# pump_trade
## 概述
当前交易器只负责发送交易，不论是买还是卖，不管确认，如果需要确认则需要自己进行改动,当前程序提供的多种购买方式都是基于官方文档进行编写的。对于jito 或者 nextblock 还有通过 grpc 方式提交交易。这些需要读者自行钻研。另外具体的交易组装过程在使用案例之后。如果有错误欢迎各位师傅指正。


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

具体的交易案例 买入
```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x4_pumpfun_trade # go run main.go -mint 34wnEj7YTYvsi2AKg8P2Xe5kWawg98D3gM59fBwGpump -op buy -amount 0.001
rootDir /root/buff/yellowstone-grpc-golang

操作信息:
Mint地址: 34wnEj7YTYvsi2AKg8P2Xe5kWawg98D3gM59fBwGpump
操作类型: buy
数量: 0.001000

已加载配置:
Private Key: Yq3tjpzzx88kGjHzKvtThmnr62fg45An63pnejVvv6rS1Y89A9VGK2M1fqgvDeVXd8m498wzSkEoQR2REWf9MzP
2025/02/14 14:24:25 开始执行买入交易...
2025/02/14 14:24:25 Mint地址: 34wnEj7YTYvsi2AKg8P2Xe5kWawg98D3gM59fBwGpump
2025/02/14 14:24:25 购买数量: 35432710885
2025/02/14 14:24:25 time: 2025-02-14 14:24:25.800
2025/02/14 14:24:25 快速模式：跳过ATA检查，直接创建ATA账户
2025/02/14 14:24:25 [Temporal] 准备发送交易...
2025/02/14 14:24:25 [Temporal] 开始发送交易: 2025-02-14 14:24:25.801
2025/02/14 14:24:25 [Temporal] 交易发送完成: 2025-02-14 14:24:25.817
2025/02/14 14:24:25 交易已发送! 交易ID: 5xSw2Fr8pxobtR2L8AXNJ4VEEry1gjCu3BnwbebrsZPEPrziSVxQG2sTjBNjGsCSCgJ3EB7snxhBCNY1YS9xhiSg
买入发送成功! 交易ID: 5xSw2Fr8pxobtR2L8AXNJ4VEEry1gjCu3BnwbebrsZPEPrziSVxQG2sTjBNjGsCSCgJ3EB7snxhBCNY1YS9xhiSg
预计获得: 35432710885 token
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x4_pumpfun_trade # 

```

卖出
```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x4_pumpfun_trade # go run main.go -mint 34wnEj7YTYvsi2AKg8P2Xe5kWawg98D3gM59fBwGpump -op sell -amount 35430
rootDir /root/buff/yellowstone-grpc-golang

操作信息:
Mint地址: 34wnEj7YTYvsi2AKg8P2Xe5kWawg98D3gM59fBwGpump
操作类型: sell
数量: 35430.000000

已加载配置:
Private Key: Yq3tjpzzx88kGjHzKvtThmnr62fg45An63pnejVvv6rS1Y89A9VGK2M1fqgvDeVXd8m498wzSkEoQR2REWf9MzP
2025/02/14 14:25:18 开始执行卖出交易...
2025/02/14 14:25:18 Mint地址: 34wnEj7YTYvsi2AKg8P2Xe5kWawg98D3gM59fBwGpump
2025/02/14 14:25:18 卖出数量: 35430000000
2025/02/14 14:25:18 最小获得: 949943 SOL
2025/02/14 14:25:18 time: 2025-02-14 14:25:18.332
2025/02/14 14:25:18 [Temporal] 准备发送交易...
2025/02/14 14:25:18 [Temporal] 开始发送交易: 2025-02-14 14:25:18.332
2025/02/14 14:25:18 [Temporal] 交易发送完成: 2025-02-14 14:25:18.354
2025/02/14 14:25:18 交易已发送! 交易ID: 5EL6gR4hXgr9uXJjzVo8tAsS8y9XvfDndKYVJSuijojKTrFkesnGAJZEctCyvMPRNiR74EZdLQE7SuLYAomXQBxn
卖出发送成功! 交易ID: 5EL6gR4hXgr9uXJjzVo8tAsS8y9XvfDndKYVJSuijojKTrFkesnGAJZEctCyvMPRNiR74EZdLQE7SuLYAomXQBxn
卖出数量: 35430000000 token
```


## 使用说明

1. 配置配置文件
请根据配置文件中的相关配置进行完整的配置，里面包括了交易的私钥、rpc url、交易模式、交易参数。

## 交易模式
当前文件一共提供了 4 中交易模式，可以同时启动. 卖出不会关闭 ata，所以需要后续自己进行关闭。
1. normal 普通模式  
2. jito 极速模式  
3. nextblock 极速模式  
4. temporal 极速模式  

