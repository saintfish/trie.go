package trie

import (
	"bytes"
	"testing"
)

var keys = []string{
	"abcdefg",
	"abcdefghijk",
	"abcd",
	"abcdxyz",
	"abXdxyz",
	"abcdefgXXX",
}

var nonKeys = []string{
	"",
	"b",
	"ab",
	"abc",
	"abcdefghijkl",
}
var content = "abcdefghijklm"
var noPrefixContent = "abcXX"
var prefixes = []string{
	"abcd",
	"abcdefg",
	"abcdefghijk",
}

func createTestTrie() *Trie {
	t := NewTrie()
	for _, k := range keys {
		t.Add([]byte(k), k)
	}
	return t
}

func TestTrieGet(t *testing.T) {
	trie := createTestTrie()
	for _, k := range keys {
		v, ok := trie.Get([]byte(k))
		if ok != true {
			t.Errorf("Unable to find key %s", k)
		}
		if v.(string) != k {
			t.Errorf("Wrong value %v, expected %v", v, k)
		}
	}
	for _, k := range nonKeys {
		v, ok := trie.Get([]byte(k))
		if ok == true {
			t.Errorf("Unexpected key %s, value %v", k, v)
		}
		if v != nil {
			t.Errorf("Value should be nil if not found, but %v", v)
		}
	}
}

func TestTrieMatchAllPrefixes(t *testing.T) {
	trie := createTestTrie()
	r := trie.MatchAllPrefixes([]byte(content))
	if len(r) != len(prefixes) {
		t.Fatalf("Wrong length of prefixes %v vs. %v)", r, prefixes)
	}
	for i, p := range prefixes {
		if !bytes.Equal(r[i].Prefix, []byte(p)) {
			t.Errorf("Wrong prefix[%d] %s vs. %s", i, r[i].Prefix, p)
		}
		if r[i].Value.(string) != p {
			t.Errorf("Wrong prefix[%d] %s vs. %s", i, r[i].Value.(string), p)
		}
	}
}

func TestTrieMatchShortestPrefix(t *testing.T) {
	trie := createTestTrie()
	v, ok := trie.MatchShortestPrefix([]byte(content))
	expected := prefixes[0]
	if !ok {
		t.Errorf("Should find shortest prefix")
	}
	if !bytes.Equal(v.Prefix, []byte(expected)) {
		t.Errorf("Wrong prefix %s vs. %s", v.Prefix, expected)
	}
	if v.Value.(string) != expected {
		t.Errorf("Wrong prefix %s vs. %s", v.Value.(string), expected)
	}
}

func TestTrieMatchLongestPrefix(t *testing.T) {
	trie := createTestTrie()
	v, ok := trie.MatchLongestPrefix([]byte(content))
	expected := prefixes[len(prefixes)-1]
	if !ok {
		t.Errorf("Should find longest prefix")
	}
	if !bytes.Equal(v.Prefix, []byte(expected)) {
		t.Errorf("Wrong prefix %s vs. %s", v.Prefix, expected)
	}
	if v.Value.(string) != expected {
		t.Errorf("Wrong prefix %s vs. %s", v.Value.(string), expected)
	}
}

func TestNoPrefixContent(t *testing.T) {
	trie := createTestTrie()
	r := trie.MatchAllPrefixes([]byte(noPrefixContent))
	if len(r) != 0 {
		t.Errorf("Unexpected result: %v", r)
	}
	v, ok := trie.MatchShortestPrefix([]byte(noPrefixContent))
	if ok {
		t.Errorf("Unexpected shortest prefix found %v", v)
	}
	v, ok = trie.MatchLongestPrefix([]byte(noPrefixContent))
	if ok {
		t.Errorf("Unexpected longest prefix found %v", v)
	}
}
