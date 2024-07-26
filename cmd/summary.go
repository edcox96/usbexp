/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"usbexp/internal/usb"

	"github.com/spf13/cobra"
)

// summaryCmd represents the summary command
var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: RunSummaryCmd,
}

func RunSummaryCmd(cmd *cobra.Command, args []string) {
	usb.SetDebugLevel()

	da := usb.GetTargetDevAddr()
	var devInfo []*usb.UsbDevInfo
	if da.Kind >= usb.DA_VidPid {
		info := usb.FindDevInfoWithTgtDevAddr(da)
		devInfo = append(devInfo, info)
	} else {
		devInfo = usb.UsbDevsInfo
	}
	for _, info := range devInfo {
		vid := info.DevDesc.Vendor
		pid := info.DevDesc.Product

		err := usb.GetDeviceDescStrings(info, vid, pid)
		if err != nil {
			fmt.Printf("GetDeviceDescStrings failed! %v", err)
			return
		}

		err = usb.DisplayDeviceSummary(info)
		if err != nil {
			fmt.Printf("usb.DisplayDeviceSummary failed! %v", err)
			return
		}
	}
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
