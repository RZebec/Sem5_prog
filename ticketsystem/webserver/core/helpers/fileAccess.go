package helpers

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
	Returns true if the file path exists, false if not.
*/
func FilePathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

/*
	Create a folder path.
*/
func CreateFolderPath(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

/*
	Creates a file if it does not already exist.
*/
func CreateFileIfNotExists(path string) error {
	exists, er := FilePathExists(path)
	if er != nil {
		return er
	} else {
		if !exists {
			var file, err = os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()
		}
	}
	return nil
}

/*
	Create a file and the path if necessary.
*/
func CreateFileWithPathIfNotExists(path string) (bool, error) {
	exists, err := FilePathExists(path)
	if err != nil {
		return false, err
	}
	if !exists {
		dir, _ := filepath.Split(path)
		if len(dir) > 0 {
			err = CreateFolderPath(dir)
			if err != nil {
				return false, err
			}
		}
		err = CreateFileIfNotExists(path)
		if err != nil {
			return false, err
		}
	}
	return exists, nil
}

/*
	Read all data from a given file.
*/
func ReadAllDataFromFile(filePath string) ([]byte, error) {
	if filePath == "" {
		return nil, errors.New("invalid file path. value:" + filePath)
	}
	var file, fileErr = os.OpenFile(filePath, os.O_RDWR, 0644)
	defer file.Close()
	if fileErr != nil {
		return nil, fileErr
	}

	return ioutil.ReadAll(file)
}

/*
	Write all data to a given file.
*/
func WriteDataToFile(filePath string, data []byte) error {
	if filePath == "" {
		return errors.New("invalid file path. value:" + filePath)
	}
	if data == nil {
		return errors.New("invalid data. value: nil")
	}
	return ioutil.WriteFile(filePath, data, 0644)
}
