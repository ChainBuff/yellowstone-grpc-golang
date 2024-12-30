# 0x0_yellow-stone-grpc-example

##  yellow-stone-grpc
官方项目地址 ： <https://github.com/rpcpool/yellowstone-grpc.git> 

在官方项目中,提供过了多个语言的 example 代码, 而我们平时使用的 golang 版本也是基于官方的 example 代码进行修改得到的。所以弄懂官方的 example 代码, 对我们理解 yellow-stone-grpc 项目有很大的帮助。

## 官方的 golang 代码测试

我们可以单独拉取官方的 golang 代码,进行简单的功能测试，了解各个功能点之后我们在进行代码的修改与拆分。

下面是抽离原本源码中的参数，我们可以根据这些参数进行自己想要的功能测试
```
var (
	// 基础连接参数
	grpcAddr           = flag.String("endpoint", "", "Solana gRPC 服务器地址，使用 URI 格式，例如 https://api.rpcpool.com")
	token              = flag.String("x-token", "", "认证令牌")
	jsonInput          = flag.String("json", "", "订阅请求的 JSON，使用 @ 前缀从文件读取")
	insecureConnection = flag.Bool("insecure", false, "使用非 TLS 连接")

	// 区块链数据订阅选项
	slots      = flag.Bool("slots", false, "订阅 slot 更新")
	blocks     = flag.Bool("blocks", false, "订阅区块更新")
	block_meta = flag.Bool("blocks-meta", false, "订阅区块元数据更新")
	signature  = flag.String("signature", "", "订阅特定交易签名")
	resub      = flag.Uint("resub", 0, "在 x 次更新后重新仅订阅 slots，0 表示禁用")

	// 账户相关订阅选项
	accounts = flag.Bool("accounts", false, "订阅账户更新")

	// 交易相关订阅选项
	transactions       = flag.Bool("transactions", false, "订阅交易，用于 tx_account_include/tx_account_exclude 和 vote/failed")
	voteTransactions   = flag.Bool("transactions-vote", false, "包含投票交易")
	failedTransactions = flag.Bool("transactions-failed", false, "包含失败的交易")

	// 过滤器数组
	accountsFilter              arrayFlags // 账户过滤器
	accountOwnersFilter         arrayFlags // 账户所有者过滤器
	transactionsAccountsInclude arrayFlags // 交易包含的账户过滤器
	transactionsAccountsExclude arrayFlags // 交易排除的账户过滤器
)

	// 设置命令行参数
	flag.Var(&accountsFilter, "accounts-account", "订阅指定账户，可多次指定")
	flag.Var(&accountOwnersFilter, "accounts-owner", "订阅指定账户所有者，可多次指定")
	flag.Var(&transactionsAccountsInclude, "transactions-account-include", "订阅包含指定账户的交易，可多次指定")
	flag.Var(&transactionsAccountsExclude, "transactions-account-exclude", "订阅不包含指定账户的交易，可多次指定")
```


### 订阅某个账号的交易
进入到 golang 目录下，这里我们使用社区节点，执行如下命令，可以订阅目标账号的交易数据:

