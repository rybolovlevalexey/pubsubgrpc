package core

import (
    "log"
    "os"
)

func NewLogger() *log.Logger {
    return log.New(os.Stdout, "[PubSub] ", log.LstdFlags)
}