## 交易参数

1. 交易数量  
2. 交易金额  
3. 交易滑点  
4. 交易手续费  


## pumpfun 交易过程中的参数获取

本文的主要的目标是实现基于 go 语言对 pump 发行的 spl 进行 buy 和 sell。
从 IDL 说起
在 Solana 中，IDL（Interface Definition Language）文件 是用来描述 Anchor 框架 中智能合约（即 Solana 程序）接口的标准化格式。它定义了程序的公共 API，包括函数（方法）、账户结构、事件和错误等信息。

IDL 文件的作用
1. 定义接口：
IDL 文件描述了程序的所有可调用方法（instruction），以及它们的参数和预期的账户结构。这使得前端应用（如 DApp）可以清晰了解如何与智能合约交互。
2. 自动化代码生成：
使用 IDL 文件，可以自动生成 TypeScript、Python 等语言的客户端代码，无需手动编写与智能合约交互的逻辑，大大减少了开发工作量。
3. 提高可读性和维护性：
对于团队开发或开源项目，IDL 文件像一个智能合约的“API 文档”，帮助新成员快速了解合约的接口规范。
4. 跨语言支持：
IDL 的标准化格式（通常是 JSON）使其可以被不同语言的 SDK 解析，方便不同技术栈的开发者使用。

