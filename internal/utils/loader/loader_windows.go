//go:build windows

package loader

import (
	"os"

	"golang.org/x/sys/windows"
)

// allow ANSI escape char for windows
func enableVirtualTermWindows() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
