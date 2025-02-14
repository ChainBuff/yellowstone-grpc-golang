package types

import (
	"github.com/gagliardetto/solana-go"
)

// TradeMode 定义交易模式
type TradeMode string

const (
	Normal    TradeMode = "normal"
	Jito      TradeMode = "jito"
	NextBlock TradeMode = "nextblock"
	Temportal TradeMode = "temportal"
)

// ValidateMode 验证交易模式是否有效
func (m TradeMode) IsValid() bool {
	switch m {
	case Normal, Jito, NextBlock, Temportal:
		return true
	default:
		return false
	}
}

// String 实现 Stringer 接口
func (m TradeMode) String() string {
	return string(m)
}

// TradeAccounts 定义交易所需的账户
type TradeAccounts struct {
	User                   solana.PublicKey
	Mint                   solana.PublicKey
	BondingCurve           solana.PublicKey
	AssociatedBondingCurve solana.PublicKey
	AssociatedTokenAccount solana.PublicKey
}

// FeeConfig 定义费用相关配置
type FeeConfig struct {
	Slippage    float64 // 滑点
	PriorityFee uint64  // 优先级费用（lamports）
}

// TradeParams 定义所有交易模式共用的基础参数
type TradeParams struct {
	Accounts  TradeAccounts
	Amount    uint64
	Cost      uint64 // 买入时是最大花费，卖出时是最小获得
	Blockhash solana.Hash
	FeeConfig FeeConfig
}

// NextBlockParams 定义 NextBlock 模式特有的参数
type NextBlockParams struct {
	BundleUrl string
	Enable    bool
	ApiKey    string
	TipAmount uint64
}

// JitoParams 定义 Jito 模式特有的参数
type JitoParams struct {
	Enable    bool
	BundleUrl string
	TipAmount uint64
}

// TemportalParams 定义 Temportal 模式特有的参数
type TemporalParams struct {
	Enable    bool
	BundleUrl string
	ApiKey    string
	TipAmount uint64
}

// OrderParams 定义完整的交易参数（适用于买入和卖出）
type OrderParams struct {
	TradeParams                      // 基础参数
	NextBlockParams *NextBlockParams // NextBlock 特有参数
	JitoParams      *JitoParams      // Jito 特有参数
	TemporalParams  *TemporalParams  // Temportal 特有参数
}
