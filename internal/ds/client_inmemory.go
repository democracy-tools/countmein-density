package ds

import "reflect"

type InMemoryClient struct {
	getByTimeDst interface{}
}

func NewInMemoryClient() Client {

	return &InMemoryClient{}
}

func (c *InMemoryClient) SetGetByTimeDst(dst interface{}) { c.getByTimeDst = dst }

func (c *InMemoryClient) GetByTime(kind Kind, from int64, dst interface{}) error {

	if c.getByTimeDst != nil {
		reflect.ValueOf(dst).Elem().Set(reflect.ValueOf(c.getByTimeDst))
	}

	return nil
}

func (c *InMemoryClient) Put(kind Kind, id string, src interface{}) error {

	return nil
}
