package main

import (
    "testing"
    "github.com/go-sonr/go-bip39/wordlists"
)

func TestEnglishMembership(t *testing.T) {
    words := []string{"diagram","matrix","fold","trip","attract","industry","torch","device","neutral","ridge","virus","attract","lizard","success","use","man","injury","anxiety","lamp","afford","happy","rich","impact","cattle"}
    set := make(map[string]bool, len(wordlists.English))
    for _, w := range wordlists.English { set[w] = true }
    for _, w := range words {
        if !set[w] {
            t.Fatalf("word not in English wordlist: %s", w)
        }
    }
}

