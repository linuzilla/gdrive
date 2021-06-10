package traversal_service

import (
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/services/ignore_service"
)

type TraversalServiceFactory interface {
	NewInstance(handler *TraversalHandler) TraversalService
}

func New() TraversalServiceFactory {
	return &traversalServiceFactorImpl{}
}

type traversalServiceFactorImpl struct {
	CryptoFactory crypto_service.CryptoFactory `inject:"*"`
	Gdrive        google_drive.API             `inject:"*"`
	IgnoreService ignore_service.IgnoreService `inject:"*"`
	//EnvSvc        environment_service.EnvironmentService `inject:"*"`
}

func (factor *traversalServiceFactorImpl) NewInstance(handler *TraversalHandler) TraversalService {
	impl := &traversalServiceImpl{
		cryptoService: factor.CryptoFactory.GetInstance(),
		gdrive:        factor.Gdrive,
		ignoreService: factor.IgnoreService,
		handler:       handler,
	}

	impl.setInterrupt()

	return impl
}
