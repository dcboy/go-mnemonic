# go-mnemonic

## Overview
- Input any string (any language)
- Generate 12 or 24 English BIP39 mnemonics
- Print wallet addresses: BTC, ETH, Solana, Sui

## Quick Start
- Requires Go 1.20+
- Example:
  - `go run . -s "Hello 世界!" -words 24`

## CLI Flags
- `-s` input string
- `-words` mnemonic length: `12` or `24`

## Output
- Mnemonic: `predict grunt tissue ...` (English wordlist)
- Addresses:
  - `BTC: 1GqUf...`
  - `ETH: 0x27E9...`
  - `SOL: FjXbU...`
  - `SUI: 0xa607...`

## Details
- String → SHA256 → 128/256-bit entropy → BIP39 mnemonic (English wordlist)
- Derivation:
  - BTC: BIP32/BIP44 `m/44'/0'/0'/0/0`, P2PKH (compressed pubkey)
  - ETH: BIP32/BIP44 `m/44'/60'/0'/0/0`, Keccak address
  - Solana: SLIP-0010 (Ed25519) `m/44'/501'/0'/0'`, address = Ed25519 pubkey Base58
  - Sui: SLIP-0010 (Ed25519) `m/44'/784'/0'/0'/0'`, address = `blake2b-256(flag||pubkey)` hex (flag=0x00)

## Security
- Use long, highly random phrases; enable wallet passphrase
- This tool is for memory aid; for production, use cryptographically secure entropy
