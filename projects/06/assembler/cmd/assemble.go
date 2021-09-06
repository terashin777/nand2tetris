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
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/terashin777/assembler/modules"
)

var dest string

// assembleCmd represents the assemble command
var assembleCmd = &cobra.Command{
	Use:   "assemble",
	Short: "assemble your assembly to hack.",
	Long:  `assemble your assembly to hack.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return assemble(args[0])
	},
}

func init() {
	rootCmd.AddCommand(assembleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// assembleCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	assembleCmd.Flags().StringVarP(&dest, "dest", "d", ".", "destination for assembled file")
}

func assemble(path, dest string) error {
	r, err := os.Open(path)
	if err != nil {
		return err
	}
	defer r.Close()
	p := modules.NewParser(r)

	w, err := os.Open(dest)
	if err != nil {
		return err
	}
	defer w.Close()

	for {
		b, err := p.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = w.Write(b)
		if err != nil {
			return err
		}
	}

	return nil
}
