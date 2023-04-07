package whatsapp

type TemplateMessageRequest struct {
	To       string   `json:"to"`
	Type     string   `json:"type"`
	Template Template `json:"template"`
}

type Language struct {
	Policy string `json:"policy"`
	Code   string `json:"code"`
}

type Parameters struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Components struct {
	Type       string       `json:"type"`
	Parameters []Parameters `json:"parameters"`
	SubType    string       `json:"sub_type,omitempty"`
	Index      string       `json:"index,omitempty"`
}

type Template struct {
	Namespace  string       `json:"namespace"`
	Language   Language     `json:"language"`
	Name       string       `json:"name"`
	Components []Components `json:"components"`
}
