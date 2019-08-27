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
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

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

			err = writeDatabase(path+".txt", filepath.Base(path), singleAS, origins)
			check(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().Bool("single-as", false, "pyasn compatible format (only a single AS from AS-SET origin)")
}

func formatSlice(s []uint32) string {
	if len(s) == 0 {
		return ""
	}
	str := fmt.Sprintf("%d", s[0])
	for _, e := range s[1:] {
		str += fmt.Sprintf(",%d", e)
	}
	return str
}

func writeDatabase(path string, source string, singleAS bool, entries []goasn.PrefixOrigin) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)

	// TODO: Diff with pyasn
	fmt.Fprintf(w, "; IP-ASN32-DAT file\n")
	fmt.Fprintf(w, "; Original source:\t%s\n", source)
	fmt.Fprintf(w, "; Converted on:\t%s\n", time.Now().Format("Mon Jan 2 15:04:05 2006"))
	fmt.Fprintf(w, "; Prefixes-v4:\t%d\n")
	fmt.Fprintf(w, "; Prefixes-v6:\t%d\n")
	fmt.Fprintf(w, ";\n")

	lastNet := ""

	_, defaultV4, _ := net.ParseCIDR("0.0.0.0/0")
	_, defaultV6, _ := net.ParseCIDR("::/0")

	// WARN if same prefix with differents ASes

	for _, entry := range entries {
		if entry.Prefix.String() == lastNet {
			continue
		}

		// TODO: Optimize
		if (entry.Prefix.String() == defaultV4.String()) || (entry.Prefix.String() == defaultV6.String()) {
			continue
		}

		lastNet = entry.Prefix.String()
		asns := entry.Origin
		if singleAS {
			asns = asns[0:1]
		}
		fmt.Fprintf(
			w,
			"%s\t%s\n",
			entry.Prefix.String(),
			formatSlice(asns),
		)
	}

	w.Flush()

	return nil
}
