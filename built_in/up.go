package built_in

import (
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/go-cmdline"
)

type UpCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*UpCommand)(nil)

func (UpCommand) Command() string {
	return `up`
}

func (cmd *UpCommand) Execute(args ...string) int {
	cmd.DriveAPI.Popd()
	return 0
}
