package modules

import (
	"bufio"
	"github.com/terashin777/vm-translator/models"
	"io"
	"os"
)

type CodeWriter struct {
	w io.WriteCloser
	bw *bufio.Writer
	t models.ITranslator
}

func NewCodeWriter(w io.WriteCloser, t models.ITranslator) *CodeWriter {
	return &CodeWriter{w, bufio.NewWriter(w), t}
}

func (w *CodeWriter) SetFileName(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}

	w.w = f
	return nil
}

func (w *CodeWriter) WriteArithmetic(c string) error {
	ar := w.t.TranslateArithmetic(c)
	if ar == "" {
		return nil
	}

	_, err := w.bw.WriteString(ar)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) WritePushPop(c, seg string, i int) error {
	p := w.t.TranslatePushPop(c)
	if p == "" {
		return nil
	}

	_, err := w.bw.Write(p)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) Close() error {
	defer w.w.Close()
	err := w.bw.Flush()
	if err != nil {
		return err
	}

	return nil
}
