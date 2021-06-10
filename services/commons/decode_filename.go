package commons

import (
	"github.com/linuzilla/gdrive/intf"
	"google.golang.org/api/drive/v3"
)

func FileNameDecoder(file *drive.File, cryptoSvc intf.CryptoService) string {
	if cryptoSvc.IsEnabled() {
		fileName, _, _ := cryptoSvc.DecryptFileNameWithMd5(file)
		return fileName
	} else {
		return file.Name
	}
}
