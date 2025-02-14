package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// ModeConfig 表示每个模式的配置
type ModeConfig struct {
	BuyTip  float64  `yaml:"buy_tip"`
	SellTip float64  `yaml:"sell_tip"`
	Enable  bool     `yaml:"enable"`
	RpcUrl  string   `yaml:"rpc_url,omitempty"`  // 可选的RPC URL
	ApiKeys []string `yaml:"api_keys,omitempty"` // NextBlock的API密钥列表
}

// Config 表示整个配置文件结构
type Config struct {
	PrivateKey      string     `yaml:"private_key"`
	HttpRpcUrl      string     `yaml:"http_rpc_url"`
	SkipATACheck    bool       `yaml:"skip_ata_check"`
	BuySlippage     float64    `yaml:"buy_slippage"`      // 全局买入滑点
	SellSlippage    float64    `yaml:"sell_slippage"`     // 全局卖出滑点
	BuyPriorityFee  float64    `yaml:"buy_priority_fee"`  // 全局买入优先级费用
	SellPriorityFee float64    `yaml:"sell_priority_fee"` // 全局卖出优先级费用
	Normal          ModeConfig `yaml:"normal"`
	Jito            ModeConfig `yaml:"jito"`
	NextBlock       ModeConfig `yaml:"nextblock"`
	Temporal        ModeConfig `yaml:"temporal"`
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	// 读取文件
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return config, nil
}

// validateConfig 验证配置是否有效
func validateConfig(config *Config) error {
	// 验证每个模式的配置
	modes := map[string]ModeConfig{
		"normal":    config.Normal,
		"jito":      config.Jito,
		"nextblock": config.NextBlock,
		"temportal": config.Temporal,
	}

	for name, mode := range modes {
		if mode.BuyTip < 0 || mode.SellTip < 0 {
			return fmt.Errorf("%s 模式的 tip 值不能为负", name)
		}
	}

	return nil
}
