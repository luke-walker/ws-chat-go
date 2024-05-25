package main

import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "log"
    "os"
    "time"

    "golang.org/x/net/websocket"
)

const ORIGIN = "http://127.0.0.1"

type Client struct {
    ws *websocket.Conn
    name string
}

func NewClient(ws *websocket.Conn, name string) *Client {
    return &Client{
        ws,
        name,
    }
}

func (c *Client) readLoop() {
    buf := make([]byte, 1024)

    for {
        n, err := c.ws.Read(buf)
        if err != nil && err != io.EOF {
            fmt.Println("Error:", err)
        }
        if n > 0 {
            fmt.Println(string(buf[:n]))
        }

        time.Sleep(500 * time.Millisecond)
    }
}

func (c *Client) writeLoop(scanner *bufio.Scanner) {
    for {
        msg := c.name + ": "

        scanner.Scan()
        msg += scanner.Text()
        c.ws.Write([]byte(msg))
    }
}

func main() {
    urlFlag := flag.String("url", "ws://127.0.0.1:5050", "specify the server url")
    flag.Parse()

    scanner := bufio.NewScanner(os.Stdin)

    fmt.Println("Connecting to server...")
    ws, err := websocket.Dial(*urlFlag, "", ORIGIN)
    if err != nil {
        log.Fatal("Fatal Error:", err)
    }
    fmt.Println("Connected to", *urlFlag)

    fmt.Print("Enter your name: ")
    scanner.Scan()
    name := scanner.Text()

    client := NewClient(ws, name)
    go client.readLoop()
    go client.writeLoop(scanner)

    select {} // surely there's a better way to do this
}
