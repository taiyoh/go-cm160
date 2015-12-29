package cm160

// FrameID is data flag.
// Live is current data, Db is stored data
const (
	FrameIDLive uint8 = 0x51
	FrameIDDb   uint8 = 0x59
)

// Record is data format from cm160
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

// NewRecord returns Record
func NewRecord(res *BulkResponse) *Record {
	return &Record{
		Year:   int(res.raw[1]) + 2000,
		Month:  int(res.raw[2] & 0x0f), // 0xcを期待してるのに0xccって返ってくることがある
		Day:    int(res.raw[3]),
		Hour:   int(res.raw[4]),
		Minute: int(res.raw[5]),
		Cost:   float32(int(res.raw[6])+(int(res.raw[7])<<8)) / 100.0,
		Amps:   float32(int(res.raw[8])+(int(res.raw[9]))) * 0.07,
		IsLive: res.raw[0] == FrameIDLive,
	}
}

// CalcWatt calculates watt by volt and amps
func (c *Record) CalcWatt(volt int) {
	c.Volt = volt
	c.Watt = float32(volt) * c.Amps
}
