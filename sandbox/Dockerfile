FROM ubuntu:14.04

ENV GOPATH /go
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
ENV GOROOT_BOOTSTRAP /usr/local/gobootstrap

RUN set -x; \
    apt-get update && apt-get install -y software-properties-common --no-install-recommends && add-apt-repository ppa:avsm/ppa && \
    apt-get update && apt-get install -y ocaml ocaml-native-compilers camlp4-extra patch opam make m4 --no-install-recommends && \
    opam init && \
    opam update && \
    opam switch 4.02.3

ENV CAML_LD_LIBRARY_PATH /root/.opam/4.02.3/lib/stublibs
ENV MANPATH :/root/.opam/4.02.3/man
ENV PERL5LIB /root/.opam/4.02.3/lib/perl5
ENV OCAML_TOPLEVEL_PATH /root/.opam/4.02.3/lib/toplevel
ENV PATH /root/.opam/4.02.3/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:$PATH

RUN set -x; \
    opam install reason

ENV REBUILD_BIN rebuild

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
