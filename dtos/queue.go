package dtos

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type DataPacket struct {
	UserID      string
	CommandName string
	MetaData    map[string]string
}

func (d *DataPacket) ToByte() ([]byte, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		logrus.Errorf("Failed to marshal message: %v", err)
		return nil, err
	}
	return bytes, nil
}

func (d *DataPacket) FromByte(bytes []byte) error {
	err := json.Unmarshal(bytes, d)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	return nil
}
