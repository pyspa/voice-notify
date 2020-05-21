package tts

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"github.com/pkg/errors"
	"github.com/pyspa/voice-notify/log"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var mutex sync.Mutex

func Speech(ctx context.Context, text string) error {
	if text == "" {
		return nil
	}
	mutex.Lock()
	defer mutex.Unlock()

	max := viper.GetInt("speech.text_max")
	spText := []rune(text)
	if len(spText) > max {
		text = string(spText[:max])
	}

	opts := make([]option.ClientOption, 0)
	cred := viper.GetString("speech.credentials")
	if cred != "" {
		opts = append(opts, option.WithCredentialsFile(cred))
	}

	client, err := texttospeech.NewClient(ctx, opts...)
	if err != nil {
		return errors.Wrap(err, "failed create client")
	}
	lang := viper.GetString("speech.lang")
	rate := viper.GetFloat64("speech.speaking_rate")
	pitch := viper.GetFloat64("speech.pitch")
	req := texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},

		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: lang,
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},

		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding:    texttospeechpb.AudioEncoding_MP3,
			SpeakingRate:     rate,
			Pitch:            pitch,
			EffectsProfileId: []string{"headphone-class-device"},
		},
	}

	log.Debug("say", log.Fields{
		"Text": text,
	})

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		return errors.Wrap(err, "failed call tts api")
	}

	out, err := ioutil.TempFile("", "tts")
	if err != nil {
		return errors.Wrap(err, "failed create tempfile")
	}

	defer func() {
		out.Close()
		os.Remove(out.Name())
	}()

	if err := ioutil.WriteFile(out.Name(), resp.AudioContent, 0644); err != nil {
		return errors.Wrap(err, "failed write contents")
	}
	cmd := viper.GetString("speech.play_cmd")
	if err := exec.Command(cmd, out.Name()).Run(); err != nil {
		return errors.Wrap(err, "failed play")
	}

	return nil
}
