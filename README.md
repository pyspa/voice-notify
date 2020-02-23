# Voice Notify

Voice-Notify is an application that reports desktop notifications by voice.

Voice-Notify tested only Linux.

### Features

- Speak smoothly with Google Cloud Text-to-Speech
- Support D-Bus notifications
- Support Pushbullet notification mirroring for Android and iOS

### Install

```
go get -u github.com/pyspa/voice-notify
```

### Requires

- Cloud Text-to-Speech API service account key
- Pushbullet API token (optional)
- mpg123

Voice-Notify uses Google [Cloud Text-to-Speech][1] API.
First, create a service account key referring to the [document][2].

Also, it supports notification mirroring via [Pushbullet's Real-time Event Stream API][3].
To use this, generate an access token according to the [Pushbullet documentation][4].

### Usage

```
Usage:
  voice-notify [flags]

Flags:
      --config string   config file (default "/home/xxxx/.config/voice-notify.toml")
  -h, --help            help for voice-notify
  -t, --toggle          help message for toggle
      --version         version for voice-notify
```

### Configuration

Voice-Notify changes the configuration in the configuration file or environment settings.
The default settings are as follows.

```
[speech]
lang = "ja-JP"
speaking_rate = 1.5
pitch = 1.5
# Maximum number of characters to speak at one time
text_max = 256

[log]
debug = false
log = stderr
log.err_log = stderr

[pushbullet]
access_token = ""

```

Voice-Notify uses [viper][5].
If you want to configure with environment values, see the [viper documentation][5].

[1]:https://cloud.google.com/text-to-speech/
[2]:https://cloud.google.com/text-to-speech/docs
[3]:https://docs.pushbullet.com/#realtime-event-stream
[4]:https://docs.pushbullet.com/
[5]:https://github.com/spf13/viper
