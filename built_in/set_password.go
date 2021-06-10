package built_in

import (
	"fmt"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type SetPasswordCommand struct {
	CryptoFactory crypto_service.CryptoFactory `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*SetPasswordCommand)(nil)

func (SetPasswordCommand) Command() string {
	return "password"
}

func (cmd *SetPasswordCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: password <password> (use - to clear password)")
	} else {
		cryptoSvc := cmd.CryptoFactory.GetInstance()

		if args[0] == `-` {
			cryptoSvc.SetPassword(``)
		} else {
			cryptoSvc.SetPassword(utils.StripQuotedString(args[0]))
		}
	}
	return 0
}
