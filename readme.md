# Mandarin Clipboard Speaker

Hanzi (汉字 / Chinese characters) are difficult to memorize; this program automatically speaks them as Mandarin when they are copied to the system clipboard on Linux.

## Running

Clone the repo then run `go run main.go`

## Installing as a systemd service

For automatic start on boot via systemd, run `make install` with the included makefile

## Techncial Details

This program runs and installs the [piper](https://rhasspy.github.io/piper-samples/) text-to-speech model for Chinese.  It currently installs the model and binary to the same location as [QuickPiperAudiobook](https://github.com/C-Loftus/QuickPiperAudiobook) but in the future once Piper is speech dispatcher, I can use that instead.