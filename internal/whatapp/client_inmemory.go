package whatsapp

type InMemoryClient struct {
}

func NewInMemoryClient() Client {

	return &InMemoryClient{}
}

func (c *InMemoryClient) Send(phone string, body string) error {
	return nil
}

func (c *InMemoryClient) SendSignupTemplate(to string, token string) error {
	return nil
}

func (c *InMemoryClient) SendVerifyTemplate(phone string) error { return nil }

func (c *InMemoryClient) SendInvitationTemplate(to string, demonstration string, userId string) error {
	return nil
}

func (c *InMemoryClient) SendDemonstrationTemplate(to string, demonstration string, userId string,
	user string, polygon string, location string) error {
	return nil
}

func (c *InMemoryClient) SendBodyParamsTemplate(template string, to string, params []string) error {
	return nil
}
