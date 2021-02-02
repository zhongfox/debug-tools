FROM alpine:3.12

MAINTAINER Kamran Azeem & Henrik Høegh (kaz@praqma.net, heh@praqma.net)

EXPOSE 80 443

# Install some tools in the container and generate self-signed SSL certificates.
# Packages are listed in alphabetical order, for ease of readability and ease of maintenance.
RUN     apk update \
    &&  apk add apache2-utils bash bind-tools busybox-extras curl ethtool git \
                iperf3 iproute2 iputils jq lftp mtr mysql-client \
                netcat-openbsd net-tools nginx nmap openssh-client openssl \
	            perl-net-telnet postgresql-client procps rsync socat tcpdump tshark wget \
    &&  mkdir /certs \
    &&  chmod 700 /certs \
    &&  openssl req \
        -x509 -newkey rsa:2048 -nodes -days 3650 \
        -keyout /certs/server.key -out /certs/server.crt -subj '/CN=localhost'

# grpcurl
RUN wget https://github.com/fullstorydev/grpcurl/releases/download/v1.7.0/grpcurl_1.7.0_linux_x86_64.tar.gz \
&& tar zxvpf grpcurl_1.7.0_linux_x86_64.tar.gz && mv grpcurl /usr/local/bin/ && rm grpcurl_1.7.0_linux_x86_64.tar.gz

# grpc-tools
# RUN wget https://github.com/bradleyjkemp/grpc-tools/releases/download/v0.2.6/grpc-tools_0.2.6_Linux_amd64.zip \
# && unzip -o ./grpc-tools_0.2.6_Linux_amd64.zip && mv grpc-dump grpc-fixture grpc-replay /usr/local/bin/ \
# && rm ./grpc-tools_0.2.6_Linux_amd64.zip 

#COPY grpcdump /usr/local/bin/grpcdump
COPY grpc-dump /usr/local/bin/grpc-dump

COPY api /api
COPY istioproto /istioproto


# Copy a simple index.html to eliminate text (index.html) noise which comes with default nginx image.
# (I created an issue for this purpose here: https://github.com/nginxinc/docker-nginx/issues/234)
COPY index.html /usr/share/nginx/html/

# Copy a custom nginx.conf with log files redirected to stderr and stdout
COPY nginx.conf /etc/nginx/nginx.conf

COPY docker-entrypoint.sh /


# Run the startup script as ENTRYPOINT, which does few things and then starts nginx.
ENTRYPOINT ["/docker-entrypoint.sh"]


# Start nginx in foreground:
CMD ["nginx", "-g", "daemon off;"]


###################################################################################################

# Build and Push (to dockerhub) instructions:
# -------------------------------------------
# docker build -t local/network-multitool .
# docker tag local/network-multitool praqma/network-multitool
# docker login
# docker push praqma/network-multitool


# Pull (from dockerhub):
# ----------------------
# docker pull praqma/network-multitool


# Usage - on Docker:
# ------------------
# docker run --rm -it praqma/network-multitool /bin/bash 
# OR
# docker run -d  praqma/network-multitool
# OR
# docker run -p 80:80 -p 443:443 -d  praqma/network-multitool
# OR
# docker run -e HTTP_PORT=1080 -e HTTPS_PORT=1443 -p 1080:1080 -p 1443:1443 -d  praqma/network-multitool


# Usage - on Kubernetes:
# ---------------------
# kubectl run multitool --image=praqma/network-multitool --replicas=1
