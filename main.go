package main

import (
	"encoding/json"
	"fmt"
	"log"
	"go-jsondiff/jsonDiff"
)

func main()  {

	c1 := map[string]string{
		"c1ab": "c111",
		"c1cd": "c211",
	}
	b1 := map[string]interface{}{
		"hoby": "dance",
		"milk": "多多",
		"disanceng": c1,
	}
	a1 := map[string]interface{}{
		"xyff": "aichifan",
		"xfy": 123,
		"xyf": b1,
		"same": "same",
		"list": []string{"a", "b", "c"},
	}


	c2 := map[string]string{
		"c2ab": "c111",
		"c2cd": "c211",
	}
	b2 := map[string]interface{}{
		"hoby": "chifan",
		"milk": "不多",
		"disanceng": c2,
	}

	a2 := map[string]interface{}{
		"xyff": "aichifan",
		"xyf": b2,
		"hahaha": "笑点低",
		"same": "same",
		"list": []string{"a", "n", "c"},
	}

	l1, err := json.Marshal(a1)
	if err != nil {
		log.Fatalln("l1", err)
	}

	l2, err := json.Marshal(a2)
	if err != nil {
		log.Fatalln("l2", err)
	}

	text, _ := jsonDiff.Compare(l1, l2)
	fmt.Println("这是diff:", string(text))
}


