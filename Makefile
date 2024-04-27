# 检查系统类型
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

# 根据系统类型和架构设置编译参数
ifeq ($(UNAME_S),Darwin)
    ifeq ($(UNAME_M),arm64)
        GOOS = darwin
        GOARCH = arm64
    else
        GOOS = darwin
        GOARCH = amd64
    endif
else ifeq ($(UNAME_S),Linux)
    ifeq ($(UNAME_M),x86_64)
        GOOS = linux
        GOARCH = amd64
    else ifeq ($(UNAME_M),aarch64)
        GOOS = linux
        GOARCH = arm64
    else
        $(error "Unsupported Linux architecture: $(UNAME_M)")
    endif
else
    GOOS = windows
    GOARCH = amd64
endif

# 检查是否安装了 Go 语言
ifeq (, $(shell which go))
    $(error "Go is not installed. Please install Go before running this Makefile.")
endif

# 默认目标
all: tidy build

# 拉取 import 包
tidy:
    go mod tidy

# 编译程序
build:
    GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o dding main.go

# 清理编译生成的文件
clean:
    rm -f dding
