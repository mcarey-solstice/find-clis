package main

import (
	"bufio"
	"path/filepath"
	"os"
	"flag"
	"fmt"
	"errors"
	"strings"
)

type Find struct {
	Filename string
	Root string
	Clis []string
}

func (f *Find) Paths() []string {
	paths := make([]string, len(f.Clis))
	p := make(map[string]bool, len(f.Clis))

	c := 0
	for _, cli := range f.Clis {
		path := filepath.Dir(cli)

		if _, ok := p[path]; !ok {
			c++
			p[path] = true
			paths = append(paths, path)
		}
	}

	return paths[c+1:]
}

func NewFind(filename string, root string) (*Find, error) {
	if filename == "" {
		return nil, errors.New("Filename cannot be empty!")
	}

	if root == "" {
		root = "."
	}

	return &Find {
		Filename: filename,
		Root: root,
	}, nil
}

func (f *Find) AddCli(path string) {
	f.Clis = append(f.Clis, path)
}

func (f *Find) AddClisFromFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		f.AddCli(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (f *Find) SearchFn() func(string, os.FileInfo, error) error {
	return func(cli string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.Name() == f.Filename {
			err = f.AddClisFromFile(cli)
			if err != nil {
				return err
			}
		}

		return nil
	}

}

func (f *Find) Walk() error {
	err := filepath.Walk(f.Root, f.SearchFn())
	if err != nil {
		return err
	}

	return nil
}

func Search(filename string, root string) (*Find, error) {
	find, err := NewFind(filename, root)
	if err != nil {
		return nil, err
	}

	err = find.Walk()
	if err != nil {
		return nil, err
	}

	return find, nil
}

func FindClis(filename string, root string) ([]string, error) {
	f, e := Search(filename, root)
	if e != nil {
		return nil, e
	}

	return f.Clis, nil
}

func FindPaths(filename string, root string) ([]string, error) {
	f, e := Search(filename, root)
	if e != nil {
		return nil, e
	}

	return f.Paths(), nil
}

func TestClis(filename string, root string) error {
	return nil
}

func handleError(err error) {
	panic(err)
}

func main() {
	var filename, root, join string
	flag.StringVar(&filename, "filename", "", "The name of the file to look for")
	flag.StringVar(&root, "root", ".", "The directory to look for these files")
	flag.StringVar(&join, "join", " ", "The joining string for list and paths commands")
	flag.Parse()

	command := flag.Arg(0)
	if command == "" {
		handleError(errors.New("Missing command!"))
	}

	switch command {
	case "paths":
		paths, err := FindPaths(filename, root)
		if err != nil {
			handleError(err)
		}

		fmt.Print(strings.Join(paths, join))
	case "test":
		err := TestClis(filename, root)
		if err != nil {
			handleError(err)
		}

		fmt.Println("All tests have passed")
	case "list":
		clis, err := FindClis(filename, root)
		if err != nil {
			handleError(err)
		}
		fmt.Println(strings.Join(clis, join))
	default:
		handleError(errors.New(fmt.Sprintf("Unknown command: %s", command)))
	}
}
