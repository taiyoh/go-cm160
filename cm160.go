package cm160

import "github.com/taiyoh/go-libusb"

// CM160 is root object of this library
type CM160 struct {
	device    *libusb.Device
	isRunning bool
	histories []*Record
}

// Open returns cm160
func Open() *CM160 {

	dev := InitializeDevice(0x0fde, 0xca05)
	return &CM160{device: dev, isRunning: true}
}

// Stop is dropping flag for stopping loop
func (c *CM160) Stop() {
	c.isRunning = false
}

// IsRunning returns flag either running or not
func (c *CM160) IsRunning() bool {
	return c.isRunning
}

// AddRecord appends Record in histories
func (c *CM160) AddRecord(r *Record) {
	c.histories = append(c.histories, r)
}

// ShiftRecord retrieve Record in histories
func (c *CM160) ShiftRecord() *Record {
	var record *Record
	if len(c.histories) > 0 {
		record = c.histories[0]
		c.histories = c.histories[1:]
	}
	return record
}

// IsEmptyRecords returns either empty or not
func (c *CM160) IsEmptyRecords() bool {
	return len(c.histories) == 0
}

func (c *CM160) Read() *Record {
	if c.IsEmptyRecords() {
		for {
			responses := c.ReadFromDevice()
			for _, res := range responses {
				if res.NeedToReply() {
					c.WriteToDevice(res.Reply())
				} else if res.IsValid() {
					c.AddRecord(res.BuildRecord())
				}
			}
			if !c.IsEmptyRecords() {
				break
			}
		}
	}
	return c.ShiftRecord()
}
