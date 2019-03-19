package mergedir

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MergeDirectories(sourceDir string, targetDir string) error {
	err := filepath.Walk(
		sourceDir,
		func(
			sourcePath string,
			sourceFInfo os.FileInfo,
			err error,
		) error {
			sharedPath := sourcePath[len(sourceDir):]
			if err != nil {
				return err
			}
			targetPath := targetDir + sharedPath
			if sourceFInfo.IsDir() {
				targetFInfo, err := os.Stat(targetPath)
				if err != nil {
					if !os.IsNotExist(err) {
						return err
					}
				} else {
					if !targetFInfo.IsDir() {
						err = os.Remove(targetPath)
						if err != nil {
							return err
						}
					}
				}
				os.MkdirAll(targetPath, sourceFInfo.Mode())
			} else {
				err = CopyFile(sourcePath, targetPath)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(sourcePath string, targetPath string) error {
	sourceFInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}
	if !sourceFInfo.Mode().IsRegular() {
		fmt.Errorf(
			"CopyFile: non-regular source file %s (%q)",
			sourceFInfo.Name(),
			sourceFInfo.Mode().String(),
		)
		return nil
	}
	targetFInfo, err := os.Stat(targetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !targetFInfo.Mode().IsRegular() {
			fmt.Errorf(
				"CopyFile: non-regular destination file %s (%q)",
				targetFInfo.Name(),
				targetFInfo.Mode().String(),
			)
		}
	}
	if os.SameFile(sourceFInfo, targetFInfo) {
		return nil
	}
	err = os.Link(sourcePath, targetPath)
	if err != nil {
		err = copyFileContents(sourcePath, targetPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFileContents(sourcePath string, targetPath string) error {
	sourceFInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}
	sourceData, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(targetPath, sourceData, sourceFInfo.Mode())
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}
