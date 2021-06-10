package plugable

import (
	"fmt"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/go-cmdline"
	"google.golang.org/api/drive/v3"
	"os"
	"os/signal"
	"strconv"
)

type TeamDriveCommand struct {
	DriveAPI     intf.DriveAPI `inject:"*"`
	terminate    bool
	interruptSet bool
}

var _ cmdline_service.CommandInterface = (*TeamDriveCommand)(nil)

func (TeamDriveCommand) Command() string {
	return "team-drives"
}

func (cmd *TeamDriveCommand) setInterrupt() {
	if !cmd.interruptSet {
		cmd.interruptSet = true
		cmd.terminate = false
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt)

		go func() {
			<-signalChan
			fmt.Printf("\nReceived an interrupt, stopping reading teamdrive\n\n")
			cmd.terminate = true
		}()
	}
}

func (cmd *TeamDriveCommand) readTeamDrive(number int) {
	counter := 0

	cmd.setInterrupt()

	if err := cmd.DriveAPI.ReadTeamDrives(func(teamDrive *drive.TeamDrive) bool {
		fmt.Printf("<Team Drive> [ %s ] %s\n", teamDrive.Id, teamDrive.Name)
		counter++

		return !cmd.terminate && (number < 0 || counter < number)
	}); err != nil {
		fmt.Println(err)
	}
}

func (cmd *TeamDriveCommand) Execute(args ...string) int {

	if len(args) == 0 {
		cmd.readTeamDrive(-1)
	} else {
		if n, err := strconv.Atoi(args[0]); err == nil {
			cmd.readTeamDrive(n)
		}
	}

	return 0
}
