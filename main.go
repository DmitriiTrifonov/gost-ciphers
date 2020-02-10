package main

import (
	"bufio"
	"github.com/DmitriiTrifonov/gost-ciphers/kuznechik"
	"github.com/DmitriiTrifonov/gost-ciphers/magma"
	"github.com/DmitriiTrifonov/gost-ciphers/util"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) > 1 {
		inLoop := false
		isDecryptor := false
		isKuzn := false
		sendDelim := false
		var delim byte = 0x0A
		keyPath := ""
		keyIndex, _ := util.ArgIndex("-key")
		if keyIndex == -1 {
			panic("Key has not found")
		} else {
			keyPath = os.Args[keyIndex + 1]
		}

		delimIndex, _ := util.ArgIndex("-delim")
		if delimIndex != -1 {
			i, _ := strconv.Atoi(os.Args[delimIndex + 1])
			delim = byte(i)
		}
		sendDelimIndex, _ := util.ArgIndex("-sd")

		if sendDelimIndex != -1 {
			sendDelim = true
		}

		for _, element := range os.Args {
			switch element {
			case "-l":
				inLoop = true
			case "-d":
				isDecryptor = true
			case "-k":
				isKuzn = true
			default:
				continue
			}
		}

		startCipher(inLoop, isDecryptor, isKuzn, keyPath, delim, sendDelim)
	} else {
		magma.SelfCheck()
		kuznechik.SelfCheck()
	}
}

func startCipher(l bool, d bool, k bool, keyPath string, delim byte, sd bool) {
	var cipher Cipher
	blockSize := 8
	if k {
		cipher = &kuznechik.Kuznechik{}
		blockSize = 16
	} else {
		cipher = &magma.Magma{}
	}
	keyFile, _ := os.Open(keyPath)
	key, _ := ioutil.ReadAll(keyFile)
	keyFile.Close()
	cipher.SetKey(key)
	cipher.SetSubKeys()
	for {
		reader := bufio.NewReader(os.Stdin)
		buffer, _ := reader.ReadBytes(delim)
		firstFlag := true
		if !sd {
			buffer = buffer[:len(buffer) - 1]
		}
		if len(buffer) % blockSize != 0 {
			for len(buffer) % blockSize != 0 {
				if firstFlag {
					buffer = append(buffer, 0x80)
					firstFlag = false
				} else {
					buffer = append(buffer, 0)
				}
			}
			padding := make([]byte, blockSize)
			buffer = append(buffer, padding...)
		} else {
			padding := make([]byte, blockSize)
			padding[0] = 0x80
			buffer = append(buffer, padding...)
		}

		log.Println("Buffer", buffer)
		splitted := split(buffer, blockSize)
		log.Println("Splitted length:", len(splitted))
		log.Println("Splitted Input:", splitted)

		wait := new(sync.WaitGroup)
		out := make([][]byte, len(splitted))

		for i, arr := range splitted {
			wait.Add(1)
			go func(j int) {
				defer wait.Done()
				if d == true {
					out[j] = cipher.Decrypt(arr)
				} else {
					out[j] = cipher.Encrypt(arr)
				}
				log.Println("j", j)
				log.Println("arr", arr)
			}(i)
			wait.Wait()
		}

		log.Println("Splitted Output:", splitted)
		for _, element := range out {
			f := bufio.NewWriter(os.Stdout)
			n, _ := f.Write(element)
			log.Println("Bytes written:", n)
			log.Println(element)
			_ = f.Flush()
		}
		if !l {
			break
		}
	}
}


func split(input []byte, size int) [][]byte {
	times := len(input) / size
	table := make([][]byte, times)
	for i := 0; i < times; i++ {
		table[i] = input[i*size:(i + 1)*size]
	}
	return table
}





