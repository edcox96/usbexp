package usb

import (
	"fmt"

	"github.com/google/gousb"
	"github.com/spf13/viper"
)

func ClaimAndDisplayIntf(info *UsbDevInfo) error {
	cfg := info.DevCfg
	if cfg == nil {
		return fmt.Errorf("in ClaimAndDisplayIntf() cfg must not be nil")
	}

	intfNum := viper.GetInt("intf")
	fmt.Printf("intfNum %d\n", intfNum)

	if intfNum >= len(cfg.Desc.Interfaces) {
		return fmt.Errorf("--intf=%d is invalid", intfNum)
	}

	altSetNum := viper.GetInt("alt")
	fmt.Printf("altSetNum %d\n", altSetNum)

	if altSetNum >= len(cfg.Desc.Interfaces[intfNum].AltSettings) {
		return fmt.Errorf("--alt=%d is invalid", altSetNum)
	}

	// the mac OS security prevents cmd line apps from claiming USB interface
	fmt.Printf("Can't open USB interface on MacOS!\n")
	/*

	intf, err := cfg.Interface(intfNum, altSetNum)
	if err != nil {
		fmt.Printf("cfg.Interfce() failed! %v", err)
		return err
	}

	intfSet := &intf.Setting
	info.DevIntf = intf

	DisplayInterfaceDesc(info, intfSet)
	*/
	return nil
}

func DisplayInterfaceDesc(info *UsbDevInfo, altSet *gousb.InterfaceSetting) error {
	fmt.Printf("    ---------------------------------\n")
	fmt.Printf("    INTERFACE %d, Alternate Setting %d\n", altSet.Number, altSet.Alternate)
	fmt.Printf("      InterfaceClass     = 0x%-2x // %s\n", uint8(altSet.Class), altSet.Class.String())
	fmt.Printf("      InterfaceSubClass  = %-4d // %s\n", uint8(altSet.SubClass), altSet.SubClass.String())
	fmt.Printf("      InterfaceProtocol  = %-6d\n", altSet.Protocol)
	fmt.Printf("      NumEndpoints       = %d\n", len(altSet.Endpoints))

	for i := 0; i < 16; i++ {
		ep := altSet.Endpoints[gousb.EndpointAddress(i)]
		fmt.Printf("ep.MaxPacketSize %d\n", ep.MaxPacketSize)
	}
	return nil
}
