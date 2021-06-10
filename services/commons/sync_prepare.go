package commons

//func PrepareBeforeSync(
//	api google_drive.API,
//	databaseFactory database_factory.DatabaseFactory,
//	cryptoService crypto_service.CryptoService,
//	folderName string,
//	folderId string,
//	callback func(connection intf.DatabaseBackendConnection, service *drive.Service, folderId, workingDir string)) {
//
//	infoDirectory := folderName + `/` + constants.GoogleDriveFolder
//
//	if fileInfo, err := os.Stat(infoDirectory); os.IsNotExist(err) {
//		if err := os.Mkdir(infoDirectory, 0700); err != nil {
//			panic(err.Error())
//		} else {
//			fmt.Printf("create directory: %s\n", infoDirectory)
//		}
//	} else if !fileInfo.IsDir() {
//		panic(fmt.Sprintf("%s : should be a directory", infoDirectory))
//	}
//
//	databaseFile := infoDirectory + `/` + constants.DatabaseFile
//	databaseBackend := databaseFactory.NewDatabase(databaseFile)
//
//	if err := databaseBackend.ConnectionEstablish(func(connection intf.DatabaseBackendConnection) error {
//		var conf models.GoogleDriveConfig
//
//		if err := connection.ReadConfig(infoDirectory, &conf); err != nil {
//
//			if folderId == `` {
//				return fmt.Errorf("ERROR!!! Specify a folder-id to upload to")
//			}
//
//			conf.Id = infoDirectory
//			conf.FolderId = folderId
//
//			if cryptoService.IsEnabled() {
//				conf.Password = sql.NullString{
//					String: base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(cryptoService.GetPassword())),
//					Valid:  true,
//				}
//			} else {
//				conf.Password = sql.NullString{Valid: false}
//			}
//
//			conf.TrashFolderId = sql.NullString{Valid: false}
//			connection.SaveConfig(&conf)
//		} else {
//			//fmt.Println(utils.Detail(&conf))
//			password := ``
//
//			if conf.Password.Valid {
//				if decodeString, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(conf.Password.String); err != nil {
//					fmt.Println(err)
//				} else {
//					password = string(decodeString)
//				}
//			}
//
//			ask := false
//			usePreviousPassword := false
//
//			if folderId == `` {
//				folderId = conf.FolderId
//
//				fmt.Printf("Upload to [ %s ]\n", folderId)
//			} else if folderId != conf.FolderId {
//				return fmt.Errorf("ERROR! upload folder was %s (vs %s)", conf.FolderId, folderId)
//			}
//
//			if cryptoService.IsEnabled() {
//				if password == `` {
//					fmt.Print("WARNING!! Folder was NOT encrypted with passwords, encrypt new files [yes/N] :")
//					ask = true
//				} else if password != cryptoService.GetPassword() {
//					return fmt.Errorf("ERROR!! Folder [ %s ] encrypted with different password", folderName)
//				} else {
//					conf.Password = sql.NullString{
//						String: base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(cryptoService.GetPassword())),
//						Valid:  true,
//					}
//					connection.SaveConfig(&conf)
//				}
//			} else {
//				if password != `` {
//					fmt.Print("WARNING!! Folder was encrypted with a password, use previous one [yes/N] :")
//					ask = true
//					usePreviousPassword = true
//				}
//			}
//
//			if ask {
//				reader := bufio.NewReader(os.Stdin)
//				text, _ := reader.ReadString('\n')
//				text = strings.TrimRight(text, " \r\n")
//
//				if text != `yes` {
//					return nil
//				}
//			}
//			if usePreviousPassword {
//				cryptoService.SetPassword(password)
//			}
//		}
//
//		return api.Connect(func(service *drive.Service) error {
//			callback(connection, service, folderId, infoDirectory)
//			return nil
//		})
//	}); err != nil {
//		panic(err.Error())
//	}
//}
