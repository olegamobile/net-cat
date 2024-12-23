package handlers

import (
	"fmt"
	"log"
	"os"
	"time"
)

var LogFile = LogFileCreate()
var MsgHistory string

func ProcessMessages(messages <-chan Request) {
	for message := range messages {
		if len(message.data) != 0 {
			// if "exit" - disconnect

			formattedMessage := formatMessage(message)
			fmt.Println(formattedMessage)
			LogWriter(formattedMessage, LogFile)

			// broadcast formatted message (separate function)
			// add formatted message to history
			// log formatted message (separtate function)

		}
	}
}

func LogFileCreate() *os.File {

	timeStamp := time.Now()
	logFileName := fmt.Sprintf("%s.txt", timeStamp.Format("20060102-150405"))
	logFile, err := os.OpenFile(LogFileDir+"/"+logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalf("Cannot create log file: %v", err)
	}

	// fmt.Println("Log file " + logFileName + " is created")
	return logFile
}

func LogWriter(formattedMessage string, logFile *os.File) {
	_, err := logFile.WriteString(formattedMessage + "\n")
	MsgHistory += formattedMessage + "\n"
	fmt.Println(formattedMessage)
	if err != nil {
		fmt.Printf("Ooops, this cannot be added to the log file.\nAnd this is the reason why: %v\n", err)
	}
}

func formatMessage(message Request) string {
	timestamp := getTimestamp()
	return fmt.Sprintf("[%v][%v]: %s", timestamp, message.client.name, message.data)
}

func getTimestamp() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
}
