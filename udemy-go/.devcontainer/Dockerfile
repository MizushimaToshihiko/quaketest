# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.194.3/containers/ubuntu/.devcontainer/base.Dockerfile

# [Choice] Ubuntu version: hirsute, bionic, focal
ARG VARIANT="focal"
FROM --platform=linux/x86_64 mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}
WORKDIR /root/
# [Optional] Uncomment this section to install additional OS packages.
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        build-essential \
        git \
        vim \
        man \
        manpages-dev \
        strace \
        gdb

FROM --platform=linux/x86_64 mcr.microsoft.com/vscode/devcontainers/go:1.18

WORKDIR /go/src/github.com/MizushimaToshihiko/gocc

COPY --from=0 /root/ ./

# Golang 環境構築(任意)
# RUN go install golang.org/x/tools/gopls@latest \
#  golang.org/x/tools/cmd/godoc@latest

RUN go install github.com/tenntenn/goplayground/cmd/gp@latest
# RUN go install -u github.com/kisielk/errcheck@latest