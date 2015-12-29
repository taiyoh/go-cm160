package cm160

import (
	"encoding/binary"
	"github.com/taiyoh/go-libusb"
	"log"
)

const (
	CP210X_IFC_ENABLE   int = 0x00
	CP210X_GET_LINE_CTL int = 0x04
	CP210X_SET_MHS      int = 0x07
	CP210X_GET_MDMSTS   int = 0x08
	CP210X_GET_FLOW     int = 0x14
	CP210X_GET_BAUDRATE int = 0x1D
	CP210X_SET_BAUDRATE int = 0x1E

	UART_ENABLE  int = 0x0001
	UART_DISABLE int = 0x0000

	IFACE_ID     int = 0
	ENDPOINT_OUT int = 0x01
	ENDPOINT_IN  int = 0x82

	USB_RECIP_INTERFACE uint8 = 0x01
	PIPE                int   = int(libusb.USB_TYPE_VENDOR | USB_RECIP_INTERFACE | libusb.USB_ENDPOINT_OUT)
)

type ctrlMsg struct {
	req   int
	value int
	data  []byte
}

// InitializeDevice : open device and send control messages
func InitializeDevice(vid, pid int) *libusb.Device {
	libusb.Init()
	dev := libusb.Open(vid, pid)
	if dev == nil {
		log.Fatalln("device not found")
	}

	dev.DetachKernelDriver(IFACE_ID)
	// if r := dev.DetachKernelDriver(IFACE_ID); r != 0 {
	// 	log.Fatalf("usb_detach_kernel_driver_np returns %d (%s)\n", r, libusb.LastError())
	// }
	if r := dev.Configuration(1); r != 0 {
		log.Fatalf("usb_set_configuration returns %d (%s)\n", r, libusb.LastError())
	}
	if r := dev.Interface(IFACE_ID); r != 0 {
		log.Fatalf("Interface cannot be claimed: %d (%s)\n", r, libusb.LastError())
	}

	brate := make([]byte, 4)
	binary.LittleEndian.PutUint32(brate, uint32(250000))
	msgs := []ctrlMsg{
		ctrlMsg{CP210X_IFC_ENABLE, UART_ENABLE, make([]byte, 0)},
		ctrlMsg{CP210X_SET_BAUDRATE, 0, brate},
		ctrlMsg{CP210X_IFC_ENABLE, UART_DISABLE, make([]byte, 0)},
	}
	for _, c := range msgs {
		if res := dev.ControlMsg(PIPE, c.req, c.value, 0, c.data); res < 0 {
			log.Fatalf("[%#v:%#v] error: %#v, %s\n", c.req, c.value, res, c.data)
			// } else {
			// log.Printf("[%#v:%#v] ok, %#v (%d)", c.req, c.value, c.data, res)
		}
	}
	return dev
}

// ReadFromDevice : read and chunk data from usb device
func (c *CM160) ReadFromDevice() []*BulkResponse {
	buf := make([]byte, 512)
	reslen := c.device.BulkRead(ENDPOINT_IN, buf)
	looptimes := int(reslen / FrameLength)

	bufptr := 0
	responses := make([]*BulkResponse, looptimes)
	for i := 0; i < looptimes; i++ {
		block := make([]byte, FrameLength)
		for j := 0; j < FrameLength; j++ {
			block[j] = buf[bufptr+j]
		}
		responses[i] = NewBulkResponse(block)
		bufptr += FrameLength
	}

	return responses
}

// WriteToDevice : send 1 byte to usb device
func (c *CM160) WriteToDevice(b uint8) int {
	return c.device.BulkWrite(ENDPOINT_OUT, []byte{b})
}

// Close : clean up
func (c *CM160) Close() {
	if r := c.device.ReleaseDevice(IFACE_ID); r < 0 {
		log.Printf("usb_release_device error: %d (%s)\n", r, libusb.LastError())
	}
	if r := c.device.Close(); r < 0 {
		log.Printf("usb_close error: %d (%s)\n", r, libusb.LastError())
	}
}
