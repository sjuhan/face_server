package face

import (
	"bufio"
	"fmt"
	"time"

	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Kagami/go-face"
)

// Path to directory with models and test images. Here it's
// assumed it points to the
// <https://github.com/Kagami/go-face-testdata> clone.
const dataDir = "testdata"

var Face face.Descriptor
var samples []face.Descriptor
var cats []int32

type vvv struct {
	v float32
	r int
}

//ByV 설명
type ByV []vvv

func (a ByV) Len() int           { return len(a) }
func (a ByV) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByV) Less(i, j int) bool { return a[i].v < a[j].v }

// This example shows the basic usage of the package: create an
// recognizer, recognize faces, classify them using few known
// ones.
func main() {
	// Init the recognizer.
	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	/*
		dirlist := directoryList("./images")
		for _, dir := range dirlist {
			jpgload(rec, dir)
		}
	*/

	f, l := NewSampleFromFile("./data/sample")
	tt, _ := NewSampleFromFile("./data/test")

	for i, t := range tt {
		j := Compare(f, t, 0.33)
		if j == -1 {
			fmt.Println("없는얼굴", i)
		} else {
			fmt.Println(l[j])
		}
	}
	setsample(rec, f, l)
	for i, t := range tt {
		j := rec.ClassifyThreshold(t, 0.33)
		if j == -1 {
			fmt.Println("없는얼굴", i)
		} else {
			fmt.Println(l[j])
		}
	}
	//faces, labels := jpgload(rec, "../images/")
	//tt, ttt := jpgload(rec, "testdata")
	//faces, labels := 얼굴열기("test")
	//tt, _ := 얼굴열기("testdata")
	//fmt.Println("test 길이:", len(tt))
	//fmt.Println("faces 길이:", len(faces))
	//setsample(rec, faces, labels)
	/*  얼굴비교부분
	for i, f := range tt {
		//j := 얼굴비교(faces, f, 0.2)
		j := rec.Classify(f)
		//j := rec.ClassifyThreshold(f, 0.9)
		if j == -1 {
			fmt.Println("없는얼굴")
		} else {
			//fmt.Println(ttt[i], labels[j])
			fmt.Println("test", i, "faces", j)
		}
	}
	*/
	// Test image with 10 faces.
	/*
			 testImagePristin := filepath.Join(dataDir, "pristin.jpg")
			 // Recognize faces on that image.
			 //faces, err := rec.RecognizeFileCNN(testImagePristin)
			 faces, err := rec.RecognizeFile(testImagePristin)
			 if err != nil {
				 log.Fatalf("Can't recognize: %v", err)
			 }
			 if len(faces) != 10 {
				 log.Fatalf("Wrong number of faces, %v", len(faces))
			 }
			 // Fill known samples. In the real world you would use a lot of
			 // images for each person to get better classification results
			 // but in our example we just get them from one big image.
			 var samples []face.Descriptor
			 var cats []int32
			 for i, f := range faces {
				 samples = append(samples, f.Descriptor)
				 // Each face is unique on that image so goes to its own
				 // category.
				 cats = append(cats, int32(i))
			 }
		 // Name the categories, i.e. people on the image.
		 // Pass samples to the recognizer.
		 now := time.Now()
		 rec.SetSamples(samples, cats)
		 fmt.Println(time.Since(now))
		 // Now let's try to classify some not yet known image.
		 testImageNayoung := filepath.Join(dataDir, ".jpg")
		 nayoungFace, err := rec.RecognizeSingleFile(testImageNayoung)
		 //nayoungFace, err := rec.RecognizeSingleFileCNN(testImageNayoung)
		 if err != nil {
			 log.Fatalf("Can't recognize: %v", err)
		 }
		 if nayoungFace == nil {
			 log.Fatalf("Not a single face on the image")
		 }
		 //var catID int
		 //catID := rec.ClassifyThreshold(nayoungFace.Descriptor, 0.26970843)
		 catID := rec.ClassifyThreshold(nayoungFace.Descriptor, 0.26)
		 //catID := rec.Classify(nayoungFace.Descriptor)
		 if catID < 0 {
			 log.Fatalf("Can't classify")
		 }
		 // Finally print the classified label. It should be "Nayoung".
		 fmt.Println(labels[catID], catID)
	*/
}

