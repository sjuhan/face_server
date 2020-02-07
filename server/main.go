package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	pb "github.com/sjuhan/face_server/proto"

	"github.com/Kagami/go-face"
	"google.golang.org/grpc"
)

const (
	port    = ":50051"
	dataDir = "./data"
)

var (
	c = make(chan chan string, 1)
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	*facestruct

	f              *pb.Face
	m              *sync.Mutex
	respool        sync.Pool
	ffpool         sync.Pool
	facestructpool sync.Pool
	recpool        sync.Pool
}

type facestruct struct {
	samples []face.Descriptor
	cats    []int32

	identify map[int32]identifypt

	cat int32
}

type identifypt struct {
	jumin string
	name  string
}

func newserverstruct() *server {
	server := new(server)

	server.facestruct = new(facestruct)

	server.samples = make([]face.Descriptor, 0, 1000000)
	server.cats = make([]int32, 0, 1000000)
	server.identify = make(map[int32]identifypt)

	server.m = new(sync.Mutex)

	server.respool = sync.Pool{
		New: func() interface{} { return new(pb.Res) },
	}
	server.ffpool = sync.Pool{
		New: func() interface{} { return [128]float32{} },
	}
	server.facestructpool = sync.Pool{
		New: func() interface{} { return new(facestruct) },
	}
	server.recpool = sync.Pool{
		New: func() interface{} {
			return func() *face.Recognizer {
				rec, _ := face.NewRecognizer(dataDir)
				return rec
			}
		},
	}

	return server
}

// SayHello implements helloworld.GreeterServer
func (s *server) Recog(ctx context.Context, in *pb.Face) (*pb.Res, error) {

	ff := s.ffpool.Get().([128]float32)
	res := s.respool.Get().(*pb.Res)
	rec := s.recpool.Get().(*face.Recognizer)
	s.facestruct = s.facestructpool.Get().(*facestruct)

	if len(in.Jumin) > 0 { //주민번호가 있을때-얼굴 저장
		for i, f := range in.Descriptor_ {
			ff[i] = f
		}
		s.m.Lock()

		s.facestruct.samples = append(s.facestruct.samples, face.Descriptor(ff)) //전송받은 얼굴특징값을 face.Descriptor로 변환하여 samples에 넣기
		s.facestruct.cat = int32(len(s.facestruct.samples) - 1)
		s.facestruct.cats = append(s.facestruct.cats, s.facestruct.cat)
		s.facestruct.identify[s.facestruct.cat] = identifypt{jumin: in.Jumin, name: in.Name}

		rec.SetSamples(s.facestruct.samples, s.cats)

		s.m.Unlock()

		res.Jumin = ""
		res.Name = ""
		ff = [128]float32{}

		s.ffpool.Put(ff)
		s.respool.Put(res)
		s.recpool.Put(rec)
		s.facestructpool.Put(s.facestruct)

	} else { //주민번호 없을때-얼굴 구분
		for i, f := range in.Descriptor_ {
			ff[i] = f
		}

		n := rec.Classify(face.Descriptor(ff))

		if n == -1 {
			res = &pb.Res{Jumin: "X", Name: "X"}
		} else {
			res = &pb.Res{Jumin: s.facestruct.identify[int32(n)].jumin, Name: s.facestruct.identify[int32(n)].name}
		}
	}
	//fmt.Println(len(samples))

	return res, nil
}

func (s *server) 얼굴저장(face face.Descriptor) {
	file, err := os.OpenFile("face.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	s.samples = append(s.samples, face)
	for _, d := range face {
		s := fmt.Sprintf("%0.9f", d) + "|"
		//fmt.Println(s)
		_, err = file.WriteString(s)
		if err != nil {
			log.Println(err)
		}
	}
	file.WriteString("\n")

	file.Close()
}

func (s *server) 얼굴열기() {
	var ff [128]float32
	var no int
	//var faced []face.Descriptor

	file, err := os.OpenFile("face.txt", os.O_RDONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		log.Println(err)
	}

	r := bufio.NewReader(file)
	for {
		d, _, err := r.ReadLine() //한줄씩 파일 읽기
		if err != nil {
			break
		}

		dd := strings.Split(string(d), "|") //실수 값 가져오기
		for i, f := range dd {
			if f == "" {
				break
			}
			//fmt.Println("f:", f)
			r, err := strconv.ParseFloat(f, 32)
			if err != nil {
				log.Println("float파싱에러", err)
			}
			ff[i] = float32(r)
		}

		s.samples = append(s.samples, face.Descriptor(ff))
		no++
	}
	file.Close()
	log.Println(no, "개 읽음")
}

func 얼굴비교(samples []face.Descriptor, comp face.Descriptor, tolerance float32) int {
	res := FaceDistance(samples, comp)
	r := -1
	v := float32(1)
	for i, s := range res {
		t := EuclideanNorm(s)
		if t < 1 {
			//log.Println(t, "\n", s)
		}
		if t < tolerance && t < v {
			//vv = append(vv, vvv{v: t, r: i})
			//v = t
			r = i
		}
	}
	/*
		if len(vv) != 0 {
			sort.Sort(ByV(vv))
			r = vv[0].r
			//t := vv[0].v
			//fmt.Printf("%v|%v|r값:%v|t값:%v\n", time.Now(), len(vv), r, t)
		}
	*/
	//r값은 얼굴특징값의 몇번째인지 반환.0부터 시작
	return r
}

// FaceDistance fadd
func FaceDistance(samples []face.Descriptor, comp face.Descriptor) []face.Descriptor {
	res := make([]face.Descriptor, len(samples))

	for i, s := range samples {
		for j := range s {
			res[i][j] = samples[i][j] - comp[j]
		}
	}

	return res
}

// EuclideanNorm ddfae
func EuclideanNorm(f face.Descriptor) float32 {
	var s float32
	for _, v := range f {
		s = s + v*v
	}

	return float32(math.Sqrt(float64(s)))
}

func main() {
	fmt.Println("시작")

	server := newserverstruct()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRecServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
