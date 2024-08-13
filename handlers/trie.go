package handlers

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

func (t *Trie) Add(v string) {
	if len(v) == 0 {
		return
	}
	current := t.root
	for _, r := range v {
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

