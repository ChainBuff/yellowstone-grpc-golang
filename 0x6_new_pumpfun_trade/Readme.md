# pump_trade
## 概述
当前交易器只负责发送**内盘**交易，不论是买还是卖，不管确认，如果需要确认则需要自己进行改动,当前程序提供的多种购买方式都是基于官方文档进行编写的。对于jito 或者 nextblock 还有通过 grpc 方式提交交易。这些需要读者自行钻研。另外具体的交易组装过程在使用案例之后。如果有错误欢迎各位师傅指正。外盘的方式需要大家自行研究。


## 使用案例
配置好配置文件之后
```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x4_pumpfun_trade # go run main.go -h
Usage of /tmp/go-build1704599728/b001/exe/main:
  -amount float
        操作数量，买入的时候是solana的数量，卖出的时候是卖出的代币数量
  -mint string
        SPL Token的Mint地址
  -op string
        操作类型 (buy/sell)
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x4_pumpfun_trade # 
```

## 更新说明
参考官方文档： https://github.com/pump-fun/pump-public-docs/blob/main/docs/PUMP_CREATOR_FEE_README.md

关键描述如下：
```
buy and sell instructions will be modified in the following way:
买入和卖出指令将按以下方式修改：

the currently unused Buy::rent account (instruction account index 10) will become Buy::creator_vault account.
目前未使用的 Buy：：rent 账户（指令账户索引 10）将变为 Buy：：creator_vault 账户。
the currently unused Sell::associated_token_program account (instruction account index 8) will become Sell::creator_vault account.
当前未使用的 Sell::associated_token_program 账户 （指示账户索引 8） 将变为 Sell：：creator_vault 账户
```

所以我们需要参考 idl 中的方式计算账户,方式如下；

```

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
```





# 交易案例

sell

```
 go run main.go -mint 代币地址  -op sell -amount 35547.760004
rootDir /root/yellowstone-grpc-golang

操作信息:
Mint地址: 代币地址
操作类型: sell
数量: 35547.760004

已加载配置:
.....
2025/05/12 13:56:17 [Jito] 准备发送交易...
2025/05/12 13:56:17 [Jito] 开始发送交易: 2025-05-12 13:56:17.242
2025/05/12 13:56:17 [Jito] 交易发送完成: 2025-05-12 13:56:17.280
2025/05/12 13:56:17 交易已发送! 交易ID: 5bdcTqXk2BKJHaaop35dmGGPPJRMS
卖出发送成功! 交易ID: 5bdcTqXk2BKJHaaop35dmG
卖出数量（token x 1e6 ）: 35547760004 token
```

buy

```
root@63882:~/yellowstone-grpc-golang/0x6_new_pumpfun_trade# go run main.go -mint 代币地址  -op buy -amount 0.001
rootDir /root/yellowstone-grpc-golang

操作信息:
Mint地址: 代币地址
操作类型: buy
数量: 0.001000

已加载配置:
Private Key: 
.....
2025/05/12 13:55:16 快速模式：跳过ATA检查，直接创建ATA账户
2025/05/12 13:55:16 [Jito] 准备发送交易...
2025/05/12 13:55:16 [Jito] 开始发送交易: 2025-05-12 13:55:16.996
2025/05/12 13:55:17 [Jito] 交易发送完成: 2025-05-12 13:55:17.031
2025/05/12 13:55:17 交易已发送! 交易ID: 
买入发送成功! 交易ID: zK5CMMxJoYoYHq
预计获得: 35547760004 token (token x 1e6)
```
