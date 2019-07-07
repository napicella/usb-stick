package download

import (
	"github.com/spf13/cobra"
	"github.com/usb-stick-client/gozip"
	"github.com/usb-stick-client/interactive"
	"github.com/usb-stick-client/model"
	"log"
	"os"
	"path/filepath"
)

func Cmd(usbService model.UsbService) (*cobra.Command, error) {
	var data, etag string
	var cmd = &cobra.Command{
		Use:   "download",
		Short: "download data from the usb cloud stick",
		RunE: func(cmd *cobra.Command, args []string) error {
			pass, e := interactive.PromptForPassword()
			if e != nil {
				return e
			}
			e = download(etag, data, pass, usbService)
			if e != nil {
				return e
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&etag, "etag", "e", "", "Etag of the data to download")
	cmd.Flags().StringVarP(&data, "data", "d", "", "Directory in which the data should be downloaded")
	_ = cmd.MarkFlagRequired("etag")
	_ = cmd.MarkFlagRequired("data")

	return cmd, nil
}

func download(etag, destFolder, password string, service model.UsbService) error {
	data, e := service.Download(etag)
	if e != nil {
		return e
	}
	log.Printf("Downloaded %d bytes \n", len(data))
	log.Println("Decrypting")
	decrypted := gozip.Decrypt(data, password)
	path, e := saveToTmpFile(decrypted)
	if e != nil {
		return e
	}
	defer os.Remove(path)

	return gozip.Unzip(path, destFolder)
}

func saveToTmpFile(data []byte) (string, error) {
	savePath := filepath.Join(os.TempDir(), "toUnzip")
	out, e := os.Create(savePath)
	if e != nil {
		return "", e
	}
	defer out.Close()
	_, e = out.Write(data)
	return savePath, e
}


