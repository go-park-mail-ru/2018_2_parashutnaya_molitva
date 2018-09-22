/*
 * Пакет для singleton логированиия ошибок и сообщений
 */
package singletoneLogger

import (
	"log"
	"sync"

	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
)

func init() {
	once.Do(func() {
		instance = initLogger()
	})
}

var (
	configInstance = config.NewLoggerConfig()
	instance       *singletonLogger // инстанс синглтона
	once           sync.Once        // Магия для реализации singleton
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
			singletonLogger.logger.Println(configInstance.ErrColor("%+v\n", err))
		case message = <-singletonLogger.messageChan:
			singletonLogger.logger.Println(configInstance.MsgColor(message))
		}
	}
}

func initLogger() *singletonLogger {
	err := configInstance.ReadConfig()

	if err != nil {
		log.Println("While reading config", err.Error()) // если не удалось, то остаются значения по умолчанию
	}

	errorChan := make(chan error, configInstance.BuffSize)
	messageChan := make(chan string, configInstance.BuffSize)
	errorLogger := singletonLogger{
		log.New(configInstance.Out, "", log.LstdFlags),
		errorChan,
		messageChan}
	go errorLogger.startLogging()
	return &errorLogger
}

// LogError - пишет в writer отдельным цветом
func LogError(err error) {
	instance.errorChan <- err
}

// LogMessage - пишет в writer отдельным цветом
func LogMessage(message string) {
	instance.messageChan <- message
}
