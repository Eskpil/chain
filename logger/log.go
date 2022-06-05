package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	Warn  *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func Init() {
	Info = log.New(os.Stdout, "\u001b[36mINFO\u001b[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn = log.New(os.Stdout, "\u001b[33mWARNING\u001b[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "\u001b[31mERROR\u001b[0m: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func PrintError(message string) {
	fmt.Fprintf(os.Stderr, message)
}
