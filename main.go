package main

import (
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-sonr/go-bip39/wordlists"
	"github.com/mr-tron/base58"
	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
)

// convertMnemonic 将助记词从源语言转换到目标语言，保持同一索引序列以保证秘钥一致
// 参数:
//   - mnemonic: 原始助记词字符串（以空格分隔）
//   - srcLang: 源语言标识("en" 或 "zh")
//   - dstLang: 目标语言标识("en" 或 "zh")
//
// 返回:
//   - 转换后的助记词字符串
//   - 错误信息
func convertMnemonic(mnemonic, srcLang, dstLang string) (string, error) {
	if srcLang == dstLang {
		return mnemonic, nil
	}
	srcList, err := getWordlist(srcLang)
	if err != nil {
		return "", err
	}
	dstList, err := getWordlist(dstLang)
	if err != nil {
		return "", err
	}

	// 构建源语言的词->索引映射
	indexMap := make(map[string]int, len(srcList))
	for i, w := range srcList {
		indexMap[w] = i
	}

	words := normalizeSplit(mnemonic)
	indices := make([]int, len(words))
	for i, w := range words {
		idx, ok := indexMap[w]
		if !ok {
			return "", fmt.Errorf("词不在源语言词表: %s", w)
		}
		indices[i] = idx
	}

	out := make([]string, len(indices))
	for i, idx := range indices {
		if idx < 0 || idx >= len(dstList) {
			return "", fmt.Errorf("索引超界: %d", idx)
		}
		out[i] = dstList[idx]
	}
	return strings.Join(out, " "), nil
}

// getWordlist 根据语言标识返回对应的BIP39词表
// 支持: en(English), zh(Chinese Simplified)
func getWordlist(lang string) ([]string, error) {
	switch strings.ToLower(lang) {
	case "en", "english":
		return wordlists.English, nil
	case "zh", "chinese", "chinese-simplified", "zh-cn":
		return wordlists.ChineseSimplified, nil
	default:
		return nil, fmt.Errorf("不支持的语言: %s", lang)
	}
}

// validateMnemonic 使用指定词表校验助记词有效性（包含词表匹配与校验和）
// 参数:
//   - mnemonic: 助记词
//   - wordlist: 词表
func validateMnemonic(mnemonic string, wordlist []string) error {
	bip39.SetWordList(wordlist)
	if !bip39.IsMnemonicValid(mnemonic) {
		return errors.New("助记词无效或与词表不匹配")
	}
	return nil
}

// seedForMnemonic 使用指定词表计算种子（PBKDF2），用于跨语言一致性校验
// 参数:
//   - mnemonic: 助记词
//   - wordlist: 词表
//   - passphrase: 额外密码（可为空字符串）
//
// 返回:
//   - 64字节种子
func seedForMnemonic(mnemonic string, wordlist []string, passphrase string) ([]byte, error) {
	if err := validateMnemonic(mnemonic, wordlist); err != nil {
		return nil, err
	}
	return bip39.NewSeedWithErrorChecking(mnemonic, passphrase)
}

// detectLanguage 粗略检测助记词语言，比较词在英文与中文词表中的匹配数量
// 返回: "en" 或 "zh"
func detectLanguage(mnemonic string) string {
	words := normalizeSplit(mnemonic)
	enSet := toSet(wordlists.English)
	zhSet := toSet(wordlists.ChineseSimplified)
	enCount, zhCount := 0, 0
	for _, w := range words {
		if enSet[w] {
			enCount++
		}
		if zhSet[w] {
			zhCount++
		}
	}
	if enCount >= zhCount {
		return "en"
	}
	return "zh"
}

// normalizeSplit 规范化空格并拆分为词数组
func normalizeSplit(s string) []string {
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	if s == "" {
		return []string{}
	}
	return strings.Split(s, " ")
}

// toSet 将词表转换为集合以便快速匹配
func toSet(list []string) map[string]bool {
	m := make(map[string]bool, len(list))
	for _, w := range list {
		m[w] = true
	}
	return m
}

