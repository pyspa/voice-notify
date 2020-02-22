package dbus

import (
	"context"
	"encoding/json"
	"time"
	"voice-notify/log"
	"voice-notify/service"
	"voice-notify/tts"

	"github.com/godbus/dbus/v5"
	"github.com/pkg/errors"
)

type notification struct {
	Type int           `json:"Type"`
	Body []interface{} `json:"Body"`
}

type DbusService struct {
}

func (s *DbusService) Start() {
	ctx := context.Background()

	conn, err := dbus.SessionBus()
	if err != nil {
		log.Error(errors.Wrap(err, "failed connect d-bus"), nil)
		return
	}
	var rules = []string{
		"type='method_call',member='Notify',path='/org/freedesktop/Notifications',interface='org.freedesktop.Notifications'",
	}
	var flag uint = 0

	call := conn.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, rules, flag)
	if call.Err != nil {
		log.Error(errors.Wrap(call.Err, "failed to become monitor"), nil)
		return
	}

	c := make(chan *dbus.Message, 64)
	conn.Eavesdrop(c)
	log.Debug("Monitoring notifications", nil)
	for v := range c {
		data, err := json.Marshal(v)
		if err != nil {
			log.Error(errors.Wrap(err, "failed marshall"), nil)
			return
		}
		// log.Debug(string(data), nil)

		var n notification
		if err := json.Unmarshal(data, &n); err != nil {
			log.Error(errors.Wrap(err, "failed unmarshall"), nil)
			return
		}

		if n.Type == 1 {
			app := n.Body[0].(string)
			title := n.Body[3].(string)
			message := n.Body[4].(string)

			log.Debug("dbus", log.Fields{
				"app":     app,
				"title":   title,
				"message": message,
			})
			if err := speech(ctx, app, title, message); err != nil {
				log.Error(err, nil)
				return
			}

		}
	}
}

func speech(ctx context.Context, app, title, msg string) error {
	if err := tts.Speech(ctx, title); err != nil {
		return errors.Wrap(err, "failed speech.")
	}
	if err := tts.Speech(ctx, msg); err != nil {
		return errors.Wrap(err, "failed speech.")
	}
	time.Sleep(time.Second * 1)
	return nil
}

func NewService() service.Service {

	return &DbusService{}
}
