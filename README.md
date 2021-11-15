## mync

The mync is simple implementation of the netcat utility in go that allows to listen and send data over TCP  protocols,as well as TCP port scanning.

This utility was created for Golang learning purposes .

### Usage:

```
mync host:port     #Open a TCP connection
mync -l -p port    #Listen on TCP port
mync -s host  #Scan TCP port
```
### Examples:

**$ mync 127.0.0.1:80**

Open a TCP connection to port 89 of localhost.

**$mync  -l -p 9999**

Listen on TCP port 9999.

**$ mync   -s 127.0.0.1**

Scan all listening TCP ports of localhost