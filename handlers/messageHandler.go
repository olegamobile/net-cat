package handlers

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func ProcessMessages(messages <-chan Request) {
	for message := range messages {
		if len(message.data) != 0 {
			if message.data == "/exit" {
				CloseConnection(message.client.conn)
				continue
			}
			formattedMessage := formatMessage(message)
			LogWriter(formattedMessage, LogFile)
			BrodcastPipe <- Request{client: message.client, data: formattedMessage + "\n"}
		}
	}
}

func BroadcastMessages(broadcastedMessages <-chan Request) {
	for message := range broadcastedMessages {
		allUsers := UserList.GetAllClients()
		for _, user := range allUsers {
			if message.client != *user {
				user.conn.Write([]byte(message.data))
			}
		}
	}
}

func CreateLogFile() *os.File {

	timeStamp := time.Now().UTC()
	logFileName := fmt.Sprintf("%s.txt", timeStamp.Format("20060102-150405"))
	logFile, err := os.OpenFile(LogFileDir+"/"+logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

	if err != nil {
		log.Fatalf("Cannot create log file: %v", err)
	}

	fmt.Println("Log file " + logFileName + " is created")
	return logFile
}

func LogWriter(message string, logFile *os.File) {
	Lock.Lock()
	defer Lock.Unlock()
	MsgHistory += message + "\n"
	fmt.Println(message)
	message = removeColors(message)
	_, err := logFile.WriteString(message + "\n")
	if err != nil {
		fmt.Printf("Ooops, this cannot be added to the log file.\nAnd this is the reason why: %v\n", err)
	}
}

func formatMessage(message Request) string {
	clearedMessage := clearInput(message.data)
	timestamp := GetTimestamp()
	return fmt.Sprintf("[%v][%v]: %s", timestamp, message.client.name, clearedMessage)
}

func clearInput(message string) string {
	message = strings.ReplaceAll(message, "\033", "")
	message = strings.TrimSpace(message)
	return message
}

func removeColors(message string) string {
	message = strings.ReplaceAll(message, Green, "")
	message = strings.ReplaceAll(message, Red, "")
	message = strings.ReplaceAll(message, Reset, "")
	return message
}

func GetTimestamp() string {
	now := time.Now().UTC()
	return now.Format("2006-01-02 15:04:05")
}
