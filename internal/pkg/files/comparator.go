package files

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/rs/zerolog/log"
)

// calcMd5 calculates the MD5 of a given file.
// TODO Improve method to avoid whole file copy into memory even if not expected to be a problem for the protos.
func calcMd5(filePath string) (*string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	md5sum := fmt.Sprintf("%x", h.Sum(nil))
	return &md5sum, nil
}

// CompareFilesAreEqual compares two files to check if they are equal.
func CompareFilesAreEqual(newPath string, oldPath string) (bool, error) {
	if _, err := os.Stat(oldPath); err != nil {
		// New file, so no match
		return false, nil
	}
	// Both files exists, so compare them
	checksumNew, err := calcMd5(newPath)
	if err != nil {
		return false, err
	}
	checksumOld, err := calcMd5(oldPath)
	if err != nil {
		return false, err
	}
	log.Debug().Str("new", *checksumNew).Str("old", *checksumOld).Msg("MD5")
	return *checksumNew == *checksumOld, nil
}

// CompareDirectoriesAreEqual compares two different set of files and return if there are changes between the two sets.
func CompareDirectoriesAreEqual(extension string, newPath string, oldPath string) (bool, error) {
	log.Debug().Str("extension", extension).Str("newPath", newPath).Str("oldPath", oldPath).Msg("comparing files")

	// Iterate on the list of new files, and compare those with the previous ones.
	fileInfo, err := ioutil.ReadDir(newPath)
	if err != nil {
		return false, err
	}

	directoriesAreEqual := true
	for _, info := range fileInfo {
		if !info.IsDir() && strings.HasSuffix(info.Name(), extension) {
			newFilePath := path.Join(newPath, info.Name())
			oldFilePath := path.Join(oldPath, info.Name())
			filesAreEqual, err := CompareFilesAreEqual(newFilePath, oldFilePath)
			if err != nil {
				return false, err
			}
			log.Debug().Str("newFilePath", newFilePath).Str("oldFilePath", oldFilePath).Bool("fileAreEqual", filesAreEqual).Msg("file comparison")
			directoriesAreEqual = directoriesAreEqual && filesAreEqual
		}
	}

	return directoriesAreEqual, nil
}
