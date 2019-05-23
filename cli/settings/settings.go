package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/winded/tyomaa/shared/api/client"
)

type Settings struct {
	Api client.Settings `json:"api"`
}

func GetFilePath() (string, error) {
	pth := os.Getenv("SETTINGS_FILE")
	if pth != "" {
		return pth, nil
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return pth, err
	}

	return path.Join(homedir, ".tyomaa-cli"), nil
}

func Load() (Settings, error) {
	pth, err := GetFilePath()
	if err != nil {
		return Settings{}, err
	}

	_, err = os.Stat(pth)
	if os.IsNotExist(err) {
		return Settings{}, nil
	} else if err != nil {
		return Settings{}, err
	}

	f, err := os.Open(pth)
	if err != nil {
		return Settings{}, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return Settings{}, err
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return settings, err
	}

	return settings, nil
}

func Save(settings Settings) error {
	pth, err := GetFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(&settings)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pth, data, 0700)
}
