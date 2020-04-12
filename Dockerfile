FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y curl build-essential libnetcdf-dev git

# Golang
RUN rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.14.1

RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz \
		| tar -C /usr/local -xz

ENV PATH /usr/local/go/bin:$PATH

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

WORKDIR /

# PyMOL
RUN curl -o pymol.tar.bz2 https://pymol.org/installers/PyMOL-2.3.4_121-Linux-x86_64-py37.tar.bz2
RUN tar -xf pymol.tar.bz2

# Fpocket
RUN git clone https://github.com/Discngine/fpocket.git
WORKDIR /fpocket
RUN make

# Node
RUN curl -sL https://deb.nodesource.com/setup_13.x | bash -
RUN apt-get install -y nodejs

# Yarn
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - 
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update
RUN apt-get install -y yarn

# VarQ build
RUN mkdir /varq
ADD . /varq/
WORKDIR /varq/

RUN make build

COPY config-example.yaml config.yaml
CMD ["/varq/varq"]
