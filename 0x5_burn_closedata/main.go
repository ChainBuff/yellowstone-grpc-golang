package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
)

// Config 配置结构体
type Config struct {
	HttpRPC    string
	PrivateKey string
	Command    string
}

// TokenInfo 结构体用于存储代币信息
type TokenInfo struct {
	Mint      string
	ATAs      []string
	Balances  []string
	UiAmounts []float64
}

func main() {
	// 解析命令行参数
	config := parseFlags()

	// 初始化RPC客户端和钱包
	rpcClient := rpc.New(config.HttpRPC)
	wallet, err := solana.WalletFromPrivateKeyBase58(config.PrivateKey)
	if err != nil {
		log.Fatalf("创建钱包失败: %v", err)
	}

	// 根据命令执行相应功能
	switch config.Command {
	case "query":
		queryTokens(rpcClient, wallet)
	case "burn":
		burnTokens(rpcClient, wallet)
	case "closed":
		closeAccounts(rpcClient, wallet)
	case "all":
		burnAndCloseAccounts(rpcClient, wallet)
	case "quick":
		quickQuery(rpcClient, wallet)
	default:
		log.Fatalf("未知命令: %s", config.Command)
	}
}

// 解析命令行参数
func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.HttpRPC, "httprpc", "", "RPC节点地址")
	flag.StringVar(&config.PrivateKey, "private_key", "", "钱包私钥")
	flag.StringVar(&config.Command, "cmd", "", "执行命令 (query/burn/closed/all/quick)")
	flag.Parse()

	if config.HttpRPC == "" || config.PrivateKey == "" || config.Command == "" {
		flag.Usage()
		log.Fatal("必须提供所有参数")
	}

	return config
}

