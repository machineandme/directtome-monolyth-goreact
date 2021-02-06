package main

import (
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})
	logFile, err := os.OpenFile("/var/log/application_backend.log.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err == nil {
		log.Out = logFile
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
	log.Info("Logs setted")

	for true {
		log.Info("Hello hello")
		time.Sleep(time.Second * 30)
	}
}