```
{
    "version": "0.1.0",
    "name": "pump",
    "instructions": [
        {
            "name": "initialize",
            "docs": [
                "Creates the global state."
            ],
            "accounts": [
                {
                    "name": "global",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "user",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": []
        },
        {
            "name": "setParams",
            "docs": [
                "Sets the global state parameters."
            ],
            "accounts": [
                {
                    "name": "global",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "user",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "eventAuthority",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "program",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": [
                {
                    "name": "feeRecipient",
                    "type": "publicKey"
                },
                {
                    "name": "initialVirtualTokenReserves",
                    "type": "u64"
                },
                {
                    "name": "initialVirtualSolReserves",
                    "type": "u64"
                },
                {
                    "name": "initialRealTokenReserves",
                    "type": "u64"
                },
                {
                    "name": "tokenTotalSupply",
                    "type": "u64"
                },
                {
                    "name": "feeBasisPoints",
                    "type": "u64"
                }
            ]
        },
        {
            "name": "create",
            "docs": [
                "Creates a new coin and bonding curve."
            ],
            "accounts": [
                {
                    "name": "mint",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "mintAuthority",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "bondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedBondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "global",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "mplTokenMetadata",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "metadata",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "user",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "tokenProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "associatedTokenProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "rent",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "eventAuthority",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "program",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": [
                {
                    "name": "name",
                    "type": "string"
                },
                {
                    "name": "symbol",
                    "type": "string"
                },
                {
                    "name": "uri",
                    "type": "string"
                }
            ]
        },
        {
            "name": "buy",
            "docs": [
                "Buys tokens from a bonding curve."
            ],
            "accounts": [
                {
                    "name": "global",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "feeRecipient",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "mint",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "bondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedBondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedUser",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "user",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "tokenProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "rent",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "eventAuthority",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "program",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": [
                {
                    "name": "amount",
                    "type": "u64"
                },
                {
                    "name": "maxSolCost",
                    "type": "u64"
                }
            ]
        },
        {
            "name": "sell",
            "docs": [
                "Sells tokens into a bonding curve."
            ],
            "accounts": [
                {
                    "name": "global",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "feeRecipient",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "mint",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "bondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedBondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedUser",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "user",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "associatedTokenProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "tokenProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "eventAuthority",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "program",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": [
                {
                    "name": "amount",
                    "type": "u64"
                },
                {
                    "name": "minSolOutput",
                    "type": "u64"
                }
            ]
        },
        {
            "name": "withdraw",
            "docs": [
                "Allows the admin to withdraw liquidity for a migration once the bonding curve completes"
            ],
            "accounts": [
                {
                    "name": "global",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "mint",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "bondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedBondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedUser",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "user",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "tokenProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "rent",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "eventAuthority",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "program",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": []
        }
    ],
    "accounts": [
        {
            "name": "Global",
            "type": {
                "kind": "struct",
                "fields": [
                    {
                        "name": "initialized",
                        "type": "bool"
                    },
                    {
                        "name": "authority",
                        "type": "publicKey"
                    },
                    {
                        "name": "feeRecipient",
                        "type": "publicKey"
                    },
                    {
                        "name": "initialVirtualTokenReserves",
                        "type": "u64"
                    },
                    {
                        "name": "initialVirtualSolReserves",
                        "type": "u64"
                    },
                    {
                        "name": "initialRealTokenReserves",
                        "type": "u64"
                    },
                    {
                        "name": "tokenTotalSupply",
                        "type": "u64"
                    },
                    {
                        "name": "feeBasisPoints",
                        "type": "u64"
                    }
                ]
            }
        },
        {
            "name": "BondingCurve",
            "type": {
                "kind": "struct",
                "fields": [
                    {
                        "name": "virtualTokenReserves",
                        "type": "u64"
                    },
                    {
                        "name": "virtualSolReserves",
                        "type": "u64"
                    },
                    {
                        "name": "realTokenReserves",
                        "type": "u64"
                    },
                    {
                        "name": "realSolReserves",
                        "type": "u64"
                    },
                    {
                        "name": "tokenTotalSupply",
                        "type": "u64"
                    },
                    {
                        "name": "complete",
                        "type": "bool"
                    }
                ]
            }
        }
    ],
    "events": [
        {
            "name": "CreateEvent",
            "fields": [
                {
                    "name": "name",
                    "type": "string",
                    "index": false
                },
                {
                    "name": "symbol",
                    "type": "string",
                    "index": false
                },
                {
                    "name": "uri",
                    "type": "string",
                    "index": false
                },
                {
                    "name": "mint",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "bondingCurve",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "user",
                    "type": "publicKey",
                    "index": false
                }
            ]
        },
        {
            "name": "TradeEvent",
            "fields": [
                {
                    "name": "mint",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "solAmount",
                    "type": "u64",
                    "index": false
                },
                {
                    "name": "tokenAmount",
                    "type": "u64",
                    "index": false
                },
                {
                    "name": "isBuy",
                    "type": "bool",
                    "index": false
                },
                {
                    "name": "user",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "timestamp",
                    "type": "i64",
                    "index": false
                },
                {
                    "name": "virtualSolReserves",
                    "type": "u64",
                    "index": false
                },
                {
                    "name": "virtualTokenReserves",
                    "type": "u64",
                    "index": false
                }
            ]
        },
        {
            "name": "CompleteEvent",
            "fields": [
                {
                    "name": "user",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "mint",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "bondingCurve",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "timestamp",
                    "type": "i64",
                    "index": false
                }
            ]
        },
        {
            "name": "SetParamsEvent",
            "fields": [
                {
                    "name": "feeRecipient",
                    "type": "publicKey",
                    "index": false
                },
                {
                    "name": "initialVirtualTokenReserves",
                    "type": "u64",
                    "index": false
                },
                {
                    "name": "initialVirtualSolReserves",
                    "type": "u64",
                    "index": false
                },
                {
                    "name": "initialRealTokenReserves",
                    "type": "u64",
                    "index": false
                },
                {
                    "name": "tokenTotalSupply",
                    "type": "u64",
                    "index": false
                },
                {
                    "name": "feeBasisPoints",
                    "type": "u64",
                    "index": false
                }
            ]
        }
    ],
    "errors": [
        {
            "code": 6000,
            "name": "NotAuthorized",
            "msg": "The given account is not authorized to execute this instruction."
        },
        {
            "code": 6001,
            "name": "AlreadyInitialized",
            "msg": "The program is already initialized."
        },
        {
            "code": 6002,
            "name": "TooMuchSolRequired",
            "msg": "slippage: Too much SOL required to buy the given amount of tokens."
        },
        {
            "code": 6003,
            "name": "TooLittleSolReceived",
            "msg": "slippage: Too little SOL received to sell the given amount of tokens."
        },
        {
            "code": 6004,
            "name": "MintDoesNotMatchBondingCurve",
            "msg": "The mint does not match the bonding curve."
        },
        {
            "code": 6005,
            "name": "BondingCurveComplete",
            "msg": "The bonding curve has completed and liquidity migrated to raydium."
        },
        {
            "code": 6006,
            "name": "BondingCurveNotComplete",
            "msg": "The bonding curve has not completed."
        },
        {
            "code": 6007,
            "name": "NotInitialized",
            "msg": "The program is not initialized."
        }
    ],
    "metadata": {
        "address": "6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P"
    }
}
```

### instruction
在 instruction 中可以看到提供了 6 个 instruction，因为主要是需要和 pump 进行交易，所以重点会关注 buy 和 sell。