// main 命令行入口（核心功能）
// 功能：
// 1) 接收任意字符串输入，按 12/24 词长度生成英文 BIP39 助记词；
// 2) 由助记词派生种子并计算 BTC/ETH（secp256k1，BIP32/BIP44）地址；
// 3) 计算 Solana/Sui（Ed25519，SLIP-0010）地址；
// 4) 将助记词与四链地址打印输出。
func main() {
	var input string
	var words int
	flag.StringVar(&input, "s", "", "input string (any language)")
	flag.IntVar(&words, "words", 24, "mnemonic length: 12 or 24")
	flag.Parse()

	if input == "" || (words != 12 && words != 24) {
		fmt.Println("Usage: go run . -s \"<string>\" -words 12|24")
		return
	}

	// 步骤1：字符串 → SHA256 → 指定长度熵
	ent, err := stringToEntropy(input, words)
	if err != nil {
		panic(err)
	}
	// 步骤2：熵 → 英文 BIP39 助记词
	mnemonic, err := bip39.NewMnemonic(ent)
	if err != nil {
		panic(err)
	}
	fmt.Println("Mnemonic:", mnemonic)

	// 步骤3：助记词 → BIP39 种子（无密码）
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		panic(err)
	}

	hard := uint32(0x80000000)
	// BTC m/44'/0'/0'/0/0
	btcPath := []uint32{44 | hard, 0 | hard, 0 | hard, 0, 0}
	btcKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		panic(err)
	}
	for _, i := range btcPath {
		btcKey, err = btcKey.NewChildKey(i)
		if err != nil {
			panic(err)
		}
	}
	btcPriv := privFromBIP32Key(btcKey)
	// 计算 BTC P2PKH 地址（压缩公钥）
	btcAddr, err := btcP2PKHFromPriv(btcPriv)
	if err != nil {
		panic(err)
	}

	// ETH m/44'/60'/0'/0/0
	ethPath := []uint32{44 | hard, 60 | hard, 0 | hard, 0, 0}
	ethKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		panic(err)
	}
	for _, i := range ethPath {
		ethKey, err = ethKey.NewChildKey(i)
		if err != nil {
			panic(err)
		}
	}
	ethPriv := privFromBIP32Key(ethKey)
	// 计算 ETH 地址（keccak 公钥后 20 字节）
	ethAddr, err := ethAddressFromPriv(ethPriv)
	if err != nil {
		panic(err)
	}

	// Solana m/44'/501'/0'/0' (Ed25519, SLIP-0010)
	solPath := []uint32{44 | hard, 501 | hard, 0 | hard, 0 | hard}
	solPriv, _, err := slip10Ed25519(seed, solPath)
	if err != nil {
		panic(err)
	}
	solKey := ed25519.NewKeyFromSeed(solPriv)
	solPub := solKey.Public().(ed25519.PublicKey)
	// Solana 地址为 Ed25519 公钥 Base58
	solAddr := base58.Encode(solPub)

	// Sui m/44'/784'/0'/0'/0' (Ed25519, SLIP-0010)
	suiPath := []uint32{44 | hard, 784 | hard, 0 | hard, 0 | hard, 0 | hard}
	suiPriv, _, err := slip10Ed25519(seed, suiPath)
	if err != nil {
		panic(err)
	}
	suiKey := ed25519.NewKeyFromSeed(suiPriv)
	suiPub := suiKey.Public().(ed25519.PublicKey)
	// Sui 地址为 blake2b-256(flag||pubkey) 的十六进制表示（flag=0x00）
	sum := blake2b.Sum256(append([]byte{0x00}, suiPub...))
	suiAddr := "0x" + hex.EncodeToString(sum[:])

	fmt.Println("BTC:", btcAddr)
	fmt.Println("ETH:", ethAddr)
	fmt.Println("SOL:", solAddr)
	fmt.Println("SUI:", suiAddr)
}

// entropyForMnemonic 根据词表计算助记词对应的原始熵（128-256位）
// 返回的熵用于验证跨语言转换的“秘钥”是否一致
func entropyForMnemonic(mnemonic string, wordlist []string) ([]byte, error) {
	bip39.SetWordList(wordlist)
	return bip39.EntropyFromMnemonic(mnemonic)
}

