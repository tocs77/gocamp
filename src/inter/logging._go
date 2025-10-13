package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("Hello, World!")
	log.SetPrefix("INFO: ")
	log.Println("Hello, World!")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Log with date, time and short file")

	infoLogger.Println("Info log")
	warningLogger.Println("Warning log")
	errorLogger.Println("Error log")
	fatalLogger.Println("Fatal log")
	panicLogger.Println("Panic log")

	localFileLogger, closeFile := makeFileLogger("log.txt")
	defer closeFile()
	localFileLogger.Println("File log")

}

var infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
var warningLogger = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
var errorLogger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
var fatalLogger = log.New(os.Stdout, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
var panicLogger = log.New(os.Stdout, "PANIC: ", log.Ldate|log.Ltime|log.Lshortfile)
var fileLogger = log.New(os.Stdout, "FILE: ", log.Ldate|log.Ltime|log.Lshortfile)

var makeFileLogger = func(filename string) (*log.Logger, func() error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return log.New(file, "FILE: ", log.Ldate|log.Ltime|log.Lshortfile), func() error {
		fmt.Println("Closing file")
		return file.Close()
	}
}
