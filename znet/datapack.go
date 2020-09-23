package znet

import (
	"bytes"
	"encoding/binary"
	"go-zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//datalen unit32 + id uint32 = 4+4
	return 8
}

//写
func (d *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个存放byte 的缓冲
	var buf []byte
	buffer := bytes.NewBuffer(buf)

	err := binary.Write(buffer, binary.LittleEndian, msg.GetMsgLen())

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, msg.GetMsgId())

	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, msg.GetMsgData())

	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

//读
func (d *DataPack) Unpack(bin []byte) (ziface.IMessage, error) {

	buf := bytes.NewReader(bin)

	msg := &Message{}

	if err := binary.Read(buf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	return msg, nil
}
