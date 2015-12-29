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

func (self *cm160) Stop() {
	self.isRunning = false
}

func (self *cm160) IsRunning() bool {
	return self.isRunning
}

func (self *cm160) Read() *Record {
	var record *Record

	if len(self.histories) == 0 {
		for {
			responses := self.ReadFromDevice()
			for _, res := range responses {
				if res.NeedToReply() {
					self.WriteToDevice(res.Reply())
				} else if res.IsValid() {
					self.histories = append(self.histories, res.BuildRecord())
				}
			}
			if len(self.histories) > 0 {
				break
			}
		}
	}
	if len(self.histories) > 0 {
		// shift操作
		record = self.histories[0]
		self.histories = self.histories[1:]
	}
	return record
}
