# goland mircorService learning


## Run docker
```
rm -rf /tmp/etcd-data.tmp && mkdir -p /tmp/etcd-data.tmp && \
  docker rmi gcr.io/etcd-development/etcd:v3.3.15 || true && \
  docker run \
  -p 2379:2379 \
  -p 2380:2380 \
  --mount type=bind,source=/tmp/etcd-data.tmp,destination=/etcd-data \
  --name etcd-gcr-v3.3.15 \
  gcr.io/etcd-development/etcd:v3.3.15 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
``` 

## Run prometheus server
```
docker run -d -p 9090:9090 -v /tmp/prometh.yaml:/etc/prometheus/prometheus.yml prom/prometheus

# 配置文件 prometh.yaml
global:
  scrape_interval:     60s
  evaluation_interval: 60s

scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets: ['0.0.0.0:9090']
        labels:
          instance: prometheus

  - job_name: linux
    metrics_path: "/metrics"
    static_configs:
      - targets: ['10.12.214.39:9100']
        labels:
          instance: localhost
```

## Run jaeger server
```
sudo docker run -d -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p5775:5775/udp -p6831:6831/udp -p6832:6832/udp -p5778:5778 -p16686:16686 -p14268:14268 -p9411:9411 jaegertracing/all-in-one:latest
```


### run server
```
go run user_agent/main.go
```

### run client
```
go test -v  user_agent/client/client_test.go
```

### monitor && opentracing
```
monitor: http://localhost:9090/targets
opentracing: http://127.0.0.1:16686/search
```