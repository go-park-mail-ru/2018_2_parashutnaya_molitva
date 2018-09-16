/*
 * Пакет для singleton логированиия ошибок и сообщений
 */
package singletoneLogger

import (
	"log"
	"sync"

	"github.com/fatih/color"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
)

func init() {
	colorFunc = map[string]func(s string, a ...interface{}) string{
		"red":   color.New(color.FgRed).SprintfFunc(),
		"green": color.New(color.FgGreen).SprintfFunc(),
	}

	once.Do(func() {
		instance = initLogger(configInstance)
	})
}

var (
	colorFunc        map[string]func(s string, a ...interface{}) string // мап для определения функции по имени
	configInstance   = &config.JsonConfig{}
	loggerInfoStruct = newLoggerInfo() // структура со всеми данными из логгера
	instance         *singletonLogger  // инстанс синглтона
	once             sync.Once         // Магия для реализации singleton
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
			singletonLogger.logger.Println(loggerInfoStruct.errColor("%+v\n", err))
		case message = <-singletonLogger.messageChan:
			singletonLogger.logger.Println(loggerInfoStruct.msgColor(message))
		}
	}
}

func initLogger(cfg config.Config) *singletonLogger {
	err := loggerInfoStruct.readFromConfig(configFilename, cfg)

	if err != nil {
		log.Println("While reading config", err.Error()) // если не удалось, то остаются значения по умолчанию
	}

	errorChan := make(chan error, loggerInfoStruct.buffSize)
	messageChan := make(chan string, loggerInfoStruct.buffSize)
	errorLogger := singletonLogger{
		log.New(loggerInfoStruct.out, "", log.LstdFlags),
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
