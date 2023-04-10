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

func CreateDemonstration() error {

	dsc := ds.NewClientWrapper(env.Project)

	id := createDemonstrationInDatastore(dsc)
	return inviteVolunteers(dsc, id)
}

func inviteVolunteers(dsc ds.Client, id string) error {

	var users []ds.User
	err := dsc.GetAll(ds.KindUser, &users)
	if err != nil {
		return err
	}

	wac := whatsapp.NewClientWrapper()
	log.Info("Sending invitations...")
	for _, currUser := range users {
		link := fmt.Sprintf("%s?user=%s&demonstration=%s", internal.JoinUrl, currUser.Id, id)
		log.Infof("%s (%s): %s, %s", currUser.Name, currUser.Id, currUser.Phone, link)
		err = wac.SendInvitationTemplate(currUser.Phone, id, currUser.Id)
		if err != nil {
			return err
		}

	}

	return nil
}

func createDemonstrationInDatastore(ds.Client) string {

	id := uuid.NewString()
	log.Infof("Creating demonstration '%s'...", id)

	return id
}
