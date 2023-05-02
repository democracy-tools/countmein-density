package whatsapp

type InMemoryClient struct {
}

func NewInMemoryClient() Client {

	return &InMemoryClient{}
}

func (c *InMemoryClient) Send(phone string, body string) error {
	return nil
}

func (c *InMemoryClient) SendOnboardingTemplate(phone string, userId string) error { return nil }

func (c *InMemoryClient) SendInvitationTemplate(to string) error {
	return nil
}

func (c *InMemoryClient) SendRegretInvitationTemplate(to string) error {
	return nil
}

func (c *InMemoryClient) SendDemonstrationTemplate(to string, userId string) error {
	return nil
}

func (c *InMemoryClient) SendBodyParamsTemplate(template string, to string, params []string) error {
	return nil
}
