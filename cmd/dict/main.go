package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	dbStore = map[string]interface{}{
		"mk2=じゅう":  []string{"銃かもしれない", "いやきっと銃に違いない"},
		"mk2=ああああ": []string{"伝説の勇者", "そういうアニメあったよね"},

		"user=taisuke.yamashita@topgate.co.jp": "やまちゃん",
	}

	cacheStore = map[string]interface{}{
		"mk2=じゅう": []string{"銃"},

		"user=taisuke.yamashita@topgate.co.jp": "これはユーザです",
	}
)

type Step interface {
	Do(fn func(interface{}), dictID, keyword string)
}

type StepFunc func(func(interface{}), string, string)

func (fn StepFunc) Do(resFn func(interface{}), dictID, keyword string) {
	fn(resFn, dictID, keyword)
}

type DatabaseStep struct{}

func (step *DatabaseStep) Do(fn func(interface{}), dictID, keyword string) { // dictID, keyword stringのあたりをどうするか？？ func(interface{}string)の形？
	log.Printf("dictID=%q keyword=%q で検索中...", dictID, keyword)
	v, _ := dbStore[dictID+"="+keyword]
	fn(v)
}

type CacheStep struct {
	Step Step
}

func (step *CacheStep) Do(fn func(interface{}), dictID, keyword string) {
	cacheKey := dictID + "=" + keyword
	if v, ok := cacheStore[cacheKey]; ok {
		log.Printf("dictID=%q keyword=%q でキャッシュにヒット！", dictID, keyword)
		fn(v)
		return
	}
	step.Step.Do(fn, dictID, keyword)
}

type UserStep struct{}

func (step *UserStep) Do(fn func(interface{}), dictID, keyword string) {
	log.Printf("userID=%q で検索中...", dictID)
	v, _ := dbStore["user="+dictID]
	fn(v)
}

func main() {
	useCache := flag.Bool("use-cache", false, "ON/OFF Switch")
	dictID := flag.String("dict-id", "", "dict-id")
	flag.Parse()
	if *dictID == "" {
		flag.Usage()
		return
	}
	keyword := flag.Arg(0)

	var step Step = &DatabaseStep{}
	if *useCache {
		step = &CacheStep{
			Step: step,
		}
	}
	step.Do(func(value interface{}) {
		fmt.Println(value)
	}, *dictID, keyword)

	step = &UserStep{}
	if *useCache {
		step = &CacheStep{
			Step: step,
		}
	}
	step.Do(func(value interface{}) {
		fmt.Println(value)
	}, *dictID, keyword)
}
