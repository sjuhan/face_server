package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"github.com/Kagami/go-face"
	pb "github.com/sjuhan/face_server/proto"
	"google.golang.org/grpc"
)

type facestruct struct {
	*pb.Face

	sendstruct pb.Face
	m          sync.Mutex
	samples    []face.Descriptor
	jumins     []string
	names      []string
}

func (s *facestruct) runRouteChat(client pb.RecClient) {
	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := client.Test(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv() //받는부분
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Println("받음", in.Descriptor_)
		}
	}()
	var test *pb.Face
	var d face.Descriptor
	s.samples = append(s.samples, d)
	s.samples = append(s.samples, d)
	s.samples = append(s.samples, d)
	test.Index = int32(len(s.samples))

	if err := stream.Send(test); err != nil { //보내는부분
		log.Fatalf("Failed to send a note: %v", err) //오류났던 부분
	}
	stream.CloseSend() //send 다 했을 때
	<-waitc
}

func main() {
	conn, err := grpc.Dial("juhan.tk:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewRecClient(conn)
	s := facestruct{}
	for {
		// RouteChat
		s.runRouteChat(client)
		time.Sleep(time.Second)
	}
}
