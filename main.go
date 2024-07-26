/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_"fmt"
	"os"

	cmd "usbexp/cmd"
	usb "usbexp/internal/usb"
)

var UsbExp = &UsbExplorer{}

type UsbExplorer struct {
//	GoUsbCtx *gousb.Context
}


func main() {
	defer usb.Cleanup()
	
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
