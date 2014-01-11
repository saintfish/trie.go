// Package trie implements a Trie data structure.
package trie

import (
	"bytes"
)

// Value can be any type. Note that Value when added to the trie and retrieved from the trie.
type Value interface{}

// Trie is an associative array where the keys are byte arrays. See http://en.wikipedia.org/wiki/Trie for details.
type Trie struct {
	Value    *Value
	prefix   []byte
	children map[byte]*Trie
}

// NewTrie creates an empty Trie.
func NewTrie() *Trie {
	return &Trie{}
}

// Add a key value to Trie. Override the value if the same key is given again.
func (this *Trie) Add(key []byte, value Value) {
	this.createNode(key).Value = &value
}

// Get the value added together with key. If no such key was added, return nil, false.
func (this *Trie) Get(key []byte) (value Value, found bool) {
	r := this.findNode(key, exactMatch)
	if len(r) == 0 {
		return Value(nil), false
	}
	return *r[0].trie.Value, true
}

// PrefixMatch is the type of returned value of Trie's prefix matching functions.
type PrefixMatch struct {
	Prefix []byte
	Value  Value
}

// Match the shortest prefix and associated value. If no prefix is found, return {nil, nil}, false.
func (this *Trie) MatchShortestPrefix(input []byte) (match PrefixMatch, found bool) {
	r := this.findNode(input, shortestPrefix)
	if len(r) == 0 {
		return PrefixMatch{}, false
	}
	return PrefixMatch{
		Prefix: input[:r[0].prefixLength],
		Value:  *r[0].trie.Value,
	}, true
}

// Match the longest prefix and associated value. If no prefix is found, return {nil, nil}, false.
func (this *Trie) MatchLongestPrefix(input []byte) (match PrefixMatch, found bool) {
	r := this.findNode(input, longestPrefix)
	if len(r) == 0 {
		return PrefixMatch{}, false
	}
	return PrefixMatch{
		Prefix: input[:r[0].prefixLength],
		Value:  *r[0].trie.Value,
	}, true
}

// Match all possible prefixes and associated values as a list. If no prefix is found, return an empty list.
func (this *Trie) MatchAllPrefixes(input []byte) []PrefixMatch {
	r := this.findNode(input, allPrefixex)
	result := make([]PrefixMatch, len(r))
	for i, v := range r {
		result[i].Prefix = input[:v.prefixLength]
		result[i].Value = *v.trie.Value
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

func (this *Trie) findNode(key []byte, mode findNodeMode) []*findNodeResult {
	result := []*findNodeResult{}
	length := 0
	for len(key) != 0 {
		if this.Value != nil && (mode == shortestPrefix || mode == allPrefixex) {
			result = append(result, &findNodeResult{length, this})
			if mode == shortestPrefix {
				return result
			}
		}
		firstByte := key[0]
		child, has := this.children[firstByte]
		has = has && bytes.HasPrefix(key, child.prefix)
		if !has {
			if this.Value != nil && mode == longestPrefix {
				result = append(result, &findNodeResult{length, this})
			}
			return result
		}
		key = key[len(child.prefix):]
		length += len(child.prefix)
		this = child
	}
	if this.Value != nil {
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
