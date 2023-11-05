package main

import (
	"context"
	"eq-expansion-switcher/internal/config"
	"eq-expansion-switcher/internal/eqassets"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"path/filepath"
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
	return &App{
		ctx:    context.Background(),
		assets: eqassets.NewEqAssets(),
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

func (a *App) PatchFilesForExpansion(expansionId int) {
	if !a.validateEqDirExists() {
		return
	}

	a.assets.PatchFilesForExpansion(expansionId)

	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "Patch",
		Message: "Patch Complete!",
	})
}

func (a *App) DumpPatchFilesForExpansion(expansionId int) {
	if !a.validateEqDirExists() {
		return
	}

	a.assets.DumpPatchFilesForExpansion(expansionId)
}

func (a *App) CloseApp() {
	dialog, _ := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   "Quit?",
		Message: "Are you sure you want to quit?",
	})

	if dialog == "Yes" {
		os.Exit(0)
	}
}

type AppInitializationCheck struct {
	IsInitialized bool `json:"is_initialized"`
}

func (a *App) AppInitializationCheck() AppInitializationCheck {
	// get count of files in a.eqassets.Basepath
	// if count > 0, then we're initialized
	count := 0
	err := filepath.Walk(a.assets.Basepath(), func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		count++
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("count:", count)

	r := AppInitializationCheck{
		IsInitialized: count > 0,
	}

	return r
}

func (a *App) AppInitialization() error {
	err := a.assets.InitPatchFiles()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) GetAssetBasepath() string {
	return a.assets.Basepath()
}

func (a *App) validateEqDirExists() bool {
	if a.config.EqDir == "" {
		_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.ErrorDialog,
			Title:   "Error",
			Message: "EverQuest Client Directory not set.",
		})
		return false
	}

	return true
}
