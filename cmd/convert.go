/*
Copyright © 2019 Maxime Mouchet <max@maxmouchet.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/maxmouchet/goasn/pkg/goasn"

	"github.com/spf13/cobra"
)

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Get bool directly
		singleAS := cmd.Flag("single-as").Value.String() == "true"

		for _, path := range args {
			// goasn.
			entries, err := goasn.RIBFromMRT(path)
			log.Println(len(entries), err)

			origins := make([]goasn.PrefixOrigin, len(entries))
			for i, entry := range entries {
				origins[i] = goasn.NewPrefixOrigin(entry)
			}

			// json.

			b, err := goasn.ASNDatabase(origins).MarshalText(singleAS)
			check(err)

			err = ioutil.WriteFile(path+".txt", b, os.FileMode(0644))
			check(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().String("format", "txt", "format: json or txt")
	convertCmd.Flags().Bool("single-as", false, "pyasn compatible format (only a single AS from AS-SET origin)")
}
