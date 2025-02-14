package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"yellowstone-grpc-golang/0x4_pumpfun_trade/core"
	"yellowstone-grpc-golang/0x4_pumpfun_trade/types"
	"yellowstone-grpc-golang/0x4_pumpfun_trade/utils"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// 定义操作类型
type Operation string

const (
	Buy  Operation = "buy"
	Sell Operation = "sell"
)

// 验证操作类型
func (o Operation) isValid() bool {
	return o == Buy || o == Sell
}

// 获取项目根目录
func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "..")
}

// validateMintAddress 验证mint地址是否有效
func validateMintAddress(mintAddress string) error {
	_, err := solana.PublicKeyFromBase58(mintAddress)
	if err != nil {
		return fmt.Errorf("无效的mint地址: %w", err)
	}
	return nil
}

func main() {
	// 定义命令行参数
	mintAddr := flag.String("mint", "", "SPL Token的Mint地址")
	operation := flag.String("op", "", "操作类型 (buy/sell)")
	amount := flag.Float64("amount", 0.0, "操作数量")

	// 解析命令行参数
	flag.Parse()

	// 验证参数
	if *mintAddr == "" {
		log.Fatal("请提供mint地址")
	}

	if err := validateMintAddress(*mintAddr); err != nil {
		log.Fatalf("无效的mint地址: %v", err)
	}

	op := Operation(strings.ToLower(*operation))
	if !op.isValid() {
		log.Fatal("无效的操作类型，请使用 buy 或 sell")
	}

	if *amount <= 0 {
		log.Fatal("amount必须大于0")
	}

	// 获取配置文件的绝对路径
	rootDir := getRootDir()
	fmt.Println("rootDir", rootDir)
	configPath := filepath.Join(rootDir, "0x4_pumpfun_trade", "config", "config.yaml")

	// 加载配置
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 打印操作信息
	fmt.Printf("\n操作信息:\n")
	fmt.Printf("Mint地址: %s\n", *mintAddr)
	fmt.Printf("操作类型: %s\n", op)
	fmt.Printf("数量: %.6f\n", *amount)

	// 打印配置内容
	fmt.Printf("\n已加载配置:\n")
	fmt.Printf("Private Key: %s\n", config.PrivateKey)

	// 创建交易实例
	privateKey, err := solana.PrivateKeyFromBase58(config.PrivateKey)
	if err != nil {
		log.Fatalf("私钥解析失败: %v", err)
	}

	trader := core.NewTrader(config.HttpRpcUrl, privateKey, config)

	// 计算所需账户
	ctx := context.Background()
	pumpAccounts, err := trader.AccountsCalc.CalculatePumpAccounts(
		ctx,
		solana.MustPublicKeyFromBase58(*mintAddr),
		privateKey.PublicKey(),
	)
	if err != nil {
		log.Fatalf("计算账户失败: %v", err)
	}

	// 准备交易参数
	accounts := core.TradeAccounts{
		Mint:                   solana.MustPublicKeyFromBase58(*mintAddr),
		User:                   privateKey.PublicKey(),
		BondingCurve:           pumpAccounts.BondingCurve,
		AssociatedBondingCurve: pumpAccounts.AssociatedBondingCurve,
		AssociatedTokenAccount: pumpAccounts.AssociatedTokenAccount,
	}

	if op == Buy {

		// 计算购买金额
		amounts, err := trader.AmountsCalc.CalculateBuyAmounts(
			ctx,
			pumpAccounts.BondingCurve,
			*amount,
			config.BuySlippage,
		)
		if err != nil {
			log.Fatalf("计算购买金额失败: %v", err)
		}

		// 获取最新的blockhash
		recent, err := trader.Client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
		if err != nil {
			log.Fatalf("获取recent blockhash失败: %v", err)
		}

		// 构建交易基础参数
		tradeParams := types.TradeParams{
			Accounts: types.TradeAccounts{
				User:                   accounts.User,
				Mint:                   accounts.Mint,
				BondingCurve:           accounts.BondingCurve,
				AssociatedBondingCurve: accounts.AssociatedBondingCurve,
				AssociatedTokenAccount: accounts.AssociatedTokenAccount,
			},
			Amount:    amounts.TokenAmount,
			Cost:      amounts.MaxCost,
			Blockhash: recent.Value.Blockhash,
			FeeConfig: types.FeeConfig{
				Slippage:    config.BuySlippage,
				PriorityFee: uint64(config.BuyPriorityFee * 1e9),
			},
		}

		// 构建完整的交易参数
		orderParams := types.OrderParams{
			TradeParams: tradeParams,
			NextBlockParams: &types.NextBlockParams{
				Enable:    config.NextBlock.Enable,
				BundleUrl: config.NextBlock.RpcUrl,
				ApiKey:    config.NextBlock.ApiKeys[0],
				TipAmount: uint64(config.NextBlock.BuyTip * 1e9),
			},
			JitoParams: &types.JitoParams{
				Enable:    config.Jito.Enable,
				BundleUrl: config.Jito.RpcUrl,
				TipAmount: uint64(config.Jito.BuyTip * 1e9),
			},
			TemporalParams: &types.TemporalParams{
				Enable:    config.Temporal.Enable,
				BundleUrl: config.Temporal.RpcUrl,
				ApiKey:    config.Temporal.ApiKeys[0],
				TipAmount: uint64(config.Temporal.BuyTip * 1e9),
			},
		}

		// 执行交易
		txid, err := trader.Buy(ctx, orderParams)
		if err != nil {
			log.Fatalf("买入失败: %v", err)
		}
		fmt.Printf("买入发送成功! 交易ID: %s\n", txid)
		fmt.Printf("预计获得: %d token (token x 1e6)\n", amounts.TokenAmount)
		// 计算并打印耗时
	} else {
		// 先计算卖出金额
		amounts, err := trader.AmountsCalc.CalculateSellAmounts(
			ctx,
			pumpAccounts.BondingCurve,
			*amount,
			config.SellSlippage,
		)
		if err != nil {
			log.Fatalf("计算卖出金额失败: %v", err)
		}

		// 获取最新的blockhash
		recent, err := trader.Client.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
		if err != nil {
			log.Fatalf("获取recent blockhash失败: %v", err)
		}

		// 构建交易基础参数
		tradeParams := types.TradeParams{
			Accounts: types.TradeAccounts{
				User:                   accounts.User,
				Mint:                   accounts.Mint,
				BondingCurve:           accounts.BondingCurve,
				AssociatedBondingCurve: accounts.AssociatedBondingCurve,
				AssociatedTokenAccount: accounts.AssociatedTokenAccount,
			},
			Amount:    amounts.TokenAmount,
			Cost:      amounts.MinOut, // 卖出时 Cost 表示最小获得数量
			Blockhash: recent.Value.Blockhash,
			FeeConfig: types.FeeConfig{
				Slippage:    config.SellSlippage,
				PriorityFee: uint64(config.SellPriorityFee * 1e9),
			},
		}

		// 构建完整的交易参数
		orderParams := types.OrderParams{
			TradeParams: tradeParams,
			NextBlockParams: &types.NextBlockParams{
				Enable:    config.NextBlock.Enable,
				BundleUrl: config.NextBlock.RpcUrl,
				ApiKey:    config.NextBlock.ApiKeys[0],
				TipAmount: uint64(config.NextBlock.SellTip * 1e9),
			},
			JitoParams: &types.JitoParams{
				Enable:    config.Jito.Enable,
				BundleUrl: config.Jito.RpcUrl,
				TipAmount: uint64(config.Jito.SellTip * 1e9),
			},
			TemporalParams: &types.TemporalParams{
				Enable:    config.Temporal.Enable,
				BundleUrl: config.Temporal.RpcUrl,
				ApiKey:    config.Temporal.ApiKeys[0],
				TipAmount: uint64(config.Temporal.SellTip * 1e9),
			},
		}

		// 执行卖出交易
		txid, err := trader.Sell(ctx, orderParams)
		if err != nil {
			log.Fatalf("卖出失败: %v", err)
		}

		fmt.Printf("卖出发送成功! 交易ID: %s\n", txid)
		fmt.Printf("卖出数量（token x 1e6 ）: %d token\n", amounts.TokenAmount)
	}

}
