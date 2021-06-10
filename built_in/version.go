package built_in

import (
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/go-cmdline"
)

type VersionCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*VersionCommand)(nil)

func (VersionCommand) Command() string {
	return "version"
}

func (cmd *VersionCommand) Execute(args ...string) int {
	cmd.DriveAPI.Version()
	return 0
}
