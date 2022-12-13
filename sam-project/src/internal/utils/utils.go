package utils

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(object any) {
	b, _ := json.MarshalIndent(object, "", "    ")

	fmt.Println(string(b))
}
