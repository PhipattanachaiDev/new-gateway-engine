package models

type Logger struct {
    MaxSize    string
    MaxBackups string
    MaxAge     string
    Compress   bool
    Level      string
}

type ServiceLog struct {
    LogPath     string
    LogFile     string
    ServiceName string
}