package modes

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"

	"yellowstone-grpc-golang/0x6_new_pumpfun_trade/types"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

// NextBlock tip 账户列表
var nextblockTipAccounts = []string{
	"NextbLoCkVtMGcV47JzewQdvBpLqT9TxQFozQkN98pE",
	"NexTbLoCkWykbLuB1NkjXgFWkX9oAtcoagQegygXXA2",
	"NeXTBLoCKs9F1y5PJS9CKrFNNLU1keHW71rfh7KgA1X",
	"NexTBLockJYZ7QD7p2byrUa6df8ndV2WSd8GkbWqfbb",
	"neXtBLock1LeC67jYd1QdAa32kbVeubsfPNTJC1V5At",
	"nEXTBLockYgngeRmRrjDV31mGSekVPqZoMGhQEZtPVG",
	"NEXTbLoCkB51HpLBLojQfpyVAMorm3zzKg7w9NFdqid",
	"nextBLoCkPMgmG8ZgJtABeScP35qLa2AMCNKntAP7Xc",
}

// API 请求相关结构
type Payload struct {
	Transaction            TransactionMessage `json:"transaction"`
	FrontRunningProtection bool               `json:"frontRunningProtection"`
}

type TransactionMessage struct {
	Content string `json:"content"`
}

// NextBlock 结构体
type NextBlock struct {
	Wallet *solana.Wallet
}

// 获取随机的 tip 账户
func getRandomNextblockTipAccount() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(nextblockTipAccounts))))
	if err != nil {
		return nextblockTipAccounts[0]
	}
	return nextblockTipAccounts[n.Int64()]
}

// SendTransaction 发送交易到 NextBlock
func SendTransactionNextBlock(wallet *solana.Wallet, instructions []solana.Instruction, params *types.NextBlockParams, blockhash solana.Hash) (string, error) {
	log.Printf("[NextBlock] 开始发送交易: %s", time.Now().Format("2006-01-02 15:04:05.000"))

	// 获取随机 tip 账户
	tipAccount := getRandomNextblockTipAccount()

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

	// 转为 Base64
	tx64Str := tx.MustToBase64()

	// 构造请求
	payload := Payload{
		Transaction: TransactionMessage{
			Content: tx64Str,
		},
		FrontRunningProtection: true, // 启用 Anti-MEV
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", params.BundleUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", params.ApiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[NextBlock] 交易发送完成: %s", time.Now().Format("2006-01-02 15:04:05.000"))
	return tx.Signatures[0].String(), nil
}
