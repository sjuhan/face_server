package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/sjuhan/face_server/proto"

	"google.golang.org/grpc"
)

const (
	address0 = "localhost:50051"
)

func main() {
	now := time.Now()

	done := make(chan bool, 5)
	done1 := make(chan bool, 5)

	var res []float32
	var ress [][]float32

	address := address0
	if len(os.Args) > 1 {
		if strings.Contains(os.Args[1], ":") {
			address = os.Args[1]
		} else {
			address = os.Args[1] + ":50051"
		}
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		/*
				if err!=nil{
				aa,err:=fmt.Scanf()
				log.Printf()
			}
		*/
	}
	defer conn.Close()
	c := pb.NewRecClient(conn)

	// Contact the server and print out its response.
	//name := defaultName
	address = address0
	if len(os.Args) > 1 {
		address = os.Args[1]
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for j := 0; j < 10; j++ {
		for i := 0; i < 128; i++ {
			res = append(res, r1.Float32())
		}
		ress = append(ress, res)
		res = nil
	}

	for i, res := range ress {
		name := strconv.Itoa(i)
		jumin := strconv.Itoa(i)
		//rname, rjumin := save(c, res, name, jumin)
		go save(c, res, name, jumin, done)
		<-done
		//go recg(c, res)
		//log.Printf("\n이름:%v \n주민:%v", name, jumin)
	}
	ress = nil
	s1 = rand.NewSource(1)
	r1 = rand.New(s1)
	s1 = rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s1)
	for j := 0; j < 500000; j++ {
		for i := 0; i < 128; i++ {
			res = append(res, (r1.Float32() + r2.Float32()/float32(r2.Intn(1000))))
		}
		ress = append(ress, res)
		res = nil
	}

	for _, res := range ress {
		//name := strconv.Itoa(i)
		//jumin := strconv.Itoa(i)
		//rname, rjumin := save(c, res, name, jumin)
		//go save(c, res, name, jumin)
		go recg(c, res, done1)
		<-done1
	}

	log.Println("끝", time.Since(now))
}

func save(c pb.RecClient, res []float32, name string, jumin string, ch chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	//rr, err := c.Recog(ctx, &pb.Face{Face: res, Name: name, Jumin: jumin})
	c.Recog(ctx, &pb.Face{Face: res, Name: name, Jumin: jumin})
	/* if err != nil {
		log.Fatalf("could not greet: %v", err)
	} */
	//log.Printf("%v,%v", rr.Jumin, rr.Name)
	ch <- true
}

func recg(c pb.RecClient, res []float32, ch chan bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	rr, err := c.Recog(ctx, &pb.Face{Face: res})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if rr.Jumin == "없는얼굴" {
		log.Printf("없음")
	} else {
		log.Printf("%v,%v", rr.Jumin, rr.Name)
	}
	ch <- true
}
