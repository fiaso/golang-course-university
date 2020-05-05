package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

type Node struct {
	Name string
	IsDir bool
	Size int64
	Children []*Node
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	startPath := path
	tree := Node{
		Name: path,
		IsDir: true,
	}
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == startPath {
				return nil
			}
			if !info.IsDir() && !printFiles {
				return nil
			}
			curPath := strings.Split(strings.TrimPrefix(path, startPath), "/")
			parent := &tree
			for i := 1; i < len(curPath)-1; i++{
				for _, node := range parent.Children {
					if node.Name == curPath[i] {
						parent = node
					}
				}
			}
			parent.Children = append(parent.Children, &Node{
				Name: info.Name(),
				IsDir: info.IsDir(),
				Size: info.Size(),
			})
			return nil
		})
	if err != nil {
		return err
	}
	sortTree(&tree)
	for i, node := range tree.Children {
		isLast := false
		if i == len(tree.Children)-1 {
			isLast = true
		}
		err = printTree(out, node,0, []bool{isLast})
		if err != nil {
			return err
		}
	}
	return nil
}

func sortTree(curNode *Node)  {
	sort.Slice(curNode.Children, func(i, j int) bool {
		return curNode.Children[i].Name < curNode.Children[j].Name
	})
	for _, node := range curNode.Children {
		sortTree(node)
	}
}

func printTree(out io.Writer, curNode *Node, level int, isLast []bool) error{
	if curNode.IsDir {
		err := printNode(out, curNode.Name, level, isLast)
		if err != nil {
			return err
		}
		for i, node := range curNode.Children {
			isLastNode := false
			if i == len(curNode.Children)-1 {
				isLastNode = true
			}
			err = printTree(out, node, level+1, append(isLast, isLastNode))
			if err != nil {
				return err
			}
		}
	} else {
		curNode.Name += " ("
		if curNode.Size > 0 {
			curNode.Name += strconv.FormatInt(curNode.Size, 10) + "b"
		} else {
			curNode.Name += "empty"
		}
		curNode.Name += ")"
		err := printNode(out, curNode.Name, level, isLast)
		if err != nil {
			return err
		}
	}
	return nil
}

func printNode(out io.Writer, nodeName string, level int, isLast []bool) error{
	strTree := ""
	for i := 0; i < level; i++ {
		if !isLast[i] {
			strTree += "│\t"
		} else {
			strTree += "\t"
		}
	}
	if isLast[level] {
		strTree += "└───"
	} else {
		strTree += "├───"
	}
	strTree += nodeName
	_, err := fmt.Fprintln(out, strTree)
	return err
}