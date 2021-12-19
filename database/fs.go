package database

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func initDataDirIfNotExists(path string) error {
	if exists(getGenesisJsonFilePath(path)) {
		return nil
	}

	if err := os.MkdirAll(getDatabaseDirPath(path), os.ModePerm); err != nil {
		return err
	}

	if err := writeGenesisToDisk(getGenesisJsonFilePath(path)); err != nil {
		return err
	}

	if err := writeEmptyBlocksDbToDisk(getBlocksDbFilePath(path)); err != nil {
		return err
	}

	return nil
}

func getDatabaseDirPath(path string) string {
	return filepath.Join(path, "database")
}

func getGenesisJsonFilePath(path string) string {
	return filepath.Join(getDatabaseDirPath(path), "genesis.json")
}

func getBlocksDbFilePath(path string) string {
	return filepath.Join(getDatabaseDirPath(path), "block.db")
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}
