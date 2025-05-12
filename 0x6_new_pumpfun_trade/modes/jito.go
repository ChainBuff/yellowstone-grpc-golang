package modes

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"time"

	"yellowstone-grpc-golang/0x6_new_pumpfun_trade/types"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	jitorpc "github.com/jito-labs/jito-go-rpc"
	"github.com/mr-tron/base58"
)

// Jito tip 账户列表
var jitoTipAccounts = []string{
	"HFqU5x63VTqvQss8hp11i4wVV8bD44PvwucfZ2bU7gRe",
	"DfXygSm4jCyNCybVYYK6DwvWqjKee8pbDmJGcLWNDXjh",
	"3AVi9Tg9Uo68tJfuvoKvqKNWKkC5wPdSSdeBnizKZ6jT",
	"Cw8CFyM9FkoMi7K7Crf6HNQqf4uEMzpKw6QNghXLvLkY",
	"96gYZGLnJYVFmbjzopPSU6QiEV5fGqZNyN9nmNhvrZU5",
	"ADuUkR4vqLUMWXxW9gh6D6L8pMSawimctcNZ5pGwDcEt",
	"DttWaMuVvTiduZRnguLF7jNxTgiMBZ1hyAumKUiL2KRL",
	"ADaUMid9yfUytqMBgopwjb2DTLSokTSzL1zt6iGPaS49",
}

// SendTransactionJito 发送交易到 Jito
func SendTransactionJito(wallet *solana.Wallet, instructions []solana.Instruction, params *types.JitoParams, blockhash solana.Hash) (string, error) {
	log.Printf("[Jito] 开始发送交易: %s", time.Now().Format("2006-01-02 15:04:05.000"))

	// 初始化 Jito 客户端
	jitoClient := jitorpc.NewJitoJsonRpcClient(params.BundleUrl, "")
	debug := false
	jitoClient.Debug = &debug

	// 获取随机 tip 账户
	tipAccount := getRandomJitoTipAccount()

	// 构建 tip 转账指令
	tipIx := system.NewTransferInstruction(
		params.TipAmount,
		wallet.PublicKey(),
		solana.MustPublicKeyFromBase58(tipAccount),
	).Build()

	// 合并所有指令
	allInstructions := []solana.Instruction{tipIx}
	allInstructions = append(allInstructions, instructions...)

	// 构建交易
	tx, err := solana.NewTransaction(
		allInstructions,
		blockhash,
		solana.TransactionPayer(wallet.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create transaction: %w", err)
	}

	// 签名交易
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(wallet.PublicKey()) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 序列化交易
	serializedTx, err := tx.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %w", err)
	}
	base58EncodedTx := base58.Encode(serializedTx)

	// 发送交易
	txnRequest := []string{base58EncodedTx}
	result, err := jitoClient.SendTxn(txnRequest, false)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	if result == nil {
		return "", fmt.Errorf("transaction result is nil")
	}

	log.Printf("[Jito] 交易发送完成: %s", time.Now().Format("2006-01-02 15:04:05.000"))
	return tx.Signatures[0].String(), nil
}

// 获取随机的 tip 账户
func getRandomJitoTipAccount() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(jitoTipAccounts))))
	if err != nil {
		return jitoTipAccounts[0]
	}
	return jitoTipAccounts[n.Int64()]
}
