/*
 * Пакет для singleton логированиия ошибок и сообщений
 */
package singletoneLogger

import (
	"log"
	"sync"

	"github.com/pkg/errors"
)

var (
	data     = NewLoggerData()
	instance *singletonLogger // инстанс синглтона
	once     sync.Once        // Магия для реализации singleton
)

type singletonLogger struct {
	logger      *log.Logger
	errorChan   chan error
	messageChan chan string
}

func (singletonLogger *singletonLogger) startLogging() {
	var err error
	var message string
	for {
		select {
		case err = <-singletonLogger.errorChan:
			singletonLogger.logger.Println(data.ErrColor("%+v\n", err))
		case message = <-singletonLogger.messageChan:
			singletonLogger.logger.Println(data.MsgColor(message))
		}
	}
}

func initLogger() *singletonLogger {
	err := data.ReadConfig()

	if err != nil {
		err = errors.Wrap(err, "While reading config")
		log.Printf("%+s", err)
	}

	errorChan := make(chan error, data.BuffSize)
	messageChan := make(chan string, data.BuffSize)
	errorLogger := singletonLogger{
		log.New(data.Out, "", log.LstdFlags),
		errorChan,
		messageChan}
	go errorLogger.startLogging()
	return &errorLogger
}

// LogError - пишет в writer отдельным цветом
func LogError(err error) {
	once.Do(func() {
		instance = initLogger()
	})
	instance.errorChan <- err
}

// LogMessage - пишет в writer отдельным цветом
func LogMessage(message string) {
	once.Do(func() {
		instance = initLogger()
	})
	instance.messageChan <- message
}
