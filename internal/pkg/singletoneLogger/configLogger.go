package singletoneLogger

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/pkg/errors"
)

func init() {
	colorFunc = map[string]func(s string, a ...interface{}) string{
		"red":   color.New(color.FgRed).SprintfFunc(),
		"green": color.New(color.FgGreen).SprintfFunc(),
	}
}

// loggerData - структура, которая предоставляет данные для инициализации логгера.
type loggerData struct {
	Out      io.Writer                               // writer для логов
	BuffSize int                                     // максимальный размер каналов
	ErrColor func(s string, a ...interface{}) string // функция окрашивающая цвет для ошибок
	MsgColor func(s string, a ...interface{}) string // функция окрашивающая цвет для сообщений
}

var (
	errIncorrectValue = errors.New("Incorrect value")
	colorFunc         map[string]func(s string, a ...interface{}) string // мап для определения функции по имени
)

// ReadConfigSpec - чтение конфига из определенного файла
func (l *loggerData) ReadConfig() error {
	dst := &config.LoggerConfig{}
	err := config.GetConfig("logger", dst)
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

// NewLoggerData - создает loggerData с значениями по умолчанию
func NewLoggerData() *loggerData {
	return &loggerData{
		Out:      os.Stdout,
		BuffSize: 100,
		ErrColor: color.New(color.FgRed).SprintfFunc(),
		MsgColor: color.New(color.FgGreen).SprintfFunc(),
	}
}
