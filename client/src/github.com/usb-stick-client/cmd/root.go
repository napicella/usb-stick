package cmd


import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/usb-stick-client/cmd/download"
	"github.com/usb-stick-client/cmd/store"
	"github.com/usb-stick-client/model"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "usb-stick",
	Short: "upload to and download from a cloud usb stick",
	Long: `Complete documentation is available at GITHUB-REPO`,
}

func Execute() {
	usbService := model.NewUsbService()
	storeCmd, e := store.Cmd(usbService)
	downloadCmd, e := download.Cmd(usbService)
	if e != nil {
		panic(e)
	}
	rootCmd.AddCommand(storeCmd, downloadCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

