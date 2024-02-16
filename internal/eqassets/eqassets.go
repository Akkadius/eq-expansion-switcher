package eqassets

import (
	"encoding/json"
	"eq-expansion-switcher/internal/config"
	"eq-expansion-switcher/internal/unzip"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/labstack/gommon/log"
	cp "github.com/otiai10/copy"
	"github.com/skratchdot/open-golang/open"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// EqAssets struct
type EqAssets struct {
	expansions []Expansion // expansions.json
	basepath   string      // base path for patch files
}

func (e *EqAssets) Basepath() string {
	return e.basepath
}

// NewEqAssets creates a new EqAssets application struct
func NewEqAssets() *EqAssets {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	basepath := filepath.Join(cacheDir, "peq-expansion-switcher")
	if _, err := os.Stat(basepath); os.IsNotExist(err) {
		err := os.MkdirAll(basepath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	e := &EqAssets{
		basepath: basepath,
	}

	// load expansions
	_ = json.Unmarshal(expansionJson, &e.expansions)

	return e
}

// Init initializes the EqAssets application struct
func (e *EqAssets) Init() error {
	// make sure ./files exists
	//if _, err := os.Stat(e.basepath); os.IsNotExist(err) {
	//	err := os.MkdirAll(e.basepath, 0755)
	//	if err != nil {
	//		return err
	//	}
	//}

	// make sure ./files/<expansion-id>-<expansion-name> exists
	//for _, s := range e.expansions {
	//	// make sure dir exists
	//	dir := filepath.Join(e.basepath, "files", fmt.Sprintf("%v-%v", s.Id, slug.Make(s.Name)))
	//	if _, err := os.Stat(dir); os.IsNotExist(err) {
	//		fmt.Println("Creating dir:", dir)
	//		err := os.MkdirAll(dir, 0755)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//}

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
			fmt.Println("err:", err)
			return files
		}

		var expansionFiles ExpansionFiles
		expansionFiles.Expansion = s

		if s.Id <= id {
			// make sure dir exists
			dir := filepath.Join(e.basepath, "files", fmt.Sprintf("%v-%v", s.Id, slug.Make(s.Name)))
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

func (e *EqAssets) PatchFilesForExpansion(id int) error {
	c := config.Get()

	for _, file := range e.GetExpansionFiles(strconv.Itoa(id)) {
		fmt.Println("Checking for files to be deleted in expansion:", file.Expansion.Name)
		for _, f := range file.Files {
			isZoneExtension := strings.Contains(f, ".s3d") ||
				strings.Contains(f, ".eqg") ||
				strings.Contains(f, ".zon")
			if isZoneExtension {
				file := strings.ReplaceAll(f, e.basepath+string(filepath.Separator), "")
				// strip two folder levels
				newFile := strings.Split(file, string(filepath.Separator))
				newFile = append(newFile[:0], newFile[2:]...)
				// path build to temp
				base := filepath.Join(c.EqDir, strings.Join(newFile, string(filepath.Separator)))

				// strip extensions for matching
				base = strings.ReplaceAll(base, ".s3d", "")
				base = strings.ReplaceAll(base, ".eqg", "")
				base = strings.ReplaceAll(base, ".zon", "")

				err := filepath.Walk(c.EqDir, func(path string, info os.FileInfo, err error) error {
					if info == nil {
						return nil
					}

					if info.IsDir() {
						return nil
					}

					// we do this to make sure we're not case sensitive since some file naming for zones
					// are not consistent
					lowerPath := strings.ToLower(path)
					lowerBase := strings.ToLower(base)
					if strings.Contains(lowerPath, lowerBase) && !strings.Contains(path, "maps") {
						fmt.Println("--- Removing file:", path)

						// remove file
						err := os.Remove(path)
						if err != nil {
							return err
						}
					}

					return nil
				})
				if err != nil {
					return err
				}
			}
		}

		fmt.Println("Patching files for expansion:", file.Expansion.Name)
		for _, f := range file.Files {
			file := strings.ReplaceAll(f, e.basepath+string(filepath.Separator), "")
			// strip two folder levels (containing files/expansion/)
			newFile := strings.Split(file, string(filepath.Separator))
			newFile = append(newFile[:0], newFile[2:]...)
			// path build to temp
			destination := filepath.Join(c.EqDir, strings.Join(newFile, string(filepath.Separator)))

			fmt.Printf("--- Copying file %v to %v\n", f, destination)

			// copy source to destination
			err := cp.Copy(f, destination)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *EqAssets) DumpPatchFilesForExpansion(id int) error {
	tmpdir := filepath.Join(os.TempDir(), "patch-dir", time.Now().Format("2006-01-02-15-04-05"))

	err := os.MkdirAll(tmpdir, 0755)
	if err != nil {
		return err
	}

	for _, file := range e.GetExpansionFiles(strconv.Itoa(id)) {
		fmt.Println("Checking for files to be deleted in expansion:", file.Expansion.Name)
		for _, f := range file.Files {
			if strings.Contains(f, ".s3d") || strings.Contains(f, ".eqg") {
				file := strings.ReplaceAll(f, e.basepath+string(filepath.Separator), "")
				// strip two folder levels
				newFile := strings.Split(file, string(filepath.Separator))
				newFile = append(newFile[:0], newFile[2:]...)
				// path build to temp
				base := filepath.Join(tmpdir, strings.Join(newFile, string(filepath.Separator)))
				base = strings.ReplaceAll(base, ".s3d", "")
				base = strings.ReplaceAll(base, ".eqg", "")
				err = filepath.Walk(tmpdir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return err
					}

					if strings.Contains(path, base) && !strings.Contains(path, "maps") {
						fmt.Println("--- Removing file:", path)

						// remove file
						err := os.Remove(path)
						if err != nil {
							return err
						}
					}

					return nil
				})
				if err != nil {
					return err
				}
			}
		}

		fmt.Println("Patching files for expansion:", file.Expansion.Name)
		for _, f := range file.Files {
			file := strings.ReplaceAll(f, e.basepath+string(filepath.Separator), "")
			// strip two folder levels (containing files/expansion/)
			newFile := strings.Split(file, string(filepath.Separator))
			newFile = append(newFile[:0], newFile[2:]...)
			// path build to temp
			destination := filepath.Join(tmpdir, strings.Join(newFile, string(filepath.Separator)))

			fmt.Printf("--- Copying file %v to %v\n", f, destination)

			// copy source to destination
			err := cp.Copy(f, destination)
			if err != nil {
				return err
			}
		}
	}

	err = open.Run(tmpdir)
	if err != nil {
		return err
	}

	return nil
}

func (e *EqAssets) InitPatchFiles() error {
	fmt.Println("InitPatchFiles")

	// download https://github.com/Akkadius/eq-expansion-switcher/releases/download/v1.0.0/files.zip
	// unzip to e.assets.Basepath()
	resp, err := http.Get("https://github.com/Akkadius/eq-expansion-switcher/releases/download/v1.0.0/files-v1.0.9.zip")
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Create the file
	source := filepath.Join(os.TempDir(), "files.zip")
	out, err := os.Create(source)
	if err != nil {
		return err
	}

	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// unzip
	err = unzip.New(source, e.basepath).Extract()
	if err != nil {
		return err
	}

	// remove zip
	_ = os.Remove(source)

	return nil
}