下方是完整的 buy instruction 的 idl 文件。里面标注了需要使用到的账号，在 accounts 列表中，同时还有需要的参数 args。
isMut 表示账户数据在当前指令执行期间是只读的，不能被更改
isSigner 表示是否需要当前账户的签名
```
            "name": "buy",
            "docs": [
                "Buys tokens from a bonding curve."
            ],
            "accounts": [
                {
                    "name": "global",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "feeRecipient",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "mint",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "bondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedBondingCurve",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "associatedUser",
                    "isMut": true,
                    "isSigner": false
                },
                {
                    "name": "user",
                    "isMut": true,
                    "isSigner": true
                },
                {
                    "name": "systemProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "tokenProgram",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "rent",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "eventAuthority",
                    "isMut": false,
                    "isSigner": false
                },
                {
                    "name": "program",
                    "isMut": false,
                    "isSigner": false
                }
            ],
            "args": [
                {
                    "name": "amount",
                    "type": "u64"
                },
                {
                    "name": "maxSolCost",
                    "type": "u64"
                }
            ]
```

需要签名的账号只有一个，那就是 user。而可以修改的参数账户有 feeRecipient、bondingCurve、associatedBondingCurve、associatedUser、user 这五个。这些账户分别是做什么的？怎么获取呢？
### 获取账户
#### 账户解析
通过观察与 pump 相关的交易，其实可以发现大部分地址都是固定的，只有5个地址需要我们自己去获取。分别是mint、bondingCurve、associatedBondingCurve、associatedUser、user
```
global：4wTV1YmiEkRvAtNtsSGPtUrqRYQMe5SKy2uB4Jjaxnjf 全局账户，这是一个固定的账号
feeRecipient：CebN5WGQ4jvEPvsVU4EoHEpgzq1VV7AbicfhtW4xC9iM 固定地址
mint：代币地址
bondingCurve：代币曲线地址
associatedBondingCurve： associated代币曲线地址
associatedUser：ata
user：用户共钥
systemProgram：11111111111111111111111111111111 系统账户
tokenProgram：	TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA tonken账户 
rent：SysvarRent111111111111111111111111111111111 固定地址
eventAuthority：Ce6TQqeHC9p8KetsN6JsjHK7UTZk7nasjjnr7XxXp9F1 固定地址
program：6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P pump 程序ID
```

#### 获取账户
mint 地址
在一笔交易中，mint 地址是我们交易的 spl 的地址。如何获取呢？通过solana.fm 观察 pump 的交易，其实不难发现在 CPI log 与 程序 log 都存在 mint 地址，只需要解析链上的数据即可。通过 grpc 对 pump 的程序ID 进行监听，即可获取到 program ID 的所有交易。
比如下面这一笔交易：https://solscan.io/tx/41rjpzBQLN4CoH2BNVzk5vSG3UV15bQQTwcVNBVKpNhdfSa4mGJ4dPdksasrY8Fi6e3kN2iBxGvn6aiMvFYKfenj


对于我们来说，有两个途径可以获取到 mint 地址，一个是Input Account这个列表，另一个则是 CPI log。
首先是 Input Account 列表，从 Grpc 获取到响应之后，取出数据进行解析。数据会存储在 resp.GetTransaction().Transaction.Transaction.Message 对象中：
```
message := resp.GetTransaction().Transaction.Transaction.Message
if message == nil || len(message.AccountKeys) == 0 {
    log.Print("invalid message data")		
    continue
}
for _, accountKey := range message.AccountKeys {
    log.Println("accountKey: ", base58.Encode(accountKey))
}
```
解析出来的数据如下所示：

根据 IDL 文件可知账号一般只有 12 个，如果多了账号，则可能存在其他的合约交互，可能是通过某些 bot 实现的交易。而我们自己进行交互的话只需要构造 12 个账号则可以进行交易。所有我们需要从 Grpc 的数据中解析出来 mint 地址。
所以我们先过滤程序日志中 Buy的交易，然后在进行解析。如何判断是 Buy 还是 Sell 呢？在 Program Logs中可以找到答案：

所以先判断是否为 Buy 然后再解析数据。
```
if resp.GetTransaction() != nil {
	log.Printf("tx_hash: %v", base58.Encode(resp.GetTransaction().Transaction.Signature))
	for _, logMessage := range resp.GetTransaction().GetTransaction().Meta.GetLogMessages() {
		if strings.Contains(logMessage, "Program log: Instruction: Buy") {
			message := resp.GetTransaction().Transaction.Transaction.Message
			if message == nil || len(message.AccountKeys) == 0 {
				log.Print("invalid message data")
				continue
			}
			for _, accountKey := range message.AccountKeys {
				log.Println("accountKey: ", base58.Encode(accountKey))
			}
		}
	}
	log.Println("================================================")
}
```

解析之后的数据为：
接下来要做的事情就是，根据 accout key 的位置来确定 mint 地址，这里我们需要注意，由于存在其他合约交互，所以我们监听所有事件的时候只是根据下标去解析这个数据可能会出错。所以最好输出一下 account key 的长度来辅助进行判断。因为大部分的 bot 需要手续费，那么会存在一次合约交互。
我们先通过 12 位长度的account keys 来解析账户，以这一笔交易为案例：31Z9QmyZBbRWJ3g8tCn54Ko3TeRgJVCTm63ddXs2rQAby4aZSEWErLwRApYEKk8xisJ1yazoCUJ7NKRXoAek9Kms

