# go-nsq-compare

run docker:
```bash
$ docker-compose up -d
```

check nsq web ui: `http://127.0.0.1:4171`

## nsqio/go-nsq

Run producer & consumer :
```sh
$ go run nsqio/main.go
```
output:
```sh
2021/02/06 06:31:29 INF    1 (127.0.0.1:4150) connecting to nsqd
2021/02/06 06:31:29 INF    2 [hello-topic/hello-channel] (127.0.0.1:4150) connecting to nsqd
2021/02/06 06:31:29 NSQ message received, msg: Hello World
```

## segmentio/nsq-go

Run producer & consumer :
```sh
$ go run segmentio/main.go
```
output:
```sh
2021/02/06 06:31:42 opening nsqd connection to localhost:4150
Hello World!
2021/02/06 06:31:43 Consumer initiating shutdown sequence
2021/02/06 06:31:43 sending CLS to all command channels
2021/02/06 06:31:43 awaiting connection waitgroup
2021/02/06 06:31:43 draining and requeueing remaining in-flight messages
2021/02/06 06:31:43 closing and cleaning up connections
2021/02/06 06:31:43 successfully flushed all connections
2021/02/06 06:31:43 Consumer exiting run
```