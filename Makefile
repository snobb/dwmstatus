TARGET := dwmstatus
MAIN := ./cmd/main.go
BIN := ./bin
INSTALL := install
INSTALL_ARGS := -o root -g root -m 755
INSTALL_DIR := /usr/local/bin/

# can do "sudo GO=$(which go) make build install clean" if Go is installed in the user folder.
GO ?= go

# autoconfiguration
# Battery status if exists.
BATPATH := ${strip ${shell find /sys -name "BAT[0-9]" -print0 -quit}}

# Infer the wifi interface name - please override here if necessary
IFNAME  := ${shell iw dev | awk '/Interface/ { print $$2 }' | tr -d '\n'}
LNKPATH := ${shell find /sys/class/net/${IFNAME}/ -name operstate -print0 -quit}

LDFLAGS := -X main.batPath=${BATPATH} -X main.wifiPath=${LNKPATH}

all: build

release: LDFLAGS += -s
release: build

build:
	${GO} build --ldflags '${LDFLAGS}' -o ${BIN}/${TARGET} ${MAIN}

install:
	${INSTALL} ${INSTALL_ARGS} ${BIN}/${TARGET} ${INSTALL_DIR}
	@echo "DONE"

clean:
	-rm -rf ${BIN}
