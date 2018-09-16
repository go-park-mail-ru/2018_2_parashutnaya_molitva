/*
 * Пакет для singleton логированиия ошибок и сообщений
 */
package singletoneLogger

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/fatih/color"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
)

const configFilename = "logger.json"

// для чтения из конфига. Имя переменной совпадает с именем в конфиге
type loggerConfigStruct struct {
	OutDist    string // возможные варианты: file, stdout
	Filename   string // название файла, если указан file
	ColorError string // цвет ошибки
	ColorMsg   string // цвет сообщения
	BuffSize   string // максимальный размер каналов

}

type loggerInfo struct {
	out      io.Writer                               // writer для логов
	buffSize int                                     // максимальный размер каналов
	errColor func(s string, a ...interface{}) string // функция окрашивающая цвет для ошибок
	msgColor func(s string, a ...interface{}) string // функция окрашивающая цвет для сообщений
}

var (
	errIncorrectValue = errors.New("Incorrect value")
)

func (l *loggerInfo) readFromConfig(filename string, conf config.Config) error {
	var dst loggerConfigStruct
	err := conf.Read(filename, &dst)
	if err != nil {
		return err
	}

	switch dst.OutDist {
	case "file":
		// TODO: Где-то закрывать файл
		f, err := os.Create(dst.Filename)
		if err != nil {
			return err
		}
		l.out = bufio.NewWriter(f)

	case "stdout":
		l.out = os.Stdout
	default:
		// логируем в stdout, если логгер не создан еще
		log.Println(errIncorrectValue.Error(), "outDst:", dst.OutDist)
	}

	fun, ok := colorFunc[dst.ColorError]
	if !ok {
		log.Println(errIncorrectValue.Error(), "colorError:", dst.ColorError)
	} else {
		l.errColor = fun
	}

	fun, ok = colorFunc[dst.ColorMsg]
	if !ok {
		log.Println(errIncorrectValue.Error(), "colorError:", dst.ColorError)
	} else {
		l.msgColor = fun
	}

	val, err := strconv.Atoi(dst.BuffSize)
	if err != nil {
		log.Println(errIncorrectValue.Error(), "bufSize:", dst.BuffSize)
	} else {
		l.buffSize = val
	}

	return nil
}

// создает loggerInfo с значениями по умолчанию
func newLoggerInfo() *loggerInfo {
	return &loggerInfo{
		out:      os.Stdout,
		buffSize: 100,
		errColor: color.New(color.FgRed).SprintfFunc(),
		msgColor: color.New(color.FgGreen).SprintfFunc(),
	}
}

var (
	colorFunc        map[string]func(s string, a ...interface{}) string
	configInstance   = &config.JsonConfig{}
	loggerInfoStruct = newLoggerInfo()
	instance         *singletonLogger // инстанс синглтона
	once             sync.Once        // Магия для реализации singleton
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
	once.Do(func() {
		instance = initLogger(configInstance)
	})
	instance.errorChan <- err
}

// LogMessage - пишет в writer отдельным цветом
func LogMessage(message string) {
	once.Do(func() {
		instance = initLogger(configInstance)
	})
	instance.messageChan <- message
}

func init() {
	colorFunc = map[string]func(s string, a ...interface{}) string{
		"red":   color.New(color.FgRed).SprintfFunc(),
		"green": color.New(color.FgGreen).SprintfFunc(),
	}
}
