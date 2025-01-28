package dtos

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type DataPacket struct {
	UserID      string
	CommandName string
	metaData    map[string]string
}

func (d *DataPacket) toByte() ([]byte, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		logrus.Errorf("Failed to marshal message: %v", err)
		return nil, err
	}
	return bytes, nil
}

func (d *DataPacket) fromByte(bytes []byte) error {
	err := json.Unmarshal(bytes, d)
	if err != nil {
		logrus.Errorf("Failed to unmarshal message: %v", err)
		return err
	}
	return nil
}