//LoadJpg LoadJpg(*face.Recognizer, 불러올 경로) (얼굴특징값 배열, 얼굴이름 배열) 불러올 경로는 인물별 폴더가 있는 폴더로 설정 ex) /images/iu/000.jpg 인경우 images만 입력
func LoadJpg(rec *face.Recognizer, loc string) ([]face.Descriptor, []string) {
	var fname string
	var faces []face.Descriptor
	var labels []string

	files, err := ioutil.ReadDir(filepath.Join(loc))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".jpg") || strings.Contains(f.Name(), ".JPG") {
			fname = f.Name()
			fmt.Println(f.Name())
			ff := filepath.Join(loc, fname)
			rfaces, err := rec.RecognizeFile(ff)
			if err != nil {
				log.Println("얼굴인식에러", err)
			}
			if len(rfaces) == 1 {
				if strings.Count(ff, "/") > 1 { //
					for _, face := range rfaces {
						FaceToFile(face.Descriptor, "sample", ff[strings.Index(ff, "/")+1:strings.LastIndex(ff, "/")]) //1. 얼굴 특징값. 2. 얼굴값 저장할 파일명. 3. 얼굴값과 함께 저장할 이름
						faces = append(faces, face.Descriptor)
					}

				} else {
					for _, face := range rfaces {
						FaceToFile(face.Descriptor, ff, ff) //1. 얼굴 특징값. 2. 얼굴값 저장할 파일명. 3. 얼굴값과 함께 저장할 이름
						faces = append(faces, face.Descriptor)
					}
				}
			} else if len(rfaces) == -1 {
				fmt.Println("얼굴없음")
			} else {
				fmt.Println("얼굴여러개")
			}
		}

	}
	return faces, labels
}

func setsample(rec *face.Recognizer, faces []face.Descriptor, labels []string) {
	for i, f := range faces {
		samples = append(samples, f)
		// Each face is unique on that image so goes to its own
		// category.
		cats = append(cats, int32(i))
	}
	rec.SetSamples(samples, cats)
}

//FaceToFile FaceToFile(얼굴특징값,저장할경로 및 파일이름,얼굴값에 저장할 이름)
func FaceToFile(face face.Descriptor, loc string, name string) {
	file, err := os.OpenFile(fmt.Sprintf("%v.csv", loc), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	//samples = append(samples, face)
	for _, d := range face {
		s := fmt.Sprintf("%f", d) + ","
		_, err = file.WriteString(s)
		if err != nil {
			log.Println(err)
		}
	}
	file.WriteString(name + "\n")
}

//NewSampleFromFile NewSampleFromFile(파일경로) (얼굴특징값 배열, 얼굴주인이름 배열)
func NewSampleFromFile(loc string) ([]face.Descriptor, []string) {
	var ff [128]float32
	var no int
	var faces []face.Descriptor
	var labels []string

	file, err := os.OpenFile(fmt.Sprintf("%v.csv", loc), os.O_RDONLY|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		log.Println(err)
	}

	r := bufio.NewReader(file)
	for {
		d, _, err := r.ReadLine() //한줄씩 파일 읽기
		if err != nil {
			break
		}

		dd := strings.Split(string(d), ",") //실수 값 가져오기
		for i, f := range dd {
			if f == "" {
				break
			}
			//fmt.Println("f:", f)
			if i == 128 { //129번째 값은 labels에 append하기
				labels = append(labels, f)
			} else {
				r, err := strconv.ParseFloat(f, 32)
				if err != nil {
					log.Panic("float파싱에러", err)
				}
				ff[i] = float32(r)
			}
		}

		faces = append(faces, face.Descriptor(ff))

		no++
	}
	file.Close()
	log.Println(no, "개 읽음")
	return faces, labels
}

//Compare Compare(얼굴특징값 배열(소스이미지), 비교할 얼굴특징값, 비교한계)
func Compare(samples []face.Descriptor, comp face.Descriptor, tolerance float32) int {
	var vv []vvv
	res := Distance(samples, comp)
	r := -1
	v := float32(1)
	for i, s := range res {
		t := EuclideanNorm(s)
		//fmt.Println(t)
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
		t := vv[0].v
		fmt.Printf("%v|%v|r값:%v|t값:%v\n", time.Now(), len(vv), r, t)
	}
	//fmt.Println(vv)
	//r값은 얼굴특징값의 몇번째인지 반환.0부터 시작
	return r
}

// Distance fadd
func Distance(samples []face.Descriptor, comp face.Descriptor) []face.Descriptor {
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

func directoryList(dir string) []string {
	var dirlist []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		if f.IsDir() {
			dirlist = append(dirlist, filepath.Join(dir, f.Name()))
			//fmt.Println(f.Name() + "\t" + "디렉토리")
		} else {
			fmt.Println(f.Name() + "\t" + "파일" + "\t" + strconv.FormatInt(f.Size(), 10) + "Byte")
		}
	}

	return dirlist
}
