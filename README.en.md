# go-mnemonic

## Overview
- Chinese → BIP39 24 words (dual CN/EN output)
- Mnemonic validation (`en`/`zh` wordlists)
- Language conversion (preserve index sequence, entropy equality)
- Strict BIP39: `ENT+CS`, MSB-first, 11-bit slicing, official wordlists

## Quick Start
- Requires Go 1.20+
- Commands:
  - CN → 24 words (dual): `go run . -cn "区块链改变世界"`
  - Validate: `go run . -validate "<mnemonic>" -lang en|zh`
  - Convert: `go run . -m "<mnemonic>" -to en|zh [-p "<passphrase>"]`

## CLI Flags
- `-cn` Chinese input, outputs 24 Chinese & English words
- `-validate` validate a mnemonic (requires `-lang en|zh`)
- `-lang` wordlist for validation (`en`/`zh`)
- `-m` original mnemonic (space-separated)
- `-to` target language (`en` or `zh`)
- `-p` optional passphrase for PBKDF2 seed checks

## Examples
- `go run . -cn "区块链改变世界"`
  - CN: `照 迎 逐 邀 四 座 抑 县 键 隐 驾 四 湿 雾 妥 忙 威 二 归 可 贵 氮 亮 术`
  - EN: `diagram matrix fold trip attract industry torch device neutral ridge virus attract lizard success use man injury anxiety lamp afford happy rich impact cattle`
- Validate English:
  - `go run . -validate "diagram ... cattle" -lang en` → `Valid`
- Convert and print consistency:
  - `go run . -m "<CN/EN mnemonic>" -to en|zh [-p "<passphrase>"]`
  - Prints: source/target language, converted mnemonic, source/target entropy, equality, source/target seed (reference)

## How It Works
- CN → entropy: `SHA256(UTF-8 CN)` full 32 bytes (256 bits)
- Checksum: `CS=ENT/32=8` bits from `SHA256(entropy)`
- Bitstream: `ENT+CS=264` bits, MSB-first slicing to 24 × 11-bit indices
- Mapping: indices → official wordlists (CN/EN 2048 words)
- Pure function: deterministic, no randomness

## Compatibility
- CN wordlist: official `bip-0039/chinese_simplified.txt` (embedded)
- EN wordlist: official `bip-0039/english.txt` (embedded)
- Wallet import: select 24-words mode; CN import requires Simplified Chinese support

## Security
- Use long, highly random sentences; enable wallet passphrase
- Convert mode prints entropy & seed (hex) — treat as sensitive output
- Standard BIP39 recommends cryptographically secure entropy; this tool is a memory aid

## Tests
- Run: `go test ./...`
- Covers: stability, fixed CN/EN outputs, wordlist validation

## Code Map
- CLI: `main.go`
- CN entropy & mnemonics: `deterministic.go`
- CN wordlist: `wordlist_zh.go`
- EN wordlist: `wordlist_en.go`
- Tests: `deterministic_test.go`
