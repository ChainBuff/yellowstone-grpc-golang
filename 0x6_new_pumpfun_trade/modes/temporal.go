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

// Temporal tip 账户列表
var temporalTipAccounts = []string{
	"TEMPaMeCRFAS9EKF53Jd6KpHxgL47uWLcpFArU1Fanq",
	"noz3jAjPiHuBPqiSPkkugaJDkJscPuRhYnSpbi8UvC4",
	"noz3str9KXfpKknefHji8L1mPgimezaiUyCHYMDv1GE",
	"noz6uoYCDijhu1V7cutCpwxNiSovEwLdRHPwmgCGDNo",
	"noz9EPNcT7WH6Sou3sr3GGjHQYVkN3DNirpbvDkv9YJ",
	"nozc5yT15LazbLTFVZzoNZCwjh3yUtW86LoUyqsBu4L",
	"nozFrhfnNGoyqwVuwPAW4aaGqempx4PU6g6D9CJMv7Z",
	"nozievPk7HyK1Rqy1MPJwVQ7qQg2QoJGyP71oeDwbsu",
	"noznbgwYnBLDHu8wcQVCEw6kDrXkPdKkydGJGNXGvL7",
	"nozNVWs5N8mgzuD3qigrCG2UoKxZttxzZ85pvAQVrbP",
	"nozpEGbwx4BcGp6pvEdAh1JoC2CQGZdU6HbNP1v2p6P",
	"nozrhjhkCr3zXT3BiT4WCodYCUFeQvcdUkM7MqhKqge",
	"nozrwQtWhEdrA6W8dkbt9gnUaMs52PdAv5byipnadq3",
	"nozUacTVWub3cL4mJmGCYjKZTnE9RbdY5AP46iQgbPJ",
	"nozWCyTPppJjRuw2fpzDhhWbW355fzosWSzrrMYB1Qk",
	"nozWNju6dY353eMkMqURqwQEoM3SFgEKC6psLCSfUne",
	"nozxNBgWohjR75vdspfxR5H9ceC7XXH99xpxhVGt3Bb",
}

// 获取随机的 tip 账户
func getRandomTemporalTipAccount() string {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(temporalTipAccounts))))
	if err != nil {
		return temporalTipAccounts[0]
	}
	return temporalTipAccounts[n.Int64()]
}

// SendTransactionTemporal 发送交易到 Temporal
func SendTransactionTemporal(wallet *solana.Wallet, instructions []solana.Instruction, params *types.TemporalParams, blockhash solana.Hash) (string, error) {
	log.Printf("[Temporal] 开始发送交易: %s", time.Now().Format("2006-01-02 15:04:05.000"))

	// 获取随机 tip 账户
	tipAccount := getRandomTemporalTipAccount()

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
	type JsonRpcRequest struct {
		Jsonrpc string        `json:"jsonrpc"`
		Id      int           `json:"id"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
	}

	payload := JsonRpcRequest{
		Jsonrpc: "2.0",
		Id:      1,
		Method:  "sendTransaction",
		Params: []interface{}{
			tx64Str,
			map[string]string{"encoding": "base64"},
		},
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 创建 HTTP 请求
	reqUrl := fmt.Sprintf("%s?c=%s", params.BundleUrl, params.ApiKey)
	req, err := http.NewRequest("POST", reqUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查响应
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.Printf("[Temporal] 交易发送完成: %s", time.Now().Format("2006-01-02 15:04:05.000"))
	return tx.Signatures[0].String(), nil
}
