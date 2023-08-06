package eqassets

import (
	"encoding/json"
	"eq-expansion-switcher/internal/config"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/labstack/gommon/log"
	cp "github.com/otiai10/copy"
	"github.com/skratchdot/open-golang/open"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	c := config.Get()

	for _, e := range e.GetExpansionFiles(strconv.Itoa(id)) {
		fmt.Println("Checking for files to be deleted in expansion:", e.Expansion.Name)
		for _, f := range e.Files {
			if strings.Contains(f, ".s3d") || strings.Contains(f, ".eqg") {
				newFile := strings.Split(f, string(filepath.Separator))
				newFile = append(newFile[:0], newFile[2:]...)
				base := strings.Join(newFile, string(filepath.Separator))
				base = strings.ReplaceAll(base, ".s3d", "")
				base = strings.ReplaceAll(base, ".eqg", "")
				_ = filepath.Walk(c.EqDir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}

					if strings.Contains(path, base) && !strings.Contains(path, "maps") {
						fmt.Println("--- Removing file:", path)

						// remove file
						err := os.Remove(path)
						if err != nil {
							fmt.Println(err)
						}
					}

					return nil
				})
			}
		}

		fmt.Println("Patching files for expansion:", e.Expansion.Name)
		for _, f := range e.Files {
			newFile := strings.Split(f, string(filepath.Separator))
			newFile = append(newFile[:0], newFile[2:]...)
			destination := filepath.Join(c.EqDir, strings.Join(newFile, string(filepath.Separator)))

			fmt.Printf("--- Copying file %v to %v\n", f, destination)

			// copy source to destination
			err := cp.Copy(f, destination)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (e *EqAssets) DumpPatchFilesForExpansion(id int) {
	tmpdir := filepath.Join(os.TempDir(), "patch-dir", time.Now().Format("2006-01-02-15-04-05"))

	err := os.MkdirAll(tmpdir, 0755)
	if err != nil {
		log.Info(err)
	}

	for _, e := range e.GetExpansionFiles(strconv.Itoa(id)) {
		fmt.Println("Checking for files to be deleted in expansion:", e.Expansion.Name)
		for _, f := range e.Files {
			if strings.Contains(f, ".s3d") || strings.Contains(f, ".eqg") {
				newFile := strings.Split(f, string(filepath.Separator))
				newFile = append(newFile[:0], newFile[2:]...)
				base := strings.Join(newFile, string(filepath.Separator))
				base = strings.ReplaceAll(base, ".s3d", "")
				base = strings.ReplaceAll(base, ".eqg", "")
				_ = filepath.Walk(tmpdir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}

					if strings.Contains(path, base) && !strings.Contains(path, "maps") {
						fmt.Println("--- Removing file:", path)

						// remove file
						err := os.Remove(path)
						if err != nil {
							fmt.Println(err)
						}
					}

					return nil
				})
			}
		}

		fmt.Println("Patching files for expansion:", e.Expansion.Name)
		for _, f := range e.Files {
			newFile := strings.Split(f, string(filepath.Separator))
			newFile = append(newFile[:0], newFile[2:]...)
			destination := filepath.Join(tmpdir, strings.Join(newFile, string(filepath.Separator)))

			fmt.Printf("--- Copying file %v to %v\n", f, destination)

			// copy source to destination
			err := cp.Copy(f, destination)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	_ = open.Run(tmpdir)
}