// 添加一个公共函数来获取代币信息
func getTokenAccounts(rpcClient *rpc.Client, wallet *solana.Wallet) map[string]*TokenInfo {
	// 查询所有代币账户
	accounts, err := rpcClient.GetTokenAccountsByOwner(
		context.Background(),
		wallet.PublicKey(),
		&rpc.GetTokenAccountsConfig{
			ProgramId: &solana.TokenProgramID,
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		log.Fatalf("获取代币账户失败: %v", err)
	}

	// 用map来存储代币信息
	tokenInfos := make(map[string]*TokenInfo)
	i := 1
	// 遍历每个ATA账户
	for _, account := range accounts.Value {
		fmt.Printf("开始查询第 %d 个代币\n", i)
		i++
		data := account.Account.Data.GetBinary()
		mint := solana.PublicKey{}
		copy(mint[:], data[0:32])
		mintAddress := mint.String()

		balance, err := rpcClient.GetTokenAccountBalance(
			context.Background(),
			account.Pubkey,
			rpc.CommitmentFinalized,
		)
		if err != nil {
			log.Printf("获取账户 %s 余额失败: %v", account.Pubkey, err)
			continue
		}

		if _, exists := tokenInfos[mintAddress]; !exists {
			tokenInfos[mintAddress] = &TokenInfo{
				Mint: mintAddress,
			}
		}

		tokenInfos[mintAddress].ATAs = append(tokenInfos[mintAddress].ATAs, account.Pubkey.String())
		tokenInfos[mintAddress].Balances = append(tokenInfos[mintAddress].Balances, balance.Value.Amount)
		if balance.Value.UiAmount != nil {
			tokenInfos[mintAddress].UiAmounts = append(tokenInfos[mintAddress].UiAmounts, *balance.Value.UiAmount)
		} else {
			tokenInfos[mintAddress].UiAmounts = append(tokenInfos[mintAddress].UiAmounts, 0)
		}
		//time.Sleep(100 * time.Millisecond)
	}

	return tokenInfos
}

// 修改查询函数
func queryTokens(rpcClient *rpc.Client, wallet *solana.Wallet) {
	tokenInfos := getTokenAccounts(rpcClient, wallet)

	fmt.Printf("钱包地址: %s\n", wallet.PublicKey().String())
	fmt.Printf("SPL代币种类总数: %d\n\n", len(tokenInfos))

	// 统计数据
	var totalAccounts int
	var zeroBalanceAccounts int
	const rentExemptBalance = 0.00203928 // SOL，每个Token账户的租金豁免额

	i := 1
	for _, info := range tokenInfos {
		fmt.Printf("SPL代币 #%d:\n", i)
		fmt.Printf("  Mint地址: %s\n", info.Mint)
		for j := 0; j < len(info.ATAs); j++ {
			totalAccounts++
			fmt.Printf("  ATA账户 #%d:\n", j+1)
			fmt.Printf("    地址: %s\n", info.ATAs[j])
			fmt.Printf("    显示余额: %.9f\n", info.UiAmounts[j])

			// 检查是否为零余额账户
			originalBalance, _ := strconv.ParseUint(info.Balances[j], 10, 64)
			if originalBalance == 0 {
				zeroBalanceAccounts++
			}
		}
		fmt.Println()
		i++
	}

	// 打印统计信息
	fmt.Printf("\n=== 租金统计 ===\n")
	fmt.Printf("总账户数: %d\n", totalAccounts)
	fmt.Printf("零余额账户数: %d\n", zeroBalanceAccounts)
	fmt.Printf("每个账户租金: %.8f SOL\n", rentExemptBalance)
	fmt.Printf("可回收租金: %.8f SOL\n", float64(zeroBalanceAccounts)*rentExemptBalance)
}

// 修改销毁函数
func burnTokens(rpcClient *rpc.Client, wallet *solana.Wallet) {
	fmt.Println("执行burnTokens功能 - 销毁余额小于10的代币")

	tokenInfos := getTokenAccounts(rpcClient, wallet)
	var instructions []solana.Instruction
	var burnCount int

	// 创建工作通道
	type burnTask struct {
		instructions []solana.Instruction
	}
	taskChan := make(chan burnTask, 100)
	doneChan := make(chan bool)
	errorChan := make(chan error, 100)

	// 启动5个工作协程
	for i := 0; i < 5; i++ {
		go func() {
			for task := range taskChan {
				if err := sendTransaction(rpcClient, wallet, task.instructions); err != nil {
					errorChan <- fmt.Errorf("发送交易失败: %v", err)
				}
				time.Sleep(100 * time.Millisecond) // 避免过快发送
			}
			doneChan <- true
		}()
	}

	// 收集需要销毁的代币指令
	for _, info := range tokenInfos {
		for j := 0; j < len(info.ATAs); j++ {
			// 解析原始余额
			fmt.Printf("开始销毁第 %d 个代币\n", j)
			originalBalance, _ := strconv.ParseUint(info.Balances[j], 10, 64)
			// 检查原始余额大于0且UI余额小于10
			if originalBalance > 0 && info.UiAmounts[j] < 1 {
				fmt.Printf("发现小额代币:\n")
				fmt.Printf("  ATA地址: %s\n", info.ATAs[j])
				fmt.Printf("  Mint地址: %s\n", info.Mint)
				fmt.Printf("  原始余额: %s\n", info.Balances[j])
				fmt.Printf("  当前余额: %.9f\n", info.UiAmounts[j])

				ataAccount := solana.MustPublicKeyFromBase58(info.ATAs[j])
				burnAmount := originalBalance // 使用解析后的原始余额

				burnIx := token.NewBurnInstruction(
					burnAmount,
					ataAccount,
					solana.MustPublicKeyFromBase58(info.Mint),
					wallet.PublicKey(),
					[]solana.PublicKey{},
				).Build()

				instructions = append(instructions, burnIx)
				burnCount++

				// 每10个指令打包成一个交易
				if burnCount == 10 {
					taskChan <- burnTask{instructions: instructions}
					instructions = make([]solana.Instruction, 0)
					burnCount = 0
				}
			}
		}
	}

	// 处理剩余的指令
	if len(instructions) > 0 {
		taskChan <- burnTask{instructions: instructions}
	}

	// 关闭任务通道
	close(taskChan)

	// 等待所有工作协程完成
	for i := 0; i < 5; i++ {
		<-doneChan
	}

	// 检查错误
	close(errorChan)
	for err := range errorChan {
		fmt.Printf("错误: %v\n", err)
	}
}

// 添加ComputeBudget相关的辅助函数
func createSetComputeUnitLimitInstruction(units uint32) solana.Instruction {
	data := make([]byte, 5)
	data[0] = 2 // 指令索引为 2
	binary.LittleEndian.PutUint32(data[1:], units)

	return solana.NewInstruction(
		solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111"),
		solana.AccountMetaSlice{},
		data,
	)
}

func createSetComputeUnitPriceInstruction(microLamports uint64) solana.Instruction {
	data := make([]byte, 9)
	data[0] = 3 // 指令索引
	binary.LittleEndian.PutUint64(data[1:], microLamports)

	return solana.NewInstruction(
		solana.MustPublicKeyFromBase58("ComputeBudget111111111111111111111111111111"),
		solana.AccountMetaSlice{},
		data,
	)
}

// 修改sendTransaction函数，添加ComputeBudget指令
func sendTransaction(rpcClient *rpc.Client, wallet *solana.Wallet, instructions []solana.Instruction) error {
	// 创建完整的指令列表，包括ComputeBudget指令
	allInstructions := []solana.Instruction{
		createSetComputeUnitLimitInstruction(200_000), // 设置CU限制
		createSetComputeUnitPriceInstruction(1_000),   // 设置优先级费用为1000 microLamports
	}
	// 添加原有指令
	allInstructions = append(allInstructions, instructions...)

	// 获取最新的blockhash
	recent, err := rpcClient.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("获取最新blockhash失败: %v", err)
	}

	// 创建交易
	tx, err := solana.NewTransaction(
		allInstructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("创建交易失败: %v", err)
	}

	// 签名交易
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if wallet.PublicKey().Equals(key) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("签名交易失败: %v", err)
	}

	// 发送交易
	sig, err := rpcClient.SendTransaction(context.Background(), tx)
	if err != nil {
		return fmt.Errorf("发送交易失败: %v", err)
	}

	fmt.Printf("交易已发送，签名: %s\n", sig.String())
	fmt.Println("等待交易确认...")

	// 等待交易确认
	for i := 0; i < 50; i++ {
		time.Sleep(time.Second)

		// 使用GetTransaction替代GetSignatureStatuses
		tx, err := rpcClient.GetTransaction(
			context.Background(),
			sig,
			&rpc.GetTransactionOpts{
				Commitment: rpc.CommitmentFinalized,
			},
		)
		if err != nil {
			continue
		}

		if tx != nil && tx.Meta != nil {
			if tx.Meta.Err != nil {
				return fmt.Errorf("交易失败: %v", tx.Meta.Err)
			}
			fmt.Println("交易已确认")
			return nil
		}
	}
	return fmt.Errorf("交易确认超时")
}

