package mergepath

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func MergePaths(rootSourcePath string, rootTargetPath string) error {
	rootSourceFInfo, err := os.Stat(rootSourcePath)
	if err != nil {
		return err
	}
	if !rootSourceFInfo.IsDir() {
		err = CopyFile(rootSourcePath, rootTargetPath)
		if err != nil {
			return err
		}
		return nil
	}
	rootTargetFInfo, err := os.Stat(rootTargetPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if !rootTargetFInfo.IsDir() {
			err = os.Remove(rootTargetPath)
			if err != nil {
				return err
			}
		}
	}
	err = filepath.Walk(
		rootSourcePath,
		func(
			sourcePath string,
			sourceFInfo os.FileInfo,
			err error,
		) error {
			sharedPath := sourcePath[len(rootSourcePath):]
			if err != nil {
				return err
			}
			targetPath := strings.ReplaceAll(rootTargetPath+"/"+sharedPath, "//", "/")
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
	} else if !targetFInfo.Mode().IsRegular() {
		if targetFInfo.IsDir() {
			err := os.RemoveAll(targetPath)
			if err != nil {
				return err
			}
		} else {
			fmt.Errorf(
				"CopyFile: non-regular destination file %s (%q)",
				targetFInfo.Name(),
				targetFInfo.Mode().String(),
			)
			return nil
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
