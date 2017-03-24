package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strings"

	"./structs"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        return true
    },
}


func main() {
	setupWebSocket()
	fmt.Println("Listening on localhost:8000")
	http.ListenAndServe(":8000", nil)
}

func setupWebSocket() {
	http.HandleFunc("/websocket", func(writer http.ResponseWriter, req *http.Request) {

		// Upgrade the connection
		conn, err := upgrader.Upgrade(writer, req, nil)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Client subscribed")

		// Read requested file
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		var reqMsg structs.ChunkRequest
		err = json.Unmarshal(msg, &reqMsg)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Chunk and write back the contents
		chunkResponse := structs.ChunkResponse{
			Chunks: wordChunk(readTxt(reqMsg.Path), reqMsg.Size),
			Page: 1,
			Max_Pages: 1,
		}
		jsonMsg, err := json.Marshal(chunkResponse)
		if err != nil {
			fmt.Println(err)
			return
		}
		conn.WriteMessage(msgType, jsonMsg)
		conn.Close()
		fmt.Println("Client Unsubscribed")
	})
}

func readTxt(path string) string {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err.Error()
	}
	return string(file)
}

/* Chunks input string into 'len' number of words.
Returns a list of strings.
*/
func wordChunk(text string, len int) []string {
	chunks := []string{}

	var tmpChunk bytes.Buffer
	for i, value := range strings.Split(text, " ") {
		if ((i % len == 0) && i != 0) {
			chunks = append(chunks, tmpChunk.String())
			tmpChunk.Reset()
		}
		tmpChunk.WriteString(value + " ")
	}
	chunks = append(chunks, tmpChunk.String()) // append remaining <len chunk

	return chunks
}
