package main

import (
  "bytes"
  "fmt"
  "net"
  "os"
  "strconv"
  "strings"
)

var srvRoot string = "."

func main() {
  port, root := loadArgs(os.Args)
  fmt.Printf("Listen port:%d, serve directory: %s\n", port, root)
  srvRoot = root

  ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
  if err != nil {
    fmt.Println("conn read error: ", err)
    return
  }
  fmt.Println(bytes.NewBuffer(buf).String())
  path := parsePath(buf[:l])
  if path == "" {
    dat400 := "HTTP/1.1 400 Bad Request\r\nContent-Type: text/html\r\nContent-Length: 12\r\n\r\n<h1>400</h1>"
    strings.NewReader(dat400).WriteTo(conn)
    return
  }
  data := loadFileData(srvRoot + path)
  if data == nil {
    dat404 := "HTTP/1.1 404 Not Found\r\nContent-Type: text/html\r\nContent-Length: 12\r\n\r\n<h1>404</h1>"
    strings.NewReader(dat404).WriteTo(conn)
    return
  }
  m := mime(path)
  fmt.Println(m)
  s := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: %s\r\nConnection: close\r\n\r\n", len(data), m)
  if _, err := strings.NewReader(s).WriteTo(conn); err == nil {
    conn.Write(data)
  }
}

func loadArgs(args []string) (int, string) {
  var err error
  port := 4000
  root := "."
  switch len(args) {
  case 3:
    port, err = strconv.Atoi(args[2])
    if port <= 0 || err != nil {
      port = 4000
      fmt.Println("port is bad format. ", err)
      os.Exit(1)
    }
    fallthrough
  case 2:
    root = args[1]
    if root[len(root)-1] == '/' {
      root = root[:len(root)-1]
    }
  default:
    fmt.Println("usage: test-go port [ root-dir ]")
    fmt.Println("e.g. test-go 8080 ~/root")
  }
  return port, root
}

func parsePath(dat []byte) (ret string) {
  l := len(dat)
  if l < 7 || dat[0] != 'G' || dat[1] != 'E' || dat[2] != 'T' || dat[3] != ' ' || dat[4] != '/' {
    return
  }
  i := 5
  state := 0
  for ;; i++ {
    r := dat[i]
    if i >= l {
      return
    } else if r == ' ' {
      break
    } else if r == '\r' || r == '\n' {
      return
    } else if state == 0 && r == '/' {
      state = 1
    } else if state == 1 {
      if r == '.' || r == '/' {
        return
      } else {
        state = 0
      }
    }
  }
  if ret = bytes.NewBuffer(dat[4:i]).String(); ret == "" {
    return
  }
  if ret[len(ret)-1] == '/' {
    ret = ret + "index.html"
  }
  return
}

func loadFileData(name string) (ret []byte) {
  f, err := os.Open(name)
  if err != nil {
    return
  }
  var l int64
  if info, err := f.Stat(); err != nil {
    return
  } else {
    l = info.Size()
  }
  dat := make([]byte, l)
  var i int64 = 0
  for {
    n, err := f.ReadAt(dat, i)
    i += int64(n)
    if i >= l || err != nil {
      break
    }
  }
  ret = dat
  return
}

func mime(path string) string {
  for i := len(path) - 1; i > 0; i-- {
    switch path[i] {
    case '.':
      return ext2mime(path[i+1:])
    case '/':
      break
    }
  }
  return "application/octet-stream"
}

func ext2mime(ext string) string {
  switch ext {
  case "html", "htm":
    return "text/html"
  case "txt", "log":
    return "text/plain"
  case "jpg", "jpeg":
    return "image/jpeg"
  case "png":
    return "image/png"
  case "gif":
    return "image/gif"
  case "pdf":
    return "application/pdf"
  case "js":
    return "text/javascript"
  case "css":
    return "text/css"
  case "json":
    return "application/json"
  }
  return "application/octet-stream"
}
