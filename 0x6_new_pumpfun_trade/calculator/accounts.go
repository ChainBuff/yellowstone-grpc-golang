package calculator

import (
	"context"
	"fmt"
	"log"

	"yellowstone-grpc-golang/0x6_new_pumpfun_trade/pump"

	ag_binary "github.com/gagliardetto/binary"
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

	// 从Solana获取BondingCurve账户数据
	accountInfo, err := c.client.GetAccountInfoWithOpts(
		ctx,
		bondingCurve,
		&rpc.GetAccountInfoOpts{
			Commitment: rpc.CommitmentConfirmed,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("获取bonding curve账户数据失败: %w", err)
	}

	if accountInfo == nil || accountInfo.Value == nil || accountInfo.Value.Data.GetBinary() == nil {
		return nil, fmt.Errorf("bonding curve账户数据为空")
	}

	// 解析BondingCurve账户数据
	// 跳过账户数据的前8个字节，这通常是锚点(Anchor)程序的鉴别符
	data := accountInfo.Value.Data.GetBinary()[8:]
	bondingCurveData := &pump.BondingCurve{}
	decoder := ag_binary.NewBorshDecoder(data)
	if err := bondingCurveData.UnmarshalWithDecoder(decoder); err != nil {
		return nil, fmt.Errorf("解析bonding curve数据失败: %w", err)
	}

	// 输出调试信息
	log.Printf("从BondingCurve获取到Creator: %s", bondingCurveData.Creator.String())

	// 验证Creator是否为零地址
	if bondingCurveData.Creator.IsZero() || bondingCurveData.Creator.Equals(solana.SystemProgramID) {
		log.Printf("警告: BondingCurve的Creator字段为零地址或系统程序ID，可能导致创建者费用功能无法正常工作")
		// 注意：我们继续执行，而不是返回错误，因为这个字段可能会由后端服务动态设置
	}

	// 在Anchor中，creator_vault是由以下种子派生的:
	// - "creator-vault" 常量字符串
	// - bonding_curve.creator

	// 直接从bonding_curve.creator计算creator_vault
	creatorVault, _, err := solana.FindProgramAddress(
		[][]byte{
			[]byte("creator-vault"),
			bondingCurveData.Creator.Bytes(),
		},
		PUMP_PROGRAM_ID,
	)
	if err != nil {
		return nil, fmt.Errorf("计算creator vault失败: %w", err)
	}

	log.Printf("计算得到的CreatorVault: %s", creatorVault.String())

	return &PumpAccounts{
		BondingCurve:           bondingCurve,
		AssociatedBondingCurve: associatedBondingCurve,
		AssociatedTokenAccount: associatedTokenAccount,
		CreatorVault:           creatorVault,
	}, nil
}

// PumpAccounts 存储所有需要的账户地址
type PumpAccounts struct {
	BondingCurve           solana.PublicKey
	AssociatedBondingCurve solana.PublicKey
	AssociatedTokenAccount solana.PublicKey
	CreatorVault           solana.PublicKey
}
