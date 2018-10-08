package config

// для чтения из конфига. Имя переменной совпадает с именем в конфиге
type LoggerConfig struct {
	OutDist    string // возможные варианты: file, stdout
	Filename   string // название файла, если указан file
	ColorError string // цвет ошибки
	ColorMsg   string // цвет сообщения
	BuffSize   string // максимальный размер каналов
}
