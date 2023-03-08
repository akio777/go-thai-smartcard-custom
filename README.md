# this repo fork from this :
https://github.com/somprasongd/go-thai-smartcard


# go-thai-smartcard

Go application read personal and nhso data from thai id card, it run in the background and wait until inserted card then send readed data to everyone via [https://socket.io/](socket.io) and [WebSockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API).

Or use like library see in [cmd/example/main.go](https://github.com/somprasongd/go-thai-smartcard/blob/main/cmd/example/main.go)

## Other Versions

- [Java](https://github.com/somprasongd/jThaiSmartCard)
- [Nodejs](https://github.com/somprasongd/thai-smartcard-nodejs)

## How to build

- Required version [Go](https://go.dev/dl/) version 1.18+
- Clone this repository
- Download all depencies with `go mod download`

> Linux install `sudo apt install build-essential libpcsclite-dev pcscd`

- Build with `go build -o bin/thai-smartcard-agent ./cmd/agent/main.go`

  > Windows `go build -o bin/thai-smartcard-agent.exe ./cmd/agent/main.go`

## How to run

Run from binary file that builded from the previous step.

### Configurations

|        ENV         | Default |                    Description                    |
| :----------------: | :-----: | :-----------------------------------------------: |
| **SMC_AGENT_PORT** | "9898"  |                    Server port                    |
| **SMC_SHOW_IMAGE** | "true"  | Enable or disable read face image from smartcard. |
| **SMC_SHOW_NHSO**  | "flase" | Enable or disable read nsho data from smartcard.  |
| **SMC_SHOW_LASER** | "flase" |  Enable or disable read laser id from smartcard.  |

### Run in daemon process

- Ubuntu

1. Clone an Build this program

```bash
cd ~

git clone https://github.com/somprasongd/go-thai-smartcard

cd go-thai-smartcard

go build -o ./bin/thai-smartcard-agent ./cmd/agent/main.go
```

2. Create System Service file

```bash
nano /lib/systemd/system/thai-smartcard-agent.service
```

```bash
[Unit]
Description=thai-smartcard-agent

[Service]
Environment="SMC_AGENT_PORT=9898"
Environment="SMC_SHOW_IMAGE=true"
Environment="SMC_SHOW_NHSO=flase"
Environment="SMC_SHOW_LASER=flase"
Type=simple
Restart=always
RestartSec=5s
ExecStart=~/go-thai-smartcard/bin/thai-smartcard-agent

[Install]
WantedBy=multi-user.target
```

3. Start the service with system service command

```bash
service thai-smartcard-agent start
```

### Run in daemon process with PM2

- Windows

```bash
npm install -g pm2 pm2-windows-startup
pm2-startup install
pm2 start .\bin\thai-smartcard-agent.exe --name smc
pm2 save
```

- Ubuntu

```bash
npm install -g pm2
pm2 start ./bin/thai-smartcard-agent --name smc
pm2 startup
pm2 save
```

- Mac

```bash
npm install -g pm2
pm2 start ./bin/thai-smartcard-agent --name smc
pm2 startup
pm2 save
```

## Example

### Client connect via socket.io

```javascript
<script>
  const socket = io.connect('http://localhost:9898');
  socket.on('connect', function () {

  });
  socket.on('smc-data', function (data) {
    console.log(data);
  });
  socket.on('smc-error', function (data) {
    console.log(data);
  });
  socket.on('smc-removed', function (data) {
    console.log(data);
  });
  socket.on('smc-inserted', function (data) {
    console.log(data);
  });
</script>
```

### Client connect via WebSokcet

```javascript
<script>
  // Connection to Websocker Server...
  if (window['WebSocket']) {
    conn = new WebSocket('ws://' + document.location.host + '/ws');
    console.log(document.location.host);
    conn.onopen = function (evt) {
      var item = document.getElementById('data-ws');
      item.innerHTML = '<b>Connected to WebSocket server</b>';
    };
    conn.onclose = function (evt) {
      var item = document.getElementById('data-ws');
      item.innerHTML = '<b>Disconnected to WebSocket server</b>';
    };
    conn.onmessage = function (evt) {
      var messages = evt.data.split('\n');
      console.log(messages);
      for (var i = 0; i < messages.length; i++) {
        var item = document.getElementById('data-ws');
        item.innerText = messages[i];
      }
    };
    conn.onerror = function (err) {
      console.log('Socket Error: ', err);
    };
  } else {
    var item = document.getElementById('data-ws');
    item.innerHTML = '<b>Your browser does not support WebSockets.</b>';
    appendLog(item);
  }
</script>
```