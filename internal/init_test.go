package internal_test

import (
	"log"
	"os"
)

func init() {

	fatalOnSetEnvErr(os.Setenv("WHATSAPP_VERIFICATION_TOKEN", "123"))
	fatalOnSetEnvErr(os.Setenv("SLACK_INFO_URL", "123"))
	fatalOnSetEnvErr(os.Setenv("SLACK_DEBUG_URL", "567"))
}

func fatalOnSetEnvErr(err error) {

	if err != nil {
		log.Fatal("failed to set env")
	}
}
