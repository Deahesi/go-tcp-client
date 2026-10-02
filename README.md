# TCP Client

A client for Lab 2. It sends files to a TCP server and can generate a test file filled with zero bytes. The CLI uses Cobra.

## Build

Requires Go 1.27.1.

```bash
go build -o client .
```

## Usage

Send a file to the server:

```bash
./client send 127.0.0.1:8080 ./file.bin
```

The client waits for a response from the server. On success, it prints `File delivered`. The file is read entirely into memory before it is sent.

Generate a test file:

```bash
./client generate ./test.bin
```

This command writes zero bytes, prints the file size every second, and stops when you press `Ctrl+C`. It overwrites the file if it already exists.

Show command help:

```bash
./client --help
./client send --help
./client generate --help
```
