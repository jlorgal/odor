FROM dockerhub.hi.inet/dcip/minimal:alpine-3.6
LABEL maintainer Jorge Lorenzo <jlorgal@gmail.com>

# Install packages for golang environment
RUN apk update && \
    apk add --no-cache git gcc make libc-dev go \
                       python2-dev py2-pip py-virtualenv \
                       zlib-dev jpeg-dev \
                       docker

USER contint

ARG GOPATH=/home/contint/go
ARG PROJECT=pdihub.hi.inet/awazza/niji-orchestrator
ENV PATH=$GOPATH/bin:${PATH} \
    GOPATH=$GOPATH \
    GOPROJECT=$GOPATH/src/$PROJECT

RUN echo "export GOPATH=$GOPATH PATH=$GOPATH/bin:${PATH}" >> /home/contint/.bashrc && \
    mkdir -p $GOPROJECT

# Install golang packages
RUN go get -v github.com/golang/lint/golint \
              github.com/golang/dep/cmd/dep \
              github.com/aktau/github-release

# Mount the source code of the project under GOPATH
WORKDIR $GOPROJECT

# Switch to root user (required by dcip)
USER root
