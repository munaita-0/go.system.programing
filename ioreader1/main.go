package main

// 参考: http://ascii.jp/elem/000/001/252/1252961/

import (
	"archive/zip"
	"bufio"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	// stdin()
	// file()
	// netcon()
	// netconParsed()
	// fileCopy()
	// randomCopy()
	// createZip()
	zipDownloader()
}

func stdin() {
	for {
		buffer := make([]byte, 5)
		size, err := os.Stdin.Read(buffer)
		if err == io.EOF {
			fmt.Println("EOF")
			break
		}
		fmt.Println("size=%d input=%s", size, string(buffer))
	}
}

func file() {
	file, err := os.Open("file.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Readした内容をStdoutしてる
	io.Copy(os.Stdout, file)
}

func netcon() {
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("GET / HTTP/1.0\r\nHost: ascii.jp\r\n\r\n"))
	io.Copy(os.Stdout, conn)
}

func netconParsed() {
	conn, err := net.Dial("tcp", "ascii.jp:80")
	if err != nil {
		panic(err)
	}
	conn.Write([]byte("GET / HTTP/1.0\r\nHost: ascii.jp\r\n\r\n"))
	// http resonseがparseされる
	res, err := http.ReadResponse(bufio.NewReader(conn), nil)
	fmt.Println(res.Header)
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
}

func fileCopy() {
	file, err := os.Open("file.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	newFile, err := os.Create("new.txt")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Readした内容をStdoutしてる
	io.Copy(newFile, file)
}

func randomCopy() {
	file, err := os.Create("random.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	io.CopyN(file, rand.Reader, 1024)
}

func createZip() {
	file, err := os.Create("zip.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()

	awriter, err := zipWriter.Create("a.txt")
	io.Copy(awriter, strings.NewReader("first"))

	bwriter, err := zipWriter.Create("b.txt")
	io.Copy(bwriter, strings.NewReader("second"))
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=ascii_sample.zip")

	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()
	awriter, _ := zipWriter.Create("a.txt")
	io.Copy(awriter, strings.NewReader("first"))
}

func zipDownloader() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
