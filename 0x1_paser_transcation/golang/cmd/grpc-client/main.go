package main

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	// 导入 protobuf 生成的 gRPC 代码
	"github.com/mr-tron/base58"
	pb "github.com/rpcpool/yellowstone-grpc/examples/golang/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// 命令行参数定义
var (
	// 基础连接参数
	grpcAddr           = flag.String("endpoint", "", "Solana gRPC 服务器地址，使用 URI 格式，例如 https://api.rpcpool.com")
	token              = flag.String("x-token", "", "认证令牌")
	jsonInput          = flag.String("json", "", "订阅请求的 JSON，使用 @ 前缀从文件读取")
	insecureConnection = flag.Bool("insecure", false, "使用非 TLS 连接")

	// 区块链数据订阅选项
	slots      = flag.Bool("slots", false, "订阅 slot 更新")
	blocks     = flag.Bool("blocks", false, "订阅区块更新")
	block_meta = flag.Bool("blocks-meta", false, "订阅区块元数据更新")
	signature  = flag.String("signature", "", "订阅特定交易签名")
	resub      = flag.Uint("resub", 0, "在 x 次更新后重新仅订阅 slots，0 表示禁用")

	// 账户相关订阅选项
	accounts = flag.Bool("accounts", false, "订阅账户更新")

	// 交易相关订阅选项
	transactions       = flag.Bool("transactions", false, "订阅交易，用于 tx_account_include/tx_account_exclude 和 vote/failed")
	voteTransactions   = flag.Bool("transactions-vote", false, "包含投票交易")
	failedTransactions = flag.Bool("transactions-failed", false, "包含失败的交易")

	// 过滤器数组
	accountsFilter              arrayFlags // 账户过滤器
	accountOwnersFilter         arrayFlags // 账户所有者过滤器
	transactionsAccountsInclude arrayFlags // 交易包含的账户过滤器
	transactionsAccountsExclude arrayFlags // 交易排除的账户过滤器
)

// gRPC 客户端保活配置
var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // 如果没有活动，每 10 秒发送一次 ping
	Timeout:             time.Second,      // ping 超时时间为 1 秒
	PermitWithoutStream: true,             // 即使没有活动的流也发送 ping
}

func main() {
	log.SetFlags(0)

	// 设置命令行参数
	flag.Var(&accountsFilter, "accounts-account", "订阅指定账户，可多次指定")
	flag.Var(&accountOwnersFilter, "accounts-owner", "订阅指定账户所有者，可多次指定")
	flag.Var(&transactionsAccountsInclude, "transactions-account-include", "订阅包含指定账户的交易，可多次指定")
	flag.Var(&transactionsAccountsExclude, "transactions-account-exclude", "订阅不包含指定账户的交易，可多次指定")

	flag.Parse()

	// 验证必需参数
	if *grpcAddr == "" {
		log.Fatalf("需要提供 GRPC 地址。请使用 --endpoint 参数。")
	}

	// 解析 gRPC 服务器地址
	u, err := url.Parse(*grpcAddr)
	if err != nil {
		log.Fatalf("提供的 GRPC 地址无效: %v", err)
	}

	// 根据 URL scheme 推断是否使用安全连接
	if u.Scheme == "http" {
		*insecureConnection = true
	}

	// 设置默认端口
	port := u.Port()
	if port == "" {
		if *insecureConnection {
			port = "80"
		} else {
			port = "443"
		}
	}
	hostname := u.Hostname()
	if hostname == "" {
		log.Fatalf("请提供 URL 格式的端点，例如 http(s)://<endpoint>:<port>")
	}

	address := hostname + ":" + port

	// 建立 gRPC 连接
	conn := grpc_connect(address, *insecureConnection)
	defer conn.Close()

	// 开始订阅
	grpc_subscribe(conn)
}

// 建立 gRPC 连接
func grpc_connect(address string, plaintext bool) *grpc.ClientConn {
	var opts []grpc.DialOption

	// 配置 TLS
	if plaintext {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		pool, _ := x509.SystemCertPool()
		creds := credentials.NewClientTLSFromCert(pool, "")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}

	// 添加保活参数
	opts = append(opts, grpc.WithKeepaliveParams(kacp))

	log.Println("启动 gRPC 客户端，连接到", address)
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	return conn
}

