TARGET		:= dwmstatus
MAIN		:= ./cmd/main.go
INSTALL		:= install
INSTALL_ARGS	:= -o root -g root -m 755
INSTALL_DIR	:= /usr/local/bin/

BATPATH	:= $(strip $(shell find /sys -name BAT0 -print0 -quit))
LDFLAGS	:= -X main.batPath=$(BATPATH)
CFLAGS	:= --ldflags '${LDFLAGS}' -o $(TARGET)

build:
	go build $(CFLAGS) $(MAIN)

clean:
	-rm -f dwmstatus

install: build
	$(INSTALL) $(INSTALL_ARGS) $(TARGET) $(INSTALL_DIR)
	@echo "DONE"