func closeAccounts(rpcClient *rpc.Client, wallet *solana.Wallet) {
	fmt.Println("执行closeAccounts功能 - 关闭余额为0的ATA账户")

	tokenInfos := getTokenAccounts(rpcClient, wallet)

	// 创建工作通道
	type closeTask struct {
		instructions []solana.Instruction
	}
	taskChan := make(chan closeTask, 100)
	doneChan := make(chan bool)
	errorChan := make(chan error, 100)

	// 启动5个工作协程
	for i := 0; i < 10; i++ {
		go func() {
			for task := range taskChan {
				if err := sendTransaction(rpcClient, wallet, task.instructions); err != nil {
					errorChan <- fmt.Errorf("发送交易失败: %v", err)
				}
				time.Sleep(100 * time.Millisecond)
			}
			doneChan <- true
		}()
	}

	var instructions []solana.Instruction
	var closeCount int

	// 收集需要关闭的账户
	for _, info := range tokenInfos {
		for j := 0; j < len(info.ATAs); j++ {
			// 检查余额是否为0
			originalBalance, _ := strconv.ParseUint(info.Balances[j], 10, 64)
			if originalBalance == 0 {
				fmt.Printf("发现余额为0的账户:\n")
				fmt.Printf("  ATA地址: %s\n", info.ATAs[j])
				fmt.Printf("  Mint地址: %s\n", info.Mint)

				ataAccount := solana.MustPublicKeyFromBase58(info.ATAs[j])

				// 创建关闭账户指令
				closeIx := token.NewCloseAccountInstruction(
					ataAccount,           // 要关闭的ATA
					wallet.PublicKey(),   // 租金接收者
					wallet.PublicKey(),   // 所有者
					[]solana.PublicKey{}, // 额外签名者
				).Build()

				instructions = append(instructions, closeIx)
				closeCount++

				// 每10个指令打包成一个交易
				if closeCount == 10 {
					taskChan <- closeTask{instructions: instructions}
					instructions = make([]solana.Instruction, 0)
					closeCount = 0
				}
			}
		}
	}

	// 处理剩余的指令
	if len(instructions) > 0 {
		taskChan <- closeTask{instructions: instructions}
	}

	// 关闭任务通道
	close(taskChan)

	// 等待所有工作协程完成
	for i := 0; i < 5; i++ {
		<-doneChan
	}

	// 检查错误
	close(errorChan)
	for err := range errorChan {
		fmt.Printf("错误: %v\n", err)
	}
}

