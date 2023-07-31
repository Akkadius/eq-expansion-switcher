package main

import (
	"context"
	"eq-expansion-switcher/internal/config"
	"eq-expansion-switcher/internal/eqassets"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"strconv"
)

// App struct
type App struct {
	ctx    context.Context
	assets *eqassets.EqAssets
	config config.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := eqassets.NewEqAssets()
	err := a.Init()
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		ctx:    context.Background(),
		assets: a,
		config: config.Get(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetExpansionFiles returns a greeting for the given name
func (a *App) GetExpansionFiles(exansionId string) []eqassets.ExpansionFiles {
	// Save the expansion to the config file
	id, _ := strconv.Atoi(exansionId)
	a.config.CurrentExpansion = id
	err := config.Save(a.config)
	if err != nil {
		log.Println(err.Error())
	}

	return a.assets.GetExpansionFiles(exansionId)
}

func (a *App) OpenFileDialogueEqDir() string {
	str, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:           "Find EverQuest Directory",
		ShowHiddenFiles: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Save the directory to the config file
	a.config.EqDir = str
	err = config.Save(a.config)
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Printf("OpenFileDialogueEqDir %s\n", str)
	return str
}

func (a *App) GetConfig() config.Config {
	fmt.Println(a.config)

	return a.config
}
