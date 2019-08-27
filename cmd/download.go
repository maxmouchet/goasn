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
	"log"
	"path/filepath"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/maxmouchet/goasn/pkg/collectors"

	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		collectorName := cmd.Flag("collector").Value.String()
		collector, err := collectors.NewCollector(collectorName)
		check(err)

		dateStr := cmd.Flag("date").Value.String()
		date, err := time.Parse("2006-01-02T15:04", dateStr)
		check(err)

		url := collector.TableURL(date)
		dst := filepath.Base(url)

		log.Printf("Downloading RIB from collector %s at %s", collector.Name(), date)
		log.Printf("Downloading %s to %s", url, dst)

		grab.Get(dst, url)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().String("collector", "route-views2.oregon-ix.net", "")
	downloadCmd.Flags().String("date", "", "")
	// TODO: Get latest by default
	downloadCmd.MarkFlagRequired("date")
}
