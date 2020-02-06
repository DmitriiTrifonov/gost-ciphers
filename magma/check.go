package magma

import (
	"fmt"
	"github.com/DmitriiTrifonov/gost-ciphers/util"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const parts = 16

const filesize = 12500000

const block = 32

// K=ffeeddccbbaa99887766554433221100f0f1f2f3f4f5f6f7f8f9fafbfcfdfeff.
var key = []byte{
	0xFF, 0xEE, 0xDD, 0xCC, //
	0xBB, 0xAA, 0x99, 0x88, //
	0x77, 0x66, 0x55, 0x44, //
	0x33, 0x22, 0x11, 0x00, //
	0xF0, 0xF1, 0xF2, 0xF3, //
	0xF4, 0xF5, 0xF6, 0xF7, //
	0xF8, 0xF9, 0xFA, 0xFB, //
	0xFC, 0xFD, 0xFE, 0xFF, //
}

var plainTest = []byte{
	0x92, 0xde, 0xf0, 0x6b, 0x3c, 0x13, 0x0a, 0x59,
	0xdb, 0x54, 0xc7, 0x04, 0xf8, 0x18, 0x9d, 0x20,
	0x4a, 0x98, 0xfb, 0x2e, 0x67, 0xa8, 0x02, 0x4c,
	0x89, 0x12, 0x40, 0x9b, 0x17, 0xb5, 0x7e, 0x41,
}

var cipherTest = []byte{
	0x2b, 0x07, 0x3f, 0x04, 0x94, 0xf3, 0x72, 0xa0,
	0xde, 0x70, 0xe7, 0x15, 0xd3, 0x55, 0x6e, 0x48,
	0x11, 0xd8, 0xd9, 0xe9, 0xea, 0xcf, 0xbc, 0x1e,
	0x7c, 0x68, 0x26, 0x09, 0x96, 0xc6, 0x7e, 0xfb,
}

func SelfCheck() {
	fmt.Println("Magma tests:")
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


	var m = Magma{}
	m.SetKey(key[:])
	m.SetSubKeys()

	pOrig, err := os.Create("magma-plaintext-original")
	if err != nil {
		panic(err)
	}

	cOrig, err := os.Create("magma-ciphertext-original")

	// Make a function from this
	var plainToFile []byte
	for i := 0; i < filesize/block; i++ {
		plainToFile = append(plainToFile, plainTest...)
	}

	// Make a function from this
	var cipherToFile []byte
	for i := 0; i < filesize/block; i++ {
		cipherToFile = append(cipherToFile, cipherTest...)
	}

	util.WriteToFile(pOrig, plainToFile, 1)
	util.WriteToFile(cOrig, cipherToFile, 1)

	var startTime = time.Now()

	prOrig, err := os.Open("magma-plaintext-original")
	//crOrig, err := os.Open("magma-plaintext-original")
	crOrig, err := os.Open("magma-ciphertext-original")
	//prOrig, err := os.Open("magma-ciphertext-original")

	bytes, err := ioutil.ReadAll(prOrig)

	var input [][]byte
	input = util.SplitToTable(bytes, 8, 1562500)

	ch := make(chan func() ([][]byte, int))

	inToChannels := util.SplitTableToParts(input, parts)

	for i := 0; i < parts; i++ {
		go processTestEncrypt(inToChannels[i], ch, m, i)
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
	err = os.Remove("magma-plaintext-original")
	err = os.Remove("magma-ciphertext-original")

	fmt.Println("Operation done in:", stopTime)
	fmt.Println()
	return stopTime
}

func processTestEncrypt(data [][]byte, c chan func() ([][]byte, int), cipher Magma, o int) {
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
