//go:build linux && cgo && gtk3 && !android && !server

package application

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestGTK3WindowStateBeforeCreation(t *testing.T) {
	if os.Getenv("WAILS_TEST_GTK3_WINDOW_STATE") != "1" {
		// GTK normally logs invalid native calls and returns false anyway. Make
		// criticals fatal in a fresh process so those calls fail the test.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGTK3WindowStateBeforeCreation$")
		cmd.Env = append(os.Environ(), "WAILS_TEST_GTK3_WINDOW_STATE=1", "G_DEBUG=fatal-criticals")
		if output, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("window state query before creation failed: %v\n%s", err, output)
		}
		return
	}

	window := &linuxWebviewWindow{}
	for _, test := range []struct {
		name  string
		query func() bool
	}{
		{"fullscreen", window.isFullscreen},
		{"maximised", window.isMaximised},
		{"minimised", window.isMinimised},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.query() {
				t.Fatal("window state before creation = true, want false")
			}
		})
	}
}
