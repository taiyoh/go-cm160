package cm160

import "bytes"

// FrameLength is data size
const FrameLength int = 11

type bulkResponse struct {
	raw   []byte
	reply uint8
}

// system message from device
var (
	IDMsg   = []byte{0xA9, 0x49, 0x44, 0x54, 0x43, 0x4D, 0x56, 0x30, 0x30, 0x31, 0x01}
	WAITMsg = []byte{0xA9, 0x49, 0x44, 0x54, 0x57, 0x41, 0x49, 0x54, 0x50, 0x43, 0x52}
)

// NewBulkResponse returns bulkResponse
func NewBulkResponse(raw []byte) *bulkResponse {
	var reply uint8 // it's magic word
	switch {
	case bytes.Compare(IDMsg, raw) == 0:
		reply = 0x5A
	case bytes.Compare(WAITMsg, raw) == 0:
		reply = 0xA5
	default:
		reply = 0x0
	}
	return &bulkResponse{raw: raw, reply: reply}
}

func (r *bulkResponse) BuildRecord() *Record {
	return NewRecord(r)
}

func (r *bulkResponse) IsValid() bool {
	checksum := 0x00
	buflen := FrameLength - 1
	for i := 0; i < buflen; i++ {
		checksum += int(r.raw[i])
	}
	checksum &= 0xff
	return checksum == int(r.raw[10])
}

func (r *bulkResponse) NeedToReply() bool {
	return r.reply != 0x0
}

func (r *bulkResponse) Reply() uint8 {
	return r.reply
}
