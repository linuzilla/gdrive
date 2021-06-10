package crypto_service

import (
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/models"
)

type CryptoFactory interface {
	GetInstance() intf.CryptoService
	Close()
	SetCoder(codec *models.Codec)
}
