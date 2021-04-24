
SRCS = $(wildcard src/*.c)
OBJS = $(SRCS:.c=.o)
CFLAGS = -Wall -O2 -ffreestanding -nostdinc -nostdlib -nostartfiles -fno-stack-protector
BUILD_DIR := build
BUILD_ABS_DIR := $(CURDIR)/$(BUILD_DIR)
LD := aarch64-linux-gnu-gcc
GO ?= go
GOROOT := $(shell $(GO) env GOROOT)
GOOS := linux
GOARCH := arm64

GOPATH := $(BUILD_ABS_DIR):$(shell pwd):$(GOPATH)

all: clean kernel8.img

arch/arm64/start.o: arch/arm64/start.S
	aarch64-linux-gnu-gcc $(CFLAGS) -c arch/arm64/start.S -o arch/arm64/start.o

main.o: main.go
	@mkdir -p $(BUILD_DIR)
	@echo "[go] compiling go sources into a standalone .o file"
	@GOARCH=$(GOARCH) GOOS=$(GOOS) GOPATH=$(GOPATH) CGO_ENABLED=0 $(GO) build -gcflags '$(GC_FLAGS)' -n 2>&1 | sed \
	    -e "1s|^|set -e\n|" \
	    -e "1s|^|export GOOS=$(GOOS)\n|" \
	    -e "1s|^|export GOARCH=$(GOARCH)\n|" \
	    -e "1s|^|export GOROOT=$(GOROOT)\n|" \
	    -e "1s|^|export CGO_ENABLED=0\n|" \
	    -e "1s|^|alias pack='$(GO) tool pack'\n|" \
	    -e "/^mv/d" \
	    -e "/\/buildid/d" \
	    -e "s|-extld=gcc|-tmpdir='$(BUILD_ABS_DIR)' -linkmode=external -extldflags='-nostartfiles -nodefaultlibs -nostdlib -r' -extld=$(LD)|g" \
	    -e 's|$$WORK|$(BUILD_ABS_DIR)|g' \
            | sh 2>&1 |  sed -e "s/^/  | /g"
		
	@echo "[objcopy] globalizing symbols {runtime.g0, main.main} in go.o"
	@aarch64-linux-gnu-objcopy \
		--globalize-symbol runtime.g0 \
		--globalize-symbol main.main \
		--globalize-symbol uart_init \
		 $(BUILD_DIR)/go.o $(BUILD_DIR)/go.o

src/%.o: src/%.c
	aarch64-linux-gnu-gcc $(CFLAGS) -c $< -o $@

kernel8.img: arch/arm64/start.o main.o
	aarch64-linux-gnu-ld -nostdlib -nostartfiles arch/arm64/start.o $(BUILD_ABS_DIR)/go.o -T arch/arm64/link.ld -o kernel8.elf
	aarch64-linux-gnu-objcopy -O binary kernel8.elf kernel8.img

clean:
	rm -rf build/ kernel8.img kernel8.elf *.o arch/arm64/*.o >/dev/null 2>/dev/null || true

run:
	qemu-system-aarch64 -M raspi3 -kernel kernel8.img -serial null -serial stdio 
