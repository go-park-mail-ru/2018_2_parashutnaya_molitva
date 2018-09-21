package config

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
)

func init() {
	colorFunc = map[string]func(s string, a ...interface{}) string{
		"red":   color.New(color.FgRed).SprintfFunc(),
		"green": color.New(color.FgGreen).SprintfFunc(),
	}
}

const configFilename = "logger.json"

// для чтения из конфига. Имя переменной совпадает с именем в конфиге
type loggerConfigStruct struct {
	OutDist    string // возможные варианты: file, stdout
	Filename   string // название файла, если указан file
	ColorError string // цвет ошибки
	ColorMsg   string // цвет сообщения
	BuffSize   string // максимальный размер каналов

}

// LoggerConfig - структура, которая предоставляет данные для инициализации логгера.
type LoggerConfig struct {
	Out      io.Writer                               // writer для логов
	BuffSize int                                     // максимальный размер каналов
	ErrColor func(s string, a ...interface{}) string // функция окрашивающая цвет для ошибок
	MsgColor func(s string, a ...interface{}) string // функция окрашивающая цвет для сообщений
}

var (
	errIncorrectValue = errors.New("Incorrect value")
	colorFunc         map[string]func(s string, a ...interface{}) string // мап для определения функции по имени
)

// ReadConfig - дефолтный Reader, который читает JsonConfig из logger.json
func (l *LoggerConfig) ReadConfig() error {
	return l.ReadConfigSpec(configFilename, jsonConfig)
}

// ReadConfigSpec - чтение конфига из определенного файла
func (l *LoggerConfig) ReadConfigSpec(filename string, conf Config) error {
	var dst loggerConfigStruct
	err := conf.Read(filename, &dst)
	if err != nil {
		return err
	}

	switch dst.OutDist {
	case "file":
		// Так как logger синглтон, то можно не закрывать файл.
		f, err := os.Create(dst.Filename)
		if err != nil {
			return err
		}
		l.Out = bufio.NewWriter(f)

	case "stdout":
		l.Out = os.Stdout
	default:
		// логируем в stdout, если логгер не создан еще
		log.Println(errIncorrectValue.Error(), "outDst:", dst.OutDist)
	}

	fun, ok := colorFunc[dst.ColorError]
	if !ok {
		log.Println(errIncorrectValue.Error(), "colorError:", dst.ColorError)
	} else {
		l.ErrColor = fun
	}

	fun, ok = colorFunc[dst.ColorMsg]
	if !ok {
		log.Println(errIncorrectValue.Error(), "colorError:", dst.ColorError)
	} else {
		l.MsgColor = fun
	}

	val, err := strconv.Atoi(dst.BuffSize)
	if err != nil {
		log.Println(errIncorrectValue.Error(), "bufSize:", dst.BuffSize)
	} else {
		l.BuffSize = val
	}

	return nil
}

// NewLoggerConfig - создает LoggerConfig с значениями по умолчанию
func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Out:      os.Stdout,
		BuffSize: 100,
		ErrColor: color.New(color.FgRed).SprintfFunc(),
		MsgColor: color.New(color.FgGreen).SprintfFunc(),
	}
}
