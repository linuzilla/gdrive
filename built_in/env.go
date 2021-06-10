package built_in

import (
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/go-logger"
	"os"
	"runtime"
)

type EnvCommand struct {
	EnvironmentService environment_service.EnvironmentService `inject:"*"`
	CryptoFactory      crypto_service.CryptoFactory           `inject:"*"`
	DriveAPI           intf.DriveAPI                          `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*EnvCommand)(nil)

func (EnvCommand) Command() string {
	return `env`
}

func (cmd *EnvCommand) Execute(args ...string) int {
	cryptoService := cmd.CryptoFactory.GetInstance()

	cmd.DriveAPI.Version()
	fmt.Println()

	fmt.Printf("Running on: %s   (Go version: %s)\n", runtime.GOOS, runtime.Version())
	if currentDir, err := os.Getwd(); err != nil {
		logger.Error(err)
	} else {
		fmt.Printf("Current Directory: [ %s ]\n", currentDir)
	}

	fmt.Printf("Root Folder Id: %s\n", cmd.EnvironmentService.RootFolder())
	fmt.Printf("Folder Id: %s\n", cmd.DriveAPI.CurrentFolder())
	fmt.Printf("Encryption: %s\n", constants.OnOffMap[cryptoService.IsEnabled()])
	codec := cmd.EnvironmentService.GetCodec()
	fmt.Printf("Encoder: %s\n", codec.Encoder)
	fmt.Printf("Decoder: %s\n", codec.Decoder)
	if credentialEmail, err := cmd.DriveAPI.GetCredentialEmail(); err == nil {
		fmt.Printf("Credential Email: %s\n", credentialEmail)
	}
	fmt.Printf("Use Domain Admin Access: %s\n", constants.OnOffMap[cmd.DriveAPI.GetDomainAdminAccess()])

	if impersonate := cmd.DriveAPI.GetImpersonate(); impersonate != `` {
		fmt.Printf("Impersonate: %s\n", impersonate)
	}

	return 0
}
