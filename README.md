# go-tcpip
goによるtcp/ipプロトコルスタックの実装  
[lectcp](https://github.com/pandax381/lectcp)と[gotcp](https://github.com/terassyi/gotcp)を参考にしています

## できること・できないこと
### [lectcp](https://github.com/pandax381/lectcp)から拡張してできること
- TCPヘッダーの生成
- flagの指定
- クライアント側のthree way handshaking
- データ送信
- 送信終了時のFINの処理

### できないこと
- 再送処理
- パケット分割
- ウィンドウ制御
- フロー制御
- サーバ側

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
- サーバー側で動くようにする
