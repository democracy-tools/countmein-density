package slack

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Send(channelHook string, message string) error {

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(struct {
		Text string `json:"text"`
	}{Text: message})
	if err != nil {
		logrus.Errorf("failed to encode message '%s' with '%v'", message, err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, channelHook, &buf)
	if err != nil {
		logrus.Errorf("failed to create slack message request '%s' with '%v'", message, err)
		return err
	}
	req.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("failed to send message to slack '%s' with '%v'", message, err)
		return err
	}
	if http.StatusOK != resp.StatusCode {
		logrus.Errorf("failed to send slack message '%s' with '%s'", message, resp.Status)
		return err
	}

	return nil
}
