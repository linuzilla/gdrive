package gdrive

import (
	"fmt"
	"github.com/kesselborn/go-getopt"
	"github.com/linuzilla/gdrive/backends/db_sqlite3"
	"github.com/linuzilla/gdrive/config"
	"github.com/linuzilla/gdrive/constants"
	"github.com/linuzilla/gdrive/intf"
	"github.com/linuzilla/gdrive/loader"
	"github.com/linuzilla/gdrive/services/crypto_service"
	"github.com/linuzilla/gdrive/services/database_factory"
	"github.com/linuzilla/gdrive/services/database_service"
	"github.com/linuzilla/gdrive/services/dot_gdrive_service"
	"github.com/linuzilla/gdrive/services/drive_api"
	"github.com/linuzilla/gdrive/services/environment_service"
	"github.com/linuzilla/gdrive/services/execution_service"
	"github.com/linuzilla/gdrive/services/google_drive"
	"github.com/linuzilla/gdrive/services/ignore_service"
	"github.com/linuzilla/gdrive/services/interrupt_service"
	"github.com/linuzilla/gdrive/services/progress_service"
	"github.com/linuzilla/gdrive/services/sync_service"
	"github.com/linuzilla/gdrive/services/traversal_service"
	"github.com/linuzilla/gdrive/utils"
	"github.com/linuzilla/go-cmdline"
	"github.com/linuzilla/go-logger"
	"github.com/linuzilla/summer"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/debug"
	"syscall"
)

func Start() {
	defaultConfigFile := `application.yml`

	if runtime.GOOS == `windows` {
		userProfile := os.Getenv("USERPROFILE")
		defaultConfigFile = path.Join(userProfile, `etc`, defaultConfigFile)
	}

	optionDefinition := getopt.Options{
		Description: constants.VERSION,
		Definitions: getopt.Definitions{
			{"debug|d|DEBUG", "debug mode", getopt.Optional | getopt.Flag, false},
			{"config|c|" + constants.ConfigFileEnv, "config file", getopt.IsConfigFile | getopt.ExampleIsDefault, defaultConfigFile},
			{"pass through", "pass through arguments", getopt.IsPassThrough | getopt.Optional, ""},
		},
	}

	options, _, passThrough, e := optionDefinition.ParseCommandLine()

	help, wantsHelp := options["help"]
	exitCode := 0

	if e != nil || wantsHelp {
		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Println("**** Error: ", e.Error(), "\n", optionDefinition.Help())
			exitCode = e.ErrorCode
		}
	} else {
		startProgram(options["config"].String, passThrough)
	}
	os.Exit(exitCode)
}

func startProgram(configYaml string, passThrough []string) {
	fmt.Println(constants.VERSION + "  (Go version: " + runtime.Version() + ")\n")
	fmt.Println("Running on: " + runtime.GOOS)

	fmt.Printf("Loading config: \"%s\"\n", configYaml)
	conf, err := config.New(configYaml)

	if err != nil {
		log.Printf("error:  #%v ", err)
		os.Exit(1)
	}

	logger.SetLogger(logger.New())
	logger.SetLevel(conf.Application.LogLevel)

	applicationContext := summer.New()

	applicationContext.Debug(false)

	applicationContext.Add(applicationContext)
	applicationContext.Add(conf)
	applicationContext.Add(environment_service.New())
	applicationContext.Add(interrupt_service.New())
	applicationContext.Add(progress_service.NewFactory())
	applicationContext.Add(google_drive.New())
	applicationContext.Add(dot_gdrive_service.New())
	applicationContext.Add(drive_api.New())
	applicationContext.Add(sync_service.New())
	applicationContext.Add(ignore_service.New())
	applicationContext.Add(execution_service.New())
	applicationContext.Add(traversal_service.New())
	applicationContext.Add(database_factory.New())
	cryptoFactory := crypto_service.New()
	applicationContext.Add(cryptoFactory)
	//
	//logger.Notice(utils.Detail(cryptoService))
	//logger.Notice(cryptoService)

	//logger.Notice(utils.Detail(&conf.Codec))

	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			fmt.Printf(">> panic: %v\n", r)
		}
		cryptoFactory.Close()
	}()

	loader.LoadStaticModules(conf, func(module interface{}) {
		applicationContext.Add(module)
	})

	if runtime.GOOS == `linux` || runtime.GOOS == `darwin` {
		if len(conf.Plugin.DatabaseBackend) > 0 {
			fmt.Printf("Login database backends from: \"%s\"\n", conf.Plugin.DatabaseBackend)
			fmt.Print("Plugins:")

			applicationContext.LoadPlugins(conf.Plugin.DatabaseBackend, func(baseName string, file string, module interface{}, err error) {
				fmt.Printf(` "%s"`, file)

				if err != nil {
					log.Printf("error:  #%v ", err)
				} else {
					if backend, ok := module.(intf.DatabaseBackend); ok {
						if err := backend.Initialize(conf.Database.File, conf.Database.Log, func(expose interface{}) {
							applicationContext.Add(expose)
						}); err != nil {
							log.Fatalf("error:  #%v ", err)
						}
					}
				}
			})
			fmt.Print("\n\n")
		} else {
			fmt.Println("Loading default database backend: SQLITE3")

			databaseBackend := db_sqlite3.New()
			applicationContext.Add(databaseBackend)
			applicationContext.Add(database_service.New())

			if err := databaseBackend.Initialize(conf.Database.File, conf.Database.Log, func(expose interface{}) {
				applicationContext.Add(expose)
			}); err != nil {
				log.Fatalf("error:  #%v", err)
			}
		}

		if len(conf.Plugin.Commands) > 0 {
			fmt.Printf("Login plugins from: \"%s\"\n", conf.Plugin.Commands)
			fmt.Print("Plugins:")

			applicationContext.LoadPlugins(conf.Plugin.Commands, func(baseName string, file string, module interface{}, err error) {
				fmt.Printf(` "%s"`, file)

				if err != nil {
					log.Printf("error:  #%v ", err)
				}
			})
			fmt.Print("\n\n")
		}
	} else {
		//db_sqlite3.New()
		fmt.Println("Loading default database backend: SQLITE3")

		databaseBackend := db_sqlite3.New()
		applicationContext.Add(databaseBackend)
		applicationContext.Add(database_service.New())

		if err := databaseBackend.Initialize(conf.Database.File, conf.Database.Log, func(expose interface{}) {
			applicationContext.Add(expose)
		}); err != nil {
			log.Fatalf("error:  #%v", err)
		}
	}

	done := applicationContext.Autowiring(func(err error) {
		if err != nil {
			utils.Println(os.Args[0])
			utils.Println(syscall.Getwd())
			utils.Println("Failed to auto wiring.")
			log.Fatalf("Error: %v\n", err)
		}
	})

	if err := <-done; err == nil {
		commandLineService := cmdline_service.New(applicationContext, "GoogleDrive")

		if len(passThrough) == 0 {
			commandLineService.Execute()
		} else {
			for _, cmd := range passThrough {
				commandLineService.RunCommand(cmd)
			}
		}
	} else {
		fmt.Println("failed")
	}
}
