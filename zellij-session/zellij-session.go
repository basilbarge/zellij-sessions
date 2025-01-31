package zellijSession

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/zellijsessions/utils"
)

type ProjectDir struct {
	Info    fs.DirEntry
	AbsPath string
}

func NewProjectDir(dirEntry fs.DirEntry, absPath string) ProjectDir {
	return ProjectDir{
		Info:    dirEntry,
		AbsPath: absPath,
	}
}

type ZellijSession struct {
	Config      *Config
	Filesystem  fs.FS
	ProjectDirs []ProjectDir
}

func NewZellijSession(filesystem fs.FS) *ZellijSession {
	sessionConfig := NewConfig(filesystem, "Documents/Projects/zellij-sessions/config.json")

	return &ZellijSession{
		Config:      sessionConfig,
		Filesystem:  filesystem,
		ProjectDirs: getProjectDirs(sessionConfig.Dirs),
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
	utils.RunShellCommand("zellij", []string{"--new-session-with-layout", "~/Documents/Projects/zellij-sessions/zession-default.kdl", "-s", sessionName})
}

func (session *ZellijSession) AttachToSession(sessionName string) {
	if !sessionExists(sessionName) {
		log.Fatal(fmt.Errorf("Can not attach to session %s. Session does not exist", sessionName))
	}

	//utils.ExecCommand("zellij", []string{"attach", strings.TrimSpace(sessionName)}, *strings.NewReader(""))

	utils.RunShellCommand("zellij", []string{"attach", sessionName})

}

func sessionExists(sessionName string) bool {
	existingSessions := utils.ExecCommand("zellij", []string{"ls"}, *strings.NewReader(""))

	outputLines := strings.Split(strings.TrimSpace(existingSessions.String()), "\n")

	var existingSessionNames []string

	for _, line := range outputLines {
		existingSessionNames = append(existingSessionNames, strings.Split(line, " ")[0])
	}

	for _, existingSessionName := range existingSessionNames {
		strippedName := strings.TrimSpace(stripAnsiiColorCodes(existingSessionName))
		if strippedName == sessionName {
			return true
		}
	}

	return false
}

func stripAnsiiColorCodes(str string) string {
	ansiiColorCodeRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)

	return ansiiColorCodeRegex.ReplaceAllString(str, "")
}

func getProjectDirs(dirPaths []string) []ProjectDir {
	var entries []ProjectDir

	for _, path := range dirPaths {
		dirEntries, err := os.ReadDir(path)

		if err != nil {
			utils.LogError(fmt.Sprintf("Could not read directory %s. Failed with err %v", path, err))
		}

		for _, entry := range dirEntries {
			if !entry.IsDir() {
				continue
			}

			entries = append(entries, NewProjectDir(entry, filepath.Join(path, entry.Name())))
		}
	}

	return entries
}
