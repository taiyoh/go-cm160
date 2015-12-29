package cm160

import "github.com/taiyoh/go-libusb"

type cm160 struct {
	device    *libusb.Device
	isRunning bool
	histories []*Record
}

func Open() *cm160 {

	dev := InitializeDevice(0x0fde, 0xca05)
	return &cm160{device: dev, isRunning: true}
}

func (c *cm160) Stop() {
	c.isRunning = false
}

func (c *cm160) IsRunning() bool {
	return c.isRunning
}

func (c *cm160) Read() *Record {
	var record *Record

	if len(c.histories) == 0 {
		for {
			responses := c.ReadFromDevice()
			for _, res := range responses {
				if res.NeedToReply() {
					c.WriteToDevice(res.Reply())
				} else if res.IsValid() {
					c.histories = append(c.histories, res.BuildRecord())
				}
			}
			if len(c.histories) > 0 {
				break
			}
		}
	}
	if len(c.histories) > 0 {
		// shift操作
		record = c.histories[0]
		c.histories = c.histories[1:]
	}
	return record
}
