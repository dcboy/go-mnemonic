package main

import (
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/go-sonr/go-bip39/wordlists"
	bip39 "github.com/tyler-smith/go-bip39"
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

// main 提供命令行接口：
//
//	-m "助记词"      必填，原始助记词
//	-to en|zh       必填，目标语言
//	-p "pass"       可选，额外密码（默认空）
//
// 自动检测源语言；执行转换并校验种子一致性
func main() {
	var mnemonic string
	var to string
	var passphrase string

	flag.StringVar(&mnemonic, "m", "", "原始助记词，使用空格分隔")
	flag.StringVar(&to, "to", "", "目标语言: en 或 zh")
	flag.StringVar(&passphrase, "p", "", "可选的额外密码")
	flag.Parse()

	if mnemonic == "" || to == "" {
		fmt.Println("用法: go run . -m \"<助记词>\" -to en|zh [-p <密码>]")
		return
	}

	src := detectLanguage(mnemonic)
	dst := strings.ToLower(to)

	srcList, err := getWordlist(src)
	if err != nil {
		fmt.Println("源语言错误:", err)
		return
	}
	dstList, err := getWordlist(dst)
	if err != nil {
		fmt.Println("目标语言错误:", err)
		return
	}

	// 源助记词校验
	if err := validateMnemonic(mnemonic, srcList); err != nil {
		fmt.Println("源助记词校验失败:", err)
		return
	}

	// 执行转换
	converted, err := convertMnemonic(mnemonic, src, dst)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}

	// 计算两种语言的熵并比对（确保秘钥一致性）
	entropySrc, err := entropyForMnemonic(mnemonic, srcList)
	if err != nil {
		fmt.Println("源熵计算失败:", err)
		return
	}
	entropyDst, err := entropyForMnemonic(converted, dstList)
	if err != nil {
		fmt.Println("目标熵计算失败:", err)
		return
	}

	sameEntropy := hex.EncodeToString(entropySrc) == hex.EncodeToString(entropyDst)

	// 同时计算种子（用于参考，注意不同语言的助记词会导致种子不同）
	seedSrc, _ := seedForMnemonic(mnemonic, srcList, passphrase)
	seedDst, _ := seedForMnemonic(converted, dstList, passphrase)

	fmt.Println("源语言:", src)
	fmt.Println("目标语言:", dst)
	fmt.Println("转换后助记词:", converted)
	fmt.Println("源熵:", hex.EncodeToString(entropySrc))
	fmt.Println("目标熵:", hex.EncodeToString(entropyDst))
	fmt.Println("熵是否一致(秘钥一致):", sameEntropy)
	fmt.Println("源种子(参考):", hex.EncodeToString(seedSrc))
	fmt.Println("目标种子(参考):", hex.EncodeToString(seedDst))
}

// entropyForMnemonic 根据词表计算助记词对应的原始熵（128-256位）
// 返回的熵用于验证跨语言转换的“秘钥”是否一致
func entropyForMnemonic(mnemonic string, wordlist []string) ([]byte, error) {
	bip39.SetWordList(wordlist)
	return bip39.EntropyFromMnemonic(mnemonic)
}
