
package main

import (
        "swcg"
	"fmt"
	"encoding/json"
)

func main() {
	bytes, e := json.Marshal(swcg.CreateDB())
	if e != nil {
		fmt.Print(e.Error())
	} else {
		fmt.Print(string(bytes))
	}
}