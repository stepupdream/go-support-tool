//go:build linux

package console

import (
	"github.com/cheggaaa/pb/v3"
)

// KeyWait waits for the user to press any key.
func KeyWait(message string) {
	// Do not stop for linux.
}

// StartProgressBar starts a progress bar.
//
//goland:noinspection GoUnusedExportedFunction
func StartProgressBar(totalCount int) *pb.ProgressBar {
	bar := pb.Simple.Start(totalCount)
	bar.SetMaxWidth(80)

	return bar
}
