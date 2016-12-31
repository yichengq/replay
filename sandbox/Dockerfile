FROM ubuntu:14.04

ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
ENV GOROOT_BOOTSTRAP /usr/local/gobootstrap

RUN set -x; \
    apt-get update && apt-get install -y curl ca-certificates git build-essential m4 tar golang --no-install-recommends && \
    curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.31.1/install.sh | bash && \
    . /root/.bashrc && \
    nvm install v6.9.2 && \

    git clone https://github.com/bloomberg/bucklescript-addons.git && \
    cd bucklescript-addons/examples/reason-demo && \
    npm install

ENV BSC_BIN /bucklescript-addons/examples/reason-demo/node_modules/bs-platform/bin/bsc
ENV REFMT_BIN /bucklescript-addons/examples/reason-demo/node_modules/reason/refmt_impl.native

COPY . /go/src/sandbox
RUN go install sandbox

EXPOSE 8080
ENTRYPOINT ["/go/bin/sandbox"]
