package installer

var expectedIgnitionServiceContents = map[string]string{
	"aro-etchosts-resolver.service": `[Unit]
Description=One shot service that appends static domains to etchosts
Before=network-online.target

[Service]
# ExecStart will copy the hosts defined in /etc/hosts.d/aro.conf to /etc/hosts
ExecStart=/bin/bash /usr/local/bin/aro-etchosts-resolver.sh

[Install]
WantedBy=multi-user.target
`,
	"fluentbit.service": `[Unit]
After=network-online.target
StartLimitIntervalSec=0

[Service]
RestartSec=1s
EnvironmentFile=/etc/sysconfig/fluentbit
ExecStartPre=-/bin/podman rm -f fluent-journal
ExecStartPre=-/bin/podman pull $FLUENTIMAGE
ExecStartPre=-mkdir -p /var/lib/fluent
ExecStart=/bin/podman run \
  --entrypoint /opt/td-agent-bit/bin/td-agent-bit \
  --net=host \
  --hostname bootstrap \
  --name fluent-journal \
  --rm \
  -v /etc/fluentbit/journal.conf:/etc/fluentbit/journal.conf \
  -v /var/lib/fluent:/var/lib/fluent:z \
  -v /var/log/journal:/var/log/journal:z,ro \
  -v /etc/machine-id:/etc/machine-id:ro \
  $FLUENTIMAGE \
  -c /etc/fluentbit/journal.conf

ExecStop=/bin/podman stop %N
Restart=always

[Install]
WantedBy=multi-user.target
`,
	"mdsd.service": `[Unit]
After=network-online.target
StartLimitIntervalSec=0

[Service]
RestartSec=1s
EnvironmentFile=/etc/sysconfig/mdsd
ExecStartPre=-/bin/podman rm -f %N
ExecStartPre=-mkdir /var/run/mdsd
ExecStartPre=-/bin/podman pull $MDSDIMAGE
ExecStart=/bin/podman run \
  --entrypoint /usr/sbin/mdsd \
  --net=host \
  --name mdsd \
  --env-file /etc/mdsd.d/mdsd.env \
  --rm \
  -v /etc/mdsd.d/:/etc/mdsd.d/:z \
  -v /var/run/mdsd:/var/run/mdsd:z \
  $MDSDIMAGE \
  -A -D -f 24224 -r /var/run/mdsd/default

ExecStop=/bin/podman stop %N
Restart=always

[Install]
WantedBy=multi-user.target
`,
	"dnsmasq.service": `
[Unit]
Description=DNS caching server.
After=network-online.target
Before=bootkube.service

[Service]
# ExecStartPre will create a copy of the customer current resolv.conf file and make it upstream DNS.
# This file is a product of user DNS settings on the VNET. We will replace this file to point to
# dnsmasq instance on the node. dnsmasq will inject certain dns records we need and forward rest of the queries to
# resolv.conf.dnsmasq upstream customer dns.
ExecStartPre=/bin/bash /usr/local/bin/aro-dnsmasq-pre.sh
ExecStart=/usr/sbin/dnsmasq -k
ExecStopPost=/bin/bash -c '/bin/mv /etc/resolv.conf.dnsmasq /etc/resolv.conf; /usr/sbin/restorecon /etc/resolv.conf'
Restart=always

[Install]
WantedBy=multi-user.target
`,
}
