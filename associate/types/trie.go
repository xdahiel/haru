package types

import "sync"

type uString []rune

type uTrieNode struct {
	isEnd    bool
	children map[rune]*uTrieNode
}

type UTrie struct {
	root *uTrieNode
}

func NewUTrie() *UTrie {
	return &UTrie{root: &uTrieNode{
		isEnd:    false,
		children: make(map[rune]*uTrieNode),
	}}
}

var (
	trie     *UTrie
	trieOnce sync.Once
)

func GetTrie() *UTrie {
	trieOnce.Do(func() {
		trie = NewUTrie()
	})
	return trie
}

func (t *UTrie) Insert(sentence string) {
	node := t.root
	us := uString(sentence)
	for _, char := range us {
		if node.children[char] == nil {
			node.children[char] = &uTrieNode{
				isEnd:    false,
				children: make(map[rune]*uTrieNode),
			}
		}
		node = node.children[char]
	}
	node.isEnd = true
}

// Seek 寻找匹配指定前缀的所有字符串
func (t *UTrie) Seek(pre string) []string {
	node := t.root
	us := uString(pre)
	var res []string

	for _, char := range us {
		if node.children[char] == nil {
			return nil
		}
		node = node.children[char]
	}
	if node == nil {
		return nil
	}

	var dfs func(*uTrieNode, uString)
	dfs = func(n *uTrieNode, cur uString) {
		if n.isEnd {
			res = append(res, pre+string(cur))
		}
		for k, v := range n.children {
			cur = append(cur, k)
			dfs(v, cur)
			cur = cur[:len(cur)-1]
		}
	}
	dfs(node, []rune{})
	return res
}
