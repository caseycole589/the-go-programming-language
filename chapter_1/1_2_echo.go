
package main 

import (
	"os"
	"fmt"
)

func main() {
	for i, arg := range os.Args[0:] {
		fmt.Println(i,arg);
	}	
}