package main

import (
	// "bufio"
	"fmt"
	// "log"
	// "os"
)

// //-----------------------------------------------------------
// // Tree
// //-----------------------------------------------------------
// type Tree struct {
// 	children []*Node
// }

// // Insert a whole string in the tree
// func (t *Tree) insert(domain string) {

// 	currentNode := t

// 	for _, c := range domain {
// 		fmt.Printf("c=%c\n", c)

// 		node := newNode(c)

// 		if currentNode.children == nil {
// 			currentNode.children = make([]*Node, 1)
// 		}
// 		currentNode.children = append(currentNode.children, node)

// 		currentNode = node
// 	}

// }

//-----------------------------------------------------------
// Node
//-----------------------------------------------------------
type Node struct {
	data     rune
	children []*Node
}

// Allocate a new node
func newNode(c rune) *Node {
	return &Node{data: c, children: nil}
}

// Length of a node is the number of children
func (n *Node) len() int {
	return len(n.children)
}

// Return the Node pointer if the children is found
// nil otherwise
func (n *Node) getNode(char rune) *Node {
	for _, c := range n.children {
		if c.data == char {
			return c
		}
	}
	return nil
}

// Attach a new Node to the current Node
func (n *Node) addNode(char rune) *Node {
	// allocate a new node
	node := newNode(char)

	// safeguard
	if n.children == nil {
		n.children = []*Node{node}
		return node
	}

	if nref := n.getNode(char); nref != nil {
		nref.children = append(n.children, node)
	} else {
		// if not, just append it
		n.children = append(n.children, node)
	}

	return node
}

// Insert a whole string in the tree
func (n *Node) Insert(domain string) {

	// we'll loop using this node
	currentNode := n

	// add each individual char
	for _, c := range domain {

		// create a new Node pointer based a char pointed by the c variable
		node := currentNode.addNode(c)

		fmt.Printf("c=%c, length=%d\n", c, len(currentNode.children))

		currentNode = node
	}
}

// // Output a tree as a graphviz .dot
// func (n *Node) toDot(name string) {
// 	f, err := os.Create(name)
// 	if err != nil {
// 		log.Fatalf("unable to create file <%s>\n", err)
// 	}
// 	defer f.Close()

// 	w := bufio.NewWriter(f)

// 	// start Graphviz DOT file
// 	f.WriteString("digraph {\n")

// 	for _, c := range n.children {
// 		c.traverse(w)
// 	}

// 	// end DOT file
// 	f.WriteString("}")
// }

// // Traverse the all tree starting from any node
// func (n *Node) traverse(w *bufio.Writer) {
// 	if n == nil {
// 		return
// 	}

// 	if n.children == nil {
// 		return
// 	}

// 	fmt.Printf("node='%c'\n", n.data)

// 	for _, c := range n.children {
// 		fmt.Printf("%c -> %c\n", n.data, c.data)
// 		fmt.Fprintf(w, "%c -> %c\n", n.data, c.data)
// 		c.traverse(w)
// 	}
// 	w.Flush()
// }
