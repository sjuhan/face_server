package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/sjuhan/face_server/proto"

	"google.golang.org/grpc"
)

const (
	address0 = "localhost:50051"
)

func main() {
	now := time.Now()
	var wg sync.WaitGroup
	var wg1 sync.WaitGroup
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

	s1 := rand.NewSource(1)
	r1 := rand.New(s1)
	for j := 0; j < 50000; j++ {
		for i := 0; i < 128; i++ {
			res = append(res, r1.Float32())
		}
		ress = append(ress, res)
		res = nil
	}
	wg.Add(len(ress))
	for i, res := range ress {
		name := strconv.Itoa(i)
		jumin := strconv.Itoa(i)
		//rname, rjumin := save(c, res, name, jumin)
		go save(c, res, name, jumin, &wg)
		//go recg(c, res)
		//log.Printf("\n이름:%v \n주민:%v", name, jumin)
	}
	wg.Wait()
	ress = nil
	s1 = rand.NewSource(1)
	r1 = rand.New(s1)
	for j := 0; j < 500; j++ {
		for i := 0; i < 128; i++ {
			res = append(res, r1.Float32())
		}
		ress = append(ress, res)
		res = nil
	}
	wg1.Add(len(ress))
	for _, res := range ress {
		//name := strconv.Itoa(i)
		//jumin := strconv.Itoa(i)
		//rname, rjumin := save(c, res, name, jumin)
		//go save(c, res, name, jumin)
		go recg(c, res, &wg1)
	}
	wg1.Wait()

	log.Println("끝", time.Since(now))
}

func save(c pb.RecClient, res []float32, name string, jumin string, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	//rr, err := c.Recog(ctx, &pb.Face{Face: res, Name: name, Jumin: jumin})
	c.Recog(ctx, &pb.Face{Face: res, Name: name, Jumin: jumin})
	/* if err != nil {
		log.Fatalf("could not greet: %v", err)
	} */
	//log.Printf("%v,%v", rr.Jumin, rr.Name)
	defer wg.Done()
}

func recg(c pb.RecClient, res []float32, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	rr, err := c.Recog(ctx, &pb.Face{Face: res})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if rr.Jumin != "없는얼굴" {
		log.Printf("%v,%v", rr.Jumin, rr.Name)
	} else {
		log.Printf("없음")
	}
	defer wg.Done()
}
