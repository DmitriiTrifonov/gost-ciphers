package main

import (
	"fmt"
	"github.com/DmitriiTrifonov/gost-ciphers/kuznechik/ecb"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	var plaintext = [0x10]byte {
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x00,
		0xFF, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
	}

	fmt.Print("\nPlaintext: ")
	for i := 0; i < len(plaintext); i++ {
		fmt.Printf("0x%X, ",plaintext[i])
	}

	var cipherText = [0x10]byte {
		0x7F, 0x67, 0x9D, 0x90, 0xBE, 0xBC, 0x24, 0x30,
		0x5A, 0x46, 0x8D, 0x42, 0xB9, 0xD4, 0xED, 0xCD,
	}

	fmt.Print("\nCiphertext: ")
	for i := 0; i < len(cipherText); i++ {
		fmt.Printf("0x%X, ", cipherText[i])
	}

	var key = [0x20]byte {
		0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 	// 00..07
		0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 	// 08..0F
		0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10, 	// 10..17
		0x01, 0x23, 0x45, 0x67, 0x89, 0xAB, 0xCD, 0xEF, 	// 18..1F
	}

	ecb.SetKey(&key)

	var encrypted = ecb.Encrypt(plaintext)
	fmt.Print("\nCipher Encryption Result: ")
	for i := 0; i < len(encrypted); i++ {
		fmt.Printf("0x%X, ",encrypted[i])
	}

	var decrypted = ecb.Decrypt(cipherText)

	fmt.Print("\nCipher Decryption Result: ")
	for i := 0; i < len(decrypted); i++ {
		fmt.Printf("0x%X, ",decrypted[i])
	}

	f, err := os.Create("plaintext-decrypted")
	if err != nil {
		panic(err)
	}

	o, errO := os.Open("plaintext-original")
	if errO != nil {
		panic(errO)
	}

	bytes, err := ioutil.ReadAll(o)


	defer f.Close()
	defer o.Close()

	var startTime = time.Now().UnixNano()

	var input [781250][]byte
	for i := 0; i < 781250; i++ {
		from := i * 16
		to := i * 16 + 16
		input[i] = bytes[from : to]
	}

	println(len(input))
	i1 := input[:len(input) / 10 ]
	i2 := input[len(input) / 10 : (len(input) / 10) * 2 ]
	i3 := input[(len(input) / 10) * 2 : (len(input) / 10) * 3 ]
	i4 := input[(len(input) / 10) * 3 : (len(input) / 10) * 4 ]
	i5 := input[(len(input) / 10) * 4 : (len(input) / 10) * 5 ]
	i6 := input[(len(input) / 10) * 5 : (len(input) / 10) * 6 ]
	i7 := input[(len(input) / 10) * 6 : (len(input) / 10) * 7 ]
	i8 := input[(len(input) / 10) * 7 : (len(input) / 10) * 8 ]
	i9 := input[(len(input) / 10) * 8 : (len(input) / 10) * 9 ]
	i10 := input[(len(input) / 10) * 9 :]
	println(len(i1))
	c1 := make(chan [][]byte)
	c2 := make(chan [][]byte)
	go processTestEncrypt(i1, c1)
	go processTestEncrypt(i2, c1)
	go processTestEncrypt(i3, c1)
	go processTestEncrypt(i4, c1)
	go processTestEncrypt(i5, c1)
	go processTestEncrypt(i6, c2)
	go processTestEncrypt(i7, c2)
	go processTestEncrypt(i8, c2)
	go processTestEncrypt(i9, c2)
	go processTestEncrypt(i10, c2)
	o1 := <-c1
	o2 := <-c1
	o3 := <-c1
	o4 := <-c1
	o5 := <-c1
	o6 := <-c2
	o7 := <-c2
	o8 := <-c2
	o9 := <-c2
	o10 := <-c2


	var endTime = time.Now().UnixNano()


	fmt.Println("\n", (endTime - startTime) / 1000000000, ":", (endTime - startTime) % 1000000000)

	println(len(o1))
	println(len(o2))
	println(len(o3))
	println(len(o4))
	println(len(o5))
	println(len(o6))
	println(len(o7))
	println(len(o8))
	println(len(o9))
	println(len(o10))
}

func processTestEncrypt(data [][]byte, c chan [][]byte) {
	r := rand.Intn(100)
	t := time.Now()
	log.Println("Started" , r ,t)
	for i := 0; i < len(data); i++ {
		var arr [16]byte
		copy(arr[:], data[i])
		data[i] = ecb.Encrypt(arr)
	}
	log.Println("Ended", r, time.Since(t))
	c <- data
}
