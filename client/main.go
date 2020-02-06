package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	pb "github.com/sjuhan/face_server/proto"

	"google.golang.org/grpc"
)

/*
캠 읽는부분 구현필요
*/
var address string

func df(now time.Time) {
	log.Println("끝", time.Since(now))
}
func main() {
	now := time.Now()
	defer df(now)

	ip := flag.String("ip", "127.0.0.1", "서버ip를 입력하세요")
	port := flag.String("port", "50051", "연결할 포트를 입력하세요")
	lensave := flag.Int("save", 1000, "저장테스트에 사용할 배열 갯수")
	lenrecog := flag.Int("recog", 1000, "얼굴인식테스트에 사용할 배열 갯수")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}
	address = fmt.Sprintf("%s:%s", *ip, *port)
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
	/*
		var s string
			fmt.Print("모드를 선택하세요 (test,real): ")
			fmt.Scan(&s)
			switch s {
			case "test":
				savetest(c, *lensave)
				if *lenrecog != 0 {
					recgtest(c, *lenrecog)
				} else {
					fmt.Println("무한!")
					looprecgtest(c, 1000)
				}
			case "real":
				//캠에서 얼굴 인식 부분 구현필요

				//구현 완료되면.. 얼굴값을 서버보내서 서버에서 연산후
				//있는 얼굴이면 index와 이름 등을 보냄
				//없는 얼굴이면 이름을 입력받음

			}
	*/
	savetest(c, *lensave)
	if *lenrecog != 0 {
		recgtest(c, *lenrecog)
	} else {
		fmt.Println("무한!")
		looprecgtest(c, 1000)
	}
}

func save(c pb.RecClient, res []float32, name string, jumin string) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	//rr, err := c.Recog(ctx, &pb.Face{Face: res, Name: name, Jumin: jumin})
	c.Recog(ctx, &pb.Face{Descriptor_: res, Name: name, Jumin: jumin})

	/* if err != nil {
		log.Fatalf("could not greet: %v", err)
	} */
	//log.Printf("%v,%v", rr.Jumin, rr.Name)
}

func recg(wg *sync.WaitGroup, ch chan []float32) {
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
	defer wg.Done()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		res, ok := <-ch
		if !ok {
			return
		}
		rr, err := c.Recog(ctx, &pb.Face{Descriptor_: res})
		if err != nil {
			//log.Fatalf("could not greet: %v", err)
			log.Printf("could not greet: %v", err)

		} else {
			if rr.Jumin == "없는얼굴" {
				//log.Printf("없음")
			} else {
				log.Printf("%v,%v", rr.Jumin, rr.Name)
			}
		}
	}
}

func realrecg(res []float32) {
	var name string
	var jumin string
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rr, err := c.Recog(ctx, &pb.Face{Descriptor_: res})
	if err != nil {
		//log.Fatalf("could not greet: %v", err)
		log.Printf("could not greet: %v", err)

	} else {
		if rr.Jumin == "없는얼굴" {
			fmt.Print("없는얼굴. 이름을 입력해주세요")
			fmt.Scan(&name)
			fmt.Print("없는얼굴. 주민번호를 입력해주세요")
			fmt.Scan(&jumin)
			save(c, res, name, jumin)
			//log.Printf("없음")
		} else {
			log.Printf("%v,%v", rr.Jumin, rr.Name)
		}
	}
}

func savetest(c pb.RecClient, i int) {
	var res []float32
	var ress [][]float32
	/*
		j6 := control.GetText("종합검진", "ThunderRT6TextBox34") //주민번호앞자리
		j7 := control.GetText("종합검진", "ThunderRT6TextBox33") //주민번호뒷자리
		jumin := j6 + j7
		if len(jumin) != 13 {
			fmt.Println(jumin)
		}
	*/
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	for j := 0; j < i; j++ {
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
		go save(c, res, name, jumin)
		//go recg(c, res)
		//log.Printf("\n이름:%v \n주민:%v", name, jumin)
	}
	ress = nil

}

func recgtest(c pb.RecClient, k int) {
	var res []float32
	var ress [][]float32
	var wg sync.WaitGroup

	s1 := rand.NewSource(1)
	r1 := rand.New(s1)
	s1 = rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s1)
	for j := 0; j < k; j++ {
		for i := 0; i < 128; i++ {
			res = append(res, (r1.Float32() + r2.Float32()/float32(r2.Intn(1000))))
		}
		ress = append(ress, res)
		res = nil
	}
	q := make(chan []float32, 30)
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go recg(&wg, q)
	}

	for i, res := range ress {
		//name := strconv.Itoa(i)
		//jumin := strconv.Itoa(i)
		//rname, rjumin := save(c, res, name, jumin)
		//go save(c, res, name, jumin)
		q <- res
		if i%150 == 0 {
			//time.Sleep(100 * time.Millisecond)
		}
	}
	close(q)
	wg.Wait() // 모든 goroutine 이 종료 될 때까지 기다린다
}

func looprecgtest(c pb.RecClient, k int) {
	var res []float32
	var ress [][]float32
	var wg sync.WaitGroup

	s1 := rand.NewSource(1)
	r1 := rand.New(s1)
	s1 = rand.NewSource(time.Now().UnixNano())
	r2 := rand.New(s1)
	for j := 0; j < k; j++ {
		for i := 0; i < 128; i++ {
			res = append(res, (r1.Float32() + r2.Float32()/float32(r2.Intn(1000))))
		}
		ress = append(ress, res)
		res = nil
	}
	q := make(chan []float32, 30)
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go recg(&wg, q)
	}
	for {
		for _, res := range ress {
			//name := strconv.Itoa(i)
			//jumin := strconv.Itoa(i)
			//rname, rjumin := save(c, res, name, jumin)
			//go save(c, res, name, jumin)
			q <- res
		}
	}
	close(q)
	wg.Wait() // 모든 goroutine 이 종료 될 때까지 기다린다
}
