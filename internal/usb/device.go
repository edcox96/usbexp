package usb

import (
	"fmt"

	"github.com/google/gousb"
)

// Process --vid, --pid, --bus, and --dev flags to create target DevAddr. The
// vid and pid are required to find target device info and to open the device.

func OpenTargetDev() (*UsbDevInfo, error) {
	da := GetTargetDevAddr()
	if da.Kind < DA_VidPid {
		return nil, fmt.Errorf("--vid and --pid flags required")
	}

	info := FindDevInfoWithTgtDevAddr(da)
	if info == nil {
		return nil, fmt.Errorf("no matching vid/pid found")
	}

	dev, err := info.GousbCtx.OpenDeviceWithVIDPID(gousb.ID(da.Vid), gousb.ID(da.Pid))
	if err != nil {
		return nil, err
	}

	//if err := dev.SetAutoDetach(true); err != nil { return nil, err }

	info.tgtAddr = da
	info.Dev = dev

	return info, nil
}

func DisplayDeviceInfo(info *UsbDevInfo) error {
	if info.Dev == nil {
		return fmt.Errorf("device must be open and the info.Dev value set")
	}
	vid := gousb.ID(info.devAddr.Vid)
	pid := gousb.ID(info.devAddr.Pid)

	// device is open so get the device desc strings
	if err := GetDeviceDescStrings(info, vid, pid); err != nil {
		fmt.Printf("GetDeviceDescStrings failed! %v", err)
		return err
	}

	if err := DisplayDeviceSummary(info); err != nil {
		fmt.Printf("DisplayDeviceSummary failed! %v", err)
		return err
	}

	if err := DisplayDeviceDesc(info); err != nil {
		fmt.Printf("DisplayDeviceDesc failed! %v", err)
		return err
	}
	return nil
}

func GetDeviceDescStrings(info *UsbDevInfo, vid, pid gousb.ID) error {
	dev := info.Dev

	man, err := dev.Manufacturer()
	if err != nil {
		fmt.Printf("dev.Manufacturer failed! %v", err)
		return err
	}
	info.Manufacturer = man

	product, err := dev.Product()
	if err != nil {
		fmt.Printf("dev.Product failed! %v", err)
		return err
	}
	info.Product = product

	sn, err := dev.SerialNumber()
	if err != nil {
		fmt.Printf("dev.SerialNumber failed! %v", err)
		return err
	}
	info.SerialNum = sn

	return nil
}

func DisplayDeviceDesc(info *UsbDevInfo) error {
	desc := info.DevDesc

	fmt.Printf("  Device Descriptor\n")
	DisplayCommonDevDesc(desc)

	fmt.Printf("    idVendor        = 0x%04x // %s\n", uint16(desc.Vendor),
		getVidStr(uint16(desc.Vendor)))
	fmt.Printf("    idProduct       = 0x%04x\n", uint16(desc.Product))
	fmt.Printf("    bcdDevice       = 0x%02x%02x // Device %s\n", desc.Device.Major(),
		desc.Device.Minor(), desc.Device.String())
	fmt.Printf("    Manufacturer    = %s\n", info.Manufacturer)
	fmt.Printf("    Product         = %s\n", info.Product)
	fmt.Printf("    SerialNumber    = %s\n", info.SerialNum)
	fmt.Printf("    NumConfigs      = %-6d\n", len(desc.Configs))
	return nil
}

func DisplayCommonDevDesc(desc *gousb.DeviceDesc) error {
	var protocolStr = ""
	if desc.SubClass == 0x2 && desc.Protocol == 1 {
		protocolStr = "// IADs - present"
	}

	fmt.Printf("    bcdUSB          = 0x%02x%02x // USB %s\n", desc.Spec.Major(),
		desc.Spec.Minor(), desc.Spec.String())
	fmt.Printf("    bDeviceClass    = 0x%-4x // %s\n", uint8(desc.Class), desc.Class.String())
	fmt.Printf("    bDeviceSubClass = %-6d // %s\n", desc.SubClass, desc.SubClass.String())
	fmt.Printf("    bDeviceProtocol = %-6d %s\n", desc.Protocol, protocolStr)

	size, str, err := GetMaxCtrlPktSize(desc)
	if err != nil {
		fmt.Printf("GetMaxPktSize failed! %v", err)
		return err
	}
	fmt.Printf("    bMaxPacketSize0 = %-6d // %s%d bytes max for EP 0\n",
		desc.MaxControlPacketSize, str, size)
	return nil
}

func GetMaxCtrlPktSize(desc *gousb.DeviceDesc) (int, string, error) {
	var maxPktSizeStr = ""
	maxPktSize := desc.MaxControlPacketSize

	//bcdval, err := strconv.Atoi(desc.Spec.String())
	if desc.Spec >= 0x0300 {
		maxPktSize = 0x1 << desc.MaxControlPacketSize
		maxPktSizeStr = fmt.Sprintf("1 << %d = ", desc.MaxControlPacketSize)
	}
	return maxPktSize, maxPktSizeStr, nil
}

var vendors = map[uint16]string{
	0x04B4: "Cypress Semiconductor",
	0x0518: "Plus More Enterprises",
	0x05ac: "Apple",
	0x2109: "Via Labs Inc",
	0x2B77: "Epiphan Systems Inc",
}

func getVidStr(vid uint16) string {
	return vendors[vid]
}
