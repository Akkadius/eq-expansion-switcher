package eqassets

import (
	"encoding/json"
	"eq-expansion-switcher/internal/config"
	"fmt"
	"github.com/gosimple/slug"
	cp "github.com/otiai10/copy"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func (e *EqAssets) PatchFilesForExpansion(id int) {
	config := config.Get()

	for _, e := range e.GetExpansionFiles(strconv.Itoa(id)) {
		fmt.Println("Patching files for expansion:", e.Expansion.Name)
		for _, f := range e.Files {
			fmt.Println("Patching file:", f)

			newFile := strings.Split(f, string(filepath.Separator))
			// remove first 2 elements from slice
			newFile = append(newFile[:0], newFile[2:]...)

			// copy file to eq dir
			destination := filepath.Join(config.EqDir, strings.Join(newFile, string(filepath.Separator)))

			basename := filepath.Base(destination)

			fmt.Println("destination:", destination)
			fmt.Println("basename:", basename)

			if strings.Contains(basename, ".s3d") {
				// check if file exists
				eqg := strings.Replace(destination, ".s3d", ".eqg", 1)
				if _, err := os.Stat(eqg); !os.IsNotExist(err) {
					fmt.Println("Removing file:", eqg)

					// remove file
					err := os.Remove(eqg)
					if err != nil {
						fmt.Println(err)
					}
				}
			}

			// copy source to destination
			err := cp.Copy(f, destination)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
