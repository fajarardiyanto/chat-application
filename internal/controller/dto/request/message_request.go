package request

type MessageRequest struct {
	Content   string `json:"content"`
	Documents []struct {
		DocumentId   string `json:"documentId"`
		DocumentName string `json:"documentName"`
	} `json:"documents"`
	ContentType string `json:"contentType"`
}
