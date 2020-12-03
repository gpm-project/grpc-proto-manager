package files

import (
	"io"
	"os"
)

// CopyFile copies the content of a file into a new path.
func CopyFile(source string, dest string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := sourceFile.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destFile.Write(buffer[:n]); err != nil {
			return err
		}
	}
	return nil
}
