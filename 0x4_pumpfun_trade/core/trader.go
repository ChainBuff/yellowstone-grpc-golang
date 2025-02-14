package core

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"yellowstone-grpc-golang/0x4_pumpfun_trade/calculator"
	"yellowstone-grpc-golang/0x4_pumpfun_trade/modes"
	"yellowstone-grpc-golang/0x4_pumpfun_trade/pump"
	"yellowstone-grpc-golang/0x4_pumpfun_trade/types"
	"yellowstone-grpc-golang/0x4_pumpfun_trade/utils"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"
)

// Trader 定义交易接口
type Trader interface {
	Buy(ctx context.Context, params BuyParams) (string, error) // 返回交易签名
	Sell(ctx context.Context, params SellParams) (string, error)
}

// TraderImpl 交易实现
type TraderImpl struct {
	Client       *rpc.Client
	Wallet       *solana.Wallet
	AccountsCalc *calculator.AccountsCalculator
	AmountsCalc  *calculator.AmountsCalculator
	Config       *utils.Config // 添加全局配置
}

// NewTrader 创建交易实例
func NewTrader(endpoint string, privateKey solana.PrivateKey, config *utils.Config) *TraderImpl {
	// 使用正确的 RPC URL 创建客户端
	client := rpc.New(endpoint)
	// 设置请求头，避免一些 RPC 节点的限制

	return &TraderImpl{
		Client:       client,
		Wallet:       &solana.Wallet{PrivateKey: privateKey},
		AccountsCalc: calculator.NewAccountsCalculator(endpoint),
		AmountsCalc:  calculator.NewAmountsCalculator(endpoint),
		Config:       config,
	}
}

