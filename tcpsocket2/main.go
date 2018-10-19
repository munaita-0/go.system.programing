package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// go server()
	// gzipClient()
	// go chunkServer()
	// chunkClient()
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}

// 青空文庫: ごんぎつねより
// http://www.aozora.gr.jp/cards/000121/card628.html
var contents = []string{
	"これは、私わたしが小さいときに、村の茂平もへいというおじいさんからきいたお話です。",
	"むかしは、私たちの村のちかくの、中山なかやまというところに小さなお城があって、",
	"中山さまというおとのさまが、おられたそうです。",
	"その中山から、少しはなれた山の中に、「ごん狐ぎつね」という狐がいました。",
	"ごんは、一人ひとりぼっちの小狐で、しだの一ぱいしげった森の中に穴をほって住んでいました。",
	"そして、夜でも昼でも、あたりの村へ出てきて、いたずらばかりしました。",
}

func chunkProsessSession(conn net.Conn) {
	fmt.Printf("Accept: %v", conn.RemoteAddr)
	defer conn.Close()

	for {
		request, err := http.ReadRequest(bufio.NewReader(conn))

		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		fmt.Fprintf(conn, strings.Join([]string{
			"HTTP/1.1 200 OK",
			"Content-Type: text/plain",
			"Transfer-Encoding: chunked",
			"", "",
		}, "\r\n"))

		for _, content := range contents {
			bytes := []byte(content)
			fmt.Fprintf(conn, "%x\r\n%s\r\n", len(bytes), content)
			fmt.Println("==endChunk waiting..==")
			time.Sleep(2 * time.Second)
			fmt.Println("==endChunk end.==")
		}
		fmt.Fprint(conn, "0\r\n%s\r\n")
		fmt.Println("==endServer==")
	}
}

func chunkServer() {
	listener, err := net.Listen("tcp", "localhost:8888")
	e(err)
	fmt.Print("running...")

	for {
		conn, err := listener.Accept()
		e(err)
		go chunkProsessSession(conn)
	}
}

func chunkClient() {
	conn, err := net.Dial("tcp", "localhost:8888")
	e(err)

	request, err := http.NewRequest("GET", "http://localhost:8888", nil)
	e(err)

	err = request.Write(conn)
	e(err)

	reader := bufio.NewReader(conn)
	response, err := http.ReadResponse(reader, request)
	e(err)

	dump, err := httputil.DumpResponse(response, false)
	e(err)

	fmt.Println(dump)

	if len(response.TransferEncoding) < 1 || response.TransferEncoding[0] != "chunked" {
		panic("wrong transfer encoding")
	}

	for {
		sizeStr, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}

		size, err := strconv.ParseInt(string(sizeStr[:len(sizeStr)-2]), 16, 64)

		if size == 0 {
			break
		}

		if err != nil {
			panic(err)
		}

		line := make([]byte, int(size))
		reader.Read(line)
		reader.Discard(2)
		fmt.Printf("  %d bytes: %s\n", size, string(line))
	}
}

func isGZipAcceptable(request *http.Request) bool {
	return strings.Index(strings.Join(request.Header["Accept-Encoding"], ","), "gzip") != -1
}

func prosessSession(conn net.Conn) {
	fmt.Printf("Accept: %v", conn.RemoteAddr)
	defer conn.Close()

	for {
		conn.SetDeadline(time.Now().Add(5 * time.Second))
		request, err := http.ReadRequest(bufio.NewReader(conn))

		if err != nil {
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				fmt.Println("Timeout")
				break
			} else if err == io.EOF {
				break
			}
			panic(err)
		}

		dump, err := httputil.DumpRequest(request, true)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(dump))

		response := http.Response{
			StatusCode: 200,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
		}

		if isGZipAcceptable(request) {
			content := "Hello World(gzipped)\n"
			// bufferに圧縮contentを書き込む
			var buffer bytes.Buffer
			writer := gzip.NewWriter(&buffer)
			io.WriteString(writer, content)
			writer.Close()
			// bufferの内容をresonse.Bodyに
			response.Body = ioutil.NopCloser(&buffer)
			response.ContentLength = int64(buffer.Len())
			response.Header.Set("Content-Encoding", "gzip")
		} else {
			content := "Hello World\n"
			response.Body = ioutil.NopCloser(strings.NewReader(content))
			response.ContentLength = int64(len(content))
		}

		response.Write(conn)
		fmt.Print("===========ServerEnd=========")
	}
}

func server() {
	listener, err := net.Listen("tcp", "localhost:8888")
	e(err)
	fmt.Print("running...")

	for {
		conn, err := listener.Accept()
		e(err)
		go prosessSession(conn)
	}
}

func gzipClient() {
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
		request.Header.Set("Accept-Encoding", "gzip")

		err = request.Write(conn)
		e(err)

		response, err := http.ReadResponse(bufio.NewReader(conn), request)

		if err != nil {
			fmt.Println("Retry")
			conn = nil
			continue
		}

		dump, err := httputil.DumpResponse(response, false)
		e(err)
		fmt.Println(string(dump))
		defer response.Body.Close()

		if response.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(response.Body)
			e(err)
			io.Copy(os.Stdout, reader)
			reader.Close()
		} else {
			io.Copy(os.Stdout, response.Body)
		}
		fmt.Print("===========ClientEnd=========")

		current++
		if current == len(sendMessages) {
			break
		}
	}
}
