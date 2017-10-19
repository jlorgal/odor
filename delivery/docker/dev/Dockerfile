FROM base/archlinux
LABEL maintainer Jorge Lorenzo <jlorgal@gmail.com>

# Install packages for golang environment
RUN pacman -Syyu --noconfirm
RUN pacman -S --noconfirm git gcc make go \
                          linux-headers iptables libpcap libnetfilter_queue \
                          pkg-config

ARG GOPATH=/root/go
ARG PROJECT=github.com/jlorgal/odor
ENV PATH=$GOPATH/bin:${PATH} \
    GOPATH=$GOPATH \
    GOPROJECT=$GOPATH/src/$PROJECT

RUN echo "export GOPATH=$GOPATH PATH=$GOPATH/bin:${PATH}" >> /root/.bashrc && \
    mkdir -p $GOPROJECT

# Install golang packages
RUN go get -v github.com/golang/lint/golint \
              github.com/golang/dep/cmd/dep \
              github.com/aktau/github-release \
              golang.org/x/net/bpf \
              golang.org/x/sys/unix

# Mount the source code of the project under GOPATH
WORKDIR $GOPROJECT
