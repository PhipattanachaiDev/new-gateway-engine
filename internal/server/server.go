package server

import (
	"bufio"
	"fmt"
	log "naturelink/internal/services/log_services"
	parse "naturelink/internal/services/parsedata_services"
	"net"
	"sync"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.ServicesInfo(fmt.Sprintf("Client connected: %s", conn.RemoteAddr().String()))

	done := make(chan struct{})

	var wg sync.WaitGroup
	var buffer string
	wg.Add(1)

	reader := bufio.NewReader(conn)
	defer func() {
		close(done)
		wg.Wait()
	}()

	for {
		msg := make([]byte, 65536)
		n, err := reader.Read(msg)
		if err != nil {
			return
		}
		hexMsg := ""
		for _, b := range msg[:n] {
			hexMsg += fmt.Sprintf("%02X ", b)
		}
		hexMsg = hexMsg[:len(hexMsg)-1]
		buffer += hexMsg + " "

		parse.ParseData(buffer)
	}
}
