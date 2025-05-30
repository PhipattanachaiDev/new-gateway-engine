package main

import (
	// "ezviewlite/controllers/v1"
	"fmt"
	"naturelink/configs"
	"naturelink/internal/server"
	log "naturelink/internal/services/log_services"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	wg             sync.WaitGroup
	tcpListener    net.Listener
	serverRunning  bool
	mutex          sync.Mutex
	stopServerChan = make(chan struct{})
)

func main() {
	err := configs.LoadEnv()
	if err != nil {
		fmt.Printf("failed to load .env: %v", err)
		// os.Exit(1)
	}

	logConfig, err := configs.GetServiceLogPath()
	if err != nil {
		log.ServicesError(fmt.Errorf("failed to get log config: %v", err))
		return
	}

	logFile := filepath.Join(logConfig.LogPath, logConfig.LogFile)

	loggerSizeConfig, err := configs.GetSizeLog()
	if err != nil {
		log.ServicesError(fmt.Errorf("failed to get logger size config: %v", err))
		return
	}

	logLevelInt, err := strconv.Atoi(loggerSizeConfig.Level)
	if err != nil {
		log.ServicesError(fmt.Errorf("failed to get logger level config: %v", err))
		return
	}

	logLevel := log.GetLevelNameByInt(logLevelInt)

	log.InitLogger(logFile, logConfig.ServiceName, logLevel)

	defer log.CloseLogger()

	log.ServicesInfo("Logger initialized Started.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start TCP Server directly
	startTCPServer()

	<-quit
	stopServerChan <- struct{}{}
	wg.Wait()
	log.ServicesInfo("Application shutdown complete")
}

func startTCPServer() {
	mutex.Lock()
	if serverRunning {
		mutex.Unlock()
		log.ServicesInfo("Server already running.")
		return
	}
	serverRunning = true
	mutex.Unlock()

	go func() {
		for {
			select {
			case <-stopServerChan:
				mutex.Lock()
				if tcpListener != nil {
					tcpListener.Close()
					tcpListener = nil
				}
				serverRunning = false
				mutex.Unlock()
				log.ServicesInfo("TCP Server stopped via signal.")
				return
			default:

				publicIP, port, err := configs.GetFromEnv()
				if err != nil {
					log.ServicesError(fmt.Errorf("failed to retrieve Gateway details: %v", err))
					time.Sleep(5 * time.Second)
					continue
				}

				tcpListener, err = net.Listen("tcp", publicIP+":"+port)
				if err != nil {
					log.ServicesError(fmt.Errorf("failed to start server: %v", err))
					time.Sleep(5 * time.Second)
					continue
				}

				log.ServicesInfo(fmt.Sprintf("Server started on IP: %s, Port: %s", publicIP, port))

				for {
					conn, err := tcpListener.Accept()
					if err != nil {
						log.ServicesError(fmt.Errorf("error accepting connection: %v", err))
						break
					}

					wg.Add(1)
					go func(c net.Conn) {
						defer wg.Done()
						server.HandleConnection(c)
					}(conn)
				}
			}
		}
	}()
}
