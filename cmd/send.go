package cmd

import (
	"fmt"
	"io"
	"net"
	"net/netip"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

func newSendCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "send <host:port> <file>",
		Short: "Send a file to the server",
		Args:  cobra.ExactArgs(2),
		Run: func(_ *cobra.Command, args []string) {
			Connect(args)
		},
	}
}

func Connect(args []string) {
	serverIP := args[0]
	fileName := args[1]
	fmt.Println(serverIP)

	addr, err := netip.ParseAddrPort(serverIP)
	if err != nil {
		fmt.Printf("Parse ip error: %v\n", err)
		return
	}

	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("File read error: %v\n", err)
		return
	}

	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		fmt.Printf("Server connection error: %v\n", err)
		return
	}
	defer conn.Close()

	header := fmt.Appendf(nil, "%s;%d;", fileName, len(file))
	if _, err := conn.Write(header); err != nil {
		fmt.Printf("File send error: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go Recv(&wg, conn)

	if _, err := conn.Write(file); err != nil {
		fmt.Printf("File send error: %v\n", err)
		return
	}
	wg.Wait()
}

func Recv(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()

	buff := make([]byte, 1)
	if _, err := io.ReadFull(conn, buff); err != nil {
		fmt.Printf("File send error: %v\n", err)
		return
	}

	switch buff[0] {
	case 0:
		fmt.Println("File delivered")
	case 1:
		fmt.Println("File send error")
	}
}
