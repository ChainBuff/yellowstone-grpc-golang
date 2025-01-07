# grpc 网络延迟测试 

在 go 中常见的网络测试方法有：

1.	使用 time 包直接记录请求的响应时间。
2.	实现基准测试进行自动化延迟测试。
3.	使用外部工具（如 ghz）模拟流量并测量延迟。
4.	集成监控工具（Prometheus + Grafana）进行实时分析。
5.	使用 gRPC 的内部调试工具捕获延迟数据。

我们可以通过 proto 文件中的 ping 和 pong 来测试网络延迟。具体的方法在:

```
type GeyserClient interface {
    // 双向流式订阅
    Subscribe(ctx context.Context, opts ...grpc.CallOption) (Geyser_SubscribeClient, error)
    
    // 普通 RPC 调用
    Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PongResponse, error)
    GetLatestBlockhash(ctx context.Context, in *GetLatestBlockhashRequest, opts ...grpc.CallOption) (*GetLatestBlockhashResponse, error)
    GetBlockHeight(ctx context.Context, in *GetBlockHeightRequest, opts ...grpc.CallOption) (*GetBlockHeightResponse, error)
    GetSlot(ctx context.Context, in *GetSlotRequest, opts ...grpc.CallOption) (*GetSlotResponse, error)
    IsBlockhashValid(ctx context.Context, in *IsBlockhashValidRequest, opts ...grpc.CallOption) (*IsBlockhashValidResponse, error)
    GetVersion(ctx context.Context, in *GetVersionRequest, opts ...grpc.CallOption) (*GetVersionResponse, error)
}
```

那么我们就可以在订阅中主动的去发送 ping 请求，等待服务器的 pong 从而测试服务器与客户端之间的抖动。

```
	go func() {
		for range ticker.C {
			lastPingTime = time.Now()
			err := stream.Send(&pb.SubscribeRequest{
				Ping: &pb.SubscribeRequestPing{},
			})
			if err != nil {
				log.Printf("发送 ping 失败: %v", err)
				return
			}
			log.Printf("发送 ping - 时间: %v", lastPingTime.Format("15:04:05.000"))
		}
	}()

	// 接收响应
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("接收响应失败: %v", err)
		}

		// 处理 pong 响应
		if resp.GetPong() != nil {
			now := time.Now()
			latency := now.Sub(lastPingTime)
			totalLatency += latency
			pingCount++
			avgLatency := totalLatency / time.Duration(pingCount)

			log.Printf("收到 pong - 时间: %v, 当前延迟: %vms, 平均延迟: %vms",
				now.Format("15:04:05.000"),
				latency.Milliseconds(),
				avgLatency.Milliseconds())
		}
	}
```

如果你使用的是 grpc.chainbuff.com:443 的地址，那么你可以使用以下命令来测试网络延迟:
```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x02_network_test/golang # go run ./cmd/ping-test/main.go -endpoint grpc.chainbuff.com:443 -interval 1s
2025/01/07 12:06:51 发送 ping - 时间: 12:06:51.295
2025/01/07 12:06:51 收到 pong - 时间: 12:06:51.303, 当前延迟: 8ms, 平均延迟: 8ms
2025/01/07 12:06:52 发送 ping - 时间: 12:06:52.295
2025/01/07 12:06:52 收到 pong - 时间: 12:06:52.298, 当前延迟: 3ms, 平均延迟: 5ms
```

如果你的 grpc 服务器不具备 TLS 认证，你需要加上 `-insecure` 参数:
```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x02_network_test/golang # go run ./cmd/ping-test/main.go -endpoint grpc.chainbuff.com:443 -interval 1s -insecure
```
