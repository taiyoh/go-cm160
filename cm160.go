package cm160

import (
	"bytes"
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

	FRAME_ID_LIVE uint8 = 0x51
	FRAME_ID_DB   uint8 = 0x59

	USB_RECIP_INTERFACE uint8 = 0x01

	FRAME_LENGTH int = 11

	// MSG_ID   int = 0x5A
	// MSG_WAIT int = 0xA5
)

const PIPE int = int(libusb.USB_TYPE_VENDOR | USB_RECIP_INTERFACE | libusb.USB_ENDPOINT_OUT)

type cm160 struct {
	device    *libusb.Device
	isRunning bool
}

type Record struct {
	Volt   int
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
	Cost   float32
	Amps   float32
	IsLive bool
	Watt   float32
}

type ctrlMsg struct {
	req   int
	value int
	data  []byte
}

type bulkResponse struct {
	buffer []byte
	length int
}

type ctrlMsgs []ctrlMsg

func (cmds ctrlMsgs) Each(cb func(c ctrlMsg)) {
	for i := 0; i < len(cmds); i++ {
		cb(cmds[i])
	}
}

func Open() *cm160 {

	libusb.Init()
	dev := libusb.Open(0x0fde, 0xca05)
	if dev == nil {
		log.Fatalln("device not found")
	}

	if r := dev.DetachKernelDriver(IFACE_ID); r != 0 {
		log.Fatalf("usb_detach_kernel_driver_np returns %d (%s)\n", r, libusb.LastError())
	}
	if r := dev.Configuration(1); r != 0 {
		log.Fatalf("usb_set_configuration returns %d (%s)\n", r, libusb.LastError())
	}
	if r := dev.Interface(IFACE_ID); r != 0 {
		log.Fatalf("Interface cannot be claimed: %d (%s)\n", r, libusb.LastError())
	}

	brate := make([]byte, 4)
	binary.LittleEndian.PutUint32(brate, uint32(250000))
	ctrlMsgs{
		ctrlMsg{CP210X_IFC_ENABLE, UART_ENABLE, make([]byte, 0)},
		ctrlMsg{CP210X_SET_BAUDRATE, 0, brate},
		ctrlMsg{CP210X_IFC_ENABLE, UART_DISABLE, make([]byte, 0)},
	}.Each(func(c ctrlMsg) {
		if res := dev.ControlMsg(PIPE, c.req, c.value, 0, c.data); res < 0 {
			log.Fatalf("[%#v:%#v] error: %#v, %s\n", c.req, c.value, res, c.data)
			// } else {
			// 	log.Printf("[%#v:%#v] ok, %s (%d)", c.req, c.value, c.data, res)
		}
	})

	return &cm160{device: dev}
}

func (self *bulkResponse) ParseFrame() *Record {
	return &Record{
		Year:   int(self.buffer[1]) + 2000,
		Month:  int(self.buffer[2] & 0x0f), // 0xcを期待してるのに0xccって返ってくることがある
		Day:    int(self.buffer[3]),
		Hour:   int(self.buffer[4]),
		Minute: int(self.buffer[5]),
		Cost:   float32(int(self.buffer[6])+(int(self.buffer[7])<<8)) / 100.0,
		Amps:   float32(int(self.buffer[8])+(int(self.buffer[9]))) * 0.07,
		IsLive: self.buffer[0] == FRAME_ID_LIVE,
	}
	// rec.Watt = float32(rec.Volt) * rec.Amps
}

func (self *Record) CalcWatt(volt int) {
	self.Volt = volt
	self.Watt = float32(volt) * self.Amps
}

var (
	ID_MSG   []byte = []byte{0xA9, 0x49, 0x44, 0x54, 0x43, 0x4D, 0x56, 0x30, 0x30, 0x31, 0x01}
	WAIT_MSG []byte = []byte{0xA9, 0x49, 0x44, 0x54, 0x57, 0x41, 0x49, 0x54, 0x50, 0x43, 0x52}
)

func (self *bulkResponse) MsgToSend() uint8 {
	switch {
	case bytes.Compare(ID_MSG, self.buffer) == 0:
		return 0x5A
	case bytes.Compare(WAIT_MSG, self.buffer) == 0:
		return 0xA5
	default:
		return 0x00
	}
}

func (self *bulkResponse) IsValid() bool {
	checksum := 0x00
	for i := 0; i < 10; i++ {
		checksum += int(self.buffer[i])
	}
	checksum &= 0xff
	return checksum == int(self.buffer[10])
}

func (self *cm160) BulkRead() ([]byte, int) {
	buf := make([]byte, FRAME_LENGTH)
	res_len := self.device.BulkRead(ENDPOINT_IN, buf)
	// fmt.Printf("cnt, err, buf: %d, %#v, %s\n", cnt, _buf, self.device.LastError())
	return buf, res_len
}

func (self *cm160) BulkWrite(b uint8) int {
	return self.device.BulkWrite(ENDPOINT_OUT, []byte{b})
}

func (self *cm160) Stop() {
	self.isRunning = false
}

func (self *cm160) Wait(cb func(buf *Record)) {

	ch1 := make(chan *bulkResponse)
	ch2 := make(chan bool)

	Read := func() {
		var res *bulkResponse = nil
		if buf, l := self.BulkRead(); l <= 0 {
			// log.Printf("[%s] length = %d, %s\n", self.device.LastError(), l, buf)
		} else {
			// log.Printf("BulkRead result (%d): %#v\n", l, buf[2])
			res = &bulkResponse{buffer: buf, length: l}
		}
		ch1 <- res
	}

	Proc := func(res *bulkResponse) {
		if msg := res.MsgToSend(); msg != 0x00 {
			self.BulkWrite(msg)
		} else if res.IsValid() {
			cb(res.ParseFrame())
		}
		ch2 <- true
	}

	self.isRunning = true
	// main loop
	for {
		go Read()
		res := <-ch1
		if res == nil {
			// time.Sleep(500 * 100000000) // 0.5sec?
			continue
		}
		go Proc(res)
		<-ch2
		if !self.isRunning {
			break
		}
	}
}

func (self *cm160) Close() {
	if r := self.device.ReleaseDevice(IFACE_ID); r < 0 {
		log.Printf("usb_release_device error: %d (%s)\n", r, libusb.LastError())
	}
	if r := self.device.Close(); r < 0 {
		log.Printf("usb_close error: %d (%s)\n", r, libusb.LastError())
	}
}