// Buy 执行买入交易
func (t *TraderImpl) Buy(ctx context.Context, params types.OrderParams) (string, error) {
	log.Printf("开始执行买入交易...")
	log.Printf("Mint地址: %s", params.TradeParams.Accounts.Mint)
	log.Printf("购买数量: %d", params.TradeParams.Amount)
	log.Printf("time: %v", time.Now().Format("2006-01-02 15:04:05.000"))

	// 构建交易指令
	instructions := []solana.Instruction{}

	// 添加计算单元限制指令（通常设置为 200,000 units）
	instructions = append(instructions, createSetComputeUnitLimitInstruction(200_000))

	// 添加优先级费用指令
	instructions = append(instructions, createSetComputeUnitPriceInstruction(params.TradeParams.FeeConfig.PriorityFee))

	if t.Config.SkipATACheck {
		// 快速模式：直接添加创建ATA指令
		log.Printf("快速模式：跳过ATA检查，直接创建ATA账户")
		createATAIx := associatedtokenaccount.NewCreateInstruction(
			params.TradeParams.Accounts.User, // payer
			params.TradeParams.Accounts.User, // owner
			params.TradeParams.Accounts.Mint, // mint
		).Build()
		instructions = append(instructions, createATAIx)
	} else {
		// 安全模式：先检查再创建
		account, err := t.Client.GetAccountInfo(ctx, params.TradeParams.Accounts.AssociatedTokenAccount)
		if err != nil || account.Value == nil {
			log.Printf("ATA账户不存在，创建新账户...")
			createATAIx := associatedtokenaccount.NewCreateInstruction(
				params.TradeParams.Accounts.User, // payer
				params.TradeParams.Accounts.User, // owner
				params.TradeParams.Accounts.Mint, // mint
			).Build()
			instructions = append(instructions, createATAIx)
		}
	}

	// 构建Buy指令
	buyInstruction := pump.NewBuyInstructionBuilder().
		SetGlobalAccount(GLOBAL_ACCOUNT).
		SetFeeRecipientAccount(FEE_RECIPIENT).
		SetMintAccount(params.TradeParams.Accounts.Mint).
		SetBondingCurveAccount(params.TradeParams.Accounts.BondingCurve).
		SetAssociatedBondingCurveAccount(params.TradeParams.Accounts.AssociatedBondingCurve).
		SetAssociatedUserAccount(params.TradeParams.Accounts.AssociatedTokenAccount).
		SetUserAccount(params.TradeParams.Accounts.User).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetRentAccount(solana.SysVarRentPubkey).
		SetEventAuthorityAccount(EVENT_AUTHORITY).
		SetProgramAccount(PUMP_PROGRAM_ID).
		SetAmount(params.TradeParams.Amount).
		SetMaxSolCost(params.TradeParams.Cost)

	instructions = append(instructions, buyInstruction.Build())

	// 创建通道
	resultChan := make(chan string, 4)
	errChan := make(chan error, 4)
	timeout := time.After(30 * time.Second)

	// 检查各个模式是否启用并准备发送
	if params.NextBlockParams != nil && params.NextBlockParams.Enable {
		log.Printf("[NextBlock] 准备发送交易...")
		go func() {
			// 构建交易并发送到 NextBlock
			sig, err := modes.SendTransactionNextBlock(
				t.Wallet,
				instructions,
				params.NextBlockParams,
				params.TradeParams.Blockhash,
			)
			if err != nil {
				errChan <- fmt.Errorf("NextBlock发送失败: %w", err)
				return
			}
			resultChan <- sig
		}()
	}
	if params.JitoParams != nil && params.JitoParams.Enable {
		log.Printf("[Jito] 准备发送交易...")
		go func() {
			sig, err := modes.SendTransactionJito(
				t.Wallet,
				instructions,
				params.JitoParams,
				params.TradeParams.Blockhash,
			)
			if err != nil {
				errChan <- fmt.Errorf("jito发送失败: %w", err)
				return
			}
			resultChan <- sig
		}()
	}
	if params.TemporalParams != nil && params.TemporalParams.Enable {
		log.Printf("[Temporal] 准备发送交易...")
		go func() {
			sig, err := modes.SendTransactionTemporal(
				t.Wallet,
				instructions,
				params.TemporalParams,
				params.TradeParams.Blockhash,
			)
			if err != nil {
				errChan <- fmt.Errorf("temporal发送失败: %w", err)
				return
			}
			resultChan <- sig
		}()
	}

	if t.Config.Normal.Enable {
		log.Printf("[Normal] 准备发送交易...")
		go func() {
			// 构建交易
			tx, err := solana.NewTransaction(
				instructions,
				params.TradeParams.Blockhash,
				solana.TransactionPayer(t.Wallet.PublicKey()),
			)
			if err != nil {
				errChan <- fmt.Errorf("创建交易失败: %w", err)
				return
			}

			// 签名交易
			_, err = tx.Sign(
				func(key solana.PublicKey) *solana.PrivateKey {
					if key.Equals(t.Wallet.PublicKey()) {
						return &t.Wallet.PrivateKey
					}
					return nil
				},
			)
			if err != nil {
				errChan <- fmt.Errorf("签名交易失败: %w", err)
				return
			}

			// 发送交易
			sig, err := t.Client.SendTransaction(ctx, tx)
			if err != nil {
				errChan <- fmt.Errorf("发送交易失败: %w", err)
				return
			}

			resultChan <- sig.String()
		}()
	}

	// 等待结果
	select {
	case sig := <-resultChan:
		log.Printf("交易已发送! 交易ID: %s", sig)
		return sig, nil
	case err := <-errChan:
		return "", err
	case <-timeout:
		return "", fmt.Errorf("transaction timeout")
	}
}

