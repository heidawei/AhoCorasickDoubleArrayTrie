package main

import (
	"fmt"

	"github.com/heidawei/AhoCorasickDoubleArrayTrie/ACDAT"
)

func main() {
	keyArray := []string{
		"hers",
		"his",
		"she",
		"he",
	}
	kvs := ACDAT.NewStringTreeMap()
	for _, key := range keyArray {
		kvs.Add(key, key)
	}
	acdat := ACDAT.NewAhoCorasickDoubleArrayTrie()
	acdat.Build(kvs)
	acdat.Dump()
	v := acdat.Get("his")
	fmt.Println("get ", v)
	hits := acdat.ParseText("u2342hers")
	for _, hit := range hits {
		fmt.Println(hit)
	}
}