经过实践，在解析buy 和 sell的时候，按照固定的顺序去排列组合，很难解析出正确的 account key 排序。所以最佳的方案还是解析 CPIlog。
同样，我们解析到 Buy 指令之后，去读取 Program data，里面的数据一定是正确的。那么我们就可以通过 mint 地址来计算 bondingCurve、associatedBondingCurve

如下所示，我们可以正确的解析出来 Program data 也就是程序的 CPI log，从中获取 mint 地址即可。
```
if resp.GetTransaction() != nil {
	for _, logMessage := range resp.GetTransaction().GetTransaction().Meta.GetLogMessages() {
		// 添加本地时间戳，精确到毫秒
		now := time.Now().Format("2006-01-02 15:04:05.000")

		if strings.Contains(logMessage, "Program log: Instruction: Buy") {
			log.Printf("[%s] tx_hash: %v", now, base58.Encode(resp.GetTransaction().Transaction.Signature))
			if resp.GetTransaction().Transaction.Meta != nil {
				for _, logMsg := range resp.GetTransaction().Transaction.Meta.LogMessages {
					if strings.Contains(logMsg, "Program data: ") {
						data := strings.Split(logMsg, "Program data: ")[1]
						decodeddata, err := base64.StdEncoding.DecodeString(data)
						if err != nil {
							log.Printf("base64 decode error: %v", err)
							continue
						}

						// 检查数据长度是否足够
						requiredLength := 8 + 32 + 8 + 8 + 1 + 32 + 8 + 8 + 8 + 8 + 8
						if len(decodeddata) < requiredLength {
							log.Printf("data too short for trade event: got %d bytes, need %d bytes", len(decodeddata), requiredLength)
							continue
						}

						offset := 8 // 跳过头部8字节

						// 解析数据...
						var mintBytes [32]byte
						copy(mintBytes[:], decodeddata[offset:offset+32])
						MintAddress := base58.Encode(mintBytes[:])
						offset += 32

						SolAmount := binary.LittleEndian.Uint64(decodeddata[offset : offset+8])
						offset += 8
						TokenAmount := binary.LittleEndian.Uint64(decodeddata[offset : offset+8])
						offset += 8
						IsBuy := decodeddata[offset] != 0
						offset += 1

						var userBytes [32]byte
						copy(userBytes[:], decodeddata[offset:offset+32])
						UserAddress := base58.Encode(userBytes[:])
						offset += 32

						Timestamp := int64(binary.LittleEndian.Uint64(decodeddata[offset : offset+8]))
						offset += 8
						VirtualSolReserves := binary.LittleEndian.Uint64(decodeddata[offset : offset+8])
						offset += 8
						VirtualTokenReserves := binary.LittleEndian.Uint64(decodeddata[offset : offset+8])
						offset += 8
						RealSolReserves := binary.LittleEndian.Uint64(decodeddata[offset : offset+8])
						offset += 8
						RealTokenReserves := binary.LittleEndian.Uint64(decodeddata[offset : offset+8])

						log.Printf("MintAddress: %v, SolAmount: %v, TokenAmount: %v, IsBuy: %v, UserAddress: %v, Timestamp: %v, VirtualSolReserves: %v, VirtualTokenReserves: %v, RealSolReserves: %v, RealTokenReserves: %v",
							MintAddress, SolAmount, TokenAmount, IsBuy, UserAddress, Timestamp, VirtualSolReserves, VirtualTokenReserves, RealSolReserves, RealTokenReserves)
					}
				}
			}
			log.Println("================================================")
		}
	}
}
```

Bonding Curve && Associated Bonding Curve
接下来就是Bonding Curve 和 Associated Bonding Curve 的计算
```
beforePDA := time.Now()
seeds := [][]byte{
	[]byte("bonding-curve"),
	solana.MustPublicKeyFromBase58(MintAddress).Bytes(),
}
bondingCurve, _, err := solana.FindProgramAddress(seeds, solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P"))
log.Printf("[%s] bonding-curve: %v", time.Now().Format("2006-01-02 15:04:05.000"), bondingCurve)
associatedBondingCurve, _, err := solana.FindAssociatedTokenAddress(
	bondingCurve,
	solana.MustPublicKeyFromBase58(MintAddress),
)
log.Printf("[%s] associatedBondingCurve-curve: %v", time.Now().Format("2006-01-02 15:04:05.000"), associatedBondingCurve)
```

