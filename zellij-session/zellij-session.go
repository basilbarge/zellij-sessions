package zellijSession

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

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
	sessionName := filepath.Base(directoryToStartIn)

	if sessionExists(sessionName) {
		session.AttachToSession(sessionName)
		return
	}

	//change directory to the directory where you want your zellij session to start
	utils.RunShellCommand("cd", []string{directoryToStartIn})

	//run the zellij command to start a session with the name of the directory you're in
	utils.RunShellCommand("zellij", []string{"--session", sessionName})
}

func (session *ZellijSession) AttachToSession(sessionName string) {
	if !sessionExists(sessionName) {
		log.Fatal(fmt.Errorf("Can not attach to session %s. Session does not exist", sessionName))
	}

	utils.ExecCommand("zellij", []string{"attach", sessionName}, *strings.NewReader(""))

}

func sessionExists(sessionName string) bool {
	fmt.Println("Connecting to ", sessionName)
	existingSessions := utils.ExecCommand("zellij", []string{"ls"}, *strings.NewReader(""))

	outputLines := strings.Split(strings.TrimSpace(existingSessions.String()), "\n")

	var existingSessionNames []string

	for _, line := range outputLines {
		existingSessionNames = append(existingSessionNames, strings.Split(line, " ")[0])
	}

	for _, existingSessionName := range existingSessionNames {
		fmt.Println("Comparing ",sessionName, " to ", existingSessionName)
		if existingSessionName == sessionName {
			return true
		}
	}

	return false
}