func burnAndCloseAccounts(rpcClient *rpc.Client, wallet *solana.Wallet) {
	fmt.Println("执行burnAndCloseAccounts功能 - 销毁代币并关闭账户")

	tokenInfos := getTokenAccounts(rpcClient, wallet)

	// 创建工作通道
	type burnAndCloseTask struct {
		instructions []solana.Instruction
	}
	taskChan := make(chan burnAndCloseTask, 100)
	doneChan := make(chan bool)
	errorChan := make(chan error, 100)

	// 启动5个工作协程
	for i := 0; i < 10; i++ {
		go func() {
			for task := range taskChan {
				if err := sendTransaction(rpcClient, wallet, task.instructions); err != nil {
					errorChan <- fmt.Errorf("发送交易失败: %v", err)
				}
				time.Sleep(100 * time.Millisecond)
			}
			doneChan <- true
		}()
	}

	var instructions []solana.Instruction
	var ixCount int

	// 收集需要处理的账户
	for _, info := range tokenInfos {
		for j := 0; j < len(info.ATAs); j++ {
			originalBalance, _ := strconv.ParseUint(info.Balances[j], 10, 64)
			if originalBalance > 0 && info.UiAmounts[j] < 1 {
				fmt.Printf("发现需要处理的账户:\n")
				fmt.Printf("  ATA地址: %s\n", info.ATAs[j])
				fmt.Printf("  Mint地址: %s\n", info.Mint)
				fmt.Printf("  原始余额: %s\n", info.Balances[j])
				fmt.Printf("  当前余额: %.9f\n", info.UiAmounts[j])

				ataAccount := solana.MustPublicKeyFromBase58(info.ATAs[j])

				// 创建销毁指令
				burnIx := token.NewBurnInstruction(
					originalBalance,
					ataAccount,
					solana.MustPublicKeyFromBase58(info.Mint),
					wallet.PublicKey(),
					[]solana.PublicKey{},
				).Build()

				// 创建关闭账户指令
				closeIx := token.NewCloseAccountInstruction(
					ataAccount,
					wallet.PublicKey(),
					wallet.PublicKey(),
					[]solana.PublicKey{},
				).Build()

				// 添加两个指令
				instructions = append(instructions, burnIx, closeIx)
				ixCount += 2

				// 每10个指令打包成一个交易（5对burn+close指令）
				if ixCount >= 10 {
					taskChan <- burnAndCloseTask{instructions: instructions}
					instructions = make([]solana.Instruction, 0)
					ixCount = 0
				}
			}
		}
	}

	// 处理剩余的指令
	if len(instructions) > 0 {
		taskChan <- burnAndCloseTask{instructions: instructions}
	}

	// 关闭任务通道
	close(taskChan)

	// 等待所有工作协程完成
	for i := 0; i < 5; i++ {
		<-doneChan
	}

	// 检查错误
	close(errorChan)
	for err := range errorChan {
		fmt.Printf("错误: %v\n", err)
	}
}

func quickQuery(rpcClient *rpc.Client, wallet *solana.Wallet) {
	fmt.Println("快速查询账户信息...")

	// 只查询账户列表，不获取余额
	accounts, err := rpcClient.GetTokenAccountsByOwner(
		context.Background(),
		wallet.PublicKey(),
		&rpc.GetTokenAccountsConfig{
			ProgramId: &solana.TokenProgramID,
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		log.Fatalf("获取代币账户失败: %v", err)
	}

	// 统计数据
	totalAccounts := len(accounts.Value)
	const rentExemptBalance = 0.00203928 // SOL，每个Token账户的租金豁免额

	fmt.Printf("\n=== 账户统计 ===\n")
	fmt.Printf("钱包地址: %s\n", wallet.PublicKey().String())
	fmt.Printf("总账户数: %d\n", totalAccounts)
	fmt.Printf("预计最大可回收租金: %.8f SOL\n", float64(totalAccounts)*rentExemptBalance)
	fmt.Printf("注意: 实际可回收租金取决于零余额账户数量\n")
}
