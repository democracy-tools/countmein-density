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
	KindObservation      Kind = "observation"
	EnvKeyDatastoreToken      = "DATASTORE_KEY"
	namespace                 = "dev"
)

func IsNoSuchEntityError(err error) bool {

	if err == nil {
		return false
	}

	return err.Error() == datastore.ErrNoSuchEntity.Error()
}

type Client interface {
	GetAll(kind Kind, order string, dst interface{}) error
	Put(kind Kind, id string, src interface{}) error
}

type ClientWrapper struct {
	ds *datastore.Client
}

func NewClientWrapper(gcpProjectID string) Client {

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

		client, err := datastore.NewClient(ctx, gcpProjectID, opt...)
		if err != nil {
			log.Fatalf("failed to create datastore client with '%v'", err)
		}

		return &ClientWrapper{ds: client}
	}

	client, err := datastore.NewClient(context.Background(), gcpProjectID)
	if err != nil {
		log.Fatalf("failed to create datastore client without token with '%v'", err)
	}

	return &ClientWrapper{ds: client}
}

func (c *ClientWrapper) GetAll(kind Kind, order string, dst interface{}) error {

	q := datastore.NewQuery(string(kind)).Namespace(namespace)
	if order != "" {
		q = q.Order(order)
	}
	_, err := c.ds.GetAll(context.Background(), q, dst)
	if err != nil {
		msg := fmt.Sprintf("failed to get all '%s' from datastore namespace '%s' with '%v'", kind, namespace, err)
		if IsNoSuchEntityError(err) {
			log.Info(msg)
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

func getKey(kind Kind, key string) *datastore.Key {

	ret := datastore.NameKey(string(kind), key, nil)
	ret.Namespace = namespace

	return ret
}
