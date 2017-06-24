package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func Sum(input interface{}) int {
	switch v := input.(type) {
	case string:
		return 0
	case float64:
		return int(v)
	case []interface{}:
		var sum int
		for _, sub := range v {
			sum += Sum(sub)
		}
		return sum
	case map[string]interface{}:
		var sum int
		for _, sub := range v {
			if sub == "red" {
				return 0
			}
			sum += Sum(sub)
		}
		return sum
	default:
		log.Fatalf("Unknown type %T %#v", v, v)
		panic("unreachable")
	}
}

func main() {
	js, err := ioutil.ReadFile("input.txt")
	if err != nil {
		log.Fatalf("failed to open input: %s", err)
	}

	var input interface{}
	if err := json.Unmarshal(js, &input); err != nil {
		log.Fatalf("failed to parse input: %s", err)
	}

	fmt.Println(Sum(input))
}
