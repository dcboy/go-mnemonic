package main

import (
    "reflect"
    "testing"

    "github.com/go-sonr/go-bip39/wordlists"
)

// TestSameInputDeterministic 验证相同输入始终得到相同的12词助记词
func TestSameInputDeterministic(t *testing.T) {
    in := "区块链改变世界"
    m1, err := ChineseToBIP39Mnemonic(in)
    if err != nil {
        t.Fatalf("转换失败: %v", err)
    }
    m2, err := ChineseToBIP39Mnemonic(in)
    if err != nil {
        t.Fatalf("转换失败: %v", err)
    }
    if !reflect.DeepEqual(m1, m2) {
        t.Fatalf("相同输入得到不同输出\n%v\n%v", m1, m2)
    }
}

// TestFixedOutput 示例用例：验证固定中文句子对应固定24词（中文）
func TestFixedOutput(t *testing.T) {
    in := "人工智能创造未来财富"
    got, err := ChineseToBIP39Mnemonic(in)
    if err != nil {
        t.Fatalf("转换失败: %v", err)
    }
    want := []string{"床", "薯", "开", "赵", "的", "除", "头", "盐", "见", "鸿", "计", "港", "借", "随", "限", "茶", "庭", "响", "励", "恒", "姑", "函", "柴", "作"}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("输出不匹配\n got=%v\nwant=%v", got, want)
    }
}

// TestSpecificSentenceOutput 验证“区块链改变世界”对应固定24词（中文）
func TestSpecificSentenceOutput(t *testing.T) {
    in := "区块链改变世界"
    got, err := ChineseToBIP39Mnemonic(in)
    if err != nil {
        t.Fatalf("转换失败: %v", err)
    }
    want := []string{"照", "迎", "逐", "邀", "四", "座", "抑", "县", "键", "隐", "驾", "四", "湿", "雾", "妥", "忙", "威", "二", "归", "可", "贵", "氮", "亮", "术"}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("输出不匹配\n got=%v\nwant=%v", got, want)
    }
}

// TestEnglishOutput 验证英文词表输出与预期一致（区块链改变世界）
func TestEnglishOutput(t *testing.T) {
    in := "区块链改变世界"
    entropy := ChineseToEntropy(in)
    got, err := EntropyToMnemonicWithWordlist(entropy, wordlists.English)
    if err != nil {
        t.Fatalf("转换失败: %v", err)
    }
    want := []string{"diagram", "matrix", "fold", "trip", "attract", "industry", "torch", "device", "neutral", "ridge", "virus", "attract", "lizard", "success", "use", "man", "injury", "anxiety", "lamp", "afford", "happy", "rich", "impact", "cattle"}
    if !reflect.DeepEqual(got, want) {
        t.Fatalf("输出不匹配\n got=%v\nwant=%v", got, want)
    }
}
