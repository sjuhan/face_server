package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	f "github.com/sjuhan/face_server/face"
	pb "github.com/sjuhan/face_server/proto"

	"github.com/Kagami/go-face"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

/*
var (
	m       = new(sync.Mutex)
	samples []face.Descriptor
	jumins  []string
	names   []string
	cats    []int32
)
*/
type vvv struct {
	v float32
	r int
}

//ByV 설명
type ByV []vvv

func (a ByV) Len() int           { return len(a) }
func (a ByV) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByV) Less(i, j int) bool { return a[i].v < a[j].v }

// server is used to implement helloworld.GreeterServer.
type server struct {
	fstruct

	samples []face.Descriptor
	jumins  []string
	names   []string
	cats    []int32

	m     *sync.Mutex
	jobch chan fstruct
}

func newserver() *server {
	server := new(server)
	server.m = new(sync.Mutex)
	server.jobch = make(chan fstruct)
	return server
}

type fstruct struct {
	face  face.Descriptor
	jumin string
	name  string
	cat   int32
	resv  chan fstruct
}

func (s *server) workers() {
	for {
		select {
		case in := <-s.jobch:
			if len(in.jumin) > 0 { //주민번호가 있을때-얼굴 저장
				s.m.Lock()
				s.samples = append(s.samples, s.face) //전송받은 얼굴특징값을 face.Descriptor로 변환하여 samples에 넣기
				s.jumins = append(s.jumins, in.jumin)
				s.names = append(s.names, in.name)
				s.m.Unlock()
				in.resv <- fstruct{jumin: "", name: ""}
				//res = &pb.Res{Jumin: "", Name: ""}

			} else { //주민번호 없을때-얼굴 구분
				n := f.Compare(s.samples, s.face, 0.6)
				if n == -1 {
					in.resv <- fstruct{jumin: "없는얼굴", name: "없는얼굴"}
					//res = &pb.Res{jumin: "없는얼굴", name: "없는얼굴"}
					//log.Printf("Received: %v", in.Face)
				} else {
					in.resv <- fstruct{jumin: s.jumins[n], name: s.names[n]}
					//res = &pb.Res{jumin: s.jumins[n], name: s.names[n]}
					//log.Printf("Received: %v", in.Face)
					//log.Printf("Received name:%v,jumin:%v,r:%v", res.Name, res.Jumin, n)
					//log.Printf("보낸주민:%v|보낸이름:%v\n", res.Jumin, res.Name)
				}
			}
		}
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) Recog(ctx context.Context, in *pb.Face) (*pb.Res, error) {
	var res *pb.Res
	resv := make(chan fstruct, 1)
	for i, f := range in.Descriptor_ {
		s.face[i] = f
	}

	struc := fstruct{face: s.face, jumin: in.Jumin, name: in.Name, resv: resv}

	s.jobch <- struc
	re := <-resv

	res = &pb.Res{Jumin: re.jumin, Name: re.name}
	/*
		if len(in.Jumin) > 0 { //주민번호가 있을때-얼굴 저장
			s.m.Lock()
			s.samples = append(s.samples, s.face) //전송받은 얼굴특징값을 face.Descriptor로 변환하여 samples에 넣기
			s.jumins = append(s.jumins, in.Jumin)
			s.names = append(s.names, in.Name)
			s.m.Unlock()
			res = &pb.Res{Jumin: "", Name: ""}

		} else { //주민번호 없을때-얼굴 구분
			n := f.Compare(s.samples, s.face, 0.6)
			if n == -1 {
				res = &pb.Res{Jumin: "없는얼굴", Name: "없는얼굴"}
				//log.Printf("Received: %v", in.Face)
			} else {
				res = &pb.Res{Jumin: s.jumins[n], Name: s.names[n]}
				//log.Printf("Received: %v", in.Face)
				//log.Printf("Received name:%v,jumin:%v,r:%v", res.Name, res.Jumin, n)
				//log.Printf("보낸주민:%v|보낸이름:%v\n", res.Jumin, res.Name)
			}
		}
		//fmt.Println(len(samples))
	*/
	return res, nil
}

/*
func (s *server) Test(in pb.Rec_TestServer) error {
	var err error
	return err
}
*/
func (s *server) setsamples(face face.Descriptor, cat int32) {

}

func 얼굴저장(face face.Descriptor) {
	file, err := os.OpenFile("face.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}

	//samples = append(samples, face)
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

func 얼굴비교(samples []face.Descriptor, comp face.Descriptor, tolerance float32, vv []vvv) int {
	res := FaceDistance(samples, comp)
	r := -1
	v := float32(1)
	for i, s := range res {
		t := EuclideanNorm(s)
		if t < 1 {
			//log.Println(t, "\n", s)
		}
		if t < tolerance && t < v {
			vv = append(vv, vvv{v: t, r: i})
			//v = t
			//r = i
		}
	}
	if len(vv) != 0 {
		sort.Sort(ByV(vv))
		r = vv[0].r
		//t := vv[0].v
		//fmt.Printf("%v|%v|r값:%v|t값:%v\n", time.Now(), len(vv), r, t)
	}
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
	server := newserver()

	fmt.Println("시작")
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRecServer(s, server)
	for i := 0; i < 10; i++ {
		go server.workers()
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
