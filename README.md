# cocoIM
this coco IM system.
This version will be built as plato

#run

go run main.go im-server

#sdk run

cd sdk
ts-node index.ts


#ipconf run 
先启动etcd集群
win下启动:
可以使用goreman,配置文件为Procfile.learner
```shell
etcd1: etcd --name jiqun1 --listen-client-urls http://127.0.0.1:12379 --advertise-client-urls http://127.0.0.1:12379 --listen-peer-urls http://127.0.0.1:12380 --initial-advertise-peer-urls http://127.0.0.1:12380 --initial-cluster-token etcd-cluster-1 --initial-cluster jiqun1=http://127.0.0.1:12380,jiqun2=http://127.0.0.1:22380,jiqun3=http://127.0.0.1:32380 --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
etcd2: etcd --name jiqun2 --listen-client-urls http://127.0.0.1:22379 --advertise-client-urls http://127.0.0.1:22379 --listen-peer-urls http://127.0.0.1:22380 --initial-advertise-peer-urls http://127.0.0.1:22380 --initial-cluster-token etcd-cluster-1 --initial-cluster jiqun1=http://127.0.0.1:12380,jiqun2=http://127.0.0.1:22380,jiqun3=http://127.0.0.1:32380 --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
etcd3: etcd --name jiqun3 --listen-client-urls http://127.0.0.1:32379 --advertise-client-urls http://127.0.0.1:32379 --listen-peer-urls http://127.0.0.1:32380 --initial-advertise-peer-urls http://127.0.0.1:32380 --initial-cluster-token etcd-cluster-1 --initial-cluster jiqun1=http://127.0.0.1:12380,jiqun2=http://127.0.0.1:22380,jiqun3=http://127.0.0.1:32380 --initial-cluster-state new --enable-pprof --logger=zap --log-outputs=stderr
```
启动完成后，本地启动ipconf
```shell
coco ipconf --config=./plato.yaml
```
可以通过curl访问
```shell
curl --location --request GET '127.0.0.1:6789/ip/list'
```
也可以浏览器直接访问上述地址 http://127.0.0.1:6789/ip/list


【2024/8/14】
目前测试的结果可得知，在笔记本上，放在docker容器里，单机能达到的连接数也就26万，离百万的距离还差的远，但是并不知道哪里出了问题