#### User 与 ATA
user 其实就是自己的公钥
ATA 就需要通过自己的公钥去创建，使用的库为 
```
import associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"


associatedTokenAccount, _, err := solana.FindAssociatedTokenAddress(
	s.wallet.PublicKey(),
	params.Mint,
)
if err != nil {
	return nil, fmt.Errorf("failed to find associated token account: %w", err)
}

var instructions []solana.Instruction

// 检查账户是否存在
account, _ := s.client.GetAccountInfo(context.Background(), associatedTokenAccount)
if account == nil || account.Value == nil {
	// 账户不存在，创建ATA账户的指令
	createATAIx := associatedtokenaccount.NewCreateInstruction(
		s.wallet.PublicKey(),
		s.wallet.PublicKey(),
		params.Mint,
	).Build()
	instructions = append(instructions, createATAIx)
}
```
这样所有的账户就获取完毕了，只需要组装交易即可。


### 交易

#### 从 IDL 到 交易
前文提到了 IDL 是 achor 生成的一个 API 接口文件。可以通过 IDL 与目标合约交互的具体细节。在 go 中，可以通过 go-anchor 来生成对应的代码文件。
首先需要获取到 IDL 文件，这个可以去 pump.fun 的前端获取。前文中有完整的文件，这里就不赘述。

#### 通过 go-anchor 生成代码
根据对应的 idl 文件，使用 anchor-go 直接生成对应的代码文件即可，文件会存储在 anchor-go 中的 generated 目录中。
```
anchor-go --src=/path/to/idl.json
```
#### 交易
在 pump 合约交互中的账号获取 一文中已经提到了交易所需的账号获取，所以我们只需要选择好代币，组装一个交易进行购买即可。
实现具体的买入交易，前文提到了需要 12 个账户，我们可以先计算账户，创建 ata，最后 build 交易。

