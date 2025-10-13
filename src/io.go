package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func readFromReader(r io.Reader) (string, error) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func writeToWriter(w io.Writer, data string) error {
	_, err := w.Write([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

func closeResource(c io.Closer) error {
	return c.Close()
}

func bufferExample() {
	buf := new(bytes.Buffer)
	writeToWriter(buf, "Hello, World!")
	data, err := readFromReader(buf)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

}

func fileExample() {
	file, err := os.Create("example.txt")
	if err != nil {
		panic(err)
	}
	defer closeResource(file)
	writeToWriter(file, "Hello, File World!")
	data, err := readFromReader(file)
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println(data)
}

func pipeExample() {
	reader, writer := io.Pipe()
	go func() {
		writeToWriter(writer, "Hello, Pipe World!")
		defer closeResource(writer)
	}()
	data, err := readFromReader(reader)
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Println(data)
}

func main() {
	bufferExample()
	fileExample()
	pipeExample()
}
