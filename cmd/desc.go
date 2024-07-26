/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	usb "usbexp/internal/usb"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CfgNum int
var IntfNum int
var AltSetNum int

// descCmd represents the desc command
var DescCmd = &cobra.Command{
	Use:   "desc",
	Short: "Get and display the device desc for specified vid and pid",
	Long:  LongDesc,
	Run:   RunDescCmd,
}

var LongDesc = `
Get and display the device descriptors for the specified
target usb device.
Usage: usbexp desc --vid=0x04b4 --pid=0x00f3 // minrray camera
` // end of LongDesc

func RunDescCmd(cmd *cobra.Command, args []string) {
	info, err := usb.OpenTargetDev()
	if err != nil {
		fmt.Printf("OpenTargetDev() failed. %v", err)
		return
	}
	defer info.Dev.Close()

	if err := usb.DisplayDeviceInfo(info); err != nil {
		fmt.Printf("DisplayDeviceInfo failed! %v", err)
		return
	}

	if err := usb.ActivateAndDisplayConfig(info); err != nil {
		fmt.Printf("ActivateAndDisplayConfig failed! %v", err)
		return
	}
	defer info.DevCfg.Close()
	
	if err := usb.ClaimAndDisplayIntf(info); err != nil {
		fmt.Printf("ClaimAndDisplayInt failed! %v", err)
		return
	}

	defer info.DevIntf.Close()
}

func init() {
	rootCmd.AddCommand(DescCmd)

	DescCmd.Flags().IntVar(&CfgNum, "cfg", 1, "--cfg=1 show config 1 descriptor)")
	viper.BindPFlag("cfg", DescCmd.Flags().Lookup("cfg"))

	DescCmd.Flags().IntVar(&IntfNum, "intf", 0, "--intf=1 show interface 1 descriptor)")
	viper.BindPFlag("intf", DescCmd.Flags().Lookup("intf"))

	DescCmd.Flags().IntVar(&AltSetNum, "alt", 0, "--alt=1 show intf alt set 1 descriptor)")
	viper.BindPFlag("alt", DescCmd.Flags().Lookup("alt"))
}
