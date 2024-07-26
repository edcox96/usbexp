package usb

import (
	"fmt"
	"strconv"

	"github.com/google/gousb"
	"github.com/spf13/viper"
)

var gousbCtx *gousb.Context
var UsbDevCnt int
var UsbDevsInfo []*UsbDevInfo

func FindDevInfoWithTgtDevAddr(da *UsbDevAddr) *UsbDevInfo {
	for _, info := range UsbDevsInfo {
		desc := info.DevDesc
		switch da.Kind {
		case DA_VidPid:
			if desc.Vendor == gousb.ID(da.Vid) && desc.Product == gousb.ID(da.Pid) {
				return info
			}
		case DA_VidPidBus:
			if desc.Vendor == gousb.ID(da.Vid) && desc.Product == gousb.ID(da.Pid) &&
				desc.Bus == da.Bus {
				return info
			}
		case DA_VidPidBusDev:
			if desc.Vendor == gousb.ID(da.Vid) && desc.Product == gousb.ID(da.Pid) &&
				desc.Bus == da.Bus && desc.Address == da.Dev {
				return info
			}
		}
	}
	return nil
}

type UsbDevInfo struct {
	// gousb state
	GousbCtx *gousb.Context
	Dev *gousb.Device
	DevCfg *gousb.Config
	DevIntf *gousb.Interface

	// top level descriptors
	DevDesc *gousb.DeviceDesc

	// descriptor strings
	Manufacturer string
	Product      string
	SerialNum    string

	// device state
	devAddr *UsbDevAddr		// device vid, pid, bus, address state
	tgtAddr *UsbDevAddr		// target vid, pid and optionally bus and address

	/* remove comments as fields are accessed
	IADpresent bool
	dev_speed  uint
	portNum    uint8
	portNums []uint8
	locID      uint32 // bus_port1_port2_port3 ...
	numCfgs    uint   // TODO - replace with len(Cfgs)
	//struct usbexp_config_info **ppconfig_info; // cfgs ptr array: cfg_info ptr for each dev cfg
	*/
}

func newDevInfo(ctx *gousb.Context, desc *gousb.DeviceDesc) *UsbDevInfo {
	devInfo := &UsbDevInfo { GousbCtx: ctx, DevDesc: desc}
	devInfo.devAddr = GetDevAddr(desc)
	return devInfo
}

type DevAddrKind byte

const (
	DA_None DevAddrKind = iota
	DA_Vid  DevAddrKind = iota
	DA_VidPid
	DA_VidPidBus
	DA_VidPidBusDev
)

type UsbDevAddr struct {
	Kind DevAddrKind
	Vid  uint16
	Pid  uint16
	Bus  int
	Dev  int
}

// get the full device address
func GetDevAddr(desc *gousb.DeviceDesc) *UsbDevAddr {
	da := &UsbDevAddr{}
	da.Vid = uint16(desc.Vendor)
	da.Pid = uint16(desc.Product)
	da.Bus = desc.Bus
	da.Dev = desc.Address
	da.Kind = DA_VidPidBusDev
	return da
}

// get the dev address flags specified on command line
func GetTargetDevAddr() *UsbDevAddr {
	da := &UsbDevAddr{Kind: DA_None}

	daFlags := GetDevAddrFlags()

	if daFlags.Vid != "" {
		vid, err := strconv.ParseInt(daFlags.Vid, 0, 32)
		if err != nil {
			fmt.Printf("ParseInt Vid failed! %v", err)
			return nil
		}
		da.Vid = uint16(vid)
		da.Kind = DA_Vid
	}

	if daFlags.Pid != "" && da.Kind == DA_Vid {
		pid, err := strconv.ParseInt(daFlags.Pid, 0, 32)
		if err != nil {
			return nil
		}
		da.Pid = uint16(pid)
		da.Kind = DA_VidPid
	}

	if daFlags.Bus != "" && da.Kind == DA_VidPid {
		bus, err := strconv.ParseInt(daFlags.Bus, 0, 16)
		if err != nil {
			return nil
		}
		da.Bus = int(bus)
		da.Kind = DA_VidPidBus
	}

	if daFlags.Dev != "" && da.Kind == DA_VidPidBus {
		dev, err := strconv.ParseInt(daFlags.Dev, 0, 16)
		if err != nil {
			return nil
		}
		da.Dev = int(dev)
		da.Kind = DA_VidPidBusDev
	}

	return da
}

type DevAddrFlags struct {
	Vid string
	Pid string
	Bus string
	Dev string
}

func GetDevAddrFlags() *DevAddrFlags {
	daFlags := &DevAddrFlags{}

	daFlags.Vid = viper.GetString("vid")
	daFlags.Pid = viper.GetString("pid")
	daFlags.Bus = viper.GetString("bus")
	daFlags.Dev = viper.GetString("dev")
	return daFlags
}

func SetDebugLevel() {
	dbgLvl := viper.GetInt("dbg_lvl")
	gousbCtx.Debug(dbgLvl)
}

// Called once by usb.init() when usb package is imported.
// Inits libusb and sets gousbCtx and usbDevCnt vars. The
// gousbCtx is kept open until app shutdown.

func getUsbDevsInfo() error {
	if gousbCtx != nil {
		// should only be called by init once so repeated call is panic
		fmt.Printf("GetUsbDevsInfo called and gousbCtx is %v\n", gousbCtx)
		panic(1)
	}

	gousbCtx = gousb.NewContext() // panics if error.

	// enumerate the attached usb devices list and get the device info 
	err := getInfoForAttachedDevices()
	if err != nil {
		fmt.Printf("GetDeviceDescs failed! %v", err)
		return err
	}
	return nil
}

// opener is called for each attached usb device by OpenDevices
func opener(devDesc *gousb.DeviceDesc) bool {
	devInfo := newDevInfo(gousbCtx, devDesc)
	UsbDevsInfo = append(UsbDevsInfo, devInfo)

	fmt.Printf("vid 0x%x pid 0x%x\n", uint16(devDesc.Vendor), uint16(devDesc.Product))
	return false // just get the UsbDevInfo for now and don't open the devices
}

func getInfoForAttachedDevices() (err error) {
	// Iterate through available Devices and call opener for each dev
	devs, err := gousbCtx.OpenDevices(opener)
	if err != nil {
		fmt.Printf("OpenDevices()failed! %v", err)
		return err
	}
	if len(devs) > 0 {
		return fmt.Errorf("unexpected dev returned from OpenDevices")
	}
	UsbDevCnt = len(UsbDevsInfo)
	return nil
}

func init() {
	// Enumerate the USB device list and get the device info and state
	// if not already done.
	if err := getUsbDevsInfo(); err != nil {
		fmt.Printf("usb.GetUsbDevsInfo failed! %v", err)
		return
	}
}

func Cleanup() {
	gousbCtx.Close()
}
