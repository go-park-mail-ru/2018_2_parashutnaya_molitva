/*
 * Пакет для singleton логированиия ошибок и сообщений
 */
package singletoneLogger

import (
	"github.com/fatih/color"
	"log"
	"os"
	"sync"
)

var (
	out      = os.Stdout // writer для логов
	buffSize = 100 // максимальный размер каналов
	red      = color.New(color.FgRed).SprintfFunc() // функция окрашивающая цвет для текста в консоли красным
	green    = color.New(color.FgGreen).SprintFunc() // функция окрашивающая цвет для текста в консоли зеленым
	instance *singletonLogger // инстанс синглтона
	once     sync.Once // Магия для реализации singleton
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
		case err = <- singletonLogger.errorChan:
			singletonLogger.logger.Println(red("%+v\n", err))
		case message = <- singletonLogger.messageChan:
			singletonLogger.logger.Println(green(message))
		}
	}
}

func initLogger() *singletonLogger {
	errorChan := make(chan error, buffSize)
	messageChan := make(chan string, buffSize)
	errorLogger := singletonLogger{
		log.New(out, "", log.LstdFlags),
		errorChan,
		messageChan}
	go errorLogger.startLogging()
	return &errorLogger
}

func LogError(err error) {
	once.Do(func() {
		instance = initLogger()
	})
	instance.errorChan <- err
}

func LogMessage(message string) {
	once.Do(func() {
		instance = initLogger()
	})
	instance.messageChan <- message
}
