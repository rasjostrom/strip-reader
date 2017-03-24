package structs

/* A chunk response contains the requested data and is decoded to JSON
before passing it to the client.
Chunks: A list of string chunks */
type ChunkResponse struct {
	Chunks []string
	Page int
	Max_Pages int

}

/* A chunk request is sent from the client in JSON and decoded into
a chunkRequest struct.
Path: Defines the path of the requested file to be read
Size: Defines the chunk size, as in number of words per chunk

Example JSON request: {"Path": "example_text.txt", "Size": 20"}*/
type ChunkRequest struct {
	Path string
	Size int
}