// btcP2PKHFromPriv 由 secp256k1 私钥计算 BTC P2PKH 地址（主网，压缩公钥）
func btcP2PKHFromPriv(privKey []byte) (string, error) {
	k, err := crypto.ToECDSA(privKey)
	if err != nil {
		return "", err
	}
	// 1) 构造压缩公钥：前缀 0x02/0x03 + 32 字节 X 坐标
	x := k.X.Bytes()
	if len(x) < 32 {
		pad := make([]byte, 32-len(x))
		x = append(pad, x...)
	}
	prefix := byte(0x02)
	if k.Y.Bit(0) == 1 {
		prefix = 0x03
	}
	comp := append([]byte{prefix}, x...)
	// 2) HASH160：RIPEMD160(SHA256(pubkey))
	h1 := sha256.Sum256(comp)
	r := ripemd160.New()
	r.Write(h1[:])
	h160 := r.Sum(nil)
	// 3) 主网 version 前缀 0x00 + payload
	payload := append([]byte{0x00}, h160...)
	// 4) 双 SHA256 取前 4 字节作为校验和
	c1 := sha256.Sum256(payload)
	c2 := sha256.Sum256(c1[:])
	// 5) Base58Check 编码
	full := append(payload, c2[:4]...)
	return base58.Encode(full), nil
}

// slip10Ed25519 派生 Ed25519 私钥（SLIP-0010，全硬化）
func slip10Ed25519(seed []byte, path []uint32) (priv []byte, chain []byte, err error) {
	// 1) Master：I = HMAC-SHA512(key="ed25519 seed", data=seed)
	I := hmacSHA512([]byte("ed25519 seed"), seed)
	k := I[:32] // 左半 32 字节作为主私钥
	c := I[32:] // 右半 32 字节作为链码
	for _, i := range path {
		// 2) 全硬化派生：非硬化索引强制 +0x80000000
		if i < 0x80000000 {
			i = i + 0x80000000
		}
		// 3) 子键派生数据: 0x00 || k_parent || index_be(4 字节)
		data := make([]byte, 0, 1+32+4)
		data = append(data, 0x00)
		data = append(data, k...)
		ib := new(big.Int).SetUint64(uint64(i)).FillBytes(make([]byte, 4))
		data = append(data, ib...)
		// 4) 计算下一层: I = HMAC-SHA512(chain, data)
		I = hmacSHA512(c, data)
		k = I[:32]
		c = I[32:]
	}
	return k, c, nil
}

// hmacSHA512 简化封装
// hmacSHA512 计算 HMAC-SHA512
// 输入：密钥与数据；输出：64 字节的 MAC 值
func hmacSHA512(key, data []byte) []byte {
	h := hmac.New(sha512.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// privFromBIP32Key 提取 BIP32 私钥原始 32 字节
// 兼容不同库：有的序列化为 33 字节（首字节 0x00），此处去掉前导 0
func privFromBIP32Key(k *bip32.Key) []byte {
	if !k.IsPrivate {
		return nil
	}
	if len(k.Key) == 32 {
		return k.Key
	}
	if len(k.Key) == 33 {
		return k.Key[1:]
	}
	return nil
}

// stringToEntropy 任意字符串 → SHA256 → 128/256 位熵
// 参数 words 指定助记词长度（12 或 24），分别对应 128/256 位熵
func stringToEntropy(s string, words int) ([]byte, error) {
	sum := sha256.Sum256([]byte(s))
	switch words {
	case 12:
		out := make([]byte, 16)
		copy(out, sum[:16])
		return out, nil
	case 24:
		out := make([]byte, 32)
		copy(out, sum[:])
		return out, nil
	default:
		return nil, errors.New("words 仅支持 12 或 24")
	}
}

// ethAddressFromPriv 由 secp256k1 私钥计算 ETH 地址
// 过程：私钥→椭圆曲线公钥→Keccak-256 哈希→后 20 字节十六进制地址
func ethAddressFromPriv(privKey []byte) (string, error) {
	k, err := crypto.ToECDSA(privKey)
	if err != nil {
		return "", err
	}
	addr := crypto.PubkeyToAddress(k.PublicKey)
	return addr.Hex(), nil
}
