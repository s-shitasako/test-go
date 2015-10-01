package main

import (
  "fmt"
  "net"
  "os"
  "strings"
)

func main() {
  ln, err := net.Listen("tcp", ":4000")
  if err != nil {
    fmt.Println("Error:")
    fmt.Println(err)
  } else {
    for {
      conn, err := ln.Accept()
      if err != nil {
        fmt.Println("Error:")
        fmt.Println(err)
        os.Exit(1)
      }
      go serve(conn)
    }
  }
}

func serve(c net.Conn) {
  buf := make([]byte, 65536)
  l, err := c.Read(buf)
  if l > 0 && err == nil {
    s := "HTTP/1.1 200 OK\r\nContent-Length: 12\r\nContent-Type: text/html;charset=utf-8\r\nConnection: close\r\n\r\n<h1>あ</h1>"
    sr := strings.NewReader(s)
    buf2 := make([]byte, sr.Len())
    sr.Read(buf2)
    c.Write(buf2)
  }
  fmt.Println("Connection closed. with err:", c.Close())
}
