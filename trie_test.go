package trie

import (
	"testing"
)

var keys = []string{
	"abcdefg",
	"abcdefghi",
	"abcdefghijk",
	"abcdefgk",
	"abcdf",
	"abcdxyz",
	"abXdxyz",
	"abcdefgXXX",
}

var nonKeys = []string{
	"",
	"b",
	"ab",
	"abc",
	"abcd",
	"abcdefghijkl",
}
var content = "abcdefghijklm"
var noPrefixContent = "abcdefXX"
var prefixes = []string{
	"abcdefg",
	"abcdefghi",
	"abcdefghijk",
}

func createTestTrie() *Trie {
	t := NewTrie()
	for _, k := range keys {
		t.Add([]byte(k), k)
	}
	return t
}

func TestTrieGetBytes(t *testing.T) {
	trie := createTestTrie()
	for _, k := range keys {
		v, ok := trie.GetBytes([]byte(k))
		if ok != true {
			t.Errorf("Unable to find key %s", k)
		}
		if v.(string) != k {
			t.Errorf("Wrong value %v, expected %v", v, k)
		}
	}
	for _, k := range nonKeys {
		v, ok := trie.GetBytes([]byte(k))
		if ok == true {
			t.Errorf("Unexpected key %s, value %v", k, v)
		}
		if v != nil {
			t.Errorf("Value should be nil if not found, but %v", v)
		}
	}
}

func TestTrieGetString(t *testing.T) {
	trie := createTestTrie()
	for _, k := range keys {
		v, ok := trie.GetString(k)
		if ok != true {
			t.Errorf("Unable to find key %s", k)
		}
		if v.(string) != k {
			t.Errorf("Wrong value %v, expected %v", v, k)
		}
	}
	for _, k := range nonKeys {
		v, ok := trie.GetString(k)
		if ok == true {
			t.Errorf("Unexpected key %s, value %v", k, v)
		}
		if v != nil {
			t.Errorf("Value should be nil if not found, but %v", v)
		}
	}
}

func TestTrieMatchAllPrefixesBytes(t *testing.T) {
	trie := createTestTrie()
	r := trie.MatchAllPrefixesBytes([]byte(content))
	if len(r) != len(prefixes) {
		t.Fatalf("Wrong length of prefixes %v vs. %v)", r, prefixes)
	}
	for i, p := range prefixes {
		prefix := content[:r[i].PrefixLength]
		if prefix != p {
			t.Errorf("Wrong prefix[%d] %s vs. %s", i, prefix, p)
		}
		if r[i].Value.(string) != p {
			t.Errorf("Wrong prefix[%d] %s vs. %s", i, r[i].Value.(string), p)
		}
	}
}

func TestTrieMatchAllPrefixesString(t *testing.T) {
	trie := createTestTrie()
	r := trie.MatchAllPrefixesString(content)
	if len(r) != len(prefixes) {
		t.Fatalf("Wrong length of prefixes %v vs. %v)", r, prefixes)
	}
	for i, p := range prefixes {
		prefix := content[:r[i].PrefixLength]
		if prefix != p {
			t.Errorf("Wrong prefix[%d] %s vs. %s", i, prefix, p)
		}
		if r[i].Value.(string) != p {
			t.Errorf("Wrong prefix[%d] %s vs. %s", i, r[i].Value.(string), p)
		}
	}
}

func TestTrieMatchShortestPrefixBytes(t *testing.T) {
	trie := createTestTrie()
	v, ok := trie.MatchShortestPrefixBytes([]byte(content))
	expected := prefixes[0]
	if !ok {
		t.Errorf("Should find shortest prefix")
	}
	prefix := content[:v.PrefixLength]
	if prefix != expected {
		t.Errorf("Wrong prefix %s vs. %s", prefix, expected)
	}
	if v.Value.(string) != expected {
		t.Errorf("Wrong prefix %s vs. %s", v.Value.(string), expected)
	}
}

func TestTrieMatchShortestPrefixString(t *testing.T) {
	trie := createTestTrie()
	v, ok := trie.MatchShortestPrefixString(content)
	expected := prefixes[0]
	if !ok {
		t.Errorf("Should find shortest prefix")
	}
	prefix := content[:v.PrefixLength]
	if prefix != expected {
		t.Errorf("Wrong prefix %s vs. %s", prefix, expected)
	}
	if v.Value.(string) != expected {
		t.Errorf("Wrong prefix %s vs. %s", v.Value.(string), expected)
	}
}

func TestTrieMatchLongestPrefixBytes(t *testing.T) {
	trie := createTestTrie()
	v, ok := trie.MatchLongestPrefixBytes([]byte(content))
	expected := prefixes[len(prefixes)-1]
	if !ok {
		t.Errorf("Should find longest prefix")
	}
	prefix := content[:v.PrefixLength]
	if prefix != expected {
		t.Errorf("Wrong prefix %s vs. %s", prefix, expected)
	}
	if v.Value.(string) != expected {
		t.Errorf("Wrong prefix %s vs. %s", v.Value.(string), expected)
	}
}

func TestTrieMatchLongestPrefixString(t *testing.T) {
	trie := createTestTrie()
	v, ok := trie.MatchLongestPrefixString(content)
	expected := prefixes[len(prefixes)-1]
	if !ok {
		t.Errorf("Should find longest prefix")
	}
	prefix := content[:v.PrefixLength]
	if prefix != expected {
		t.Errorf("Wrong prefix %s vs. %s", prefix, expected)
	}
	if v.Value.(string) != expected {
		t.Errorf("Wrong prefix %s vs. %s", v.Value.(string), expected)
	}
}

func TestNoPrefixContentBytes(t *testing.T) {
	trie := createTestTrie()
	r := trie.MatchAllPrefixesBytes([]byte(noPrefixContent))
	if len(r) != 0 {
		t.Errorf("Unexpected result: %v", r)
	}
	v, ok := trie.MatchShortestPrefixBytes([]byte(noPrefixContent))
	if ok {
		t.Errorf("Unexpected shortest prefix found %v", v)
	}
	v, ok = trie.MatchLongestPrefixBytes([]byte(noPrefixContent))
	if ok {
		t.Errorf("Unexpected longest prefix found %v", v)
	}
}

func TestNoPrefixContentString(t *testing.T) {
	trie := createTestTrie()
	r := trie.MatchAllPrefixesString(noPrefixContent)
	if len(r) != 0 {
		t.Errorf("Unexpected result: %v", r)
	}
	v, ok := trie.MatchShortestPrefixString(noPrefixContent)
	if ok {
		t.Errorf("Unexpected shortest prefix found %v", v)
	}
	v, ok = trie.MatchLongestPrefixString(noPrefixContent)
	if ok {
		t.Errorf("Unexpected longest prefix found %v", v)
	}
}
