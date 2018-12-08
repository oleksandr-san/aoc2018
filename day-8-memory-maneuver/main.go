package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type node struct {
	children []*node
	metadata []int
}

func walkNodes(node *node, handler func(*node)) {
	handler(node)
	for _, child := range node.children {
		walkNodes(child, handler)
	}
}

type lexer struct {
	reader io.Reader
}

func (l *lexer) getInt() (int, error) {
	var value int
	_, err := fmt.Fscanf(l.reader, "%d", &value)
	return value, err
}

type parser struct {
	lexer lexer
}

type parserState int

const (
	readHeader parserState = iota
	readChildren
	readMetadata
	commitNode
)

type parserContext struct {
	node          *node
	state         parserState
	childrenCount int
	metadataSize  int
}

func (p *parser) parse() (*node, error) {
	rootNode := &node{}
	ctxStack := []*parserContext{&parserContext{node: rootNode}}

parseLoop:
	for len(ctxStack) != 0 {
		ctx := ctxStack[len(ctxStack)-1]

		switch ctx.state {
		case readHeader:
			childrenCnt, err := p.lexer.getInt()
			if err == io.EOF {
				break parseLoop
			} else if err != nil {
				return nil, err
			}

			metadataSize, err := p.lexer.getInt()
			if err != nil {
				return nil, err
			}

			ctx.childrenCount = childrenCnt
			ctx.metadataSize = metadataSize
			ctx.state = readChildren

		case readChildren:
			if len(ctx.node.children) < ctx.childrenCount {
				ctxStack = append(ctxStack, &parserContext{node: &node{}})
			} else {
				ctx.state = readMetadata
			}

		case readMetadata:
			for i := 0; i < ctx.metadataSize; i++ {
				entry, err := p.lexer.getInt()
				if err != nil {
					return nil, err
				}

				ctx.node.metadata = append(ctx.node.metadata, entry)
			}
			ctx.state = commitNode

		case commitNode:
			if len(ctxStack) > 1 {
				parentCtx := ctxStack[len(ctxStack)-2]
				parentCtx.node.children = append(parentCtx.node.children, ctx.node)
			}
			ctxStack = ctxStack[:len(ctxStack)-1]
		}
	}

	return rootNode, nil
}

func calculateMetadataSum(n *node) (sum int) {
	walkNodes(n, func(n *node) {
		for _, entry := range n.metadata {
			sum += entry
		}
	})
	return
}

func calculateSpecificSum(n *node) (sum int) {
	if len(n.children) == 0 {
		for _, entry := range n.metadata {
			sum += entry
		}
	} else {
		for _, entry := range n.metadata {
			if entry > 0 && entry < len(n.children)+1 {
				sum += calculateSpecificSum(n.children[entry-1])
			}
		}
	}
	return
}

func readNodeTree(path string) (*node, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	parser := parser{lexer{f}}
	node, err := parser.parse()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	node, err := readNodeTree(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Sum of all metadata entries: %d\n", calculateMetadataSum(node))
	fmt.Printf("Specific sum of root node: %d\n", calculateSpecificSum(node))
}
