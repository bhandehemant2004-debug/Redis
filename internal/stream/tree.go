package stream

import "bytes"

type node struct {
	prefix   []byte
	children []*node 
	value    any
	hasValue bool
}

type Tree struct{ root *node }

func New() *Tree { return &Tree{root: &node{}} }

func commonPrefix(a, b []byte) int {
	i := 0
	for i < len(a) && i < len(b) && a[i] == b[i] {
		i++
	}
	return i
}

func (t *Tree) Insert(key []byte, val any) {
	n := t.root
	for {
		if len(key) == 0 {
			n.value, n.hasValue = val, true
			return
		}
		idx := -1
		for i, c := range n.children {
			if c.prefix[0] == key[0] {
				idx = i
				break
			}
		}
		if idx == -1 {
			child := &node{prefix: append([]byte(nil), key...), value: val, hasValue: true}
			n.addChild(child)
			return
		}
		c := n.children[idx]
		cp := commonPrefix(c.prefix, key)
		if cp < len(c.prefix) { 
			split := &node{
				prefix:   c.prefix[cp:],
				children: c.children,
				value:    c.value,
				hasValue: c.hasValue,
			}
			c.prefix = c.prefix[:cp]
			c.children = []*node{split}
			c.value, c.hasValue = nil, false
		}
		key = key[cp:]
		n = c
	}
}

func (n *node) addChild(c *node){
	i := 0
	for i < len(n.children) && n.children[i].prefix[0] < c.prefix[0] {
		i++
	}
	n.children = append(n.children, nil)
	copy(n.children[i+1:], n.children[i:])
	n.children[i] = c
}


func (t *Tree) Ascend(from []byte, fn func(key []byte, val any) bool) {
	t.root.walk(nil, from, true, fn)
}

func (n *node) walk(acc, bound []byte, bounded bool, fn func([]byte, any) bool) bool {
	acc = append(append([]byte(nil), acc...), n.prefix...)
	if bounded {
		m := min(len(n.prefix), len(bound))
		switch c := bytes.Compare(n.prefix[:m], bound[:m]); {
		case c < 0:
			return true // whole subtree is below the bound
		case c > 0:
			bounded = false // whole subtree is above
		default:
			if len(n.prefix) >= len(bound) {
				bounded = false
			} else {
				bound = bound[m:]
			}
		}
	}
	if n.hasValue && !bounded {
		if !fn(acc, n.value) {
			return false
		}
	}
	for _, c := range n.children {
		if !c.walk(acc, bound, bounded, fn) {
			return false
		}
	}
	return true
}

func (t *Tree) Floor(k []byte) (key []byte, val any, ok bool) {
	var best []byte
	var bestVal any
	found := false
	t.Ascend(nil, func(kk []byte, v any) bool {
		if bytes.Compare(kk, k) > 0 {
			return false
		}
		best, bestVal, found = kk, v, true
		return true
	})
	return best, bestVal, found
}
func (t *Tree) Last() ([]byte, any, bool) {
    n := t.root
    var key []byte

    for {
        key = append(key, n.prefix...)

        if len(n.children) > 0 {
            n = n.children[len(n.children)-1]
            continue
        }
        if n.hasValue {
            return key, n.value, true
        }

        return nil, nil, false
    }
}