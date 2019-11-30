# murmur-cli

murmur-cli provides an interface to a grpc-enabled murmur server.

This is a fork of the [original project](https://github.com/layeh/murmur-cli)
as the build didn't work for me, and it was missing some features I wanted.

## Installation

    env GOBIN=$PWD go get -u github.com/MckayJT/murmur-cli

## Precompiled downloads

New and fresh, supporting features that murmur doesn't even have yet, client
side TLS certificate authentication.

[Grab them today from my website!](https://octet.space/murmur-cli/)

All the packages come with a signature file included using
[minisign](https://jedisct1.github.io/minisign/). My key is

```
minisign public key 6125BC1D7F42FCBB
RWS7/EJ/HbwlYZWZHMuSK11XUJvhTHnWzFP3mz0nb9WWQt4j98ozJ5ou
```

You can verify these with minisign or they also verify with signify if you're
on OpenBSD.


## Usage Tips

gRPC will do host name checking against the server certificate, and it's
unlikely you will have a valid certificate that matches an internal
IP address. Since there is no authentication in murmur for gRPC connections
opening it up on a public address that you can get a certificate for
is a security risk.

You can securely set up communications via a unix socket easily.
If you are using systemd, just add a drop-in file in
`/etc/systemd/system/murmur.service.d/override.conf` such as:

```
[Service]
User=murmur
RuntimeDirectory="murmur/"
RuntimeDirectoryMode=0770
```

and change your murmur.ini like

```
grpc="unix:///run/murmur/grpc.sock"
```

This will create a socket that only the user or group murmur can access.
Remember to use the -address or set $MURMUR\_ADDRESS to this value
and set -insecure=true.

You can also now use the -hostoverride flag to expect a different name
from the server than what you are connecting to. So if you use your regular
certificate for 'murmur.foo.io' in the gRPC settings for murmur, you can
use -hostoverride="murmur.foo.io' and it will accept the certificate as valid
if you connect using a local socket or a loopback address.

I am working on client authentication for murmur. To use this, you need to have
a certificate and key that can be used when connecting.

## Syntax
    usage: murmur-cli [flags] [command... [arguments...]]

    Flags:
      -address="127.0.0.1:50051"   address and port of murmur's grpc endpoint
                                    (can also be set via $MURMUR_ADDRESS).
      -timeout="10s"               duration to wait for connection.
      -template=""                 Go text/template template to use when outputing
                                    data. By default, JSON objects are printed.
      -hostoverride=""             Expect a different hostname from the server
      -cert=""                     Client certificate (pem format)
      -key=""                      Client certeficate key (pem format), unencrypted
      -insecure=true               Disable TLS encryption.
      -help                        Print command list.

    Commands:
      acl get <server id> <channel id>
      acl get-effective-permissions <server id> <session> <channel id>

      ban get <server id>

      channel query <server id>
      channel get <server id> <channel id>
      channel add <server id> <parent channel id> <name>
      channel remove <server id> <channel id>

      config get <server id>
      config get-field <server id> <key>
      config set-field <server id> <key> <value>
      config get-defaults

      contextaction add <server id> <context> <action> <text> <session>
        Context is a comma seperated list of the following:
          Server
          Channel
          User
      contextaction remove <server id> <action> [session]
      contextaction events <server id> <action>

      database query <server id> [filter]
      database get <server id> <user id>
      database add <server id> <user id> <password>  

      log query <server id> (<min> <max>)

      meta uptime
      meta version
      meta events

      server create
      server query
      server get <server id>
      server start <server id>
      server stop <server id>
      server remove <server id>
      server events <server id>

      textmessage send <server id> [sender:<session>] [targets...] <text>
        Valid targets:
          user:<session>
          channel:<id>
          tree:<id>
      textmessage filter <server id> <program> [args...]

      tree query <server id>

      user query <server id>
      user get <server id> <session>
      user kick <server id> <session> [reason]


## Original Author

Tim Cooper (<tim.cooper@layeh.com>)
