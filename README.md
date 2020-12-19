# nsq-go

run docker:
```bash
$ docker-compose up -d
```

check nsq web ui: 127.0.0.1:4171

run consumer:
```bash
$ go run consume/consume.go
```

run publisher:
```bash
$ go run publish/publish.go
```