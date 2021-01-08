package preferences

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/iac-io/myiac/internal/util"
	"gopkg.in/ini.v1"
)

const (
	preferencesFile = "/.myiac/prefs"
)

type Preferences interface {
	Set(name string, value string)
	Get(name string) string
	Del(name string)
}

type configPreferences struct {
	fileName       string
	propertiesFile *ini.File
}

func DefaultConfig() *configPreferences {
	homeDir, _ := os.UserHomeDir()
	prefsFilePath := homeDir + preferencesFile
	return NewConfig(prefsFilePath)
}

func NewConfig(prefsFilePath string) *configPreferences {
	var errFile error = nil

	if !util.FileExists(prefsFilePath) {
		_, errFile = createPath(prefsFilePath)
		if errFile != nil {
			fmt.Printf("Error creating file %v\n", errFile)
			panic(errFile)
		}
	}

	propertiesFile, err := ini.Load(prefsFilePath)
	if err == nil {
		log.Printf("Preferences file created at %s\n", prefsFilePath)
	} else {
		panic(err)
	}

	return &configPreferences{fileName: prefsFilePath, propertiesFile: propertiesFile}
}

func (prefs configPreferences) Set(name string, value string) {
	key, _ := prefs.propertiesFile.Section("").NewKey(name, value)
	err := prefs.propertiesFile.SaveTo(prefs.fileName)
	if err != nil {
		log.Printf("[WARN] error saving preferences %v", err)
	}
	log.Printf("Saved preference %s\n", key.Name())
}

func (prefs configPreferences) Get(name string) string {
	value, err := prefs.propertiesFile.Section("").GetKey(name)
	if err != nil {
		return ""
	}

	return value.String()
}

func (prefs configPreferences) Del(name string) {
	prefs.propertiesFile.Section("").DeleteKey(name)
}

func createPath(p string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(p), 0770); err != nil {
		return nil, err
	}
	return os.Create(p)
}
