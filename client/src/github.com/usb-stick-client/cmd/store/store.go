package store

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/usb-stick-client/gozip"
	"github.com/usb-stick-client/interactive"
	"github.com/usb-stick-client/model"
	"log"
	"os"
	"path/filepath"
)



func Cmd(usbService model.UsbService) (*cobra.Command, error) {
	var cmd = &cobra.Command{
		Use:   "upload",
		Aliases: []string{"up", "store"},
		Short: "upload file or directory to the the cloud usb",
		Long: ``,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1  {
				return errors.New("too many params")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, e := getOrDefault(args, os.Getwd)
			if e != nil {
				return e
			}

			log.Printf("Going to upload %s\n", dir)

			pass, e := interactive.PromptForPassword()
			if e != nil {
				return e
			}
			etag, e := store(dir, pass, usbService)
			if e != nil {
				return e
			}
			fmt.Println(etag)
			return nil
		},
	}


	return cmd, nil
}

// directory or file
func store(directory string, password string, service model.UsbService) (string, error) {
	log.Println("Encrypting the content")

	destinationFilePath := filepath.Join(os.TempDir(), "toUpload")
	destinationFile, e := os.Create(destinationFilePath)
	if e != nil {
		return "", e
	}
	defer os.Remove(destinationFile.Name())

	e = gozip.ZipAndEncrypt(directory, destinationFilePath, password)
	if e != nil {
		return "", e
	}

	return service.Store(destinationFilePath)
}

func getOrDefault(args []string, supplier func() (string, error)) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}

	return supplier()
}
