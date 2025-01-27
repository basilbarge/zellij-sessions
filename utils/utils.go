package utils

import (
	"fmt"
	"log"
	"os/exec"
	"slices"
	"strings"
	"syscall"
	"unsafe"
)

func RunShellCommand(command string, args []string) {

	shellCmdString := fmt.Sprintf("%s %s\n", command, strings.Join(args, " "))

	shellCmdBytes, err := syscall.ByteSliceFromString(shellCmdString)
	if err != nil {
		log.Fatalln(err)
	}

	var eno syscall.Errno
	for _, c := range shellCmdBytes {
		_, _, eno = syscall.Syscall(syscall.SYS_IOCTL,
			0,
			syscall.TIOCSTI,
			uintptr(unsafe.Pointer(&c)),
		)
		if eno != 0 {
			log.Fatalln(eno)
		}
	}
}

func ExecCommand(cmdString string, cmdArgs []string, stdinData strings.Reader) strings.Builder {

	cmd := exec.Command(cmdString, cmdArgs...)

	cmd.Stdin = &stdinData

	var stdOutBuilder strings.Builder
	cmd.Stdout = &stdOutBuilder

	if err := cmd.Run(); err != nil {
		if (slices.Contains(cmdArgs, "ls")) {
			return stdOutBuilder
		}

		log.Fatal(fmt.Sprintf("There was a problem running command %s. %s\n", strings.Join(cmd.Args, " "), err))
	}

	return stdOutBuilder
}

func LogError(errString string) {
	log.Fatal(errString)
}
