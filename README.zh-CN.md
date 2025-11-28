# go-mnemonic

## 功能概览

- 中文句子 → BIP39 24 词（中文/英文双语）
- 助记词有效性校验（支持 `en`/`zh` 词表）
- 助记词中英转换（保持索引序列一致，熵一致）
- 完全遵循 BIP39：`ENT+CS` 拼接、MSB 优先、11 位分片、官方词表映射

## 安装与运行

- 要求：Go 1.20+
- 常用命令：
  - 中文 → 24 词（双语）：`go run . -cn "区块链改变世界"`
  - 助记词校验：`go run . -validate "<mnemonic>" -lang en|zh`
  - 语言转换：`go run . -m "<mnemonic>" -to en|zh [-p "<密码>"]`

## CLI 参数

- `-cn` 中文输入，输出 24 个中文与英文助记词
- `-validate` 校验助记词有效性（需配合 `-lang en|zh`）
- `-lang` 指定校验词表语言（`en`/`zh`）
- `-m` 原始助记词（空格分隔）
- `-to` 目标语言（`en` 或 `zh`）
- `-p` 可选密码，用于 PBKDF2 种子派生验证

## 示例

- `go run . -cn "区块链改变世界"`
  - 中文: `照 迎 逐 邀 四 座 抑 县 键 隐 驾 四 湿 雾 妥 忙 威 二 归 可 贵 氮 亮 术`
  - 英文: `diagram matrix fold trip attract industry torch device neutral ridge virus attract lizard success use man injury anxiety lamp afford happy rich impact cattle`
- 校验英文助记词：
  - `go run . -validate "diagram ... cattle" -lang en` → `校验通过`
- 中英转换并打印一致性信息：
  - `go run . -m "<中文/英文助记词>" -to en|zh [-p "<密码>"]`
  - 输出包含：源/目标语言、转换后助记词、源/目标熵、是否一致、源/目标种子（参考）

## 原理

- 中文 → 熵：`SHA256(UTF-8中文)` 取完整 32 字节（256 位）
- 校验和：`CS=ENT/32=8` 位（取 `SHA256(entropy)` 前 `CS` 位）
- 比特流：`ENT+CS=264` 位，MSB 优先切分为 24 × 11 位索引
- 词表映射：索引映射至官方词表（中文/英文各 2048 词）
- 纯函数：无随机性；同一输入始终得到相同输出

## 兼容性

- 中文词表：官方 `bip-0039/chinese_simplified.txt`（项目内嵌）
- 英文词表：官方 `bip-0039/english.txt`（项目内嵌）
- 钱包导入：确保选择 24 词模式；英文默认可导入；中文需钱包支持中文简体词表

## 安全提醒

- 安全性取决于中文输入的熵；建议使用长且随机性高的句子
- 建议启用钱包附加密码（passphrase）
- 转换模式会打印熵与种子（十六进制），实际使用中请谨慎对待敏感输出
- 标准 BIP39 推荐使用密码学安全随机熵；本工具主要用于记忆辅助

## 测试

- 执行：`go test ./...`
- 覆盖：同输入稳定性、固定句子中/英输出、词表有效性校验

## 代码结构

- 命令行入口：`main.go`
- 中文熵与助记词：`deterministic.go`
- 中文词表：`wordlist_zh.go`
- 英文词表：`wordlist_en.go`
- 测试：`deterministic_test.go`
