package ds

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/datastore"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/sabik/common/env"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Kind string

const (
	KindObservation     Kind = "observation"
	KindRegisterRequest Kind = "register_request"
	KindUser            Kind = "user"
	KindVolunteer       Kind = "volunteer"
	KindPreference      Kind = "preference"

	EnvKeyDatastoreToken = "DATASTORE_KEY"
	namespace            = "v1"
)

type FilterField struct {
	Name     string
	Operator string
	Value    interface{}
}

type Client interface {
	Get(kind Kind, id string, dst interface{}) error
	GetAll(kind Kind, dst interface{}) error
	GetFilter(kind Kind, filters []FilterField, dst interface{}) error
	Put(kind Kind, id string, src interface{}) error
}

type ClientWrapper struct {
	ds *datastore.Client
}

func NewClientWrapper(project string) Client {

	if key := env.GetEnvSensitive(EnvKeyDatastoreToken); key != "" {
		conf, err := google.JWTConfigFromJSON([]byte(key), datastore.ScopeDatastore)
		if err != nil {
			log.Fatalf("failed to config datastore JWT from JSON key with '%v'", err)
		}

		ctx := context.Background()
		opt := []option.ClientOption{option.WithTokenSource(conf.TokenSource(ctx))}

		dataStoreEndPoint := os.Getenv("DATASTORE_ENDPOINT")
		if dataStoreEndPoint != "" {
			opt = append(opt, option.WithEndpoint(dataStoreEndPoint))
		}

		client, err := datastore.NewClient(ctx, project, opt...)
		if err != nil {
			log.Fatalf("failed to create datastore client with '%v'", err)
		}

		return &ClientWrapper{ds: client}
	}

	client, err := datastore.NewClient(context.Background(), project)
	if err != nil {
		log.Fatalf("failed to create datastore client without token with '%v'", err)
	}

	return &ClientWrapper{ds: client}
}

func (c *ClientWrapper) Get(kind Kind, id string, dst interface{}) error {

	err := c.ds.Get(context.Background(), getKey(kind, id), dst)
	if err != nil {
		log.Errorf("failed to get '%s' id '%s' from datastore namespace '%s' with '%v'", kind, id, namespace, err)
	}

	return err
}

func (c *ClientWrapper) GetAll(kind Kind, dst interface{}) error {

	q := datastore.NewQuery(string(kind)).Namespace(namespace)
	_, err := c.ds.GetAll(context.Background(), q, dst)
	if err != nil {
		log.Errorf("failed to get all '%s' from datastore namespace '%s' with '%v'", kind, namespace, err)
	}

	return err
}

func (c *ClientWrapper) GetFilter(kind Kind, filters []FilterField, dst interface{}) error {

	q := datastore.NewQuery(string(kind)).Namespace(namespace)
	for _, currFilter := range filters {
		q = q.FilterField(currFilter.Name, currFilter.Operator, currFilter.Value)
	}
	_, err := c.ds.GetAll(context.Background(), q, dst)
	if err != nil {
		msg := fmt.Sprintf("failed to get with filter '%+v' kind '%s' from datastore namespace '%s' with '%v'",
			filters, kind, namespace, err)
		if IsNoSuchEntityError(err) {
			log.Debug(msg)
		} else {
			log.Error(msg)
		}
	}

	return err
}

func (c *ClientWrapper) Put(kind Kind, id string, src interface{}) error {

	_, err := c.ds.Put(context.Background(), getKey(kind, id), src)
	if err != nil {
		log.Errorf("failed to create '%s/%s' item '%+v' type: '%T' with '%v'", namespace, kind, src, src, err)
	} else {
		log.Debugf("Item '%s/%s' created: '%+v' type: '%T'", namespace, kind, src, src)
	}

	return err
}

func getKey(kind Kind, id string) *datastore.Key {

	res := datastore.NameKey(string(kind), id, nil)
	res.Namespace = namespace

	return res
}
