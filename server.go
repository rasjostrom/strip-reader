package main

import (
	"strings"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"encoding/json"
	"bytes"
	"strconv"
)


func main() {
	fmt.Println("StripReader Server v0.1")
	startListen()
}


func startListen() {
	listener, err := net.Listen("tcp", "localhost:3333")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on localhost:3333")

	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		go handleRequest(conn)
	}
}


/*
Accepts a request formatted as "PATH;CHUNK_SIZE;" where
PATH = path to a file on disk
CHUNK_SIZE = the number of words per chunk

The file is then read, split into chunks and sent back on the
socket as a JSON object.
*/
func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)

	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		conn.Write([]byte("Failed to read request, closing..."))
		conn.Close()
		return
	}

	req := strings.Split(string(buf), ";")
	path := req[0]
	chunkSize,_ := strconv.Atoi(req[1])

	fmt.Println("PATH:", path, "WORDS/CHUNK:", chunkSize)

	conn.Write(wordChunk(path, chunkSize))
	conn.Close()
}


func readTextFile(path string) string {
	file, err := ioutil.ReadFile(path)
	if err != nil {  // FIXME: deal with the error
		fmt.Println(err.Error())
		return ""
	}
	return string(file)
}


func wordChunk(path string, len int) []byte {
	text := readTextFile(path)
	chunkMap := make(map[string][]string)

	chunks := []string{}
	var tmpChunk bytes.Buffer

	for i, v := range strings.Split(text, " ") {
		if ((i % len == 0) && i != 0) {
			chunks = append(chunks, tmpChunk.String())
			tmpChunk.Reset()
		}
		tmpChunk.WriteString(v + " ")
	}
	chunks = append(chunks, tmpChunk.String()) // append remaining <len chunk

	chunkMap["chunks"] = chunks
	jsonResponse, err := json.Marshal(chunkMap)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return jsonResponse
}
