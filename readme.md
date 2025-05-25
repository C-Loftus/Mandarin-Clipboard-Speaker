# Mandarin Clipboard Speaker

Hanzi (汉字 / Chinese characters) are difficult to memorize for non-native speakers; this program automatically speaks them as Mandarin when they are copied to the system clipboard on Linux.

## Running

Clone the repo then run `go run main.go`

## Installing as a systemd service

For automatic start on boot via systemd, run `make install` with the included makefile

## Technical Details

This program runs and installs the [Piper](https://rhasspy.github.io/piper-samples/) text-to-speech model for Mandarin Chinese.  It currently installs the model and binary to the same location as [QuickPiperAudiobook](https://github.com/C-Loftus/QuickPiperAudiobook) but in the future once Piper is included within [speech dispatcher](https://github.com/brailcom/speechd), I can use that instead.


### Quickstart Video

[![Mandarin-Clipboard-Speaker Quickstart](./docs/video_thumbnail.jpg)](https://youtu.be/Ax4buJ-f4Jg "Mandarin-Clipboard-Speaker Quickstart")