package cm160

import "bytes"

const (
	FRAME_LENGTH int = 11
)

type bulkResponse struct {
	raw   []byte
	reply uint8
}

var (
	ID_MSG   []byte = []byte{0xA9, 0x49, 0x44, 0x54, 0x43, 0x4D, 0x56, 0x30, 0x30, 0x31, 0x01}
	WAIT_MSG []byte = []byte{0xA9, 0x49, 0x44, 0x54, 0x57, 0x41, 0x49, 0x54, 0x50, 0x43, 0x52}
)

func NewBulkResponse(raw []byte) *bulkResponse {
	var reply uint8
	switch {
	case bytes.Compare(ID_MSG, raw) == 0:
		reply = 0x5A
	case bytes.Compare(WAIT_MSG, raw) == 0:
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
	buflen := FRAME_LENGTH - 1
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
