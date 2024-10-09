package helper

import (
	"github.com/3rd_rec/air_api_tool/consts"
	"os"
	"path/filepath"
)

func BuildModuleDir(module string) string {
	dataDir := os.Getenv(consts.EnvNameDataDir)
	if len(dataDir) == 0 {
		dataDir = consts.DefaultDataDir
	}
	return filepath.Join(dataDir, module)
}
