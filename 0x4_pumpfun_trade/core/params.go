package core

import (
	"yellowstone-grpc-golang/0x4_pumpfun_trade/types"

	"github.com/gagliardetto/solana-go"
)

// TradeAccounts 定义交易所需的账户
type TradeAccounts struct {
	Mint                   solana.PublicKey
	User                   solana.PublicKey
	BondingCurve           solana.PublicKey
	AssociatedBondingCurve solana.PublicKey
	AssociatedTokenAccount solana.PublicKey
}

// BuyParams 定义购买交易所需的参数
type BuyParams struct {
	Mode      types.TradeMode
	Accounts  TradeAccounts
	Amount    uint64      // token数量
	MaxCost   uint64      // 最大SOL花费
	Blockhash solana.Hash // 添加blockhash参数
}

// SellParams 定义卖出交易所需的参数
type SellParams struct {
	Mode      types.TradeMode
	Accounts  TradeAccounts
	Amount    uint64      // token数量
	MinOut    uint64      // 最小SOL获得
	Blockhash solana.Hash // 添加blockhash参数
}
