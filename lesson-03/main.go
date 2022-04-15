package main

import (
	"context"
	"fmt"
	"gb-go-best-practices/lesson-03/dirscan"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"go.uber.org/zap"
)

func main() {
	// logger, err := zap.NewProduction()
	logger, err := zap.NewDevelopment()

	if err != nil {
		log.Fatalln(err)
	}

	defer logger.Sync()

	logger.Info("Process ID:", zap.Int("pid", os.Getpid()))

	dirPath := "."

	DS, err := dirscan.Create(dirPath, dirscan.Recursive, time.Second*0, logger)
	if err != nil {
		logger.Fatal("Can't read directory:", zap.String("path", dirPath), zap.Error(err))
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := DS.Find(ctx, dirscan.ExtFilter(".go", ".csv")); err != nil {
		logger.Error("File search error: ", zap.Error(err))
	} else {
		result := new(strings.Builder)
		for _, f := range DS.Result() {
			fmt.Fprintf(result, "path: %s ->\tname: %s\n", f.Path(), f.Name())
		}
		fmt.Fprintln(os.Stdout, "\nScan results:\n-------------")
		fmt.Fprintln(os.Stdout, result)
	}
}
