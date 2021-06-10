package built_in

import (
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/go-cmdline"
)

type DirsCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DirsCommand)(nil)

func (DirsCommand) Command() string {
	return "dirs"
}

func (cmd *DirsCommand) Execute(args ...string) int {
	cmd.DriveAPI.Dirs()
	return 0
}
