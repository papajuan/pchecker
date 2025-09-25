package pchecker

import (
	"fmt"
	"sync"
)

/**
 * @author  papajuan
 * @date    10/10/2024
 **/

// SafeTrie represents the concurrently safe prefix tree
type SafeTrie[K comparable] struct {
	root       *node[K]
	lock       sync.RWMutex
	comparator func(K) K
	strFunc    func([]K) string
}

// node represents a node in the Trie
type node[K comparable] struct {
	children map[K]*node[K]
	isEnd    bool // Marks the end of a word
}

func NewSafeTrie[K comparable](length int) *SafeTrie[K] {
	return &SafeTrie[K]{
		root: &node[K]{
			children: make(map[K]*node[K], length),
		},
	}
}

func (t *SafeTrie[K]) WithStrFunc(f func([]K) string) *SafeTrie[K] {
	t.strFunc = f
	return t
}

func (t *SafeTrie[K]) WithComparator(f func(K) K) *SafeTrie[K] {
	t.comparator = f
	return t
}

// Insert adds a word to the Trie
func (t *SafeTrie[K]) Insert(arr []K) {
	t.lock.Lock()
	defer t.lock.Unlock()
	n := t.root
	for _, key := range arr {
		// If the symbol does not exist, create a new node
		if _, exists := n.children[key]; !exists {
			n.children[key] = &node[K]{
				children: make(map[K]*node[K]),
			}
		}
		n = n.children[key]
	}
	n.isEnd = true
}

// Exists checks if a word exists in the Trie
func (t *SafeTrie[K]) Exists(arr []K) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	n := t.root
	for _, ch := range arr {
		if t.comparator != nil {
			ch = t.comparator(ch)
		}
		if _, exists := n.children[ch]; !exists {
			return false
		}
		n = n.children[ch]
	}
	return n.isEnd
}

// StartsWith checks if any words in the Trie start with the given prefix
func (t *SafeTrie[K]) StartsWith(arr []K) (bool, bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	n := t.root
	for _, ch := range arr {
		if t.comparator != nil {
			ch = t.comparator(ch)
		}
		// If the character doesn't exist, return false
		if _, exists := n.children[ch]; !exists {
			return false, false
		}
		n = n.children[ch]
	}
	return true, n.isEnd
}

func (t *SafeTrie[K]) IsPrefixInTrie(arr []K) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()
	n := t.root
	for _, ch := range arr {
		if t.comparator != nil {
			ch = t.comparator(ch)
		}
		if _, exists := n.children[ch]; !exists {
			break
		}
		n = n.children[ch]
	}
	return n.children == nil || len(n.children) == 0
}

func (t *SafeTrie[K]) PrintAll() {
	if t.strFunc == nil {
		panic("string function is not set")
	}
	t.lock.RLock()
	defer t.lock.RUnlock()
	var dfs func(n *node[K], prefix []K)
	dfs = func(n *node[K], prefix []K) {
		if n.isEnd {
			fmt.Println(t.strFunc(prefix))
		}
		for r, child := range n.children {
			dfs(child, append(prefix, r))
		}
	}
	dfs(t.root, nil)
}
