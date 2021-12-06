package log_monitor

import (
	"container/list"
	"fmt"
	"sort"
	"time"
	"unicode/utf8"
)

type NextType []*Node

func (n NextType) Len() int { return len(n) }

func (n NextType) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

func (n NextType) Less(i, j int) bool { return n[i].Val < n[j].Val }

type Node struct {
	Val   rune
	Depth int
	Next  NextType
	Fail  *Node
	Type  int
}

// 获取指定val的子节点
func (N *Node) GetChildNodeByVal(val rune) *Node {
	for i := 0; i <= len(N.Next)-1; i++ {
		if N.Next[i].Val == val {
			return N.Next[i]
		}
	}
	return nil
}

// 插入指定val的子节点
func (N *Node) InsertChildNodeByVal(val rune) *Node {
	node := new(Node)
	node.Val = val
	node.Depth = N.Depth + 1
	N.Next = append(N.Next, node)
	return node
}

func (N *Node) BinGetChildNodeByVal(val rune) *Node {
	right := len(N.Next) - 1
	left := 0
	mid := 0
	var midnode *Node
	for left <= right {
		mid = (left + right) / 2
		midnode = N.Next[mid]
		if midnode.Val == val {
			return midnode
		} else if midnode.Val < val {
			left = mid + 1
		} else if midnode.Val > val {
			right = mid - 1
		}
	}
	return nil
}

type Trie struct {
	Root *Node
	Time int64
}

func InitTree(dictionaries map[int][]string) *Trie {
	var tree *Trie
	tree = new(Trie)
	tree.Root = new(Node)
	for typ, words := range dictionaries {
		if len(words) == 0 {
			continue
		}
		for _, word := range words {
			tree.Put(word, typ)
		}
	}
	tree.Build()
	tree.Time = time.Now().Unix()
	return tree
}

// 在trie树中添加word，区分匹配词还是过滤词
func (t *Trie) Put(word string, typ int) {
	if len(word) == 0 {
		return
	}
	parent := t.Root
	for len(word) > 0 {
		char, length := utf8.DecodeRuneInString(word)
		if length <= 0 {
			break
		}
		child := parent.GetChildNodeByVal(char)
		if child == nil {
			child = parent.InsertChildNodeByVal(char)
		}
		parent = child
		word = word[length:]
	}
	parent.Type = typ
}

func (t *Trie) Dump() {

	lst := new(list.List)
	lst.PushBack(t.Root)

	for lst.Len() > 0 {
		node := lst.Remove(lst.Front()).(*Node)
		snode := node.Fail

		var adr *Node = nil
		var sadr *Node = nil
		var cadr *Node = nil

		var val rune = 0
		var sval rune = 0
		var cval rune = 0

		val = node.Val
		adr = node

		if snode != nil {
			sval = snode.Val
			sadr = snode
		}

		fmt.Printf("adr:%p  val:%c  depth:%d  sadr:%p  sval:%c  eow:%v\n", adr, val, node.Depth, sadr, sval, node.Type)

		for _, child := range node.Next {
			cadr = child
			cval = child.Val
			fmt.Printf("-------------->cadr:%p  cval:%c\n", cadr, cval)
			lst.PushBack(child)
		}
	}
}

// 使用BFS构造trie图
func (t *Trie) Build() {
	q := new(list.List)
	q.PushBack(t.Root)
	var p, c *Node

	for q.Len() > 0 {
		node := q.Remove(q.Front()).(*Node) // 设置该节点的子节点的Fail指针
		sort.Sort(node.Next)
		for _, child := range node.Next {
			q.PushBack(child)
			if node == t.Root { // 直接和root相连的节点，它们的fail指针指向root
				child.Fail = t.Root
			} else { // 其他节点：根据父节点的fail指针来设置子节点的fail指针
				p = node.Fail
				for p != nil {
					c = p.GetChildNodeByVal(child.Val)
					if c != nil {
						child.Fail = c
						break
					}
					p = p.Fail
				}
				if p == nil { // 没找到，直接将子节点fail指针指向root
					child.Fail = t.Root
				}
			}
		}
	}
}

func (t *Trie) Match(text []rune, seps Seps) map[int][]string {
	walkNode := t.Root
	var res = make(map[int][]string)
	var nextNode *Node

	for i, char := range text {
		if FindSep(seps, char) {
			continue
		}

		for nextNode = walkNode.BinGetChildNodeByVal(char); nextNode == nil && walkNode != t.Root; {
			walkNode = walkNode.Fail
			nextNode = walkNode.BinGetChildNodeByVal(char)
		}

		walkNode = nextNode
		if walkNode == nil {
			walkNode = t.Root
		}

		temp := walkNode
		for temp != t.Root {
			if temp.Type > 0 {
				start := startIndex(text, i, temp.Depth, seps)
				res[temp.Type] = append(res[temp.Type], string(text[start:i+1]))
			}
			temp = temp.Fail
		}

	}
	return res
}

type Seps []rune

func (s Seps) Len() int {
	return len(s)
}
func (s Seps) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Seps) Less(i, j int) bool {
	return s[i] < s[j]
}

func FindSep(seps Seps, char rune) bool {
	i := sort.Search(len(seps), func(i int) bool { return seps[i] >= char })
	return i < len(seps) && seps[i] == char
}

// 忽略seps中的任意字符
func startIndex(text []rune, end, depth int, seps Seps) int {
	for depth > 0 {
		for end >= 0 && FindSep(seps, text[end]) {
			end--
		}
		if end >= 0 {
			end--
		}
		depth--
	}
	return end + 1
}

//func main() {
//	dict := make(map[int][]string)
//	dict[Warning] = []string{
//		"机器学习",
//		"机器人",
//		"统计学习",
//		"学习时间",
//		"时间机器",
//	}
//	dict[Ignored] = []string{
//		"基础",
//	}
//
//	tree := InitTree(dict)
//	txt := []rune("统计基础的时间机用了机器人学习时间机器")
//	res := tree.Match(txt, []rune{'人'})
//	fmt.Println(res)
//}
