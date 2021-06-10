package plugable

import (
	"bufio"
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/commons"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/go-cmdline"
	"os"
)

type CatCommand struct {
	DriveAPI      intf.DriveAPI                          `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory           `inject:"*"`
	EnvSvc        environment_service.EnvironmentService `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*CatCommand)(nil)

func (CatCommand) Command() string {
	return `cat`
}

func (cmd *CatCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: cat <fileId>")
	} else {
		writer := bufio.NewWriter(os.Stdout)

		if _, _, _, err := commons.PipeDownloadFile(cmd.CryptoFactory.GetInstance(), cmd.DriveAPI, args[0], writer); err != nil {
			fmt.Println(err)
		}
		writer.Flush()
	}
	return 0
}
