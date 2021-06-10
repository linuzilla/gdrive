package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type RenameCommand struct {
	DriveAPI      intf.DriveAPI                `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*RenameCommand)(nil)

func (RenameCommand) Command() string {
	return `rename`
}

func (cmd *RenameCommand) Execute(args ...string) int {
	if len(args) != 2 {
		fmt.Println("usage: rename <fileId> <fileName>")
	} else {
		cryptoSvc := cmd.CryptoFactory.GetInstance()
		sourceFileNameOrId := utils.StripQuotedString(args[0])
		newFileName := utils.StripQuotedString(args[1])

		if cryptoSvc.IsEnabled() {
			if _, file, err := cmd.DriveAPI.FileInfo(sourceFileNameOrId, nil); err != nil {
				fmt.Println(err)
			} else {
				if _, md5, err := cryptoSvc.DecryptFileNameWithMd5(file); err != nil {
					fmt.Println(err)
				} else {
					encryptFileName := cryptoSvc.EncryptFileNameWithMd5(newFileName, md5)
					cmd.DriveAPI.Rename(utils.StripQuotedString(sourceFileNameOrId), encryptFileName)
				}
			}
		} else {
			cmd.DriveAPI.Duplicate(sourceFileNameOrId, newFileName)
		}
	}
	return 0
}
