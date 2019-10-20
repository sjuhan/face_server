package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/Kagami/go-face"
	pb "github.com/sjuhan/face_server/proto"
	"google.golang.org/grpc"
)

type server struct {
	faceStruct

	faces *pb.Face

	mu sync.Mutex

	setsamplesChan chan face.Descriptor
}

type faceStruct struct {
	samples []face.Descriptor
	jumins  []string
	names   []string
}

func (s *server) StreamTest(stream pb.Rec_StreamTestServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if in.GetHash() == "" {
			continue
		} else if in.GetHash() == makehash(s.samples) {

		} else {
			sendstruct := &pb.Face{Jumin: ""}
			if err := stream.Send(sendstruct); err != nil {
				return err
			}
		}

	}
}

func (s *server) Recog(context context.Context, in *pb.Face) (out *pb.Face, err error) {
	return
}

func (s *server) settamples(face1 *pb.Face) face.Descriptor {
	var vFace [128]float32
	for i, f := range face1.FaceDescriptor {
		vFace[i] = f
	}
	v := s.faces
	v.GetFaceDescriptor()
	s.mu.Lock()
	s.samples = append(s.samples, face.Descriptor(vFace)) //전송받은 얼굴특징값을 face.Descriptor로 변환하여 samples에 넣기
	s.jumins = append(s.jumins, face1.Jumin)
	s.names = append(s.names, face1.Name)
	s.mu.Unlock()
	return face.Descriptor(vFace)
}

func makehash(in interface{}) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(fmt.Sprintf("%v", in))))
}

func main() {
	fmt.Println("시작")
	lis, err := net.Listen("tcp", "50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRecServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
