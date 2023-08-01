## 1. Use

### 1. create client.ini

```ini
[common]
tunnel_addr = ":9998"
tunnel_url = "/tunnel1"

[http]
local_addr = ":8080"
remote_addr = ":30001"

[ssh]
local_addr = ":8081"
remote_addr = "172.16.1.104:22"
tunnel_addr = ":9999"
tunnel_url = "/tunnel2"
```

### 2. create server.ini

```ini
[default1]
tunnel_addr = ":9998"
tunnel_url = "/tunnel1"

[default2]
tunnel_addr = ":9999"
tunnel_url = "/tunnel2"
```

### 3. start client and server

```shell
./http-tunnel-server -c server.ini

./http-tunnel-client -c client.ini
```