// 处理 gRPC 订阅
func grpc_subscribe(conn *grpc.ClientConn) {
	var err error
	client := pb.NewGeyserClient(conn)

	var subscription pb.SubscribeRequest

	// 处理 JSON 输入
	if *jsonInput != "" {
		var jsonData []byte

		// 从文件或直接字符串读取 JSON
		if (*jsonInput)[0] == '@' {
			jsonData, err = os.ReadFile((*jsonInput)[1:])
			if err != nil {
				log.Fatalf("读取 JSON 文件错误: %v", err)
			}
		} else {
			jsonData = []byte(*jsonInput)
		}
		err := json.Unmarshal(jsonData, &subscription)
		if err != nil {
			log.Fatalf("解析 JSON 错误: %v", err)
		}
	} else {
		// 如果没有提供 JSON，创建空的订阅请求
		subscription = pb.SubscribeRequest{}
	}

	// 配置 slots 订阅
	if *slots {
		if subscription.Slots == nil {
			subscription.Slots = make(map[string]*pb.SubscribeRequestFilterSlots)
		}
		subscription.Slots["slots"] = &pb.SubscribeRequestFilterSlots{}
	}

	// 配置区块订阅
	if *blocks {
		if subscription.Blocks == nil {
			subscription.Blocks = make(map[string]*pb.SubscribeRequestFilterBlocks)
		}
		subscription.Blocks["blocks"] = &pb.SubscribeRequestFilterBlocks{}
	}

	// 配置区块元数据订阅
	if *block_meta {
		if subscription.BlocksMeta == nil {
			subscription.BlocksMeta = make(map[string]*pb.SubscribeRequestFilterBlocksMeta)
		}
		subscription.BlocksMeta["block_meta"] = &pb.SubscribeRequestFilterBlocksMeta{}
	}

	// 配置账户订阅
	if (len(accountsFilter)+len(accountOwnersFilter)) > 0 || (*accounts) {
		if subscription.Accounts == nil {
			subscription.Accounts = make(map[string]*pb.SubscribeRequestFilterAccounts)
		}
		subscription.Accounts["account_sub"] = &pb.SubscribeRequestFilterAccounts{}

		if len(accountsFilter) > 0 {
			subscription.Accounts["account_sub"].Account = accountsFilter
		}

		if len(accountOwnersFilter) > 0 {
			subscription.Accounts["account_sub"].Owner = accountOwnersFilter
		}
	}

	// 配置交易订阅
	if subscription.Transactions == nil {
		subscription.Transactions = make(map[string]*pb.SubscribeRequestFilterTransactions)
	}

	// 配置特定签名的交易订阅
	if *signature != "" {
		tr := true
		subscription.Transactions["signature_sub"] = &pb.SubscribeRequestFilterTransactions{
			Failed: &tr,
			Vote:   &tr,
		}
		subscription.Transactions["signature_sub"].Signature = signature
	}

	// 配置通用交易订阅
	if *transactions {
		subscription.Transactions["transactions_sub"] = &pb.SubscribeRequestFilterTransactions{
			Failed: failedTransactions,
			Vote:   voteTransactions,
		}
		subscription.Transactions["transactions_sub"].AccountInclude = transactionsAccountsInclude
		subscription.Transactions["transactions_sub"].AccountExclude = transactionsAccountsExclude
	}

	// 打印订阅请求
	subscriptionJson, err := json.Marshal(&subscription)
	if err != nil {
		log.Printf("序列化订阅请求失败: %v", subscriptionJson)
	}
	log.Printf("订阅请求: %s", string(subscriptionJson))

	// 设置上下文和认证信息
	ctx := context.Background()
	if *token != "" {
		md := metadata.New(map[string]string{"x-token": *token})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// 创建订阅流
	stream, err := client.Subscribe(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = stream.Send(&subscription)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// 处理订阅更新
	var i uint = 0
	for {
		i += 1
		// 重新订阅示例
		if i == *resub {
			subscription = pb.SubscribeRequest{}
			subscription.Slots = make(map[string]*pb.SubscribeRequestFilterSlots)
			subscription.Slots["slots"] = &pb.SubscribeRequestFilterSlots{}
			stream.Send(&subscription)
		}

		// 接收更新
		resp, err := stream.Recv()
		//timestamp := time.Now().UnixNano()

		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("接收更新时发生错误: %v", err)
		}
		if resp.GetTransaction() != nil {
			tx := resp.GetTransaction()
			for _, logMessage := range tx.GetTransaction().Meta.GetLogMessages() {
				// 只处理买入指令
				if strings.Contains(logMessage, "Program log: Instruction: Buy") {
					// 查找并解析Program data
					for _, msg := range tx.GetTransaction().Meta.GetLogMessages() {
						if strings.Contains(msg, "Program data: ") {
							data := strings.Split(msg, "Program data: ")[1]
							decoded, err := base64.StdEncoding.DecodeString(data)
							if err != nil {
								log.Printf("解码失败: %v", err)
								continue
							}

							// 验证数据长度
							if len(decoded) < 8+32+8+8+1+32+8 {
								continue
							}

							offset := 8 // 跳过魔数和版本

							// 解析Mint地址
							var mintBytes [32]byte
							copy(mintBytes[:], decoded[offset:offset+32])
							mintAddress := base58.Encode(mintBytes[:])
							offset += 32

							// 跳过sol和token数量
							offset += 16 // 8+8

							// 跳过isBuy标志
							offset += 1

							// 解析用户地址
							var userBytes [32]byte
							copy(userBytes[:], decoded[offset:offset+32])
							userAddress := base58.Encode(userBytes[:])
							offset += 32

							// 获取本地时间戳（毫秒）
							milliseconds := time.Now().UnixMilli()
							tradeTime := time.Now().Format("2006-01-02 15:04:05.000")

							// 输出交易信息
							log.Printf("\n===================== 买入交易 =====================\n"+
								"用户地址: %s\n"+
								"Mint地址: %s\n"+
								"交易时间: %s\n"+
								"时间戳(毫秒): %d\n"+
								"================================================\n",
								userAddress,
								mintAddress,
								tradeTime,
								milliseconds,
							)
						}
					}
				}
			}
		}
	}
}