demo 代码
```
package main

import (
	"context"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"time"
	"trade/pump"

	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/spf13/viper"
)

var (
	EVENT_AUTHORITY = solana.MustPublicKeyFromBase58("Ce6TQqeHC9p8KetsN6JsjHK7UTZk7nasjjnr7XxXp9F1")
	PUMP_PROGRAM_ID = solana.MustPublicKeyFromBase58("6EF8rrecthR5Dkzon8Nwu78hRvfCKubJ14M5uBEwF6P")
	GLOBAL_ACCOUNT  = solana.MustPublicKeyFromBase58("4wTV1YmiEkRvAtNtsSGPtUrqRYQMe5SKy2uB4Jjaxnjf")
	FEE_RECIPIENT   = solana.MustPublicKeyFromBase58("CebN5WGQ4jvEPvsVU4EoHEpgzq1VV7AbicfhtW4xC9iM")
)

// Config 配置结构体
type Config struct {
	HttpRpcUrl string `mapstructure:"httpRpcUrl"`
	PrivateKey string `mapstructure:"privateKey"`
}

// Wallet 钱包结构体
type Wallet struct {
	privateKey solana.PrivateKey
	publicKey  solana.PublicKey
}

// NewWalletFromPrivateKey 从私钥创建钱包
func NewWalletFromPrivateKey(privateKeyStr string) (*Wallet, error) {
	privateKey, err := solana.PrivateKeyFromBase58(privateKeyStr)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %w", err)
	}

	return &Wallet{
		privateKey: privateKey,
		publicKey:  privateKey.PublicKey(),
	}, nil
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 搜索配置文件的路径

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig 验证配置是否有效
func validateConfig(config *Config) error {
	if config.HttpRpcUrl == "" {
		return fmt.Errorf("httpRpcUrl 不能为空")
	}
	return nil
}

// validateSolAmount 验证SOL数量是否有效
func validateSolAmount(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("SOL数量必须大于0")
	}
	if amount > 100 { // 设置一个上限以防止意外
		return fmt.Errorf("SOL数量不能大于100")
	}
	return nil
}

// validateMintAddress 验证mint地址是否有效
func validateMintAddress(mintAddress string) error {
	_, err := solana.PublicKeyFromBase58(mintAddress)
	if err != nil {
		return fmt.Errorf("无效的mint地址: %w", err)
	}
	return nil
}

func main() {
	// 定义命令行参数
	mintAddr := flag.String("mint", "", "代币的Mint地址")
	solAmount := flag.Float64("amount", 0, "要购买的SOL数量")

	// 解析命令行参数
	flag.Parse()

	// 验证必要参数
	if *mintAddr == "" {
		log.Fatal("请提供mint地址 (-mint)")
	}
	if *solAmount == 0 {
		log.Fatal("请提供SOL数量 (-amount)")
	}

	// 验证参数有效性
	if err := validateMintAddress(*mintAddr); err != nil {
		log.Fatal(err)
	}
	if err := validateSolAmount(*solAmount); err != nil {
		log.Fatal(err)
	}

	// 加载配置
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 打印配置和参数信息
	fmt.Printf("已加载配置:\n")
	fmt.Printf("RPC URL: %s\n", config.HttpRpcUrl)
	if config.PrivateKey != "" {
		fmt.Println("私钥已配置")
	} else {
		fmt.Println("私钥未配置")
	}

	fmt.Printf("\n交易参数:\n")
	fmt.Printf("Mint地址: %s\n", *mintAddr)
	fmt.Printf("SOL数量: %f\n", *solAmount)

	// 创建钱包
	wallet, err := NewWalletFromPrivateKey(config.PrivateKey)
	if err != nil {
		log.Fatalf("创建钱包失败: %v", err)
	}

	// 打印钱包信息
	fmt.Printf("\n钱包信息:\n")
	fmt.Printf("公钥: %s\n", wallet.publicKey.String())

	// 获取bonding curve地址
	seeds := [][]byte{
		[]byte("bonding-curve"),
		solana.MustPublicKeyFromBase58(*mintAddr).Bytes(),
	}
	bondingCurve, _, err := solana.FindProgramAddress(seeds, PUMP_PROGRAM_ID)
	if err != nil {
		fmt.Printf("failed to find bonding curve: %w", err)
	}
	fmt.Printf("bonding curve地址: %s\n", bondingCurve)

	// 获取associated bonding curve地址
	associatedBondingCurve, _, err := solana.FindAssociatedTokenAddress(
		bondingCurve,
		solana.MustPublicKeyFromBase58(*mintAddr),
	)
	if err != nil {
		fmt.Printf("failed to find associated bonding curve: %s", err)
	}
	fmt.Printf("associated bonding curve地址: %s\n", associatedBondingCurve)

	// 假设我们是第一次购买，所以我们需要在购买的时候创建 ata
	associatedTokenAccount, _, err := solana.FindAssociatedTokenAddress(
		wallet.publicKey,
		solana.MustPublicKeyFromBase58(*mintAddr),
	)
	if err != nil {
		fmt.Printf("failed to find associated token account: %s", err)
	}
	// 创建指令集
	var instructions []solana.Instruction

	// 检查账户是否存在
	// 由于检查账户是否存在需要使用到 rpc 所以需要新建一个 rpc 客户端
	client := rpc.New(config.HttpRpcUrl)
	// 从 rpc 获取账户是否存在，如果不存在则创建 并将创建的指令添加到指令集中
	account, _ := client.GetAccountInfo(context.Background(), associatedTokenAccount)
	if account == nil || account.Value == nil {
		// 账户不存在，创建ATA账户的指令
		createATAIx := associatedtokenaccount.NewCreateInstruction(
			wallet.publicKey,
			wallet.publicKey,
			solana.MustPublicKeyFromBase58(*mintAddr),
		).Build()
		instructions = append(instructions, createATAIx)
	}

	// 构建指令之前我们要确定购买的参数，即我们花费多少 sol 可以购买到多少 spl token
	// 这里有一个重点 sol 的 lanports 是 1e9 而 spl token 的 lanports 是 1e6 需要注意计算单位
	// 第二个是如何具体计算单位 sol 能购买多少 spl token 在这里不解释具体的算法
	bondingaccount, err := client.GetAccountInfo(context.Background(), bondingCurve)
	if err != nil {
		fmt.Printf("failed to get bonding curve account: %s", err)
	}

	if bondingaccount == nil || bondingaccount.Value == nil {
		fmt.Println("account not found or empty")
	}

	data := bondingaccount.Value.Data.GetBinary()
	if len(data) < 8+40 {
		fmt.Println("invalid account data length")
	}

	// 跳过discriminator
	data = data[8:]

	VirtualTokenReserves := binary.LittleEndian.Uint64(data[0:8])
	VirtualSolReserves := binary.LittleEndian.Uint64(data[8:16])
	fmt.Printf("VirtualTokenReserves: %d\n", VirtualTokenReserves)
	fmt.Printf("VirtualSolReserves: %d\n", VirtualSolReserves)
	// bonding curve 曲线公式 计算
	base := float64(1073000191)
	coef := float64(32190005730)
	oneBillion := float64(1000000000)

	virtualSol := float64(VirtualSolReserves)
	virtualSol2 := virtualSol + (*solAmount * oneBillion)

	pay1 := base - (coef / (virtualSol / oneBillion))
	pay2 := base - (coef / (virtualSol2 / oneBillion))

	// 计算最后可以购买多少 spl 代币 并打印
	tokenAmount := pay2 - pay1
	tokenAmountWithDecimals := tokenAmount * 1e6
	fmt.Printf("tokenAmountWithDecimals: %f\n", tokenAmountWithDecimals)

	// 如果存在滑点的话，则需要计算滑点
	// 买入的话点是计算 最大可以花费多少 sol 则意味着 我们设置 滑点为 20% 则意味着 我们最大可以花费的 sol 为 1.2 * solAmount

	// 计算滑点
	slippage := 0.2
	maxSolAmount := (*solAmount * (1 + slippage)) * 1e9
	fmt.Printf("maxSolAmount: %f\n", maxSolAmount)

	// 构建交易指令
	instruction := pump.NewBuyInstructionBuilder().
		SetGlobalAccount(GLOBAL_ACCOUNT).
		SetFeeRecipientAccount(FEE_RECIPIENT).
		SetMintAccount(solana.MustPublicKeyFromBase58(*mintAddr)).
		SetBondingCurveAccount(bondingCurve).
		SetAssociatedBondingCurveAccount(associatedBondingCurve).
		SetAssociatedUserAccount(associatedTokenAccount).
		SetUserAccount(wallet.publicKey).
		SetSystemProgramAccount(solana.SystemProgramID).
		SetTokenProgramAccount(solana.TokenProgramID).
		SetRentAccount(solana.SysVarRentPubkey).
		SetEventAuthorityAccount(EVENT_AUTHORITY).
		SetProgramAccount(PUMP_PROGRAM_ID).
		SetAmount(uint64(tokenAmountWithDecimals)).
		SetMaxSolCost(uint64(maxSolAmount))

	ix, err := instruction.ValidateAndBuild()
	if err != nil {
		fmt.Printf("failed to build instruction: %s", err)
	}

	instructions = append(instructions, ix)
	// 获取最新的区块hash
	recent, err := client.GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		fmt.Printf("failed to get recent blockhash: %s", err)
	}

	// 构建普通交易
	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(wallet.publicKey),
	)
	if err != nil {
		fmt.Printf("failed to create transaction: %s", err)
	}

	// 签名交易
	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(wallet.publicKey) {
			pk := wallet.privateKey
			return &pk
		}
		return nil
	})

	if err != nil {
		fmt.Printf("failed to sign transaction: %s", err)
	}

	txSig, err := client.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Printf("failed to send transaction: %s", err)
	}

	fmt.Printf("txSig: %s\n", txSig)

	// 最后可以通过 http rpc 来轮询确认链上的交易状态
	for i := 0; i < 60; i++ {
		time.Sleep(time.Second * 1)
		status, err := client.GetTransaction(
			context.Background(),
			txSig,
			&rpc.GetTransactionOpts{
				Commitment: rpc.CommitmentFinalized,
			},
		)
		if err != nil {
			continue
		}
		if status != nil {
			fmt.Printf("\n交易确认结果:\n")
			fmt.Printf("签名: %s\n", txSig)
			if status.Meta.Err == nil {
				fmt.Println("状态: 成功")
			} else {
				fmt.Printf("状态: 失败\n")
				fmt.Printf("错误信息: %v\n", status.Meta.Err)
			}
			return
		}
	}
	fmt.Println("交易确认超时")
}
```

