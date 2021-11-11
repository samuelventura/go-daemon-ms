# go-daemon-ms

Daemon manager with RESTish API.

## API

```bash
#explicit /api/daemon to make it proxy aggregatable
#daemon actions
curl -X POST http://127.0.0.1:31600/api/daemon/install/:name?path=/path/to/executable
curl -X POST http://127.0.0.1:31600/api/daemon/uninstall/:name
curl -X POST http://127.0.0.1:31600/api/daemon/enable/:name
curl -X POST http://127.0.0.1:31600/api/daemon/disable/:name
curl -X POST http://127.0.0.1:31600/api/daemon/stop/:name
curl -X POST http://127.0.0.1:31600/api/daemon/start/:name
#daemon queries
curl -X GET http://127.0.0.1:31600/api/daemon/list
curl -X GET http://127.0.0.1:31600/api/daemon/info/:name
#environ management
curl -X POST http://127.0.0.1:31600/api/daemon/env/:name \
    -H "DaemonEnviron: VAR1=VALUE1" \
    -H "DaemonEnviron: VAR2=VALUE2"
curl -X GET http://127.0.0.1:31600/api/daemon/env/:name
curl -X DELETE http://127.0.0.1:31600/api/daemon/env/:name
```

## Test Drive

```bash
go install && ~/go/bin/go-daemon-ms 
go install github.com/samuelventura/go-state/sample
~/go/bin/sample 
curl -X GET http://127.0.0.1:31600/api/daemon/list
curl -X GET http://127.0.0.1:31600/api/daemon/info/sample
curl -X POST "http://127.0.0.1:31600/api/daemon/install/sample?path=$HOME/go/bin/sample"
curl -X POST http://127.0.0.1:31600/api/daemon/uninstall/sample
curl -X POST http://127.0.0.1:31600/api/daemon/enable/sample
curl -X POST http://127.0.0.1:31600/api/daemon/disable/sample
curl -X POST http://127.0.0.1:31600/api/daemon/stop/sample
curl -X POST http://127.0.0.1:31600/api/daemon/start/sample
curl -X POST http://127.0.0.1:31600/api/daemon/env/sample \
    -H "DaemonEnviron: VAR1=VALUE1" \
    -H "DaemonEnviron: VAR2=VALUE2"
curl -X GET http://127.0.0.1:31600/api/daemon/env/sample
curl -X DELETE http://127.0.0.1:31600/api/daemon/env/sample
```

## Ports

- 31600 go-daemon-ms
- 31607 go-echo-ms
- 31622 go-dock-ms
- 31625 go-mail-ms
- 31651 go-auth-ms
- 00000 go-ship-ms
- 00000 go-proxy-ms
- 00000 go-pay-ms
