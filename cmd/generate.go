package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func newGenerateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate <file>",
		Short: "Generate a file filled with zero bytes until interrupted",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			Generate(args)
		},
	}
}

func Generate(args []string) {
	if args[0] == "" {
		fmt.Println("Invalid args")
		return
	}
	fileName := args[0]

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("File create fail: %v\n", err)
		return
	}
	defer file.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	generator := Generator{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		generator.generateLoop(file, ctx)
	}()
	go func() {
		defer wg.Done()
		generator.logLoop(ctx)
	}()

	<-ctx.Done()
	wg.Wait()
}

type Generator struct {
	SizeBytes int
	mu        sync.Mutex
}

func (g *Generator) generateLoop(file *os.File, ctx context.Context) {
	for {
		select {
		default:
			text := make([]byte, 1000)
			file.Write(text)
			g.mu.Lock()
			g.SizeBytes += len(text)
			g.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (g *Generator) logLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			g.mu.Lock()
			fmt.Printf("Final file size: %dMB\n", g.SizeBytes/1000000)
			g.mu.Unlock()
			return
		case <-ticker.C:
			g.mu.Lock()
			fmt.Printf("File size: %dMB\n", g.SizeBytes/1000000)
			g.mu.Unlock()
		}
	}
}
