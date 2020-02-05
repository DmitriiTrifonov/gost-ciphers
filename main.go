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
)

func main() {
	if len(os.Args) > 1 {
		inLoop := false
		isDecryptor := false
		isKuzn := false
		var delim byte
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

		startCipher(inLoop, isDecryptor, isKuzn, keyPath, delim)
	} else {
		magma.SelfCheck()
		kuznechik.SelfCheck()
	}
}

func startCipher(l bool, d bool, k bool, keyPath string, delim byte) {
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
		buffer, _ := reader.ReadBytes(0x0A )
		firstFlag := true
		for len(buffer) % blockSize != 0 {
			if firstFlag {
				buffer = append(buffer, 0x80)
				firstFlag = false
			} else {
				buffer = append(buffer, 0)
			}
		}
		log.Println("Buffer", buffer)
		splitted := split(buffer, blockSize)
		log.Println("Splitted length", len(splitted))
		log.Println("Splitted Input:", splitted)
		for i := 0; i < len(splitted); i++ {

				if d == true {
					splitted[i] = cipher.Decrypt(splitted[i])
				} else {
					splitted[i] = cipher.Encrypt(splitted[i])
				}

		}
		log.Println("Splitted Output:", splitted)
		for _, element := range splitted {
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

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}


