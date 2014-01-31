// Package trie implements a Trie data structure.
package trie

import (
	"bytes"
	"strings"
)

// Value can be any type. Note that Value when added to the trie and retrieved from the trie.
type Value interface{}

// Trie is an associative array where the keys are byte arrays. See http://en.wikipedia.org/wiki/Trie for details.
type Trie struct {
	value    *Value
	prefix   []byte
	children map[byte]*Trie
}

// NewTrie creates an empty Trie.
func NewTrie() *Trie {
	return &Trie{}
}

// Add a key value to Trie. Override the value if the same key is given again.
func (this *Trie) Add(key []byte, value Value) {
	this.createNode(key).value = &value
}

// Get the value associated with the key. If no such key was added, return nil, false.
func (this *Trie) GetBytes(key []byte) (value Value, found bool) {
	return this.get(&inputBytes{key})
}

// Same as GetBytes but works for string.
func (this *Trie) GetString(key string) (value Value, found bool) {
	return this.get(&inputString{key})
}

func (this *Trie) get(in input) (value Value, found bool) {
	r := this.findNode(in, exactMatch)
	if len(r) == 0 {
		return Value(nil), false
	}
	return *r[0].trie.value, true
}

// PrefixMatch is the type of returned value of Trie's prefix matching functions.
type PrefixMatch struct {
	PrefixLength int
	Value        Value
}

// Match the shortest prefix and associated value. If no prefix is found, return {nil, nil}, false.
func (this *Trie) MatchShortestPrefixBytes(input []byte) (match PrefixMatch, found bool) {
	return this.matchPrefix(&inputBytes{input}, shortestPrefix)
}

// Same as MatchShortestPrefixBytes but works for string.
func (this *Trie) MatchShortestPrefixString(input string) (match PrefixMatch, found bool) {
	return this.matchPrefix(&inputString{input}, shortestPrefix)
}

// Match the longest prefix and associated value. If no prefix is found, return {nil, nil}, false.
func (this *Trie) MatchLongestPrefixBytes(input []byte) (match PrefixMatch, found bool) {
	return this.matchPrefix(&inputBytes{input}, longestPrefix)
}

// Same as MatchLongestPrefixBytes but works for string.
func (this *Trie) MatchLongestPrefixString(input string) (match PrefixMatch, found bool) {
	return this.matchPrefix(&inputString{input}, longestPrefix)
}

func (this *Trie) matchPrefix(in input, mode findNodeMode) (match PrefixMatch, found bool) {
	r := this.findNode(in, mode)
	if len(r) == 0 {
		return PrefixMatch{}, false
	}
	return PrefixMatch{
		PrefixLength: r[0].prefixLength,
		Value:        *r[0].trie.value,
	}, true
}

// Match all possible prefixes and associated values as a list. If no prefix is found, return an empty list.
func (this *Trie) MatchAllPrefixesBytes(in []byte) []PrefixMatch {
	return this.matchAllPrefixes(&inputBytes{in})
}

// Same as MatchAllPrefixesBytes but works for string input.
func (this *Trie) MatchAllPrefixesString(in string) []PrefixMatch {
	return this.matchAllPrefixes(&inputString{in})
}

func (this *Trie) matchAllPrefixes(in input) []PrefixMatch {
	r := this.findNode(in, allPrefixex)
	result := make([]PrefixMatch, len(r))
	for i, v := range r {
		result[i].PrefixLength = v.prefixLength
		result[i].Value = *v.trie.value
	}
	return result
}

func (this *Trie) createNode(key []byte) *Trie {
	for len(key) != 0 {
		firstByte := key[0]
		child, has := this.children[firstByte]
		if !has {
			child = &Trie{prefix: make([]byte, len(key))}
			copy(child.prefix, key)
			if this.children == nil {
				this.children = make(map[byte]*Trie)
			}
			this.children[firstByte] = child
			return child
		}
		commonPrefixLen := longestCommonPrefix(child.prefix, key)
		if commonPrefixLen < len(child.prefix) {
			newChild := &Trie{
				prefix:   child.prefix[:commonPrefixLen],
				children: map[byte]*Trie{child.prefix[commonPrefixLen]: child},
			}
			child.prefix = child.prefix[commonPrefixLen:]
			this.children[firstByte] = newChild
			this = newChild
			key = key[commonPrefixLen:]
		} else {
			this = child
			key = key[commonPrefixLen:]
		}
	}
	return this
}

type findNodeMode int

const (
	exactMatch findNodeMode = iota
	shortestPrefix
	longestPrefix
	allPrefixex
)

type findNodeResult struct {
	prefixLength int
	trie         *Trie
}

type input interface {
	end() bool
	char() byte
	hasPrefix([]byte) bool
	advance(int)
}

type inputBytes struct {
	b []byte
}

func (i *inputBytes) end() bool {
	return len(i.b) == 0
}

func (i *inputBytes) char() byte {
	return i.b[0]
}

func (i *inputBytes) hasPrefix(prefix []byte) bool {
	return bytes.HasPrefix(i.b, prefix)
}

func (i *inputBytes) advance(n int) {
	i.b = i.b[n:]
}

type inputString struct {
	s string
}

func (i *inputString) end() bool {
	return len(i.s) == 0
}

func (i *inputString) char() byte {
	return i.s[0]
}

func (i *inputString) hasPrefix(prefix []byte) bool {
	return strings.HasPrefix(i.s, string(prefix))
}

func (i *inputString) advance(n int) {
	i.s = i.s[n:]
}

func (this *Trie) findNode(key input, mode findNodeMode) []*findNodeResult {
	result := []*findNodeResult{}
	length := 0
	for !key.end() {
		if this.value != nil && (mode == shortestPrefix || mode == allPrefixex) {
			result = append(result, &findNodeResult{length, this})
			if mode == shortestPrefix {
				return result
			}
		}
		firstByte := key.char()
		child, has := this.children[firstByte]
		has = has && key.hasPrefix(child.prefix)
		if !has {
			if this.value != nil && mode == longestPrefix {
				result = append(result, &findNodeResult{length, this})
			}
			return result
		}
		key.advance(len(child.prefix))
		length += len(child.prefix)
		this = child
	}
	if this.value != nil {
		result = append(result, &findNodeResult{length, this})
	}
	return result
}

func longestCommonPrefix(a, b []byte) int {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return minLen
}
