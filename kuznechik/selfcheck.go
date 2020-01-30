package kuznechik

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"time"
)

var plaintext = [0x10]byte{
	0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
	0xFF, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
}

var cipherText = [0x10]byte{
	0x7F, 0x67, 0x9D, 0x90, 0xBE, 0xBC, 0x24, 0x30,
	0x5A, 0x46, 0x8D, 0x42, 0xB9, 0xD4, 0xED, 0xCD,
}

var key = [0x20]byte{
	0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, // 00..07
	0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, // 08..0F
	0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10, // 10..17
	0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, // 18..1F
}

func SelfCheck() {
	var average time.Duration
	for i := 0; i < 1; i++ {
		fmt.Println("Test No.", i+1)
		average += run()
	}
	fmt.Println("Average time is:", average/10)

}

func run() time.Duration {
	runtime.GOMAXPROCS(10)
	fmt.Println("Self Test Started:", time.Now())

	fmt.Print("\nPlaintext: ")
	for i := 0; i < len(plaintext); i++ {
		fmt.Printf("0x%X, ", plaintext[i])
	}

	fmt.Print("\nCiphertext: ")
	for i := 0; i < len(cipherText); i++ {
		fmt.Printf("0x%X, ", cipherText[i])
	}

	var k = Kuznechik{}
	k.SetKey(key[:])
	k.SetSubKeys()

	var encrypted = k.Encrypt(plaintext)
	fmt.Print("\nCipher Encryption Result: ")
	for i := 0; i < len(encrypted); i++ {
		fmt.Printf("0x%X, ", encrypted[i])
	}

	var decrypted = k.Decrypt(cipherText)

	fmt.Print("\nCipher Decryption Result: ")
	for i := 0; i < len(decrypted); i++ {
		fmt.Printf("0x%X, ", decrypted[i])
	}
	fmt.Println()

	pOrig, err := os.Create("plaintext-original")
	if err != nil {
		panic(err)
	}

	cOrig, err := os.Create("ciphertext-original")
	if err != nil {
		panic(err)
	}

	writeToFile(pOrig, plaintext[:], 781250)
	writeToFile(cOrig, cipherText[:], 781250)

	var startTime = time.Now()

	prOrig, err := os.Open("plaintext-original")
	if err != nil {
		panic(err)
	}

	crOrig, err := os.Open("ciphertext-original")
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(prOrig)

	var input [][]byte
	input = splitToTable(bytes, 16, 781250)

	ch := make(chan [][]byte)

	inToChannels := splitTableToParts(input, 10)

	for i := 0; i < 10; i++ {
		go processTestEncrypt(inToChannels[i], ch, k)
	}

	outFromChannels := make([][][]byte, 10)

	for i := 0; i < 10; i++ {
		out := <-ch
		outFromChannels[i] = out
	}

	stopTime := time.Since(startTime)
	log.Println("Cipher goroutines were done in:", stopTime)

	resultSlice := combineBytesFrom(outFromChannels, 10)

	bytesOc, err := ioutil.ReadAll(crOrig)

	log.Println(float32(len(resultSlice))/float32(1000000), "Mb")
	log.Println(isEqual(resultSlice, bytesOc))
	pOrig.Close()
	prOrig.Close()
	cOrig.Close()
	crOrig.Close()
	os.Remove("plaintext-original")
	os.Remove("ciphertext-original")

	fmt.Println("Operation done in:", stopTime)
	fmt.Println()
	return stopTime
}

func processTestEncrypt(data [][]byte, c chan [][]byte, cipher Kuznechik) {
	t := time.Now()
	log.Println("Started", &c, time.Since(t))
	for i := 0; i < len(data); i++ {
		var arr [16]byte
		copy(arr[:], data[i])
		data[i] = cipher.Encrypt(arr)
	}
	log.Println("Ended", &c, time.Since(t))
	c <- data
}

func writeToFile(f *os.File, data []byte, c int) {
	for i := 0; i < c; i++ {
		_, err := f.Write(data)
		if err != nil {
			panic(err)
		}
	}
	f.Close()
}

func splitToTable(in []byte, s byte, c int) [][]byte {
	input := make([][]byte, c)
	for i := 0; i < c; i++ {
		from := i * int(s)
		to := i*int(s) + int(s)
		input[i] = in[from:to]
	}
	return input
}

func splitTableToParts(in [][]byte, c int) [][][]byte {
	s := make([][][]byte, c)
	for i := 1; i <= c; i++ {
		if i == c {
			s[i-1] = in[(len(in)/c)*(i-1):]
		} else {
			s[i-1] = in[(len(in)/c)*(i-1) : (len(in)/c)*i]
		}
	}
	return s
}

func combineBytesFrom(in [][][]byte, c int) []byte {
	encSlice := make([][]byte, c)
	for i := 0; i < c; i++ {
		encSlice = append(encSlice, in[i]...)
	}

	var resultSlice []byte

	for i := 0; i < len(encSlice); i++ {
		resultSlice = append(resultSlice, encSlice[i]...)
	}
	return resultSlice
}

func isEqual(r []byte, c []byte) string {
	for j := 0; j < len(r); j++ {
		if c[j] != r[j] {
			log.Println("Byte", c[j], "is not equal to", r[j])
			log.Println("Byte address is:", j)
			panic("Stopped")
		}
	}
	return "Everything is OK!"
}
