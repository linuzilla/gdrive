package commons

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/utils"
	"google.golang.org/api/drive/v3"
	"io"
)

func DownloadAnDecodeRemoteFile(api google_drive.API,
	driveSvc *drive.Service,
	cryptoService intf.CryptoService,
	remoteFile *drive.File,
	localFilePath string,
	md5sum string,
	encrypted bool) error {

	if encrypted {
		if err := api.RetrieveFile(driveSvc, remoteFile.Id, func(file *drive.File, reader io.Reader) error {
			return cryptoService.DecodeAndDownloadFromReader(localFilePath, reader)
		}); err != nil {
			return err
		}

		if md5, err := utils.Md5sum(localFilePath); err != nil {
			return nil
		} else if md5 != md5sum {
			return fmt.Errorf("checksum mismatch (%s vs %s)", md5, md5sum)
		}
		return nil
	} else {
		api.DownloadFile(driveSvc, remoteFile.Id, localFilePath, false)
		return nil
	}
}
