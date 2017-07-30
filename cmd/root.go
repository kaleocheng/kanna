// Copyright © 2017 Kaleo Cheng <kaleocheng@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kaleocheng/kanna/run"
	"github.com/kaleocheng/kanna/watcher"
	"github.com/spf13/cobra"
)

var (
	watchPath string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kanna COMMAND",
	Short: "Monitor for any changes in your project and automatically restart",
	Long: `Kanna watch the files in the directory in which Kanna was started or your specify, 
and if any files change, Kanna will automatically restart with your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("kanna --help")
			return
		}

		c := args[0]
		a := args[1:]

		if len(args) == 1 {
			cs := strings.Split(args[0], " ")
			c = cs[0]
			a = cs[1:]
		}

		start(c, a, watchPath, 1*time.Second)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().StringVarP(&watchPath, "path", "p", "./", "path which you want to wath")
}

func start(cmd string, args []string, path string, interval time.Duration) {
	watcher, err := watcher.NewRecursiveWatcher(path, interval)
	if err != nil {
		log.Fatal(err)
	}
	go watcher.Run()
	defer watcher.Close()

	run.Start(cmd, args)
	for {
		select {
		case <-watcher.Files:
			run.Restart(cmd, args)
		}
	}
}
