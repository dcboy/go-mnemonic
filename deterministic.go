package main

import (
	"crypto/sha256"
	"errors"
)

// ChineseToEntropy 中文字符串 → SHA256 → 32字节熵（256位）
// 说明：
// 1) 输入按UTF-8编码；
// 2) 对编码后的字节做SHA256；
// 3) 取完整32字节作为熵，满足24词（256位熵）需求；
// 安全提醒：此熵完全由中文输入决定，安全性取决于输入的随机性与长度。
func ChineseToEntropy(chineseStr string) []byte {
	data := []byte(chineseStr)
	sum := sha256.Sum256(data)
	entropy := make([]byte, 32)
	copy(entropy, sum[:])
	return entropy
}

// EntropyToMnemonicWithWordlist 通用BIP39熵到助记词转换（支持16/20/24/28/32字节熵）
// 步骤：
// 1) 计算校验位长度 CS=ENT/32（位）；取SHA256(entropy)的前CS位；
// 2) 拼接得到总位长 ENT+CS；按11位分割为 N 个索引；
// 3) 使用提供的词表映射索引得到助记词数组；
// 注：比特读取采用MSB优先，严格符合BIP39。
func EntropyToMnemonicWithWordlist(entropy []byte, wordlist []string) ([]string, error) {
	entBytes := len(entropy)
	switch entBytes {
	case 16, 20, 24, 28, 32:
		// ok
	default:
		return nil, errors.New("熵长度必须为16/20/24/28/32字节")
	}
	entBits := entBytes * 8
	csBits := entBits / 32
	totalBits := entBits + csBits

	h := sha256.Sum256(entropy)

	getBit := func(pos int) uint8 {
		if pos < entBits {
			b := entropy[pos/8]
			shift := 7 - (pos % 8)
			return (b >> shift) & 1
		}
		k := pos - entBits // 0..csBits-1
		// 仅需首字节即可覆盖最多8位校验
		return (h[0] >> (7 - k)) & 1
	}

	wordsCount := totalBits / 11
	out := make([]string, wordsCount)
	for i := 0; i < wordsCount; i++ {
		var idx uint16 = 0
		for j := 0; j < 11; j++ {
			idx = (idx << 1) | uint16(getBit(i*11+j))
		}
		out[i] = wordlist[idx]
	}
	return out, nil
}

// EntropyToMnemonic 中文简体词表的便捷封装（默认中文）
func EntropyToMnemonic(entropy []byte) ([]string, error) {
	return EntropyToMnemonicWithWordlist(entropy, zhCNWordlist)
}

// ChineseToBIP39Mnemonic 主函数：中文 → 助记词（中文简体，纯函数，无随机性）
// 安全提醒：标准BIP39建议使用密码学安全的随机熵。此方法仅作记忆辅助，
// 如需实际资金安全，请选择长且难以猜测的中文句子，并在钱包中使用附加密码。
func ChineseToBIP39Mnemonic(chineseStr string) ([]string, error) {
	entropy := ChineseToEntropy(chineseStr)
	return EntropyToMnemonic(entropy)
}
