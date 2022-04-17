package hee

import "strings"

type tireNode struct {
	pattern  string      // 待匹配路由，例如 /p/:lang,存放在前缀树的最底端
	part     string      //路由的一部分
	children []*tireNode //前缀树节点
	isWild   bool        // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *tireNode) matchChild(part string) *tireNode {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *tireNode) matchChildren(part string) []*tireNode {
	var nodes []*tireNode
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入路由规则
func (n *tireNode) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &tireNode{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//进入下一层
	child.insert(pattern, parts, height+1)
}

//查找路由
func (n *tireNode) search(parts []string, height int) *tireNode {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
