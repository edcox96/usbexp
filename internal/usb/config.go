package usb

import (
	"fmt"

	"github.com/google/gousb"
	"github.com/spf13/viper"
)

func ActivateAndDisplayConfig(info *UsbDevInfo) error {
	cfgNum := viper.GetInt("cfg")
	if !IsCfgNumValid(info, cfgNum) {
		return fmt.Errorf("invalid cfgNum %d", cfgNum)
	}

	actCfg, err := info.Dev.ActiveConfigNum()
	if err != nil {
		fmt.Printf("Dev.ActiveConfigNum failed! %v\n", err)
		return err
	}
	fmt.Printf("actCfg %d, requested cfg %d\n", actCfg, cfgNum)

	cfg, err := info.Dev.Config(cfgNum)
	if err != nil {
		fmt.Printf("Config failed! %v\n", err)
		return err
	} else if cfg == nil {
		return fmt.Errorf("no config found with cfgNum %d", cfgNum)
	}

	info.DevCfg = cfg

	DisplayConfigDesc(info, cfg.Desc)
	return nil
}

func IsCfgNumValid(info *UsbDevInfo, cfgNum int) bool {
	for _, cfg := range info.DevDesc.Configs {
		if cfgNum == cfg.Number {
			return true
		}
	}
	return false
}

func DisplayConfigDesc(info *UsbDevInfo, cfgDesc gousb.ConfigDesc) error {
	var selfPwrStr string
	var remoteWakeupStr string

	//devDesc := info.DevDesc
	//cfgDesc := devDesc.Configs[cfgNum]

	if cfgDesc.SelfPowered {
		selfPwrStr = "Yes"
	} else {
		selfPwrStr = "No"
	}

	if cfgDesc.RemoteWakeup {
		remoteWakeupStr = "Yes"
	} else {
		remoteWakeupStr = "No"
	}

	fmt.Printf("  ---------------------------\n")
	fmt.Printf("  Configuration %d Descriptor:\n", cfgDesc.Number)
	fmt.Printf("      NumInterfaces = %-6d\n", len(cfgDesc.Interfaces))
	fmt.Printf("      ConfigNum     = %-6d\n", cfgDesc.Number)
	fmt.Printf("      SelfPowered   = %s\n", selfPwrStr)
	fmt.Printf("      RemoteWakeup  = %s\n", remoteWakeupStr)
	fmt.Printf("      MaxPower      = %d mA\n", cfgDesc.MaxPower*2) // FIXME - * 8 for SS
	return nil
}
