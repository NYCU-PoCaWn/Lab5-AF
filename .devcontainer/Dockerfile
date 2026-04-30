FROM golang:1.24

# Install common tools
RUN apt-get update && apt-get install -y \
    git \
    make \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Install Golangci-lint 
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin

RUN apt-get update && apt-get install -y net-tools iputils-ping curl tcpdump jq


WORKDIR /workspace