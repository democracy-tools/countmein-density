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
