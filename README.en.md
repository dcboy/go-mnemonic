# go-mnemonic

## Overview
- Deterministically converts any Chinese sentence to BIP39 mnemonic words
- Generates 24 Chinese words and prints the equivalent English words from the same entropy
- Fully compliant with BIP39: `ENT + CS` bit concat, 11-bit slicing, official wordlist mapping
- Uses official Simplified Chinese BIP39 wordlist (2048 words), importable into wallets that support CN wordlist

## Features
- Pure function: same input → same 24-word mnemonic (both CN/EN from identical entropy)
- BIP39 compliance: `CS=ENT/32`, MSB-first bit order, 11-bit indices
- Dual output: Chinese and English mnemonics derived from the same entropy

## Install & Run
- Requires Go 1.20+
- Run the tool:
  - `go run . -cn "your Chinese sentence"`

## Examples
- Example 1:
  - Command: `go run . -cn "区块链改变世界"`
  - Output:
    - CN: `照 迎 逐 邀 四 座 抑 县 键 隐 驾 四 湿 雾 妥 忙 威 二 归 可 贵 氮 亮 术`
    - EN: `diagram matrix fold trip attract industry torch device neutral ridge virus attract lizard success use man injury anxiety lamp afford happy rich impact cattle`
- Example 2:
  - Command: `go run . -cn "人工智能创造未来财富"`
  - Output:
    - CN: `床 薯 开 赵 的 除 头 盐 见 鸿 计 港 借 随 限 茶 庭 响 励 恒 姑 函 柴 作`
    - EN: `lock surprise artwork jar abandon curious box floor broken stomach bronze pave layer disagree fatigue gorilla limit cross raw toilet ocean march plate addict`

## CLI Flags
- `-cn "Chinese sentence"`: required, outputs 24 Chinese and English words
- Legacy conversion mode (kept for compatibility):
  - `-m "<mnemonic>" -to en|zh [-p <passphrase>]` for language conversion & entropy/seed checks

## Algorithm & Determinism
- CN → Entropy: `entropy = SHA256(UTF-8 Chinese)` (full 32 bytes = 256 bits)
- Checksum: `CS = ENT/32 = 8` bits; take first `CS` bits of `SHA256(entropy)`
- Concat to `ENT+CS = 264` bits; slice MSB-first into 24 groups of 11 bits
- Map indices into wordlists (CN/EN 2048 words) to get equivalent dual-language mnemonics
- Pure function: no randomness; same input always produces the same output

## BIP39 Compatibility
- Wordlist: official BIPs repository `bip-0039/chinese_simplified.txt`
- Strict bit semantics & slicing as per BIP39; resulting mnemonics can be imported into wallets that support CN wordlist (e.g., MetaMask, Trust Wallet)
- Note: some wallets default to English wordlist; ensure you switch to Simplified Chinese when importing CN mnemonics

## Security Notes
- Security depends entirely on the entropy of your Chinese input; use a long, highly random sentence
- Consider using a wallet passphrase to strengthen overall security
- Standard BIP39 recommends cryptographically secure random entropy; this tool is intended as a memory aid only

## Tests
- Run: `go test ./...`
- Includes:
  - Same-input stability (two runs produce identical outputs)
  - Fixed output cases for the above examples (both CN & EN)

## Code Map
- CN entropy & mnemonics: `deterministic.go`
- CN wordlist loading: `wordlist_zh.go`
- CLI entry: `main.go`
- Tests: `deterministic_test.go`

