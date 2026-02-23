FROM fedora:43

# Qt6 dev libs, ImageMagick, build tools
RUN dnf install -y \
    qt6-qtbase-devel \
    ImageMagick \
    gcc gcc-c++ \
    git \
    && dnf clean all

# Install Go 1.26.0 manually (much faster CGO builds than distro Go)
RUN curl -fsSL https://go.dev/dl/go1.26.0.linux-amd64.tar.gz | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:/root/go/bin:${PATH}"
ENV GOPATH="/root/go"

# Install rugo
RUN go install github.com/rubiojr/rugo@v0.26.3

ENV QT_QPA_PLATFORM=offscreen
ENV LC_ALL=C.UTF-8

WORKDIR /cute
