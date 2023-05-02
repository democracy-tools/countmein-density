package whatsapp

type TextMessageRequest struct {
	MessagingProduct string      `json:"messaging_product"`
	RecipientType    string      `json:"recipient_type"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Text             MessageText `json:"text"`
}

type MessageText struct {
	PreviewURL bool   `json:"preview_url"`
	Body       string `json:"body"`
}
