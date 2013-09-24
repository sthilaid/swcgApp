
package main

import (
        "swcg"
	//"fmt"
	//"encoding/json"
)

type toto struct {
	x int
	y int
}

func main() {
	swcg.AnalyzeDB(swcg.CreateDB())

	// bytes, e := json.Marshal(cards)
	// if e != nil {
	// 	fmt.Print(e.Error())
	// } else {
	// 	fmt.Print(string(bytes))
	// }
	// fmt.Print("\n")
}