package main

// 参考: http://ascii.jp/elem/000/001/276/1276572/

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

func main() {
	// go server()
	go keepaliveServer()
	// client()
	keepaliveClient()
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}

// TCPでの実装
// HTTPでの実装は実際もっと簡易にできる http://ascii.jp/elem/000/001/243/1243667/
// http.HandleFuncやhttp.ListenAndServeなどが一通りやってくれる
func server() {
	fmt.Println("==inserver==")
	listener, err := net.Listen("tcp", "localhost:8888")
	e(err)

	fmt.Println("running..")

	for {
		conn, err := listener.Accept()
		e(err)

		go func() {
			fmt.Print("accept %v", conn.RemoteAddr())

			request, err := http.ReadRequest(bufio.NewReader(conn))
			e(err)

			dump, err := httputil.DumpRequest(request, true)
			e(err)

			fmt.Println(string(dump))

			response := http.Response{
				StatusCode: 200,
				ProtoMajor: 1,
				ProtoMinor: 0,
				Body:       ioutil.NopCloser(strings.NewReader("Hello World")),
			}
			response.Write(conn)
			conn.Close()
		}()
	}
}

func keepaliveServer() {
	fmt.Println("==inKeepAviceServer==")
	listener, err := net.Listen("tcp", "localhost:8888")
	e(err)

	fmt.Println("running..")

	for {
		conn, err := listener.Accept()
		e(err)

		go func() {
			fmt.Print("accept %v", conn.RemoteAddr())

			for {
				conn.SetDeadline(time.Now().Add(5 * time.Second))

				// 次のリクエストが来るのを待つ
				// timeoutした場合はerrを返す
				// 5秒以内のリクエストであれば再接続しない分オーバーヘッド減る
				request, err := http.ReadRequest(bufio.NewReader(conn))

				if err != nil {
					neterr, ok := err.(net.Error)
					if ok && neterr.Timeout() {
						fmt.Println("Timeout")
						break
					} else if err == io.EOF {
						fmt.Println("???")
						break
					}
					panic(err)
				}

				dump, err := httputil.DumpRequest(request, true)
				if err != nil {
					panic(err)
				}
				fmt.Println("SSSSSSS" + string(dump))
				content := "Hello World"

				response := http.Response{
					StatusCode:    200,
					ProtoMajor:    1,
					ProtoMinor:    1,
					ContentLength: int64(len(content)),
					Body:          ioutil.NopCloser(strings.NewReader(content)),
				}
				response.Write(conn)
			}
			conn.Close()
		}()
	}
}

// TCPでの実装
func client() {
	fmt.Println("==inclient==")
	conn, err := net.Dial("tcp", "localhost:8888")
	e(err)

	request, err := http.NewRequest("GET", "http://localhost:8888", nil)
	e(err)

	request.Write(conn)
	response, err := http.ReadResponse(bufio.NewReader(conn), request)
	e(err)

	dump, err := httputil.DumpResponse(response, true)
	e(err)

	fmt.Println(string(dump))
}

// TCPでの実装
func keepaliveClient() {
	sendMessages := []string{
		"ASCII",
		"PROGRAMMING",
		"PLUS",
	}

	current := 0
	var conn net.Conn = nil

	for {
		var err error

		if conn == nil {
			conn, err = net.Dial("tcp", "localhost:8888")
			e(err)
			fmt.Printf("Access: %d", current)
		}

		request, err := http.NewRequest("POST", "http://localhost:8888", strings.NewReader(sendMessages[current]))
		e(err)

		err = request.Write(conn)
		e(err)

		response, err := http.ReadResponse(bufio.NewReader(conn), request)

		if err != nil {
			fmt.Println("Retry")
			conn = nil
			continue
		}

		dump, err := httputil.DumpResponse(response, true)
		e(err)

		fmt.Println("cccccc" + string(dump))
		current++
		if current == len(sendMessages) {
			break
		}
	}
}
