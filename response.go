package cm160

import "bytes"

// FrameLength : chunked data size
const FrameLength int = 11

// BulkResponse : chunked response from usb device
type BulkResponse struct {
	raw   []byte
	reply uint8
}

// system message from device
var (
	IDMsg   = []byte{0xA9, 0x49, 0x44, 0x54, 0x43, 0x4D, 0x56, 0x30, 0x30, 0x31, 0x01}
	WAITMsg = []byte{0xA9, 0x49, 0x44, 0x54, 0x57, 0x41, 0x49, 0x54, 0x50, 0x43, 0x52}
)

// NewBulkResponse : returns new BulkResponse
func NewBulkResponse(raw []byte) *BulkResponse {
	var reply uint8 // it's magic word
	switch {
	case bytes.Compare(IDMsg, raw) == 0:
		reply = 0x5A
	case bytes.Compare(WAITMsg, raw) == 0:
		reply = 0xA5
	default:
		reply = 0x0
	}
	return &BulkResponse{raw: raw, reply: reply}
}

// BuildRecord : returns new Record
func (r *BulkResponse) BuildRecord() *Record {
	return NewRecord(r)
}

// IsValid : returns that response is either valid or not
func (r *BulkResponse) IsValid() bool {
	checksum := 0x00
	buflen := FrameLength - 1
	for i := 0; i < buflen; i++ {
		checksum += int(r.raw[i])
	}
	checksum &= 0xff
	return checksum == int(r.raw[10])
}

// NeedToReply : returns that response is either system message or not
func (r *BulkResponse) NeedToReply() bool {
	return r.reply != 0x0
}

// Reply : returns next message to send to usb device
func (r *BulkResponse) Reply() uint8 {
	return r.reply
}
