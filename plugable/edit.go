package plugable

import (
	"bufio"
	"fmt"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/services/commons"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type EditCommand struct {
	DriveAPI      intf.DriveAPI                          `inject:"*"`
	CryptoFactory crypto_service.CryptoFactory           `inject:"*"`
	EnvSvc        environment_service.EnvironmentService `inject:"*"`
}

var _ cmdline_service.CommandInterface = (*EditCommand)(nil)

func (EditCommand) Command() string {
	return `edit`
}

func (cmd *EditCommand) Execute(args ...string) int {
	if len(args) != 1 {
		fmt.Println("usage: edit <fileId>")
	} else {
		cryptoSvc := cmd.CryptoFactory.GetInstance()
		fileId := utils.StripQuotedString(args[0])

		tempFile, err := ioutil.TempFile(cmd.EnvSvc.GetWorkingDirectory(), `edit`)
		if err != nil {
			fmt.Println(err)
		}

		file, fileName, md5sum, err := commons.PipeDownloadFile(cryptoSvc, cmd.DriveAPI, fileId, tempFile)
		tempFile.Close()

		defer func() {
			fmt.Printf("remove temp file %s\n", tempFile.Name())
			os.Remove(tempFile.Name())
		}()

		if err != nil {
			fmt.Println(err)
		} else {
			parentFolder := ``
			inTrash := false

			if file.Parents != nil && len(file.Parents) > 0 {
				parentFolder = file.Parents[0]

				if parentFileName, _, err := cmd.DriveAPI.FileInfo(parentFolder, func(file *drive.File) string {
					return commons.FileNameDecoder(file, cryptoSvc)
				}); err != nil {
					fmt.Println(err)
					return 0
				} else if parentFileName == constants.TrashFolderName {
					inTrash = true
					fmt.Printf("WARNING!! file already in Trash folder, continue anyway? [yes/No]: ")
					reader := bufio.NewReader(os.Stdin)
					text, _ := reader.ReadString('\n')

					if strings.TrimSpace(text) != `yes` {
						return 0
					}
				}
			} else {
				fmt.Println("parent folder not found")
				return 0
			}

			fmt.Printf("Edit file: %s, md5sum: %s\n", fileName, md5sum)

			execCmd := exec.Command(cmd.EnvSvc.GetEditor(), tempFile.Name())
			execCmd.Stdout = os.Stdout
			execCmd.Stderr = os.Stderr
			execCmd.Stdin = os.Stdin

			if err := execCmd.Run(); err != nil {
				fmt.Println(err)
			} else {
				if sum, err := utils.Md5sum(tempFile.Name()); err != nil {
					fmt.Println(err)
				} else if sum != md5sum {
					if inTrash {
						fmt.Printf("file %s change, but ignored (in trash)\n", fileName)
					} else {
						fmt.Printf("file changed (%s -> %s)\n", md5sum, sum)

						if err := cmd.DriveAPI.EncryptedUpload(cryptoSvc, tempFile.Name(), parentFolder, fileName, sum); err != nil {
							fmt.Println(err)
						} else {
							cmd.DriveAPI.DropFile(file.Id)
						}
					}
				} else {
					fmt.Printf("file %s not change\n", fileName)
				}
			}
		}
	}
	return 0
}
