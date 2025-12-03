# go-mnemonic

## 项目简介

- 输入任意字符串（任何语言）
- 生成 12 或 24 个英文 BIP39 助记词
- 同步计算并打印钱包地址：BTC、ETH、Solana、Sui

## 快速开始

- 环境要求：Go 1.20+
- 示例：
  - `go run . -s "Hello 世界!" -words 24`

## CLI 参数

- `-s` 输入任意字符串（任何语言）
- `-words` 助记词长度，支持 `12` 或 `24`

## 输出示例

- 助记词：`predict grunt tissue ...`（英文词表）
- 地址：
  - `BTC: 1GqUf...`
  - `ETH: 0x27E9...`
  - `SOL: FjXbU...`
  - `SUI: 0xa607...`

## 原理与兼容

- 字符串 → SHA256 → 128/256 位熵 → BIP39 助记词（英文词表）
- 地址派生：
  - BTC：BIP32/BIP44 路径 `m/44'/0'/0'/0/0`，P2PKH（压缩公钥、Base58Check）
  - ETH：BIP32/BIP44 路径 `m/44'/60'/0'/0/0`，Keccak-256 公钥后 20 字节地址
  - Solana：SLIP-0010（Ed25519）路径 `m/44'/501'/0'/0'`，地址为 Ed25519 公钥 Base58
  - Sui：SLIP-0010（Ed25519）路径 `m/44'/784'/0'/0'/0'`，地址为 `blake2b-256(flag||pubkey)` 十六进制（flag=0x00）

## 安全提醒

- 安全性取决于输入字符串的熵；建议使用长且随机性高的短语
- 建议启用钱包附加密码（passphrase）
- 本项目用于记忆辅助；生产用途建议使用密码学安全随机熵

## English README

- See [README.en.md](README.en.md)
