package pb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/pyspa/voice-notify/log"
	"github.com/pyspa/voice-notify/service"
	"github.com/pyspa/voice-notify/tts"
	"github.com/spf13/viper"
)

type PbService struct {
	AccessToken string
}

type pushMessage struct {
	Type string          `json:"type"`
	Push json.RawMessage `json:"push"`
}

type mirrorMessage struct {
	Type            string `json:"type"`
	Title           string `json:"title"`
	Body            string `json:"body"`
	ApplicationName string `json:"application_name"`
	PackageName     string `json:"package_name"`
}

func (s *PbService) Start() {
	url := fmt.Sprintf("wss://stream.pushbullet.com/websocket/%s", s.AccessToken)

	ctx := context.Background()
	c, _, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		log.Error(errors.Wrap(err, "failed connect pushbullet"), nil)
		return
	}
	defer c.Close()

	log.Info("Monitoring Pushbullet real-time event stream", nil)
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Error(errors.Wrap(err, "failed read."), nil)
			return
		}

		var pmsg pushMessage
		if err := json.Unmarshal(msg, &pmsg); err != nil {
			log.Error(errors.Wrap(err, "failed json unmarshal."), nil)
			return
		}

		switch pmsg.Type {
		case "push":

			var mm mirrorMessage
			if err := json.Unmarshal(pmsg.Push, &mm); err != nil {
				log.Error(errors.Wrap(err, "failed json unmarshal."), nil)
				return
			}
			log.Debug("recv push", log.Fields{
				"app":   mm.ApplicationName,
				"type":  mm.Type,
				"title": mm.Title,
				"body":  mm.Body,
			})

			switch mm.Type {
			case "mirror":
				if err := speechMirror(ctx, mm); err != nil {
					log.Error(err, nil)
					return
				}

			}
		}
	}

}

func speechMirror(ctx context.Context, mm mirrorMessage) error {
	if err := tts.Speech(ctx, mm.ApplicationName); err != nil {
		return errors.Wrap(err, "failed speech.")
	}
	if err := tts.Speech(ctx, mm.Title); err != nil {
		return errors.Wrap(err, "failed speech.")
	}
	if err := tts.Speech(ctx, mm.Body); err != nil {
		return errors.Wrap(err, "failed speech.")
	}
	return nil
}

func NewService() service.Service {

	token := viper.GetString("pushbullet.access_token")
	log.Debug("pushbullet service", log.Fields{
		"token": token,
	})

	return &PbService{
		AccessToken: token,
	}

}
