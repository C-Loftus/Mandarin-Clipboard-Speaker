# Mandarin Clipboard Speaker

- Automatically speak Chinese characters as Mandarin when they are copied to the system clipboard on Linux.
- Runs offline and works in any application with the goal to make 汉字 easier to read for non-native speakers.

## Running

- Clone the repo then run `go run main.go`
- If the local model is not present, it will be installed automatically
- For automatic start on boot via systemd, run `make install` with the included makefile

## Technical Details

This program runs and installs the [Piper](https://rhasspy.github.io/piper-samples/) text-to-speech model for Mandarin Chinese.  It currently installs the model and binary to the same location as [QuickPiperAudiobook](https://github.com/C-Loftus/QuickPiperAudiobook) but in the future once Piper is included within [speech dispatcher](https://github.com/brailcom/speechd), I can use that instead.


### Quickstart Video

[![Mandarin-Clipboard-Speaker Quickstart](./docs/video_thumbnail.jpg)](https://youtu.be/Ax4buJ-f4Jg "Mandarin-Clipboard-Speaker Quickstart")