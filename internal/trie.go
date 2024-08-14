package internal

type node struct {
	children map[rune]*node
	isEnd    bool
}

type Trie struct {
	root *node
}

func NewTrie() *Trie {
	return &Trie{
		root: &node{
			children: make(map[rune]*node),
			isEnd:    false,
		},
	}
}

func (t *Trie) Add(s string) {
	if len(s) == 0 {
		return
	}
	current := t.root
	for _, r := range s {
		if _, exists := current.children[r]; !exists {
			current.children[r] = &node{
				children: make(map[rune]*node),
				isEnd:    false,
			}
		}
		current = current.children[r]
	}
	current.isEnd = true
}

func (t *Trie) Exists(s string) bool {
	if len(s) == 0 {
		return false
	}
	current := t.root
	for _, r := range s {
		if _, exists := current.children[r]; !exists {
			return false
		}
		current = current.children[r]
	}
	return current.isEnd
}
