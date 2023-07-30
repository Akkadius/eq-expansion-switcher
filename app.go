package main

import (
	"context"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/fs"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {

	// filepath walk ~/code/
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	_ = filepath.WalkDir(filepath.Join(homedir, "code/shared-task-server/code"), func(path string, d fs.DirEntry, err error) error {
		runtime.EventsEmit(a.ctx, "terminal-echo", path)
		return nil
	})
	if err != nil {
		return ""
	}

	// runtime.EventsEmit("terminal-echo", text)

	return fmt.Sprintf("Hello %s, It's show time!", name)
}
