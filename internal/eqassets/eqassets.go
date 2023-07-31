package eqassets

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type EqAssets struct {
}

func NewEqAssets() *EqAssets {
	return &EqAssets{}
}

func (e *EqAssets) Init() error {
	// make sure ./files exists
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		err := os.MkdirAll("./files", 0755)
		if err != nil {
			return err
		}
	}

	fmt.Println(string(expansionJson))

	// unmarshal expansionJson to Expansion struct
	var expansions []Expansion
	err := json.Unmarshal(expansionJson, &expansions)
	if err != nil {
		return err
	}

	for _, s := range expansions {
		// make sure dir exists
		dir := filepath.Join("./files", fmt.Sprintf("%v", s.Id))
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Println("Creating dir:", dir)
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *EqAssets) GetAsset(name string) ([]byte, error) {
	return nil, nil
}
