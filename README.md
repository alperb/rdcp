# rdcp
Remote Daemon Control Protocol

## Introduction
rdcp is a protocol for controlling a daemon remotely. It is designed to be simple and easy to implement. It is also designed to be secure, and to be able to be used over an insecure connection.

Initial motivation of implementation is to use it with a remote docker daemon control tool ([drun](https://github.com/alperb/drun)). However, it is not limited to docker, and can be used with any daemon, if tailored to that daemon.

## Protocol
RDCP presents 4 different methods of request. 
- ORDER
- PROXY
- OPTIONS
- HEARTBEAT

## Methods

### `ORDER`

ORDER method is used to send a command to the daemon. It is a simple string, and some optional parameters. Daemon is expected to execute a command/function based on the string, and return a response. It is similar to RPC.

#### Request
```
ORD
Version: 0.1
Time: 1686604453

OrderedAction(parameter1=value1,parameter2=value2)
```

#### Response
```
R/ORD
Version: 0.1
From: 192.168.0.1
Received: 1686604453
Replied: 1686604700

Hello, World!
```

### `PROXY`

PROXY method is used to proxy a connection to the daemon. It is used to connect to a daemon, and send/receive data to/from it. It is similar to SSH tunneling.

Daemon is expected to proxy the request to the target, and return the response. It is similar to a SOCKS proxy.

#### Request
```
PRX
Version: 0.1
Time: 1686604453
Host: example.com

GET / HTTP/1.1
Host: example.com
```

#### Response
```
R/PRX
Version: 0.1
From: 192.168.0.1
Proxied-To: example.com
Received: 1686604453
Replied: 1686604700

Hello, World!
```

### `OPTIONS`

OPTIONS method is used to get the options of the daemon. It is used to get the options of the daemon, and return the response. It is similar to HTTP OPTIONS. It can be used to decide on the capabilities of the daemon.

#### Request
```
OPT
Version: 0.1
Time: 1686604453

localhost?
```

#### Response
```
R/OPT
Version: 0.1
From: 192.168.0.1
Received: 1686604453
Replied: 1686604700

192.168.0.7
```

### `HEARTBEAT`

HEARTBEAT method is used to keep the connection alive. It is used to keep the connection alive, and return the response. It is similar to HTTP HEAD. It can be used to keep the connection alive.

#### Request
```
HBT
Version: 0.1
Time: 1686604453
```

#### Response
```
R/HBT
Version: 0.1
From: 192.168.0.1
Received: 1686604453
Replied: 1686604700
```

## Security

### Encryption

RDCP is designed to be used over an insecure connection. However, it is also designed to be able to be used over a secure connection. It is recommended to use a secure connection, if possible. A RDCP over TLS is available. It is recommended to use so. 

### Authentication

There's no authentication mechanism available in RDCP since it is a stateless protocol like HTTP. It is possible and recommended to use Authorization header for authentication like used in HTTP. It is also possible to use a custom header for authentication.


