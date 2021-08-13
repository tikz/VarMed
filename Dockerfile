FROM ubuntu:20.04

RUN apt-get update
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends tzdata
RUN apt-get install -y curl build-essential git python2

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

# Node
RUN curl -sL https://deb.nodesource.com/setup_14.x | bash -
RUN apt-get install -y nodejs

# Yarn
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - 
RUN echo "deb https://dl.yarnpkg.com/debian/ stable main" | tee /etc/apt/sources.list.d/yarn.list
RUN apt-get update
RUN apt-get install -y yarn

# Fpocket
RUN apt-get install -y libnetcdf-dev 
WORKDIR /
RUN git clone https://github.com/Discngine/fpocket.git
RUN cd fpocket && make && make install


# FreeSASA
RUN apt-get install -y autoconf
WORKDIR /
RUN git clone https://github.com/mittinatten/freesasa.git
RUN apt-get install -y libc++-dev libc++abi-dev
RUN cd freesasa && git submodule init && git submodule update && autoreconf -i && ./configure --disable-json --disable-xml && make && make install

# DSSP
RUN apt-get install -y dssp

# HMMER
RUN apt-get install -y hmmer

# VarMed build
RUN mkdir /varmed
ADD . /varmed/
WORKDIR /varmed/

RUN make build

COPY config-example.yaml /varmed/config.yaml

COPY pipeline-bins.tar.gz /
RUN tar -C /varmed/ -xvf /pipeline-bins.tar.gz

WORKDIR /varmed/

CMD ["/varmed/varmed"]