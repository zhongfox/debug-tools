FROM alpine
RUN apk update
# telnet
RUN apk add busybox-extras
RUN apk add curl
# arp, hostname, ifconfig, ipmaddr, iptunnel, mii-tool, nameif, netstat, plipconfig, rarp, route, slattach
RUN apk add net-tools
# bridge, ss, ip
RUN apk add iproute2
# delv, dig, host, nslookup, nsupdate
RUN apk add bind-tools

CMD ["/bin/sh"]