// Sell 实现卖出交易
func (t *TraderImpl) Sell(ctx context.Context, params types.OrderParams) (string, error) {
	log.Printf("开始执行卖出交易...")
	log.Printf("Mint地址: %s", params.TradeParams.Accounts.Mint)
	log.Printf("卖出数量: %d", params.TradeParams.Amount)
	log.Printf("最小获得: %d SOL", params.TradeParams.Cost)
	log.Printf("time: %v", time.Now().Format("2006-01-02 15:04:05.000"))

	// 构建交易指令
	instructions := []solana.Instruction{}

	// 添加计算单元限制指令（通常设置为 200,000 units）
	instructions = append(instructions, createSetComputeUnitLimitInstruction(200_000))

	// 添加优先级费用指令
	instructions = append(instructions, createSetComputeUnitPriceInstruction(params.TradeParams.FeeConfig.PriorityFee))

	// 构建Sell指令
	sellInstruction := pump.NewSellInstructionBuilder().
		SetGlobalAccount(GLOBAL_ACCOUNT).
		SetFeeRecipientAccount(FEE_RECIPIENT).
		SetMintAccount(params.TradeParams.Accounts.Mint).
		SetBondingCurveAccount(params.TradeParams.Accounts.BondingCurve).
		SetAssociatedBondingCurveAccount(params.TradeParams.Accounts.AssociatedBondingCurve).
		SetAssociatedUserAccount(params.TradeParams.Accounts.AssociatedTokenAccount).
		SetUserAccount(params.TradeParams.Accounts.User).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetAssociatedTokenProgramAccount(associatedtokenaccount.ProgramID).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetEventAuthorityAccount(EVENT_AUTHORITY).
		SetProgramAccount(PUMP_PROGRAM_ID).
		SetAmount(params.TradeParams.Amount).
		SetMinSolOutput(params.TradeParams.Cost)

	ix, err := sellInstruction.ValidateAndBuild()
	if err != nil {
		return "", fmt.Errorf("构建指令失败: %w", err)
	}

	instructions = append(instructions, ix)

	// 创建通道
	resultChan := make(chan string, 4)
	errChan := make(chan error, 4)
	timeout := time.After(30 * time.Second)

	// 检查各个模式是否启用并准备发送
	if params.NextBlockParams != nil && params.NextBlockParams.Enable {
		log.Printf("[NextBlock] 准备发送交易...")
		go func() {
			sig, err := modes.SendTransactionNextBlock(
				t.Wallet,
				instructions,
				params.NextBlockParams,
				params.TradeParams.Blockhash,
			)
			if err != nil {
				errChan <- fmt.Errorf("NextBlock发送失败: %w", err)
				return
			}
			resultChan <- sig
		}()
	}

	if params.JitoParams != nil && params.JitoParams.Enable {
		log.Printf("[Jito] 准备发送交易...")
		go func() {
			sig, err := modes.SendTransactionJito(
				t.Wallet,
				instructions,
				params.JitoParams,
				params.TradeParams.Blockhash,
			)
			if err != nil {
				errChan <- fmt.Errorf("jito发送失败: %w", err)
				return
			}
			resultChan <- sig
		}()
	}

	if params.TemporalParams != nil && params.TemporalParams.Enable {
		log.Printf("[Temporal] 准备发送交易...")
		go func() {
			sig, err := modes.SendTransactionTemporal(
				t.Wallet,
				instructions,
				params.TemporalParams,
				params.TradeParams.Blockhash,
			)
			if err != nil {
				errChan <- fmt.Errorf("temporal发送失败: %w", err)
				return
			}
			resultChan <- sig
		}()
	}

	if t.Config.Normal.Enable {
		log.Printf("[Normal] 准备发送交易...")
		go func() {
			// 构建交易
			tx, err := solana.NewTransaction(
				instructions,
				params.TradeParams.Blockhash,
				solana.TransactionPayer(t.Wallet.PublicKey()),
			)
			if err != nil {
				errChan <- fmt.Errorf("创建交易失败: %w", err)
				return
			}

			// 签名交易
			_, err = tx.Sign(
				func(key solana.PublicKey) *solana.PrivateKey {
					if key.Equals(t.Wallet.PublicKey()) {
						return &t.Wallet.PrivateKey
					}
					return nil
				},
			)
			if err != nil {
				errChan <- fmt.Errorf("签名交易失败: %w", err)
				return
			}

			// 发送交易
			sig, err := t.Client.SendTransaction(ctx, tx)
			if err != nil {
				errChan <- fmt.Errorf("发送交易失败: %w", err)
				return
			}

			resultChan <- sig.String()
		}()
	}

	// 等待结果
	select {
	case sig := <-resultChan:
		log.Printf("交易已发送! 交易ID: %s", sig)
		return sig, nil
	case err := <-errChan:
		return "", err
	case <-timeout:
		return "", fmt.Errorf("transaction timeout")
	}
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

// createSetComputeUnitLimitInstruction 创建设置计算单元限制的指令
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
