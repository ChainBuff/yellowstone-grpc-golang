package main

import (
	"context"
	"crypto/x509"
	"flag"
	"log"
	"time"

	pb "github.com/rpcpool/yellowstone-grpc/examples/golang/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	endpoint     = flag.String("endpoint", "", "gRPC endpoint")
	token        = flag.String("x-token", "", "认证令牌")
	interval     = flag.Duration("interval", 5*time.Second, "ping 间隔")
	insecureFlag = flag.Bool("insecure", false, "使用非 TLS 连接")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption

	// 根据 insecure 参数选择是否使用 TLS
	if *insecureFlag {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	// 建立连接
	conn, err := grpc.Dial(*endpoint, opts...)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 gRPC 客户端
	client := pb.NewGeyserClient(conn)

	// 设置认证信息
	ctx := context.Background()
	if *token != "" {
		md := metadata.New(map[string]string{"x-token": *token})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// 创建空的订阅流
	stream, err := client.Subscribe(ctx)
	if err != nil {
		log.Fatalf("创建流失败: %v", err)
	}

	// 发送空订阅请求以建立流
	err = stream.Send(&pb.SubscribeRequest{})
	if err != nil {
		log.Fatalf("发送订阅请求失败: %v", err)
	}

	// 启动 ping 计时器
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	var (
		totalLatency time.Duration
		pingCount    int
		lastPingTime time.Time
	)

	// 开始 ping/pong 循环
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
}
