package main

import (
	"flag"
	"fmt"
	"os"
)

type animal interface {
	cry()
}
type dog struct {
}
type cat struct {
}

func (d dog) cry() {
	fmt.Println("わん！")
}
func (c cat) cry() {
	fmt.Println("にゃー")
}

func main() {
	var name string
	flag.StringVar(&name, "animal", "", "動物名")
	flag.Parse()
	flag.Usage = func() {
		p := func(format string, args ...interface{}) {
			fmt.Fprintf(os.Stderr, format, args...)
		}
		p("Usage:\n")
		// p("  {prog} [options]\n")
		// p("\n")
		p("Available Options:\n")
		flag.PrintDefaults()
	}
	// operation := flag.Arg(0)
	// if operation == "" {
	// 	flag.Usage()
	// 	return
	// }

	switch name {
	case "dog":
		// dogの場合
		fmt.Println("犬の場合の処理を記載します")
	case "cat":
		// catの場合
		fmt.Println("ネコの場合の処理を記載します")
	default:
		fmt.Println("dogかcatで選択してください")
		flag.Usage()
	}
	dog := dog{}
	cat := cat{}
	dog.cry()
	cat.cry()
}
