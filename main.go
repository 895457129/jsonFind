package main

import (
	"flag"
	"fmt"
	"github.com/atotto/clipboard"
	"os"
)

var (
	key string
	num int
)

func init()  {
	flag.IntVar(&num, "n", 20, "最大匹配条数，默认20")
	flag.StringVar(&key, "k", "", "要查找的关键字")
}

func main()  {
	content, err := clipboard.ReadAll()
	flag.Parse()
	if err != nil {
		panic(err)
	}
	if len(content) == 0 {
		fmt.Println("粘贴板没有内容")
		os.Exit(0)
	}
	if len(key) == 0 {
		fmt.Println("请输入要查找的关键字")
		os.Exit(0)
	}
	FindPath(content, key, num)
}


