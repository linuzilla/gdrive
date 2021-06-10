package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/dot_gdrive_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type LoadFolderCommand struct {
	DotService    dot_gdrive_service.DotGDriveService    `inject:"*"`
	EnvSvc        environment_service.EnvironmentService `inject:"*"`
	DriveAPI      intf.DriveAPI                          `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory           `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*LoadFolderCommand)(nil)

func (LoadFolderCommand) Command() string {
	return `load-folder-env`
}

func (cmd *LoadFolderCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: load-folder-env <folder-id>")
	} else {
		folderId := utils.StripQuotedString(args[0])

		environment := cmd.DotService.LoadTopLevelEnvByFolder(folderId)

		if environment != nil {
			cmd.DriveAPI.Pushd()
			cmd.DriveAPI.ChangeFolder(environment.FolderId)
			cmd.CryptoFactory.GetInstance().SetPassword(environment.Password)
		} else {
			fmt.Printf("%s: folder not in inventory\n", folderId)
		}
	}
	return 0
}
