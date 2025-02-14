package calculator

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/shopspring/decimal"
)

// AmountsCalculator 金额计算器
type AmountsCalculator struct {
	client *rpc.Client
}

func NewAmountsCalculator(endpoint string) *AmountsCalculator {
	return &AmountsCalculator{
		client: rpc.New(endpoint),
	}
}

// CalculateBuyAmounts 计算购买所需的金额参数
func (c *AmountsCalculator) CalculateBuyAmounts(ctx context.Context, bondingCurve solana.PublicKey, solAmount float64, slippage float64) (*BuyAmounts, error) {
	// 获取bonding curve账户数据
	account, err := c.client.GetAccountInfo(ctx, bondingCurve)
	if err != nil {
		return nil, fmt.Errorf("获取bonding curve账户失败: %w", err)
	}

	data := account.Value.Data.GetBinary()
	data = data[8:] // 跳过discriminator

	//virtualTokenReserves := binary.LittleEndian.Uint64(data[0:8])
	virtualSolReserves := binary.LittleEndian.Uint64(data[8:16])

	// 计算bonding curve
	base := float64(1073000191)
	coef := float64(32190005730)
	oneBillion := float64(1000000000)

	virtualSol := float64(virtualSolReserves)
	virtualSol2 := virtualSol + (solAmount * oneBillion)

	pay1 := base - (coef / (virtualSol / oneBillion))
	pay2 := base - (coef / (virtualSol2 / oneBillion))

	tokenAmount := pay2 - pay1
	tokenAmountWithDecimals := uint64(tokenAmount * 1e6)

	// 计算基础金额
	solAmountWithDecimals := uint64(solAmount * 1e9)

	// 计算最大花费（考虑滑点）
	maxCost := uint64(float64(solAmountWithDecimals) * (1 + slippage))

	return &BuyAmounts{
		TokenAmount: tokenAmountWithDecimals,
		SolAmount:   solAmountWithDecimals,
		MaxCost:     maxCost,
	}, nil
}

// BuyAmounts 存储购买相关的金额
type BuyAmounts struct {
	TokenAmount uint64 // 获得的token数量
	SolAmount   uint64 // 支付的SOL数量
	MaxCost     uint64 // 最大花费（含滑点）
}

// CalculateSellAmounts 计算卖出所需的金额参数
func (c *AmountsCalculator) CalculateSellAmounts(ctx context.Context, bondingCurve solana.PublicKey, tokenAmount float64, slippage float64) (*SellAmounts, error) {
	// 获取bonding curve账户数据
	account, err := c.client.GetAccountInfo(ctx, bondingCurve)
	if err != nil {
		return nil, fmt.Errorf("获取bonding curve账户失败: %w", err)
	}

	data := account.Value.Data.GetBinary()
	data = data[8:] // 跳过discriminator

	virtualTokenReserves := binary.LittleEndian.Uint64(data[0:8])
	virtualSolReserves := binary.LittleEndian.Uint64(data[8:16])

	// 计算价格
	virtualSol := decimal.NewFromInt(int64(virtualSolReserves)).Div(decimal.NewFromInt(1e9))
	virtualToken := decimal.NewFromInt(int64(virtualTokenReserves)).Div(decimal.NewFromInt(1e6))
	price := virtualSol.Div(virtualToken).Round(12)

	// 计算token数量（需要乘以1e6转换为正确精度）
	tokenAmountWithDecimals := uint64(tokenAmount * 1e6)

	// 计算预期获得的SOL数量
	expectedSol, _ := price.Float64()
	solAmount := expectedSol * tokenAmount // 将token数量转换为原始单位
	solAmountWithDecimals := uint64(solAmount * 1e9)

	// 计算最小获得数量（考虑滑点）
	minOut := uint64(float64(solAmountWithDecimals) * (1 - slippage))

	return &SellAmounts{
		TokenAmount: tokenAmountWithDecimals,
		SolAmount:   solAmountWithDecimals,
		MinOut:      minOut,
	}, nil
}

// SellAmounts 存储卖出相关的金额
type SellAmounts struct {
	TokenAmount uint64 // 卖出的token数量
	SolAmount   uint64 // 预期获得的SOL数量
	MinOut      uint64 // 最小获得数量（考虑滑点）
}
