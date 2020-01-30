package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		log.Fatal("Not Implemented")
	} else {
		// Magma
		/*key := []byte {
			0xFF, 0xEE, 0xDD, 0xCC,   //
			0xBB, 0xAA, 0x99, 0x88,   //
			0x77, 0x66, 0x55, 0x44,   //
			0x33, 0x22, 0x11, 0x00,   //
			0xF0, 0xF1, 0xF2, 0xF3,   //
			0xF4, 0xF5, 0xF6, 0xF7,   //
			0xF8, 0xF9, 0xFA, 0xFB,   //
			0xFC, 0xFD, 0xFE, 0xFF,   //
		}
		m := magma.Magma{}
		m.SetKey(key)
		m.SetSubKeys()
		plaintext := []byte {
			0xFE, 0xDC, 0xBA, 0x98, 0x76, 0x54, 0x32, 0x10,
		}
		ciphertext := []byte {
			0x4E, 0xE9, 0x01, 0xE5, 0xC2, 0xD8, 0xCA, 0x3D,
		}
		enc := m.Encrypt(plaintext)
		for i := 0; i < len(enc); i++ {
			fmt.Printf("0x%X ", enc[i])
		}
		fmt.Println()
		dec := m.Decrypt(ciphertext)
		for i := 0; i < len(dec); i++ {
			fmt.Printf("0x%X ", dec[i])
		}*/

		kuznechik.SelfCheck()
	}

}
