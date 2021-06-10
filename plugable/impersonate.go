package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type ImpersonateCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*ImpersonateCommand)(nil)

func (ImpersonateCommand) Command() string {
	return `impersonate`
}

func (cmd *ImpersonateCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: impersonate <email-address>")
	} else {
		cmd.DriveAPI.SetImpersonate(utils.StripQuotedString(args[0]))
	}
	return 0
}
