package kuznechik

import (
	"fmt"
	"github.com/DmitriiTrifonov/gost-ciphers/util"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

const parts = 10

const filesize = 12500000

const block = 64


var key = [0x20]byte{
	0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, // 00..07
	0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, // 08..0F
	0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10, // 10..17
	0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, // 18..1F
}

var plainTest = []byte{
	0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
	0xFF, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
	0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
	0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xEE, 0xFF, 0x0A,
	0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
	0x99, 0xAA, 0xBB, 0xCC, 0xEE, 0xFF, 0x0A, 0x00,
	0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99,
	0xAA, 0xBB, 0xCC, 0xEE, 0xFF, 0x0A, 0x00, 0x11,
}

var cipherTest = []byte{
	0x7F, 0x67, 0x9D, 0x90, 0xBE, 0xBC, 0x24, 0x30,
	0x5A, 0x46, 0x8D, 0x42, 0xB9, 0xD4, 0xED, 0xCD,
	0xB4, 0x29, 0x91, 0x2C, 0x6E, 0x00, 0x32, 0xF9,
	0x28, 0x54, 0x52, 0xD7, 0x67, 0x18, 0xD0, 0x8B,
	0xF0, 0xCA, 0x33, 0x54, 0x9D, 0x24, 0x7C, 0xEE,
	0xF3, 0xF5, 0xA5, 0x31, 0x3B, 0xD4, 0xB1, 0x57,
	0xD0, 0xB0, 0x9C, 0xCD, 0xE8, 0x30, 0xB9, 0xEB,
	0x3A, 0x02, 0xC4, 0xC5, 0xAA, 0x8A, 0xDA, 0x98,
}

func SelfCheck() {
	fmt.Println("Kuznechik tests:")
	var average time.Duration
	for i := 0; i < 10; i++ {
		fmt.Println("Test No.", i+1)
		time := run()
		fmt.Printf("Time;%v\n", time)
		util.PrintMemoryUsage()
		fmt.Println()
		average += time
	}
	fmt.Printf("Average;%v\n", average/10)
	util.PrintMemoryUsage()

}

func run() time.Duration {
	runtime.GOMAXPROCS(10)


	var k = Kuznechik{}
	k.SetKey(key[:])
	k.SetSubKeys()

	fmt.Println()

	pOrig, err := os.Create("kuzn-plaintext-original")
	if err != nil {
		panic(err)
	}

	cOrig, err := os.Create("kuzn-ciphertext-original")
	if err != nil {
		panic(err)
	}

	// Make a function from this
	var plainToFile []byte
	for i := 0; i < filesize/block; i++ {
		plainToFile = append(plainToFile, plainTest...)
	}
	plainToFile = append(plainToFile, plainTest[:0x20]...)

	// Make a function from this
	var cipherToFile []byte
	for i := 0; i < filesize/block; i++ {
		cipherToFile = append(cipherToFile, cipherTest...)
	}
	cipherToFile = append(cipherToFile, cipherTest[:0x20]...)

	util.WriteToFile(pOrig, plainToFile, 1)
	util.WriteToFile(cOrig, cipherToFile, 1)

	var startTime = time.Now()

	prOrig, err := os.Open("kuzn-plaintext-original")
	//crOrig, err := os.Open("kuzn-plaintext-original")

	crOrig, err := os.Open("kuzn-ciphertext-original")
	//prOrig, err := os.Open("kuzn-ciphertext-original")

	bytes, err := ioutil.ReadAll(prOrig)

	var input [][]byte
	input = util.SplitToTable(bytes, 16, 781250)

	ch := make(chan func() ([][]byte, int))

	inToChannels := util.SplitTableToParts(input, parts)

	for i := 0; i < parts; i++ {
		go processTestEncrypt(inToChannels[i], ch, k, i)
	}

	outFromChannels := make([][][]byte, parts)

	for i := 0; i < parts; i++ {
		out, c := (<-ch)()
		outFromChannels[c] = out
	}

	stopTime := time.Since(startTime)
	log.Println("Cipher goroutines were done in:", stopTime)

	resultSlice := util.CombineBytesFrom(outFromChannels, parts)

	bytesOc, err := ioutil.ReadAll(crOrig)

	log.Println(float32(len(resultSlice))/float32(1000000), "Mb")
	log.Println(util.IsEqual(resultSlice, bytesOc))
	pOrig.Close()
	prOrig.Close()
	cOrig.Close()
	crOrig.Close()
	err = os.Remove("kuzn-plaintext-original")
	err = os.Remove("kuzn-ciphertext-original")

	fmt.Println("Operation done in:", stopTime)
	fmt.Println()
	return stopTime
}

func processTestEncrypt(data [][]byte, c chan func() ([][]byte, int), cipher Kuznechik, o int) {
	t := time.Now()
	log.Println("Started", &c, o, time.Since(t))
	for i := 0; i < len(data); i++ {
		data[i] = cipher.Encrypt(data[i])
	}
	log.Println("Ended", &c, time.Since(t))
	c <- func() (bytes [][]byte, i int) {
		return data, o
	}
}
