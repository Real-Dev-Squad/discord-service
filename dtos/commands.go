package dtos

type CommandNameTypes struct {
	Hello     string
	Listening string
	Verify    string
}

type CommandError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Success bool   `json:"success"`
}
