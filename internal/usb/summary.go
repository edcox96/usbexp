package usb

import (
	"fmt"
	"os"
)

const frameWidth = 50 //128
var numDevs = 0       // track num_devs across mult invocations of this func
var frame = newFrame()

func DisplayDeviceSummary(info *UsbDevInfo) error {
	desc := info.DevDesc
	if numDevs == 0 {
		fmt.Printf("%s", frame)
	}
	numDevs++

	if numDevs == 1 { /*|| papp_settings->show_descriptors*/
		//fmt.Printf("|     Manufacturer    |       Product      |  VID    PID   | Location ID |    Bus    | Dev Addr | Port | Hub | CID |   Speed   |\n")
		fmt.Printf("| Bus  |       Device Info        |\n")
		fmt.Printf("%s", frame)
	}

    var cnt int
	cnt, _ = fmt.Printf("|      | Manufacturer: %s", info.Manufacturer) //   |       Product      |  VID    PID   | Location ID |    Bus    | Dev Addr | Port | Hub | CID |   Speed   |\n")
	fmt.Printf("%s|\n", padCol(34, cnt))
    if info.DevDesc.Bus > 100 {
        cnt, _ = fmt.Printf("| %d  | Product: %s", info.DevDesc.Bus, info.Product)
    } else {
	    cnt, _ = fmt.Printf("|  %d  | Product: %s", info.DevDesc.Bus, info.Product)
    }
	fmt.Printf("%s|\n", padCol(34, cnt))
    var sn = info.SerialNum
    if len(sn) > 12 {
        for i := 0; i < 12; i++ {
            sn = info.SerialNum[:12]
        }
    }
	cnt, _ = fmt.Printf("|      | Serial Num: %s", sn)
	fmt.Printf("%s|\n", padCol(34, cnt))
	fmt.Printf("|      | VID: 0x%04x, PID: 0x%04x |\n", uint16(desc.Vendor), uint16(desc.Product))
	fmt.Printf("%s", frame)
	/*
		        cidIdxStr := "   "

		        if (pdev_info->container_id_idx >= 0) {
		            snprintf(cid_idx_str, 3, "%d", pdev_info->container_id_idx);
		        } else {
		            strcpy(cid_idx_str, "   ");
		        }

			fmt.Printf("| %-19s | %-18s | 0x%04X 0x%04X | 0x%08x  | 0x%x (%d)%s |    %2d    |  %2d  | %3d |%3s  | %3s bps |\n",
				info.Manufacturer, info.Product, uint16(desc.Vendor), uint16(desc.Product) pdev_info->loc_id, 0x00001234, desc.Bus, desc.Bus,
				" ", desc.Address, desc.Port plog_dev_info->parent_hub_addr, 0, cidIdxStr, desc.Speed.String())
	*/

	return nil
}

func newFrame() string {
	frame := "|"
	for i := 1; i < frameWidth-2; i++ {
		frame += "-"
	}
	frame += "|"
	frame += "\n"
	return frame
}

func padCol(colWidth, curWidth int) string {
	if curWidth >= colWidth {
		fmt.Printf("padCol: colWidth %d curWidth %d\n", colWidth, curWidth)
		os.Exit(1)
	}
	pad := ""
	for i := curWidth; i < colWidth; i++ {
		pad += " "
	}
	return pad
}

/*
void display_device_summary(struct usbexp_log_usb_dev_info *plog_dev_info)
{
    struct libusb_device_descriptor *pdev_desc;
    struct usbexp_dev_info *pdev_info;
    usbexp_dev_addr_t *pdev_addr;
    //char man_str[20];
    //char prod_str[18];
    int num_matching_devs;
    int frame_width = 129;
    char frame[frame_width+2]; // leave space for newline and NULL char
    int cur_bus_num;
    static int num_devs = 0; // track num_devs across mult invocations of this func
    static int bus_num = -1; // track the locID bus_num and display frame line when it changes

    pdev_info = plog_dev_info->pdev_info;
    pdev_desc = &pdev_info->dev_desc;
    pdev_addr = &pdev_info->dev_addr;

    num_matching_devs = plog_dev_info->pdev_info->pusbexp->num_matching_devs;
    num_devs++;

    usbexp_get_manufacturer_str(plog_dev_info, pdev_desc, man_str, sizeof(man_str) - 1);
    usbexp_get_product_str(plog_dev_info, pdev_desc, prod_str, sizeof(prod_str) - 1);
    usbexp_get_frame_str(frame, frame_width);

    cur_bus_num = (pdev_info->loc_id >> 24) & 0xff;
    if (cur_bus_num != bus_num) {
        bus_num = cur_bus_num;
        if (!papp_settings->show_descriptors) {
            printf("%s", frame);
        }
    }

    if (papp_settings->show_descriptors) {
        printf("%s", frame);
    }

    if (num_devs == 1 || papp_settings->show_descriptors) {
        printf("|     Manufacturer    |       Product      |  VID    PID   | Location ID |    Bus    | Dev Addr | Port | Hub | CID |   Speed   |\n");
        printf("%s", frame);
    }

    char cid_idx_str[4];
    if (pdev_info->container_id_idx >= 0) {
        snprintf(cid_idx_str, 3, "%d", pdev_info->container_id_idx);
    } else {
        strcpy(cid_idx_str, "   ");
    }

    printf("| %-19s | %-18s | 0x%04X 0x%04X | 0x%08x  | 0x%x (%u)%s |    %2u    |  %2u  | %3u |%3s  | %3s bps |\n",
           man_str, prod_str, pdev_desc->idVendor, pdev_desc->idProduct, pdev_info->loc_id, pdev_addr->bus, pdev_addr->bus,
           (pdev_addr->bus < 10) ? "  " : "", pdev_addr->dev, pdev_info->port_num, plog_dev_info->parent_hub_addr, cid_idx_str, plog_dev_info->speed_str);

    if (num_devs == num_matching_devs || papp_settings->show_descriptors) {
        printf("%s", frame);
    }
}
*/
