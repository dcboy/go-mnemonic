# go-mnemonic

一个使用 Go 编写的英文与中文（简体）BIP39 助记词互转工具。

## 项目简介

- 在英文（`en`）与中文简体（`zh`）词表间进行助记词转换
- 依据 BIP39 官方词表进行索引映射，确保底层熵（秘钥）保持一致
- 在转换前后进行有效性校验，并输出参考种子

## 功能特性

- 自动检测输入助记词的语言
- 验证转换前后熵是否一致（秘钥不变）
- 输出参考 PBKDF2 种子（说明：不同语言的助记词字符串本身不同，种子可能不同，这是标准行为）

## 环境要求

- Go 1.18+（在 macOS 上测试）

## 快速开始

```bash
# 直接运行
go run . -m "<你的助记词>" -to en|zh [-p <密码>]
```

## 命令行使用

- `-m` 原始助记词（以空格分隔）
- `-to` 目标语言：`en` 或 `zh`
- `-p` 可选的附加密码

示例：

```bash
# 英文 -> 中文简体
go run . -m "all hour make first leader extend hole alien behind guard gospel lava path output census museum junior mass reopen famous sing advance salt reform" -to zh

# 中文简体 -> 英文
go run . -m "而 怕 夏 客 盖 古 松 面 解 谓 鲜 唯 障 烯 共 吴 永 丁 赤 副 醒 分 猛 埔" -to en
```

## 输出说明

- `转换后助记词` 目标语言的助记词
- `源熵 / 目标熵` 十六进制的底层熵；相等即表示秘钥一致
- `源种子 / 目标种子` 参考 PBKDF2 种子；跨语言可能不同

## 原理说明

- 使用源语言词表构建“词 → 索引”映射；按相同索引在目标词表取词
- 词表来源：`github.com/go-sonr/go-bip39/wordlists`
- 验证与推导：`github.com/tyler-smith/go-bip39`

## 内部函数

- `convertMnemonic(mnemonic, srcLang, dstLang string) (string, error)` 按词表索引进行转换
- `validateMnemonic(mnemonic string, wordlist []string) error` 校验助记词在词表下的有效性
- `entropyForMnemonic(mnemonic string, wordlist []string) ([]byte, error)` 计算底层熵
- `seedForMnemonic(mnemonic string, wordlist []string, passphrase string) ([]byte, error)` 计算 PBKDF2 种子
- `detectLanguage(mnemonic string) string` 粗略语言检测

## 注意事项

- 熵一致即秘钥保持一致；这是判定转换正确的核心标准
- 种子受“助记词字符串+密码”影响，跨语言不必相同
- 请妥善保管助记词与密码，不要在公开仓库中泄露任何秘钥信息

## 贡献

- 欢迎提交 Issue 与 PR，请附上可复现的示例说明

## 许可证

- 如需开源发布，请在仓库中添加 `LICENSE` 文件（例如 MIT）。
