package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func inList(input []string, check string) bool {
	for _, name := range input {
		if name == check {
			return true
		}
	}
	return false
}

func untartar(tarName, xpath string, allowedFiles []string) (err error) {
	tarFile, err := os.Open(tarName)

	defer tarFile.Close()
	absPath, err := filepath.Abs(xpath)

	gz, _ := gzip.NewReader(tarFile)
	tr := tar.NewReader(gz)

	// untar each segment
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// determine proper file path info
		finfo := hdr.FileInfo()
		fileName := hdr.Name
		absFileName := filepath.Join(absPath, fileName)
		s := strings.Split(fileName, "/")

		JustFile := s[len(s)-1]

		fmt.Println("file name here ", s[len(s)-1])
		// if a dir, create it, then go to next segment
		if finfo.Mode().IsDir() {

			continue
		}
		// create new file with original file mode

		if JustFile != "" {
			if inList(allowedFiles, JustFile) {
				finalName := "data/" + JustFile
				file, err := os.OpenFile(
					finalName,
					os.O_RDWR|os.O_CREATE|os.O_TRUNC,
					finfo.Mode().Perm(),
				)

				if err != nil {
					return err
				}
				fmt.Printf("x %s\n", absFileName)

				n, cpErr := io.Copy(file, tr)
				if closeErr := file.Close(); closeErr != nil {
					return err
				}

				if cpErr != nil {
					return cpErr
				}
				if n != finfo.Size() {
					return fmt.Errorf("wrote %d, want %d", n, finfo.Size())
				}
			}
		}
	}
	return nil
}
func main() {
	fmt.Println("hello world write files to data directory ")

	untartar("data/example.tar.gz", "data", []string{"file2.txt", "file3.txt"})

}
