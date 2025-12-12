
package logger

import (
    "github.com/rs/zerolog"
    "os"
)

var Logger zerolog.Logger


func Init() {
    Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}


func Sync() {
    
}
