package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/dot_gdrive_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"os"
	"strconv"
)

type LoadInventoryCommand struct {
	DotService    dot_gdrive_service.DotGDriveService    `inject:"*"`
	EnvSvc        environment_service.EnvironmentService `inject:"*"`
	DriveAPI      intf.DriveAPI                          `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory           `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*LoadInventoryCommand)(nil)

func (LoadInventoryCommand) Command() string {
	return `load-inventory`
}

func (cmd *LoadInventoryCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: load-inventory <directory|number>")
	} else {
		targetDir := utils.StripQuotedString(args[0])

		index, err := strconv.Atoi(targetDir)

		if err == nil {
			if model := cmd.EnvSvc.LoadCache(index); model != nil {
				targetDir = model.Id
			}
		}

		if args[0] == `.` {
			if currentDir, err := os.Getwd(); err == nil {
				targetDir = currentDir
			}
		}

		environment := cmd.DotService.LoadTopLevelEnvironment(targetDir)

		if environment != nil {
			cmd.DriveAPI.Pushd()
			cmd.DriveAPI.ChangeFolder(environment.FolderId)
			cmd.CryptoFactory.GetInstance().SetPassword(environment.Password)
		} else {
			fmt.Printf("%s: inventory not found\n", targetDir)
		}
	}
	return 0
}
