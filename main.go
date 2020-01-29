package main

import (
	"github.com/DmitriiTrifonov/gost-ciphers/kuznechik"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		log.Fatal("Not Implemented")
	} else {
		kuznechik.SelfCheck()
	}

}
