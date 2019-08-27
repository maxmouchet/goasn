package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/maxmouchet/goasn/pkg/goasn"
	"github.com/maxmouchet/goasn/pkg/peeringdb"

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

		// TMP

		var db peeringdb.DB
		err := db.FromAPI()
		check(err)

		tree, err := goasn.NewIXPTree(db)
		check(err)

		fmt.Println(tree.LookupStr("8.8.8.8"))
		fmt.Println(tree.LookupStr("2001:7f8:1::64"))

		for _, path := range args {
			entries, err := goasn.RIBFromMRT(path)
			log.Println(len(entries), err)

			origins := make([]goasn.PrefixOrigin, len(entries))
			for i, entry := range entries {
				origins[i] = goasn.NewPrefixOrigin(entry)
			}

			b, err := goasn.ASNDatabase{origins}.MarshalText(singleAS)
			check(err)

			err = ioutil.WriteFile(path+".txt", b, os.FileMode(0644))
			check(err)

			b, err = ioutil.ReadFile(path + ".txt")
			check(err)

			db := goasn.ASNDatabase{}
			err = db.UnmarshalText(b)
			check(err)
			fmt.Println(len(db.Entries))

			tree, err := goasn.NewASNTreeFromFile(path + ".txt")
			check(err)
			fmt.Println(tree.LookupStr("2001:660:7302:5:6153:a38:4c2d:5b3b"))
			fmt.Println(tree.LookupStrMultiple([]string{
				"2001:660:7302:5:6153:a38:4c2d:5b3b",
				"2001:4860:4860::8888",
				"2405:9800:b000::1",
			}))
		}
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)
	convertCmd.Flags().String("format", "txt", "format: json or txt")
	convertCmd.Flags().Bool("single-as", false, "pyasn compatible format (only a single AS from AS-SET origin)")
}
