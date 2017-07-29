package main

import (
    "io/ioutil"
    "os"
    "fmt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: embedder resourceName < input > output\n %v", os.Args)
		os.Exit(1)
	}

    b, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(2)
	}
	
    fmt.Printf("package main \n\nvar (\n")
    fmt.Printf("\t%s = %#v\n", os.Args[1], b)
    fmt.Println(")\n")
}
