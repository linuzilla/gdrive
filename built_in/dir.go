package built_in

import (
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/go-cmdline"
	"google.golang.org/api/drive/v3"
)

type DirCommand struct {
	DriveAPI      intf.DriveAPI                `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DirCommand)(nil)

func (cmd *DirCommand) PostSummerConstruct() {
}

func (DirCommand) Command() string {
	return "dir"
}

func (cmd *DirCommand) Execute(args ...string) int {
	cryptoSvc := cmd.CryptoFactory.GetInstance()

	cmd.DriveAPI.ListWithDecoder(func(file *drive.File) string {
		if cryptoSvc.IsEnabled() {
			fileName, _, _ := cryptoSvc.DecryptFileNameWithMd5(file)
			return fileName
		} else {
			return file.Name
		}
	})
	return 0
}