测试结果
```
root@kvm12191 ~/trade # go run main.go -mint AMinB8yJJRCV7w8GZsqWAM3r1M8xPtn6RxHotCZHpump -amount 0.001
已加载配置:
RPC URL: https://mainnet.chainbuff.com
私钥已配置

交易参数:
Mint地址: AMinB8yJJRCV7w8GZsqWAM3r1M8xPtn6RxHotCZHpump
SOL数量: 0.001000

钱包信息:
公钥: 5zwZscF4YfvvfKyu5ijmUv3vRwLei3yT5ZShLMYXw18B
bonding curve地址: Gje1fcSnA8FNyFNHqbQoVDHUdr7UmMwxv7i75iJ5S762
associated bonding curve地址: BZmCnYZvnnr6a2c23HX8Z5cJBPU6oBkd9uThy2QFXMfj
VirtualTokenReserves: 848844566834666
VirtualSolReserves: 37922139437
tokenAmountWithDecimals: 22383290971.517563
maxSolAmount: 1200000.000000
txSig: 2H3y7uSL3pej4DBSzvXJjij9Q927pTR7vy8eBoDha8y2H2zDHU63EggH7PR6cSvd6183ammJmKbEbDfyyWRXRdPX

交易确认结果:
签名: 2H3y7uSL3pej4DBSzvXJjij9Q927pTR7vy8eBoDha8y2H2zDHU63EggH7PR6cSvd6183ammJmKbEbDfyyWRXRdPX
状态: 成功
```