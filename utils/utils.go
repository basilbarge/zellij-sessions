package utils

import (
	"fmt"
	"log"
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
