package calculator

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	PUMP_PROGRAM_ID = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
)

// AccountsCalculator 账户计算器
type AccountsCalculator struct {
	client *rpc.Client
}

func NewAccountsCalculator(endpoint string) *AccountsCalculator {
	return &AccountsCalculator{
		client: rpc.New(endpoint),
	}
}

// CalculatePumpAccounts 计算Pump所需的所有账户
func (c *AccountsCalculator) CalculatePumpAccounts(ctx context.Context, mint, user solana.PublicKey) (*PumpAccounts, error) {
	// 计算bonding curve地址
	bondingCurve, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("bonding-curve"),
			mint.Bytes(),
		},
		PUMP_PROGRAM_ID,
	)
	if err != nil {
		return nil, fmt.Errorf("计算bonding curve失败: %w", err)
	}

	// 计算associated bonding curve地址
	associatedBondingCurve, _, err := solana.FindAssociatedTokenAddress(
		bondingCurve,
		mint,
	)
	if err != nil {
		return nil, fmt.Errorf("计算associated bonding curve失败: %w", err)
	}

	// 计算用户的associated token account
	associatedTokenAccount, _, err := solana.FindAssociatedTokenAddress(
		user,
		mint,
	)
	if err != nil {
		return nil, fmt.Errorf("计算associated token account失败: %w", err)
	}

	return &PumpAccounts{
		BondingCurve:           bondingCurve,
		AssociatedBondingCurve: associatedBondingCurve,
		AssociatedTokenAccount: associatedTokenAccount,
	}, nil
}

// PumpAccounts 存储所有需要的账户地址
type PumpAccounts struct {
	BondingCurve           solana.PublicKey
	AssociatedBondingCurve solana.PublicKey
	AssociatedTokenAccount solana.PublicKey
}
