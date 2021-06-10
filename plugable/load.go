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

type LoadCommand struct {
	DotService    dot_gdrive_service.DotGDriveService    `inject:"*"`
	EnvSvc        environment_service.EnvironmentService `inject:"*"`
	DriveAPI      intf.DriveAPI                          `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory           `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*LoadCommand)(nil)

func (LoadCommand) Command() string {
	return `load`
}

func (cmd *LoadCommand) Execute(args ...string) int {
	//if topLevel, err := cmd.EnvSvc.FetchTopLevel(); err != nil {
	//	logger.Error(err)
	//} else {
	//	for i, entry := range topLevel {
	//		fmt.Printf("%3d. %s (%s)\n", i+1, entry.Id, entry.FolderId)
	//	}
	//}

	if len(args) != 1 {
		fmt.Println("usage: load <path>")
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

		if err := os.Chdir(targetDir); err != nil {
			fmt.Println(err)
		} else {
			environment := cmd.DotService.LoadEnvironment()

			if environment != nil {
				cmd.DriveAPI.Pushd()
				cmd.DriveAPI.ChangeFolder(environment.FolderId)
				cmd.CryptoFactory.GetInstance().SetPassword(environment.Password)
			} else {
				fmt.Printf("%s: record not found\n", targetDir)
			}
		}
	}
	return 0
}
