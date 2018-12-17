package fileStorage

import (
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/config"
	"github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/internal/pkg/singletoneLogger"
	"net/http"
	"path/filepath"
)

const (
	configFilename = "file_storage.json"
)

var (
	jsonConfigReader  = config.JsonConfigReader{}
	fileStorageConfig         fileStorageConfigData
	StoragePath      string
	StorageHandler  http.Handler
)

type fileStorageConfigData struct {
	Path string
}

func init() {
	jsonConfigReader = config.JsonConfigReader{}
	err := jsonConfigReader.Read(configFilename, &fileStorageConfig)
	if err != nil {
		singletoneLogger.LogError(err)
	}
	StoragePath = filepath.Join(fileStorageConfig.Path, "storage")
	StorageHandler = http.StripPrefix("/storage/", http.FileServer(http.Dir(StoragePath)))
}
