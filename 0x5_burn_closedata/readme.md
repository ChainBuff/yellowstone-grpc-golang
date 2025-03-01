# burn and closeata
*本文仅适用于代码流程学习，不建议在生产环境使用，不建议直接用来进行代币销毁和关闭帐户操作*  
*本文仅适用于代码流程学习，不建议在生产环境使用，不建议直接用来进行代币销毁和关闭帐户操作*   
*本文仅适用于代码流程学习，不建议在生产环境使用，不建议直接用来进行代币销毁和关闭帐户操作*  
*本文仅适用于代码流程学习，不建议在生产环境使用，不建议直接用来进行代币销毁和关闭帐户操作*  
*本文仅适用于代码流程学习，不建议在生产环境使用，不建议直接用来进行代币销毁和关闭帐户操作*  

## burn 
根据 SDK 的代码，我们可以查看 burn 函数的实现 <https://github.com/gagliardetto/solana-go/tree/main/programs/token>
```
// NewBurnInstruction declares a new Burn instruction with the provided parameters and accounts.
func NewBurnInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	source ag_solanago.PublicKey,
	mint ag_solanago.PublicKey,
	owner ag_solanago.PublicKey,
	multisigSigners []ag_solanago.PublicKey,
) *Burn {
	return NewBurnInstructionBuilder().
		SetAmount(amount).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetOwnerAccount(owner, multisigSigners...)
}
```
这里需要的参数就 5 个，一个是 amount 销毁数量，第二个是 ata ，第三个是 mint 地址，第四个是从什么账户销毁，第五个是多签账户   
```
[]solana.PublicKey{} 用来传入额外的签名者（signers）的列表。在 Solana 的 Token Program 中，如果目标账户是一个多签（multisig）账户，则需要提供一个额外签名者的列表来验证操作的合法性。如果账户不是多签账户，或者当前操作不需要额外签名者，就可以传入一个空的切片（即 []solana.PublicKey{}）。这种方式确保了函数接口的一致性，同时也提供了扩展性，以便在需要时可以轻松传入额外的签名者。
```

怎么获取这些信息呢？从链上可以通过` GetTokenAccountsByOwner` 方法来获取，然后从 data 中切片取出 mint 数据 ` copy(mint[:], data[0:32]) ` , 然后通过 ` GetTokenAccountBalance ` 方法来获取余额

```
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// TokenInfo 结构体用于存储代币信息
type TokenInfo struct {
	Mint      string
	ATAs      []string
	Balances  []string
	UiAmounts []float64 // 添加UI余额字段
}

func main() {
	// 初始化RPC客户端
	rpcClient := rpc.New("https://go.getblock.io/xxx")
	private_key := ""
	wallet, err := solana.WalletFromPrivateKeyBase58(private_key)
	if err != nil {
		log.Fatalf("创建钱包失败: %v", err)
	}

	// 查询所有代币账户
	accounts, err := rpcClient.GetTokenAccountsByOwner(
		context.Background(),
		wallet.PublicKey(),
		&rpc.GetTokenAccountsConfig{
			ProgramId: &solana.TokenProgramID,
		},
		&rpc.GetTokenAccountsOpts{
			Encoding: solana.EncodingBase64,
		},
	)
	if err != nil {
		log.Fatalf("获取代币账户失败: %v", err)
	}

	// 用map来存储代币信息
	tokenInfos := make(map[string]*TokenInfo)

	// 遍历每个ATA账户
	for _, account := range accounts.Value {
		// 从账户数据中获取Mint地址
		data := account.Account.Data.GetBinary()
		mint := solana.PublicKey{} // 32字节的Mint地址
		copy(mint[:], data[0:32])  // Mint地址在数据的前32字节
		mintAddress := mint.String()

		// 获取代币余额
		balance, err := rpcClient.GetTokenAccountBalance(
			context.Background(),
			account.Pubkey,
			rpc.CommitmentFinalized,
		)
		if err != nil {
			log.Printf("获取账户 %s 余额失败: %v", account.Pubkey, err)
			continue
		}

		// 如果是新的Mint地址，创建TokenInfo
		if _, exists := tokenInfos[mintAddress]; !exists {
			tokenInfos[mintAddress] = &TokenInfo{
				Mint: mintAddress,
			}
		}

		// 添加ATA和余额信息
		tokenInfos[mintAddress].ATAs = append(tokenInfos[mintAddress].ATAs, account.Pubkey.String())
		tokenInfos[mintAddress].Balances = append(tokenInfos[mintAddress].Balances, balance.Value.Amount)
		if balance.Value.UiAmount != nil {
			tokenInfos[mintAddress].UiAmounts = append(tokenInfos[mintAddress].UiAmounts, *balance.Value.UiAmount)
		} else {
			tokenInfos[mintAddress].UiAmounts = append(tokenInfos[mintAddress].UiAmounts, 0)
		}
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Printf("钱包地址: %s\n", wallet.PublicKey().String())
	fmt.Printf("SPL代币种类总数: %d\n\n", len(tokenInfos))

	// 打印所有SPL代币信息
	i := 1
	for _, info := range tokenInfos {
		fmt.Printf("SPL代币 #%d:\n", i)
		fmt.Printf("  Mint地址: %s\n", info.Mint)
		for j := 0; j < len(info.ATAs); j++ {
			fmt.Printf("  ATA账户 #%d:\n", j+1)
			fmt.Printf("    地址: %s\n", info.ATAs[j])
			fmt.Printf("    原始余额: %s\n", info.Balances[j])
			fmt.Printf("    显示余额: %.9f\n", info.UiAmounts[j])
		}
		fmt.Println()
		i++


	}
}
```

