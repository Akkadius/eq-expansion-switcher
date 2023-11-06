#############################################
# go workspace
#############################################
FROM debian:12-slim

ARG PUID=1000
ENV PUID ${PUID}
ARG PGID=1000
ENV PGID ${PGID}

ENV GO_VERSION 1.19.1

RUN apt-get update -yqq && \
    groupadd -g ${PGID} go && \
    useradd -u ${PUID} -g go -m go -G go && \
    usermod -p "*" go

#############################################
# install dependencies
#############################################
RUN apt-get update && apt-get install -y \
    bash \
    curl \
    git \
    tree \
    jq \
    inetutils-ping \
    redis-tools \
    wget \
    cron \
    procps \
    make \
    sudo \
    inotify-tools \
    mariadb-client \
    zip \
    unzip \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
 && rm -rf /var/lib/apt/lists/*

#############################################
# install go
#############################################
RUN cd /tmp && wget --quiet https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz \
	&& tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && rm -rf /tmp/*

#######################################################################
# setup
#######################################################################
RUN echo "go ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/user && \
    chmod 0440 /etc/sudoers.d/user

#############################################
# set go env vars
# https://golang.org/doc/code.html
#############################################
ENV GOPATH=/home/go
ENV GOROOT=/usr/local/go/
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
ENV CGO_ENABLED 0

#######################################################################
# default shell
#######################################################################
RUN chsh -s /bin/bash go

#############################################
# download migrate utility (migrations)
#############################################
RUN mkdir -p /tmp/ && chmod 1777 /tmp && wget --quiet https://github.com/golang-migrate/migrate/releases/download/v4.12.2/migrate.linux-amd64.tar.gz \
    && tar -xvf migrate.linux-amd64.tar.gz -C /tmp/ \
    && mv /tmp/migrate.linux-amd64 /bin/migrate \
    && chmod +x /bin/migrate

USER go

#############################################
# download go utilities
#############################################
RUN go install github.com/google/wire/cmd/wire@latest
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest

#############################################
# node
#############################################
RUN cd /tmp/ && curl -SLO https://deb.nodesource.com/nsolid_setup_deb.sh && \
    chmod 500 nsolid_setup_deb.sh && \
    sudo ./nsolid_setup_deb.sh 18 && \
    sudo apt-get install nodejs -y

RUN sudo npm install -g gh-release && \
    sudo npm install -g @vue/cli

RUN sudo npm install -g vue-tsc

RUN sudo apt-get install build-essential -y

#############################################
# bash
#############################################
ENV PS1 '\[\e]0;\w\a\]\n\[\e[32m\]\u@\h \[\e[33m\]\w\[\e[0m\]\n\$ '

WORKDIR /home/go/src
