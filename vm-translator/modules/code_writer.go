package modules

import (
	"bufio"
	"io"
	"os"

	"github.com/terashin777/vm-translator/models"
)

type CodeWriter struct {
	w  io.WriteCloser
	bw *bufio.Writer
	t  models.ITranslator
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

func (w *CodeWriter) WritePushPop(c models.CommandType, seg string, i int) error {
	p := ""
	switch c {
	case models.C_PUSH:
		p = w.t.TranslatePush(seg, i)
	case models.C_POP:
		p = w.t.TranslatePop(seg, i)
	}
	if p == "" {
		return nil
	}

	_, err := w.bw.WriteString(p)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) WriteLabel(l string) error {
	ar := w.t.TranslateLabel(l)
	if ar == "" {
		return nil
	}

	_, err := w.bw.WriteString(ar)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) WriteGoto(l string) error {
	ar := w.t.TranslateGoto(l)
	if ar == "" {
		return nil
	}

	_, err := w.bw.WriteString(ar)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) WriteIf(l string) error {
	ar := w.t.TranslateIf(l)
	if ar == "" {
		return nil
	}

	_, err := w.bw.WriteString(ar)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) WriteCall(fn string, n int) error {
	ar := w.t.TranslateCall(fn, n)
	if ar == "" {
		return nil
	}

	_, err := w.bw.WriteString(ar)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) WriteReturn() error {
	ar := w.t.TranslateReturn()
	if ar == "" {
		return nil
	}

	_, err := w.bw.WriteString(ar)
	if err != nil {
		return err
	}

	return nil
}

func (w *CodeWriter) WriteFunction(fn string, n int) error {
	ar := w.t.TranslateFunction(fn, n)
	if ar == "" {
		return nil
	}

	_, err := w.bw.WriteString(ar)
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