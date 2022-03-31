package main

import (
	"context"
	"fmt"
	"gb-go-best-practice/lesson-02/dirscan"
	"log"
	"os"
	"os/signal"
)

func main() {
	ds := dirscan.Create()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	fileList, err := ds.FindFiles(ctx, ".", -1, dirscan.ExtFilter(".csv", ".go"))
	if err != nil {
		log.Fatalln(err)
	}

	for _, f := range fileList {
		fmt.Printf("Name: %s\tPath: %s\n", f.Name(), f.Path())
	}
}
