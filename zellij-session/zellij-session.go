package zellijSession

import (
	"io/fs"
	"path/filepath"
	"github.com/zellijsessions/utils"
)

type ZellijSession struct {
	Config     *Config
	Filesystem fs.FS
}

func NewZellijSession(filesystem fs.FS) *ZellijSession {
	return &ZellijSession{
		Config:     NewConfig(filesystem, "Documents/Projects/zellij-sessions/config.json"),
		Filesystem: filesystem,
	}
}

func (session *ZellijSession) StartSession(directoryToStartIn string) {
	//change directory to the directory where you want your zellij session to start
	utils.RunShellCommand("cd", []string{directoryToStartIn})

	//run the zellij command to start a session with the name of the directory you're in
	utils.RunShellCommand("zellij", []string{"--session", filepath.Base(directoryToStartIn)})
}
