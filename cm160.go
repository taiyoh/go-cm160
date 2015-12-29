package cm160

import "github.com/taiyoh/go-libusb"

// CM160 : client for device handling
type CM160 struct {
	device    *libusb.Device
	isRunning bool
	records   []*Record
}

// Open : returns new CM160
func Open() *CM160 {
	dev := InitializeDevice(0x0fde, 0xca05)
	return &CM160{device: dev, isRunning: true}
}

// Stop : drops flag for stopping loop
func (c *CM160) Stop() {
	c.isRunning = false
}

// IsRunning : returns that process shoud either run or not
func (c *CM160) IsRunning() bool {
	return c.isRunning
}

// AddRecord : appends Record to records
func (c *CM160) AddRecord(r *Record) {
	c.records = append(c.records, r)
}

// ShiftRecord : retrieves Record from records
func (c *CM160) ShiftRecord() *Record {
	var record *Record
	if len(c.records) > 0 {
		record = c.records[0]
		c.records = c.records[1:]
	}
	return record
}

// IsEmptyRecords : returns that records is either empty or not
func (c *CM160) IsEmptyRecords() bool {
	return len(c.records) == 0
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
