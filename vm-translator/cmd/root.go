/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/terashin777/vm-translator/models"
	"github.com/terashin777/vm-translator/modules"
)

var (
	cfgFile        string
	dest           string
	assemblyExt    = ".asm"
	vmExt          = ".vm"
	defaultDestDir = "same dir as source file"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "translate [file]",
	Short: "translate your vm code to assembly",
	Long:  `translate your vm code to assembly.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := translate(args[1], dest)
		if err != nil {
			fmt.Printf("assemble is failed because: %s", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vm-translator.yaml)")
	rootCmd.Flags().StringVarP(&dest, "dest", "d", defaultDestDir, "destination for translated file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".vm-translator")
	}

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func translate(src, dest string) error {
	fns, isDir, err := getFiles(src)
	if err != nil {
		return err
	}
	if isDir {
		dest = src
	}

	cw, err := newCodeWriter(src, dest)
	if err != nil {
		return err
	}
	defer cw.Close()

	p, err := modules.NewParser(fns)
	if err != nil {
		return err
	}
	return parseAll(p, cw, isDir)
}

func getFiles(src string) ([]string, bool, error) {
	fi, err := os.Stat(src)
	if err != nil {
		return nil, false, err
	}
	if fi.IsDir() {
		fs, err := getFilesInDir(src)
		if err != nil {
			return nil, false, err
		}

		return fs, true, nil
	}

	return []string{
		src,
	}, false, nil
}

func getFilesInDir(dir string) ([]string, error) {
	fs, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fns := []string{}
	for _, f := range fs {
		fn := f.Name()
		if isVM(fn) {
			fns = append(fns, filepath.Join(dir, fn))
		}
	}

	return fns, nil
}

func isVM(fn string) bool {
	return filepath.Ext(fn) == vmExt
}

func newCodeWriter(src, dest string) (*modules.CodeWriter, error) {
	if dest == defaultDestDir {
		dest = filepath.Join(filepath.Dir(src), makeSameFileName(src))
	} else {
		dest = filepath.Join(dest, makeSameFileName(src))
	}
	w, err := os.Create(dest)
	if err != nil {
		return nil, err
	}

	return modules.NewCodeWriter(
		w,
		modules.NewTranslator(),
	), nil
}

func parseAll(p *modules.Parser, w *modules.CodeWriter, isDir bool) error {
	err := nextFile(p, w)
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	if isDir {
		w.WriteInit()
	}
	for {
		err := p.Advance()
		if err == io.EOF {
			err := nextFile(p, w)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}

			continue
		}
		if err != nil {
			return err
		}

		ct := p.CommandType()
		switch ct {
		case models.C_ARITHMETIC:
			err = w.WriteArithmetic(p.Arg1())
		case models.C_PUSH, models.C_POP:
			err = w.WritePushPop(ct, p.Arg1(), p.Arg2())
		case models.C_LABEL:
			err = w.WriteLabel(p.Arg1())
		case models.C_GOTO:
			err = w.WriteGoto(p.Arg1())
		case models.C_IF:
			err = w.WriteIf(p.Arg1())
		case models.C_CALL:
			err = w.WriteCall(p.Arg1(), p.Arg2())
		case models.C_FUNCTION:
			err = w.WriteFunction(p.Arg1(), p.Arg2())
		case models.C_RETURN:
			err = w.WriteReturn()
		}
		if err != nil {
			return err
		}
	}
}

func makeSameFileName(path string) string {
	fn := filepath.Base(path)
	return fmt.Sprintf("%s%s", strings.Split(fn, ".")[0], assemblyExt)
}

func nextFile(p *modules.Parser, w *modules.CodeWriter) error {
	fn, err := p.NextFile()
	if err != nil {
		return err
	}

	w.SetFunctionName(fn)
	return nil
}
