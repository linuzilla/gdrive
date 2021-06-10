package built_in

import (
	"fmt"
	"github.com/linuzilla/gdrive/services/dot_gdrive_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"os"
)

type ChdirCommand struct {
	DotService  dot_gdrive_service.DotGDriveService `inject:"*"`
	previousDir string
}

var _ cmdline_service.CommandInterface = (*ChdirCommand)(nil)

func (ChdirCommand) Command() string {
	return `lcd`
}

func (cmd *ChdirCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: lcd <path> (use - to clear password)")
	} else {
		currentDir, _ := os.Getwd()

		if args[0] == `-` {
			if err := os.Chdir(cmd.previousDir); err != nil {
				fmt.Println(err)
			} else {
				cmd.previousDir = currentDir
				cmd.DotService.LoadEnvironment()
			}
		} else {
			targetDir := utils.StripQuotedString(args[0])

			if targetDir != currentDir {
				if err := os.Chdir(targetDir); err != nil {
					fmt.Println(err)
				} else {
					cmd.previousDir = currentDir
					cmd.DotService.LoadEnvironment()
				}
			}
		}
	}
	return 0
}
