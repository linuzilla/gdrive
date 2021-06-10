package loader

import (
	"github.com/linuzilla/gdrive/built_in"
	"github.com/linuzilla/gdrive/config"
	"github.com/linuzilla/gdrive/plugable"
)

var modules []interface{}

//var dbBackend intf.DatabaseBackend

func init() {
	//dbBackend = db_sqlite3.New()

	modules = []interface{}{
		&plugable.DomainAdminAccessCommand{},
		&plugable.CatCommand{},
		&plugable.CopyCommand{},
		&plugable.DeleteCommand{},
		&plugable.DownloadCommand{},
		&plugable.DropCommand{},
		&plugable.EditCommand{},
		&plugable.ImpersonateCommand{},
		&plugable.LoadCommand{},
		&plugable.LoadFolderCommand{},
		&plugable.LoadInventoryCommand{},
		&plugable.NewFolderCommand{},
		&plugable.PermissionInfoCommand{},
		&plugable.PullFolderCommand{},
		&plugable.PushFolderCommand{},
		&plugable.RemovePermissionCommand{},
		&plugable.RenameCommand{},
		&plugable.ShareFileCommand{},
		&plugable.ShareInfoCommand{},
		&plugable.TeamDriveCommand{},
		&plugable.UnshareFileCommand{},
		&plugable.UploadCommand{},

		&built_in.ListCommand{},
		&built_in.LsCommand{},
		&built_in.PwdCommand{},
		&built_in.ChdirCommand{},
		&built_in.DirCommand{},
		&built_in.DirsCommand{},
		&built_in.InfoCommand{},
		&built_in.UpCommand{},
		&built_in.FolderCommand{},
		&built_in.VersionCommand{},
		&built_in.SetPasswordCommand{},
		&built_in.SetCoderCommand{},
		&built_in.EnvCommand{},

		//&dbBackend,
	}
}

func LoadStaticModules(conf *config.Config, callback func(module interface{})) {
	for _, m := range modules {
		callback(m)
	}

	//dbBackend := db_sqlite3.New()
	////applicationContext.Add(databaseBackend)
	////applicationContext.Add(database_service.New())
	//
	//if err := dbBackend.Initialize(conf.Database.File, conf.Database.Log, func(expose interface{}) {
	//	callback(expose)
	//}); err != nil {
	//	log.Fatal(err)
	//}
}