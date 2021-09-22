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
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/terashin777/assembler/modules"

	"github.com/spf13/viper"
)

var (
	cfgFile        string
	dest           string
	assembledExt   = ".hack"
	defaultDestDir = "same dir as source file"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "assemble [file]",
	Short: "assemble your assembly to hack",
	Long:  `assemble your assembly to hack.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := assemble(args[1], dest)
		if err != nil {
			fmt.Printf("assemble is failed because: %s", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.assembler.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&dest, "dest", "d", defaultDestDir, "destination for assembled file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".assembler" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".assembler")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func assemble(path, dest string) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	if dest == defaultDestDir {
		dest = filepath.Join(filepath.Dir(path), makeSameFileName(path))
	}
	fw, err := createDestFile(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	i, err := fw.Stat()
	if err != nil {
		return err
	}
	size := i.Size()

	p, err := modules.NewParser(r, size)
	if err != nil {
		return err
	}
	err = r.Close()
	if err != nil {
		return err
	}

	w := bufio.NewWriter(fw)
	for {
		s, err := p.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if s == "" {
			continue
		}

		_, err = w.WriteString(s + "\n")
		if err != nil {
			return err
		}
	}
	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

func createDestFile(dest string) (*os.File, error) {
	var w *os.File
	w, err := os.Create(dest)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func makeSameFileName(path string) string {
	fn := filepath.Base(path)
	return fmt.Sprintf("%s%s", strings.Split(fn, ".")[0], assembledExt)
}
