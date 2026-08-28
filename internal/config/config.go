package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = "/.gatorconfig.json"

type Config struct {
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}


func Read() (Config, error) {
	file , err := getConfigFilePath()
	if err != nil {
        return Config{}, err
    }
	data, err := os.ReadFile(file)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
        return Config{}, err
    }
	return config, nil
}


func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username
	return write(*cfg)

}


func write(cfg Config) error{
	path , err := getConfigFilePath()
	if err!=nil{
		return err
	}
	file,err := os.Create(path)
	if err!=nil{
		return err
	}
	defer file.Close()
	data, err := json.Marshal(cfg)
	if err!=nil{
		return err
	}
	return os.WriteFile(path, data, 0600)
}


func getConfigFilePath() (string, error){
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}