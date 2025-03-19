package dtos

type DataPacket struct {
	UserID      string
	CommandName string
	MetaData    map[string]string
}
