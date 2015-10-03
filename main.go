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
    fmt.Println("Error:", err)
  } else {
    for {
      conn, err := ln.Accept()
      if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
      }
      go serve(conn)
    }
  }
}

func serve(conn net.Conn) {
  defer func(){
    if err := conn.Close(); err != nil {
      fmt.Println("conn closed with err: ", err)
    }
  }()
  buf := make([]byte, 65536)
  l, err := conn.Read(buf)
  if l > 0 && err == nil {
    data := []byte{'<', 'h', '1', '>', 0xe3, 0x81, 0x84, '<', '/', 'h', '1', '>'}
    s := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html;charset=utf-8\r\nConnection: close\r\n\r\n", len(data))
    sr := strings.NewReader(s)
    buf2 := make([]byte, sr.Len())
    sr.Read(buf2)
    if _, err := conn.Write(buf2); err == nil {
      conn.Write(data)
    }
  }
}
