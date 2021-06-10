// +build darwin

package loader
//
//import (
//	"github.com/linuzilla/gdrive/built_in"
//	"github.com/linuzilla/gdrive/config"
//)
//
//var modules []interface{}
//
////var dbBackend db_bolt_impl.BoltDbBackend
//
//func init() {
//	modules = []interface{}{
//		&built_in.ListCommand{},
//		&built_in.LsCommand{},
//		&built_in.PwdCommand{},
//		&built_in.ChdirCommand{},
//		&built_in.DirCommand{},
//		&built_in.DirsCommand{},
//		&built_in.InfoCommand{},
//		&built_in.UpCommand{},
//		&built_in.FolderCommand{},
//		&built_in.VersionCommand{},
//		&built_in.SetPasswordCommand{},
//		&built_in.SetCoderCommand{},
//		&built_in.EnvCommand{},
//
//		//&dbBackend,
//	}
//}
//
//func LoadStaticModules(conf *config.Config, callback func(module interface{})) {
//	for _, m := range modules {
//		callback(m)
//	}
//	//
//	//if err :=dbBackend.Initialize(conf.Database.File, conf.Database.Log, func(expose interface{}) {
//	//	callback(expose)
//	//}); err != nil {
//	//	log.Fatal(err)
//	//}
//}