```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x0_yellow-stone-grpc-example/golang # go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -transactions -transactions-account-include 696969Y6orZEjp4gZtwcCZS7TNVuhMVE6G5mFvdJ4mYq


启动 gRPC 客户端，连接到 grpc.chainbuff.com:443
订阅请求: {"transactions":{"transactions_sub":{"vote":false,"failed":false,"account_include":["696969Y6orZEjp4gZtwcCZS7TNVuhMVE6G5mFvdJ4mYq"]}}}
1735569493008873325 filters:"transactions_sub"  transaction:{transaction:{signature:"{\x12\xc3& \xd9~\r<\xa38\xd9CdN\x1e\x9b\xf14m?\xb6Fu\x85u\xcc\x12\x11\xe2\xe7ō\xd2\x1cou\xbf\x8d\\\x11\xa5m\x1c\xdd\xfcBBS\xea?\xd8#_\xa0\xb4\xf5bc_C\xad\xf9\x05"  transaction:{signatures:"{\x12\xc3& \xd9~\r<\xa38\xd9CdN\x1e\x9b\xf14m?\xb6Fu\x85u\xcc\x12\x11\xe2\xe7ō\xd2\x1cou\xbf\x8d\\\x11\xa5m\x1c\xdd\xfcBBS\xea?\xd8#_\xa0\xb4\xf5bc_C\xad\xf9\x05"  message:{header:{num_required_signatures:1  num_readonly_unsigned_accounts:10}  account_keys:"L\\\xe4n\x87\xe3t\x15'\xb9\xa1\x0e2\xb4\xf9\xe5\x96g\xdf\xc7\xd23UJ\xc0\x01\x8f\xee\xeawƮ"  account_keys:"\x92\x8c\xf6\x1c\xf1\x8d\x15\xe5\xa5\x11:/I\xad<\x91\x17`-\x8a\x04\x03\x14\x9e\x18\n\x16eu\xc7~\t"  account_keys:"\xad\x11\xe6\xa4\xfc)D\xa4\xfa\x82Q\xbe\xf8\x15Bn\x1b\xfb(ƶdfw`|j\xd9\xf5f\xa6F"  account_keys:"\x88\xe7\xe
```


### 订阅其他内容
既然我们能订阅一个账号的交易，那么我们就可以订阅其他内容，比如区块，账户，slot 等。

从代码给出的功能点不难看出主要的订阅类型为：
1. slots - 区块链槽位更新  
2. blocks - 完整区块信息  
3. blocks-meta - 区块元数据  
4. transactions - 交易信息  
5. accounts - 账户更新  

订阅 slots 的命令如下：
```
go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -slots

root@kvm12191 ~/buff/yellowstone-grpc-golang/0x0_yellow-stone-grpc-example/golang # go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go \
  -endpoint https://grpc.chainbuff.com \
  -slots
启动 gRPC 客户端，连接到 grpc.chainbuff.com:443
订阅请求: {"slots":{"slots":{}}}
1735569879016257324 filters:"slots"  slot:{slot:310788730  parent:310788729}
1735569879038823226 filters:"slots"  slot:{slot:310788699  parent:310788698  status:FINALIZED}
1735569879279981521 filters:"slots"  slot:{slot:310788730  status:CONFIRMED}
1735569879383190946 filters:"slots"  slot:{slot:310788731  parent:310788730}
```

订阅 blocks-meta 的命令如下：
```
go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -blocks-meta

root@kvm12191 ~/buff/yellowstone-grpc-golang/0x0_yellow-stone-grpc-example/golang # go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -blocks-meta
启动 gRPC 客户端，连接到 grpc.chainbuff.com:443
订阅请求: {"blocks_meta":{"block_meta":{}}}
1735569917068401560 filters:"block_meta"  block_meta:{slot:310788823  blockhash:"N3BYKQDdezsde6xmSbDRMWTPfkYAhXKVqGKBPqpgwwL"  rewards:{rewards:{pubkey:"DRpbCBMxVnDK7maPM5tGv6MvB3v1sRMC86PZ8okm21hy"  lamports:36849615  post_balance:4029935466060  reward_type:Fee}}  block_time:{timestamp:1735569916}  block_height:{block_height:289118211}  parent_slot:310788822  parent_blockhash:"B7EvTc8DcNSSiXYq3URDUzhprDfnGZs3t4ddncLCPQZL"  executed_transaction_count:1696  entries_count:501}
```

订阅 blocks 的命令如下：
```
go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -blocks
```


### 认证令牌

一般的私人节点，需要使用认证令牌，这时候你可以加上你的认证令牌来进行认证，比如：

```
go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -blocks -x-token "xToken"
```

### 订阅多个内容

我们可以同时订阅多个内容，比如：
```
go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -blocks -blocks-meta -slots
```

或者订阅多个账号的交易

```
go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -transactions -transactions-account-include 696969Y6orZEjp4gZtwcCZS7TNVuhMVE6G5mFvdJ4mYq -transactions-account-include 696969Y6orZEjp4gZtwcCZS7TNVuhMVE6G5mFvdJ4mYq
```


### json 来订阅账号或者交易

我们可以使用 json 来订阅账号或者交易，比如：
```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x0_yellow-stone-grpc-example/golang # go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go \
  -endpoint https://grpc.chainbuff.com \
  -json '{"accounts":{"my_sub":{"account":["696969Y6orZEjp4gZtwcCZS7TNVuhMVE6G5mFvdJ4mYq"]}}}' 

启动 gRPC 客户端，连接到 grpc.chainbuff.com:443
订阅请求: {"accounts":{"my_sub":{"account":["696969Y6orZEjp4gZtwcCZS7TNVuhMVE6G5mFvdJ4mYq"]}}}
1735570543788920488 filters:"my_sub"  account:{account:{pubkey:"L\\\xe4n\x87\xe3t\x15'\xb9\xa1\x0e2\xb4\xf9\xe5\x96g\xdf\xc7\xd23UJ\xc0\x01\x8f\xee\xeawƮ"  lamports:14262359193  owner:"\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00"  rent_epoch:18446744073709551615  write_version:1567321535342  txn_signature:"\xbe7\xac_k\x8a\xacX\xbf\x1a짒b\xd4s\xeb\xd6\xe4\xf9\xf0\x05\x08SE\xe7x\xde8\x0c\xd4\xd5N\xcds\x7f\xe5\"㪚\xed\x84\x01\xa0]O\xbdN*q\xe4蒚m\xca\xc8!\xe5v/f\x0b"}  slot:310790353}
```

通过 json文件 订阅交易

```
jsonfile:

{
  "transactions": {
    "tx_sub": {
      "accountInclude": [
        "696969Y6orZEjp4gZtwcCZS7TNVuhMVE6G5mFvdJ4mYq"
      ],
      "vote": false,
      "failed": true
    }
  }
}
```

执行命令：
```
root@kvm12191 ~/buff/yellowstone-grpc-golang/0x0_yellow-stone-grpc-example/golang # go run ./cmd/grpc-client/main.go ./cmd/grpc-client/array-flag.go -endpoint https://grpc.chainbuff.com -json @./transactions.json
```


### 总结

通过以上测试，我们可以了解到 yellow-stone-grpc 的基本使用方法，以及各个功能点的使用方法。后续我们就只需要把自己需要的功能拆解出来，然后根据 proto 文件夹中的结构体，来找到对应的数据进行解析即可。