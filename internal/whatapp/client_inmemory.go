package whatsapp

type InMemoryClient struct {
}

func NewInMemoryClient() Client {

	return &InMemoryClient{}
}

func (c *InMemoryClient) Send(phone string, body string) error {
	return nil
}
