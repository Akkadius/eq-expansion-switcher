package main

import (
	"context"
	"eq-expansion-switcher/internal/eqassets"
	"fmt"
	"log"
)

// App struct
type App struct {
	ctx    context.Context
	assets *eqassets.EqAssets
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
	fmt.Printf("GetExpansionFiles %s\n", exansionId)

	return a.assets.GetExpansionFiles(exansionId)
}
