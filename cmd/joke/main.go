package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/241m/jokes"
)

func main() {
	j := jokes.NewRequest()

	flag.IntVar(&j.Amount, "amount", 1, "Get `n` number of jokes.")
	flag.StringVar(&j.Contains, "contains", "", "Get jokes containing `text`")
	flag.BoolVar(&j.Safe, "safe", false, "Set safe-mode on")
	flag.Func("flag", "Add blacklist `flag`", j.Blacklist.Add)
	flag.Func("category", "Add category `cat`", j.Category.Add)
	flag.Func("lang", "Set languge to `lang`", j.Lang.Set)
	flag.Func("type", "Set type to `type`", j.Type.Set)
	flag.Parse()

	if jokes, e := j.Get(); e != nil {
		fmt.Println(e)
		os.Exit(1)
	} else {
		n := len(jokes)

		for i, j := range jokes {
			fmt.Println(j)

			if i < n-1 {
				fmt.Println("---")
			}
		}
	}
}
