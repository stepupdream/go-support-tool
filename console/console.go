//go:build !linux

package console

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"
)

// KeyWait waits for the user to press any key.
func KeyWait(message string) {
	fmt.Println(message)
	bufio.NewScanner(os.Stdin).Scan()
}

// StartProgressBar starts a progress bar.
//
//goland:noinspection GoUnusedExportedFunction
func StartProgressBar(totalCount int) *pb.ProgressBar {
	bar := pb.Simple.Start(totalCount)
	bar.SetMaxWidth(80)

	return bar
}
