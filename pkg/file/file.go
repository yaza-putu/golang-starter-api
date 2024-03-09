package file

import (
	"fmt"
	"github.com/yaza-putu/golang-starter-api/pkg/unique"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

// ToPublic folder
func ToPublic(file *multipart.FileHeader, dest string, randomName bool) (string, error) {
	src, err := file.Open()
	defer src.Close()

	if err != nil {
		return "", err
	}

	// Destination
	fileName := file.Filename
	if randomName {
		split := strings.Split(file.Filename, ".")
		fileName = fmt.Sprintf("%s.%s", unique.Uid(13), split[len(split)-1])
	}

	destPath := fmt.Sprintf("public/%s/%s", dest, fileName)
	_, err = os.Stat(fmt.Sprintf("public/%s"))
	if err != nil {
		err = os.Mkdir(fmt.Sprintf("public/%s", dest), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	dst, err := os.Create(destPath)
	defer dst.Close()
	if err != nil {
		return "", err
	}

	// store to destination
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dest, nil
}

// ToPrivate folder
func ToPrivate(file *multipart.FileHeader, dest string, randomName bool) (string, error) {
	src, err := file.Open()
	defer src.Close()

	if err != nil {
		return "", err
	}

	// Destination
	fileName := file.Filename
	if randomName {
		split := strings.Split(file.Filename, ".")
		fileName = fmt.Sprintf("%s.%s", unique.Uid(13), split[len(split)-1])
	}

	destPath := fmt.Sprintf("storage/%s/%s", dest, fileName)
	_, err = os.Stat(fmt.Sprintf("storage/%s"))
	if err != nil {
		err = os.Mkdir(fmt.Sprintf("storage/%s", dest), os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	dst, err := os.Create(destPath)
	defer dst.Close()
	if err != nil {
		return "", err
	}

	// store to destination
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return dest, nil
}

// DetectContentType to make sure the mimes of data
func DetectContentType(file multipart.File, allowMimes []string) bool {
	// get type MIME of file
	buffer := make([]byte, 512) // read 512 bytes to make sure MIME
	_, err := file.Read(buffer)
	if err != nil {
		return false
	}
	fileType := http.DetectContentType(buffer)
	// compare mime type
	for _, validType := range allowMimes {
		fmt.Println(fileType, validType)
		if strings.HasPrefix(fileType, validType) {
			return true
		}
	}
	return false
}
