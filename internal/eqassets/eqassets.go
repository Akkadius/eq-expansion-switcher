package eqassets

import (
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"os"
	"path/filepath"
	"strconv"
)

// EqAssets struct
type EqAssets struct {
	expansions []Expansion
}

// NewEqAssets creates a new EqAssets application struct
func NewEqAssets() *EqAssets {
	return &EqAssets{}
}

// Init initializes the EqAssets application struct
func (e *EqAssets) Init() error {
	// make sure ./files exists
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		err := os.MkdirAll("./files", 0755)
		if err != nil {
			return err
		}
	}

	// load expansions
	err := json.Unmarshal(expansionJson, &e.expansions)
	if err != nil {
		return err
	}

	// make sure ./files/<expansion-id>-<expansion-name> exists
	for _, s := range e.expansions {
		// make sure dir exists
		dir := filepath.Join("./files", fmt.Sprintf("%v-%v", s.Id, slug.Make(s.Name)))
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

type ExpansionFiles struct {
	Expansion Expansion `json:"expansion"`
	Files     []string  `json:"files"`
}

func (e *EqAssets) GetExpansionFiles(expansionId string) []ExpansionFiles {
	var files []ExpansionFiles

	fmt.Println("GetExpansionFiles", expansionId)

	for _, s := range e.expansions {
		// convert expansionId to int
		id, err := strconv.Atoi(expansionId)
		if err != nil {
			return files
		}

		var expansionFiles ExpansionFiles
		expansionFiles.Expansion = s

		if s.Id <= id {
			// make sure dir exists
			dir := filepath.Join("./files", fmt.Sprintf("%v-%v", s.Id, slug.Make(s.Name)))
			fmt.Println("dir:", dir)
			if _, err := os.Stat(dir); !os.IsNotExist(err) {
				_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}

					expansionFiles.Files = append(expansionFiles.Files, path)
					return nil
				})
			}
			if len(expansionFiles.Files) > 0 {
				files = append(files, expansionFiles)
			}
		}
	}

	fmt.Println(files)

	return files
}
