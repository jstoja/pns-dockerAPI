# Use phusion/baseimage as base image. To make your builds
# reproducible, make sure you lock down to a specific version, not
# to `latest`! See
# https://github.com/phusion/baseimage-docker/blob/master/Changelog.md
# for a list of version numbers.
FROM phusion/baseimage:0.9.11
MAINTAINER Julien Bordellier <me@julienbordellier.com>

# Set correct environment variables.
ENV HOME /root

# Regenerate SSH host keys. baseimage-docker does not contain any, so you
# have to do that yourself. You may also comment out this instruction; the
# init system will auto-generate one during boot.
RUN /etc/my_init.d/00_regen_ssh_host_keys.sh

# Use baseimage-docker's init system.
CMD ["/sbin/my_init"]

RUN apt-get update
RUN apt-get -y install wget

#
# Install Go to build platform specific application server
# (just avoids to build before creating the docker image)
#
#WORKDIR /usr/local/
#RUN ["wget", "http://golang.org/dl/go1.3.linux-amd64.tar.gz"]
#RUN ["tar", "xfvz", "go1.3.linux-amd64.tar.gz"]
#ENV PATH $PATH:/usr/local/go/bin
#ENV GOPATH /usr/local/gopackages

#
# Installing Go packages
#
#RUN ["/usr/local/go/bin/go", "get", "github.com/gorilla/mux"]

#
# RTMP server waiting for streams
#
WORKDIR /tmp
RUN ["wget", "http://rtmpd.com/assets/binaries/784/crtmpserver-1.1_beta-x86_64-Ubuntu_12.04.tar.gz"]
RUN ["tar", "xfz", "crtmpserver-1.1_beta-x86_64-Ubuntu_12.04.tar.gz", "-C", "/etc/"]
RUN ["mv", "/etc/crtmpserver-1.1_beta-x86_64-Ubuntu_12.04", "/etc/crtmpserver"]
RUN ["rm", "-rf", "/tmp/crtmpserver-1.1_beta-x86_64-Ubuntu_12.04"]

WORKDIR /
RUN ["mkdir", "-p", "/var/log/crtmpserver/"]
VOLUME /var/log/crtmpserver

RUN ["mkdir", "-p", "/etc/service/crtmpd"]
ADD ./scripts/crtmpd.sh /etc/service/crtmpd/run
RUN ["chmod", "+x", "/etc/service/crtmpd/run"]

#EXPOSE 1935
#EXPOSE 1234

#
# PlugAndStream server manager
#
#ADD ./server_manager /etc/pns
#RUN ["/usr/local/go/bin/go", "build", "-o", "/etc/pns/pns", "/etc/pns/app.go"]
#
#RUN mkdir /etc/service/pns
#ADD pns.sh /etc/service/pns/run
#RUN ["chmod", "755", "/etc/service/pns/run"]
#
#RUN ["mkdir", "-p", "/var/log/pns/"]
#RUN ["touch", "/var/log/pns/pns.log"]
#
#EXPOSE 3001

# Clean up APT when done.
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
