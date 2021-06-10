package crypto_service

import (
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/services/execution_service"
	"github.com/linuzilla/go-logger"
)

type cryptoFactoryImpl struct {
	EnvSvc           environment_service.EnvironmentService    `inject:"*"`
	ExecutionService execution_service.CommandExecutionService `inject:"*"`
	currentInstance  intf.CryptoService
}

func New() CryptoFactory {
	return &cryptoFactoryImpl{}
}

func (factory *cryptoFactoryImpl) initializeInstance() {
	oldInstance := factory.currentInstance
	factory.currentInstance = newCryptoService(factory.EnvSvc.GetCodec(), factory.EnvSvc, factory.ExecutionService)
	if oldInstance != nil {
		logger.Notice("Load new Crypto Service")
		oldInstance.Close()
	}
}

func (factory *cryptoFactoryImpl) SetCoder(codec *models.Codec) {
	if *codec != *factory.EnvSvc.GetCodec() {
		factory.EnvSvc.SetCoder(codec)
		factory.initializeInstance()
	}
}

func (factory *cryptoFactoryImpl) GetInstance() intf.CryptoService {
	if factory.currentInstance == nil {
		factory.initializeInstance()
	}
	return factory.currentInstance
}

func (factory *cryptoFactoryImpl) Close() {
	saved := factory.currentInstance
	factory.currentInstance = nil

	if saved != nil {
		saved.Close()
	}
}
