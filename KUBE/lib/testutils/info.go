package testutils

import (
	"encoding/json"
	"fmt"
)

func PrintJSON(i interface{}) {
	result, _ := json.Marshal(i)
	fmt.Println(string(result))
}

func PrintJSONpretty(i interface{}) {
	result, _ := json.MarshalIndent(i, "", "\t")
	fmt.Println(string(result))
}
