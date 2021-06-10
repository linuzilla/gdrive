package built_in

import (
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/go-cmdline"
)

type ListCommand struct {
	EnvSvc environment_service.EnvironmentService `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*ListCommand)(nil)

func (ListCommand) Command() string {
	return "list"
}

func (cmd *ListCommand) Execute(args ...string) int {
	cmd.EnvSvc.ListSyncInfos()
	return 0
}
