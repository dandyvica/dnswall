package main

import (
	//"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNode(t *testing.T) {
	assert := assert.New(t)

	// create a sample slice of nodes
	n := newNode('\x00')
	n.children = make([]*Node, 4)

	n.children[0] = &Node{data: 'a', children: nil}
	n.children[1] = &Node{data: 'b', children: nil}
	n.children[2] = &Node{data: 'c', children: nil}
	n.children[3] = &Node{data: 'd', children: nil}

	assert.Equal(len(n.children), 4)
	assert.NotNil(n.getNode('a'))
	assert.Equal(n.getNode('a').data, 'a')
	assert.Nil(n.getNode('e'))
}

func TestAddNode(t *testing.T) {
	assert := assert.New(t)

	root := Node{data: '\x00', children: nil}
	a := root.addNode('a')
	b := root.addNode('b')

	assert.Equal(len(root.children), 2)
	assert.Equal(a.data, 'a')
	assert.Nil(a.children)
	assert.Equal(b.data, 'b')
	assert.Nil(b.children)	

	a = root.addNode('a')
	assert.Equal(len(root.children), 2)
	assert.Equal(a.data, 'a')
	assert.Nil(a.children)	

	c := root.addNode('c')
	assert.Equal(len(root.children), 3)
	assert.Equal(c.data, 'c')
	assert.Nil(c.children)	

	a.addNode('c')
	a.addNode('d')
	assert.Equal(len(a.children), 2)	
}

func TestInsert(t *testing.T) {
	assert := assert.New(t)

	root := Node{data: '\x00', children: nil}

	// add one string: a ⭢ b ⭢ c
	root.Insert("abc")
	assert.Equal(len(root.children), 1)

	assert.Equal(root.getNode('a').data, 'a')
	assert.Equal(len(root.getNode('a').children), 1)

	assert.Equal(root.getNode('a').getNode('b').data, 'b')
	assert.Equal(len(root.getNode('a').getNode('b').children), 1)

	assert.Equal(root.getNode('a').getNode('b').getNode('c').data, 'c')
	assert.Equal(len(root.getNode('a').getNode('b').getNode('c').children), 0)

	// add another string having nothing in common with the first one
	// a ⭢ b ⭢ c
	// d ⭢ e ⭢ f
	root.Insert("def")
	assert.Equal(len(root.children), 2)

	assert.Equal(root.getNode('d').data, 'd')
	assert.Equal(len(root.getNode('d').children), 1)

	assert.Equal(root.getNode('d').getNode('e').data, 'e')
	assert.Equal(len(root.getNode('d').getNode('e').children), 1)

	assert.Equal(root.getNode('d').getNode('e').getNode('f').data, 'f')
	assert.Equal(len(root.getNode('d').getNode('e').getNode('f').children), 0)

	// add another one having something in common
	// a ⭢ b ⭢ c
	//        ⭨ d
	// d ⭢ e ⭢ f	
	root.Insert("abc")
	assert.Equal(len(root.children), 2)

	a := root.getNode('a')
	assert.Equal(a.data, 'a')
	assert.Equal(a.len(), 1)
	//assert.Equal(len(root.getNode('a').getNode('b').children), 2)
	//assert.Equal(len(root.getNode('a').getNode('b').getNode('c').children), 0)
	//assert.Equal(len(root.getNode('a').getNode('b').getNode('d').children), 0)


	// assert.Equal(len(root.getNode('d').children), 1)
	// assert.Equal(len(root.getNode('e').children), 1)
	// assert.Equal(len(root.getNode('f').children), 0)


}
