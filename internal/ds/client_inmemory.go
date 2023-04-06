package ds

type InMemoryClient struct {
	getFilterDelegate func(kind Kind, filterFieldName string, filterOperator string, filterValue interface{}, dst interface{}) error
}

func NewInMemoryClient() Client {

	return &InMemoryClient{getFilterDelegate: func(kind Kind, filterFieldName string, filterOperator string, filterValue interface{}, dst interface{}) error {
		return nil
	}}
}

func (c *InMemoryClient) Get(kind Kind, id string, dst interface{}) error { return nil }

func (c *InMemoryClient) GetAll(kind Kind, dst interface{}) error { return nil }

func (c *InMemoryClient) SetGetFilterDelegate(delegate func(kind Kind, filterFieldName string, filterOperator string, filterValue interface{}, dst interface{}) error) {

	c.getFilterDelegate = delegate
}

func (c *InMemoryClient) GetFilter(kind Kind, filterFieldName string, filterOperator string, filterValue interface{}, dst interface{}) error {

	return c.getFilterDelegate(kind, filterFieldName, filterOperator, filterValue, dst)
}

func (c *InMemoryClient) Put(kind Kind, id string, src interface{}) error { return nil }
