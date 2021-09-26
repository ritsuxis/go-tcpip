# go-tcpip
goによるtcp/ipプロトコルスタックの実装  
[lectcp](https://github.com/pandax381/lectcp)と[gotcp](https://github.com/terassyi/gotcp)を参考にしています

## Prepare
### tap device
```
sudo ip tuntap add mode tap user $USER name tap0
sudo ip addr add 192.0.2.1/24 dev tap0
sudo ip link set tap0 up
```

### nc tcp server
```
nc -l -p 10381 -s 192.0.2.1
```

### build
```
go build -o tcpc.exe
```

### execute examples
```
sudo ./tcpc.exe -name tap0
``` 

## TODO
- arpを外に対してできるようにする
- 外に対して通信できるようにする
- arpが未解決の時はキューにメッセージを溜める
- Close処理を適切にする
- サーバー側で動くようにする
