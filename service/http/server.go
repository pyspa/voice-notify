package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/pyspa/voice-notify/log"
	"github.com/pyspa/voice-notify/service"
	"github.com/pyspa/voice-notify/tts"
	"github.com/spf13/viper"
)

type HttpService struct {
}

func (s *HttpService) Start() {
	http.HandleFunc("/", handler)
	addr := viper.GetString("http.addr")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Error(errors.Wrapf(err, "failed listen addr=%s", addr), nil)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Error(errors.Wrap(err, "failed parse form"), nil)
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	text := r.FormValue("text")
	ctx := context.Background()
	if err := tts.Speech(ctx, text); err != nil {
		log.Error(errors.Wrap(err, "failed speech."), nil)
		fmt.Fprintf(w, "Speech err: %v", err)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "OK")
}

func NewService() service.Service {
	return &HttpService{}
}