结果

```
root@vm48863:~/# go run account/main.go 
钱包地址: 5k---
SPL代币种类总数: 10

SPL代币 #1:
  Mint地址: 9s6---
  ATA账户 #1:
    地址: Dxo---
    原始余额: 0
    显示余额: 0.000000000
```

值得一提的时候这里有 原始余额 与 显示余额 ，
显示余额一般是 `原始余额 / spl 代币的精度` ，目前程序只测试了 pumpfun 内盘的代币，精度为 6 ，所以显示余额是原始余额除以 1000000 可以得到你的原始余额。

这也是为什么程序不支持你进行直接销毁的原因，在 `main.go ` 中，你可以看到程序回销毁 显示余额 < 1 的代币，如果你的单个代币价格高，就去执行销毁会造成你本金亏损。
```
		for j := 0; j < len(info.ATAs); j++ {
			// 解析原始余额
			fmt.Printf("开始销毁第 %d 个代币\n", j)
			originalBalance, _ := strconv.ParseUint(info.Balances[j], 10, 64)
			// 检查原始余额大于0且UI余额小于1
			if originalBalance > 0 && info.UiAmounts[j] < 1 {
				fmt.Printf("发现小额代币:\n")
				fmt.Printf("  ATA地址: %s\n", info.ATAs[j])
				fmt.Printf("  Mint地址: %s\n", info.Mint)
				fmt.Printf("  原始余额: %s\n", info.Balances[j])
				fmt.Printf("  当前余额: %.9f\n", info.UiAmounts[j])

				ataAccount := solana.MustPublicKeyFromBase58(info.ATAs[j])
				burnAmount := originalBalance // 使用解析后的原始余额

				burnIx := token.NewBurnInstruction(
					burnAmount,
					ataAccount,
					solana.MustPublicKeyFromBase58(info.Mint),
					wallet.PublicKey(),
					[]solana.PublicKey{},
				).Build()

```

我们拿到了相关的参数，接下来只需要构建销毁指令，然后进行ata的关闭即可。对于关闭 ata 同样的参考 sdk 给的到代码,提供余下参数进行关闭即可:

```
				closeIx := token.NewCloseAccountInstruction(
					ataAccount,           // 要关闭的ATA
					wallet.PublicKey(),   // 租金接收者
					wallet.PublicKey(),   // 所有者
					[]solana.PublicKey{}, // 额外签名者
				).Build()
```

## 使用
代码中提供了多种模式，包括 quick 快速查询能关闭多少账号，burn 销毁代币，burnAndClose 销毁代币并关闭账号。
```
root@vm48863:~/buff/yellowstone-grpc-golang/0x5_burn_closedata# go run main.go 
Usage of /tmp/go-build2285630972/b001/exe/main:
  -cmd string
        执行命令 (query/burn/closed/all/quick)
  -httprpc string
        RPC节点地址
  -private_key string
        钱包私钥
2025/03/01 09:09:24 必须提供所有参数
exit status 1
```

```
root@vm48863:~/buff/yellowstone-grpc-golang/0x5_burn_closedata# go run main.go -cmd all  -private_key pk -httprpc https://mainnet.chainbuff.com

执行burnAndCloseAccounts功能 - 销毁代币并关闭账户
开始查询第 1 个代币
发现需要处理的账户:
  ATA地址: CsCkd --- 
  Mint地址: DTqm ---
  原始余额: 739352
  当前余额: 0.739352000

交易已发送，签名: 5WAkS4X
等待交易确认...
交易已确认
```