package ds

import "cloud.google.com/go/datastore"

func IsNoSuchEntityError(err error) bool {

	if err == nil {
		return false
	}

	return err.Error() == datastore.ErrNoSuchEntity.Error()
}
