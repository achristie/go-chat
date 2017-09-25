package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"

	"github.com/ardanlabs/kit/tcp"
)

type connHandler struct{}

func (connHandler) Bind(conn net.Conn) (io.Reader, io.Writer) {
	return bufio.NewReader(conn), bufio.NewWriter(conn)
}

type reqHandler struct{}

func (reqHandler) Read(ipAddress string, reader io.Reader) ([]byte, int, error) {
	bufReader := reader.(*bufio.Reader)

	line, err := bufReader.ReadString('\n')
	if err != nil {
		log.Printf("read : IP[ %s ] : %s", ipAddress, err)
		return nil, 0, err
	}
	log.Printf("read : IP[ %s ] : Length[%d] Data[%s]", ipAddress, len(line), line)
	return []byte(line), len(line), nil
}

func (reqHandler) Process(r *tcp.Request) {
	log.Printf("read : IP[ %s ] : %s\n", r.TCPAddr.IP.String(), string(r.Data))

	resp := tcp.Response{
		TCPAddr: r.TCPAddr,
		Data:    []byte("GOT IT\n"),
		Length:  7,
	}

	r.TCP.Send(context.TODO(), &resp)
}

type respHandler struct{}

func (respHandler) Write(r *tcp.Response, writer io.Writer) error {
	bufWriter := writer.(*bufio.Writer)
	if _, err := bufWriter.WriteString(string(r.Data)); err != nil {
		return err
	}
	return bufWriter.Flush()
}
