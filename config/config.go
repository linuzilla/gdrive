package config

import (
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/go-logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Application struct {
		Name       string `yaml:"name"`
		LogLevel   string `yaml:"log-level"`
		WorkingDir string `yaml:"working-dir"`
		Editor     string `yaml:"editor"`
	}
	GSuite struct {
		Impersonate          string `yaml:"impersonate"`
		UseDomainAdminAccess bool   `yaml:"use-domain-admin-access"`
	} `yaml:"gsuite"`
	GoogleDrive struct {
		Credential string `yaml:"credential"`
		FolderId   string `yaml:"folder-id"`
	} `yaml:"google-drive"`
	Database struct {
		File string `yaml:"file"`
		Log  bool   `yaml:"log"`
	} `yaml:"database"`
	Plugin struct {
		DatabaseBackend string `yaml:"db-backend"`
		Commands        string `yaml:"commands"`
	} `yaml:"plugin"`
	Codec models.Codec `yaml:"codec"`
}

func New(fileName string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		logger.Error("yamlFile.Get err   #%v ", err)
		return nil, err
	}
	conf := new(Config)
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		logger.Error("Unmarshal: %v", err)
		return nil, err
	}
	return conf, nil
}
