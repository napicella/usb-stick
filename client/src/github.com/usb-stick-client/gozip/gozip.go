package gozip

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"log"

	"os"
	"path/filepath"
	"strings"
)

const (
	encrypted              = "encrypted"
	decrypted              = "decrypted"
	toEncryptedZipFileName = "to-encrypt.zip"
	decryptedZipFileName   = "decrypted.zip"

)

func ZipAndEncrypt(pathToZip, destFilePath, password string) error {
	// TODO dir and files should be deleted afterwards
	tmpFolder, e := ioutil.TempDir("", encrypted)
	if e != nil {
		return e
	}
	encryptedFilePath := filepath.Join(tmpFolder, toEncryptedZipFileName)

	e = zipit(pathToZip, encryptedFilePath)
	if e != nil {
		return  e
	}

	return encryptFile(encryptedFilePath, destFilePath, password)
}

func DecryptAndUnzip(sourceZip, destFolder, password string) error {
	decryptDestFolder, e := ioutil.TempDir("", decrypted)
	if e != nil {
		return e
	}
	decryptDestFilePath := filepath.Join(decryptDestFolder, decryptedZipFileName)
	_, e = decryptFile(sourceZip, decryptDestFilePath, password)
	if e != nil {
		return e
	}
	log.Println("File decrypted")
	return Unzip(decryptDestFilePath, destFolder)
}

func zipit(source, target string) error {
	zipFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	_ = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

func Unzip(archive, destFolder string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(destFolder, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(destFolder, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		targetFile, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if e != nil {
			return err
		}

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
		targetFile.Close()
		fileReader.Close()
	}

	return nil
}