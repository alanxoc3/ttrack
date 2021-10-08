package config

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetAppDir(appName, appDirVar, xdgDirVar, homeFallback string) string {
	if val, present := os.LookupEnv(appDirVar); present {
		return val
	} else if val, present := os.LookupEnv(xdgDirVar); present {
		return filepath.Join(val, appName)
	} else if usr, err := user.Current(); err == nil && usr.HomeDir != "" {
		return filepath.Join(usr.HomeDir, homeFallback, appName)
	} else {
		return "."
	}
}
