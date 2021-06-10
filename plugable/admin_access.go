package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/go-cmdline"
)

type DomainAdminAccessCommand struct {
	DriveAPI intf.DriveAPI `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*DomainAdminAccessCommand)(nil)

func (DomainAdminAccessCommand) PostSummerConstruct() {
}

func (DomainAdminAccessCommand) Command() string {
	return `use-domain-admin-access`
}

func (DomainAdminAccessCommand) useDomainAdminAccessUsage() {
	fmt.Println("usage: use-domain-admin-access [on|off]")
}

func (cmd *DomainAdminAccessCommand) Execute(args ...string) int {
	if len(args) != 1 {
		cmd.useDomainAdminAccessUsage()
	} else {
		switch args[0] {
		case `on`:
			cmd.DriveAPI.SetDomainAdminAccess(true)
		case `off`:
			cmd.DriveAPI.SetDomainAdminAccess(false)
		default:
			cmd.useDomainAdminAccessUsage()
		}
	}
	return 0
}
