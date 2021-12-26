package modules

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/terashin777/jack-analyzer/models"
)

type JackAnalyzer struct {
	e  *CompilationEngine
	fs []string
	fi int
}

func NewJackAnalyzer(src, dest string) (*JackAnalyzer, error) {
	fs, err := getFileNames(src, dest)
	if err != nil {
		return nil, err
	}

	e, err := NewCompilationEngine(src, dest, NewJackTokenizer())
	if err != nil {
		return nil, err
	}

	return &JackAnalyzer{
		e:  e,
		fs: fs,
		fi: -1,
	}, nil
}

func getFileNames(src, dest string) ([]string, error) {
	stat, err := os.Stat(src)
	if err != nil {
		return nil, err
	}

	fs := []string{}
	if stat.IsDir() {
		fs, err = getFileNamesFromDir(src)
		if err != nil {
			return nil, err
		}
	} else {
		fs = []string{src}
	}
	if len(fs) == 0 {
		return nil, fmt.Errorf("no jack file")
	}

	return fs, nil
}

func getFileNamesFromDir(dir string) ([]string, error) {
	es, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fns := []string{}
	for _, e := range es {
		f := e.Name()
		if isJack(f) {
			fns = append(fns, filepath.Join(dir, f))
		}
	}

	return fns, nil
}

func isJack(fn string) bool {
	return filepath.Ext(fn) == models.JackExt
}

func (a *JackAnalyzer) Compile() error {
	f, err := a.nextFile()
	if err == io.EOF {
		return nil
	}
	a.e.SetFile(f)

	err = a.e.Open()
	if err != nil {
		return err
	}

	err = a.e.CompileClass()
	if err == io.EOF {
		return a.Compile()
	}
	if err != nil {
		return err
	}

	return nil
}

func (a *JackAnalyzer) nextFile() (string, error) {
	a.fi++
	if !a.hasNextFiles() {
		return "", io.EOF
	}

	return a.fs[a.fi], nil
}

func (a *JackAnalyzer) hasNextFiles() bool {
	return a.fi < len(a.fs)
}
