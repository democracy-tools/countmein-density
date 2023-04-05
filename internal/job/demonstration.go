package job

import (
	"fmt"

	"github.com/democracy-tools/countmein-density/internal"
	"github.com/democracy-tools/countmein-density/internal/ds"
	"github.com/democracy-tools/countmein-density/internal/env"
	whatsapp "github.com/democracy-tools/countmein-density/internal/whatapp"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func CreateDemonstration() {

	dsc := ds.NewClientWrapper(env.Project)

	id := createDemonstrationInDatastore(dsc)
	inviteVolunteers(dsc, id)
}

func inviteVolunteers(dsc ds.Client, id string) error {

	var users []ds.User
	err := dsc.GetAll(ds.KindUser, &users)
	if err != nil {
		return err
	}

	wac := whatsapp.NewClientWrapper()
	for _, currUser := range users {
		link := fmt.Sprintf("%s?user=%s&demonstration=%s", internal.JoinUrl, currUser.Id, id)
		log.Infof("sending invitation... '%s (%s): %s' link '%s'", currUser.Name, currUser.Id, currUser.Phone, link)
		wac.Send(currUser.Phone, fmt.Sprintf("היי, רוצה להתנדב לספירה ביום שבת? לחץ על הלינק כדי להצטרף אלינו\n%s", link))
	}

	return nil
}

func createDemonstrationInDatastore(ds.Client) string {

	id := uuid.NewString()
	log.Infof("Creating demonstration '%s'...", id)

	return id
}
