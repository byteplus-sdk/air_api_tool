package helper

import (
	"github.com/3rd_rec/air_api_tool/consts"
	"os"
	"os/user"
	"path/filepath"
)

func BuildModuleDir(module string) string {
	dataDir := os.Getenv(consts.EnvNameDataDir)
	if len(dataDir) == 0 {
		usr, _ := user.Current()
		return filepath.Join(usr.HomeDir, consts.DefaultDataDirPrefix, module)
	}
	return filepath.Join(dataDir, module)
}
