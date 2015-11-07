package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)


func main() {
	counts := make(map[string]int)
	//s_map maps each line to all the files it occurs in 
	s_map := make(map[string]string)
	files := os.Args[1:]
  	

	for _, arg := range files {
		fmt.Println(arg)
		f, err := os.Open(arg)
		if err != nil {;''
			fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
			continue
		}
		countLines(f, counts, s_map, arg)
		f.Close()
	}

	//print number of occurences files in occures in and the line
	//for each occurence that was seen more than once
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\t%s\n", n, s_map[line], line);
		}
	}
}

//does not need a return value because maps are passed my reference
//eg there variables are not copied
func countLines(f *os.File, counts map[string]int, s_map map[string]string, arg string) {
	
	// create a scanner object
	input := bufio.NewScanner(f)
	
	//while input
	for input.Scan() {
		counts[input.Text()]++
		//if there is more the one occurents of the current line
		if counts[input.Text()] > 1 {
			//if the current list of file names does not contain the current argument
			if !strings.Contains(s_map[input.Text()], arg) {
				//I know a join would be faster here
				s_map[input.Text()] += arg + " "
			}
		}
	}
}