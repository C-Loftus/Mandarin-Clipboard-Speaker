BINARY_NAME := mandarin-clipboard-speaker
INSTALL_DIR := $(HOME)/.local/bin
SERVICE_FILE := $(BINARY_NAME).service
USER_SYSTEMD_DIR := $(HOME)/.config/systemd/user

.PHONY: all build install service uninstall clean logs

all: build

# build the go binary
build:
	go build -o $(BINARY_NAME) .

# install the systemd service
install: build
	systemd-analyze verify mandarin-clipboard-speaker.service
	mkdir -p $(INSTALL_DIR)

	# stop the systemd service if it is already installed and running
	# that way the service can be reinstalled
	systemctl --user stop $(SERVICE_FILE) || true
	
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	mkdir -p $(USER_SYSTEMD_DIR)
	cp $(SERVICE_FILE) $(USER_SYSTEMD_DIR)/$(SERVICE_FILE)
	systemctl --user daemon-reexec
	systemctl --user enable $(SERVICE_FILE)
	systemctl --user start $(SERVICE_FILE)

# remove the systemd service
uninstall:
	systemctl --user stop $(SERVICE_FILE) || true
	systemctl --user disable $(SERVICE_FILE) || true
	rm -f $(USER_SYSTEMD_DIR)/$(SERVICE_FILE)
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	systemctl --user daemon-reexec

# remove the binary
clean:
	rm -f $(BINARY_NAME)

# tail the logs of the systemd service
logs:
	journalctl --user -u $(SERVICE_FILE) -f
