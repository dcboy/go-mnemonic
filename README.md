# go-mnemonic

English–Chinese BIP39 mnemonic converter written in Go.

## Overview
- Converts mnemonics between English (`en`) and Simplified Chinese (`zh`) using BIP39 official wordlists
- Preserves the underlying entropy (secret) across languages to ensure the same cryptographic key
- Validates mnemonics against the selected wordlist and shows reference seeds

## Highlights
- Language auto-detection for the input mnemonic
- Entropy equality check to prove the secret stays identical
- Reference seed outputs for both languages (PBKDF2). Note: Seeds differ across languages by design, while entropy stays the same.

## Requirements
- Go 1.18+ (tested on macOS)

## Quick Start
```bash
# Clone and run
go run . -m "<your mnemonic>" -to en|zh [-p <passphrase>]
```

## CLI Usage
- `-m` mnemonic string separated by spaces
- `-to` target language: `en` or `zh`
- `-p` optional passphrase

Examples:
```bash
# English -> Chinese
go run . -m "all hour make first leader extend hole alien behind guard gospel lava path output census museum junior mass reopen famous sing advance salt reform" -to zh

# Chinese -> English
go run . -m "而 怕 夏 客 盖 古 松 面 解 谓 鲜 唯 障 烯 共 吴 永 丁 赤 副 醒 分 猛 埔" -to en
```

## Output Explanation
- `转换后助记词` Converted mnemonic in the target language
- `源熵 / 目标熵` Hex-encoded entropy; must match for the same secret
- `源种子 / 目标种子` Reference PBKDF2 seeds; may differ across languages

## How It Works
- Builds an index map from the source wordlist and re-maps to the target wordlist by index
- Uses `github.com/go-sonr/go-bip39/wordlists` for official BIP39 wordlists
- Uses `github.com/tyler-smith/go-bip39` for validation, entropy and seed derivation

## Internal API
- `convertMnemonic(mnemonic, srcLang, dstLang string) (string, error)` Convert by word indices
- `validateMnemonic(mnemonic string, wordlist []string) error` Validate mnemonic against a wordlist
- `entropyForMnemonic(mnemonic string, wordlist []string) ([]byte, error)` Get raw entropy
- `seedForMnemonic(mnemonic string, wordlist []string, passphrase string) ([]byte, error)` Derive PBKDF2 seed
- `detectLanguage(mnemonic string) string` Heuristic language detection

## Notes
- Entropy equality guarantees the same cryptographic secret
- Seeds depend on the final mnemonic string and passphrase; they may differ between languages
- Always keep mnemonics and passphrases secure; never commit secrets

## Contributing
- Issues and PRs are welcome. Please include reproducible examples.

## License
- Choose an open-source license (e.g., MIT) and add a `LICENSE` file as needed.

