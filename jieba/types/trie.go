package types

import (
	"fmt"
)

type Rune []rune

const (
	maxWordLength = 512
)

// 有向无环图节点
type DictUnit struct {
	word   Rune    // 汉字
	weight float64 // 权重
	tag    Rune    // 标记
}

// 有向无环图
type Dag struct {
	r rune // 汉字

	// 下一个汉字，可能有多个
	nexts []struct {
		offset uint
		info   *DictUnit
	}

	info    *DictUnit
	weight  float64
	nextPos uint
}

func (d *Dag) appendNext(o uint, info *DictUnit) {
	d.nexts = append(d.nexts, struct {
		offset uint
		info   *DictUnit
	}{offset: o, info: info})
}

func NewDag() *Dag {
	return &Dag{
		r: 0,
		nexts: make([]struct {
			offset uint
			info   *DictUnit
		}, 0),
		info:    nil,
		weight:  0,
		nextPos: 0,
	}
}

type TrieNode struct {
	next map[rune]*TrieNode // 下一个树上节点
	val  *DictUnit          // 当前节点权值
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		next: nil,
		val:  nil,
	}
}

type Trie struct {
	root *TrieNode // 根节点
}

func NewTrie(keys []Rune, vals []*DictUnit) *Trie {
	root := NewTrieNode()
	trie := &Trie{root: root}
	trie.createTrie(keys, vals)
	return trie
}

func (t *Trie) findWord(word Rune) *DictUnit {
	if len(word) == 0 {
		panic(word)
	}

	p := t.root
	for _, v := range word {
		if x, ok := p.next[v]; ok {
			p = x
		} else {
			return nil
		}
	}

	return p.val
}

func (t *Trie) findAsDag(dst Rune) []Dag {
	return t.findAsDagWithMaxWordLength(dst, maxWordLength)
}

func (t *Trie) findAsDagWithMaxWordLength(dst Rune, length int) []Dag {
	if t.root == nil {
		panic(t.root)
	}

	check := func(r rune) bool {
		if t.root.next == nil {
			return false
		}
		_, ok := t.root.next[r]
		return ok
	}

	dags := make([]Dag, len(dst))
	var p *TrieNode
	for i, v := range dst {
		dags[i].r = v
		if check(v) {
			p = t.root.next[v]
		} else {
			p = nil
		}

		if p == nil {
			dags[i].appendNext(uint(i), nil)
		} else {
			dags[i].appendNext(uint(i), p.val)
		}

		for j := i + 1; j < len(dst) && (j-i+1) <= length; j++ {
			if p == nil || p.next == nil {
				break
			}
			if x, ok := p.next[dst[j]]; ok {
				p = x
				if p.val != nil {
					dags[i].appendNext(uint(j), p.val)
				}
			} else {
				break
			}
		}
	}

	return dags
}

func (t *Trie) createTrie(keys []Rune, vals []*DictUnit) {
	if len(keys) == 0 || len(vals) == 0 {
		return
	}

	if len(keys) != len(vals) {
		panic(fmt.Sprintln(len(keys), len(vals)))
	}

	for i, v := range keys {
		t.insertNode(v, vals[i])
	}
}

func (t *Trie) insertNode(key Rune, val *DictUnit) {
	if len(key) == 0 {
		panic(key)
	}

	p := t.root
	for _, v := range key {
		if p.next == nil {
			p.next = make(map[rune]*TrieNode)
		}

		if x, ok := p.next[v]; ok {
			p = x
		} else {
			node := NewTrieNode()
			p.next[v] = node
			p = node
		}
	}

	if p == nil {
		panic(p)
	}
	p.val = val
}
