package tinyGin

import "strings"

//注册路由，映射handler；匹配路由，执行handler

type node struct {
	pattern  string  // 待匹配路由
	part     string  // 路由中的一部分
	children []*node // 子节点
	isWild   bool    // 是否精确匹配,part含有:或*时为true
}

//第一个匹配成功的节点，用于插入

func (n *node) matchChild(pattern string) *node {
	for _, child := range n.children {
		if child.part == pattern || child.isWild {
			return child
		}
	}
	return nil
}

//所有匹配成功的节点，用于查询

func (n *node) matchChildren(pattern string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == pattern || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 递归

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
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
