# kitten

***A tiny open source workflow engine written in Go***

## Tutorial

### Init

```go
dnsOption := db.SetDSN("root@tcp(127.0.0.1:3306)/flow_test?charset=utf8")

traceOption := db.SetTrace(true)

kitten.Init(dnsOption, traceOption)
```

