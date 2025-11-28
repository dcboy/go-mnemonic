# go-mnemonic

## 项目简介
- 通过纯函数将任意中文句子确定性地转换为 BIP39 助记词
- 支持生成 24 个中文助记词，并同步输出等价的英文助记词
- 完全遵循 BIP39 标准：`ENT + CS` 比特拼接，11 位索引切分，官方词表映射
- 使用官方 BIP39 中文简体词表（2048 词），可在支持中文词表的钱包导入

## 特性
- 纯函数：同一中文输入 → 同一 24 词助记词（中/英一致映射至同一熵）
- BIP39 合规：`CS=ENT/32`，MSB 优先比特顺序，11 位索引切分
- 双语输出：中文与英文助记词基于同一熵生成，跨语言一致可验证

## 安装与运行
- 运行环境：Go 1.20+
- 拉取代码并运行：
  - `go run . -cn "你的中文句子"`

## 使用示例
- 示例一：
  - 命令：`go run . -cn "区块链改变世界"`
  - 输出：
    - 中文: `照 迎 逐 邀 四 座 抑 县 键 隐 驾 四 湿 雾 妥 忙 威 二 归 可 贵 氮 亮 术`
    - 英文: `diagram matrix fold trip attract industry torch device neutral ridge virus attract lizard success use man injury anxiety lamp afford happy rich impact cattle`
- 示例二：
  - 命令：`go run . -cn "人工智能创造未来财富"`
  - 输出：
    - 中文: `床 薯 开 赵 的 除 头 盐 见 鸿 计 港 借 随 限 茶 庭 响 励 恒 姑 函 柴 作`
    - 英文: `lock surprise artwork jar abandon curious box floor broken stomach bronze pave layer disagree fatigue gorilla limit cross raw toilet ocean march plate addict`

## 命令行参数
- `-cn "中文句子"`：必填，输入中文字符串，输出 24 个中文与英文助记词
- 兼容原工具的转换模式：
  - `-m "<助记词>" -to en|zh [-p <密码>]` 用于中英词表索引一致性转换与验证（保留）

## 算法与确定性保证
- 中文 → 熵：`entropy = SHA256(UTF-8中文)` 的完整 32 字节（256 位）
- 校验和：`CS = ENT/32 = 8` 位，取 `SHA256(entropy)` 的前 `CS` 位
- 拼接为 `ENT+CS = 264` 位比特，按 MSB 优先从左往右切分为 24 个 11 位索引
- 索引映射至词表（中文/英文词表各 2048 词），得到等价的双语助记词
- 由于熵由中文输入唯一确定，整个过程无随机性；同输入始终生成相同输出

## BIP39 兼容性
- 中文词表来源：Bitcoin 官方 BIPs 仓库 `bip-0039/chinese_simplified.txt`
- 比特顺序与切片规则严格遵循 BIP39，生成的助记词可在支持中文词表的钱包（如 MetaMask、Trust Wallet）导入
- 注意：部分钱包默认英文词表，导入中文时需切换到中文简体词表

## 安全提醒
- 此方法的安全性完全取决于中文输入的熵，请使用长且随机性高、不易被猜测的中文句子
- 建议在钱包使用附加密码（passphrase），提升整体安全性
- 标准 BIP39 建议使用密码学安全的随机熵；本工具仅作为“记忆辅助”方案

## 运行测试
- 执行：`go test ./...`
- 已包含以下测试：
  - 同一输入稳定性（两次生成的助记词完全一致）
  - 指定句子（如上示例）固定输出验证（中文与英文）

## 代码位置
- 中文熵与助记词：`deterministic.go`
- 中文词表加载：`wordlist_zh.go`
- 命令行入口：`main.go`
- 测试用例：`deterministic_test.go`

