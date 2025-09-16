package dtos

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type DataPacket struct {
	UserID      string `json:"userId"`
	CommandName string `json:"commandName"`
	MetaData    map[string]string `json:"metaData"`
}

var ToByte = func (d *DataPacket) ([]byte, error) {
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

