/*
Copyright © 2022 Gabe Joseph

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
	"os"

	"time"

	"github.com/spf13/cobra"

	speedclip "github.com/gjoseph92/speedclip/pkg"
)

var Start time.Duration
var End time.Duration

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "speedclip",
	Short: "Trim speedscope files",
	Long: `Speedclip crops speedscope files by timestamp.

Examples:

speedclip --start 33s --end 36.5s profile.json > clipped.json

speedclip --start 5m --end 10m profile.json clipped.json`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		out := "-"
		if len(args) == 2 {
			out = args[1]
		}
		if err := speedclip.Clip(path, out, Start, End); err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().DurationVarP(&Start, "start", "s", time.Second, "Start timestamp")
	rootCmd.Flags().DurationVarP(&End, "end", "e", time.Second, "End timestamp")
	rootCmd.MarkFlagRequired(("start"))
	rootCmd.MarkFlagRequired(("end"))
}
