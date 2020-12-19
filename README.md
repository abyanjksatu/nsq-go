# nsq-go

run docker:
```bash
$ docker-compose up -d
```

check nsq web ui:
```bash
$ curl 127.0.0.1:4171/ping
```

run consumer:
```bash
$ go run consume/consume.go
```

run publisher:
```bash
$ go run publish/publish.go
```