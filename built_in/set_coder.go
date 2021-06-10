package built_in

import (
	"fmt"
	"github.com/linuzilla/gdrive/models"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
)

type SetCoderCommand struct {
	environment_service.EnvironmentService `inject:"*"`
	crypto_service.CryptoFactory           `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*SetCoderCommand)(nil)

func (SetCoderCommand) Command() string {
	return `set-coder`
}

func (cmd *SetCoderCommand) Execute(args ...string) int {
	if len(args) != 2 {
		fmt.Println("usage: set-coder \"encoder\" \"decoder\"\n")
		fmt.Println(`Example:
    encoder: "/usr/bin/openssl aes-256-cbc -e -pbkdf2 -pass file:#{password-file}"
    decoder: "/usr/bin/openssl aes-256-cbc -d -pbkdf2 -pass file:#{password-file}"`)
	} else {
		cmd.CryptoFactory.SetCoder(&models.Codec{
			Encoder: utils.StripQuotedString(args[0]),
			Decoder: utils.StripQuotedString(args[1]),
		})
	}
	return 0
}
