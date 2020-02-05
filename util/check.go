package util

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func WriteToFile(f *os.File, data []byte, c int) {
	for i := 0; i < c; i++ {
		_, err := f.Write(data)
		if err != nil {
			panic(err)
		}
	}
	f.Close()
}

func SplitToTable(in []byte, s byte, c int) [][]byte {
	input := make([][]byte, c)
	for i := 0; i < c; i++ {
		from := i * int(s)
		to := i*int(s) + int(s)
		input[i] = in[from:to]
	}
	return input
}

func SplitTableToParts(in [][]byte, c int) [][][]byte {
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

func CombineBytesFrom(in [][][]byte, c int) []byte {
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

func IsEqual(r []byte, c []byte) string {
	for j := 0; j < len(r); j++ {
		if c[j] != r[j] {
			log.Println("Byte", c[j], "is not equal to", r[j])
			log.Println("Byte address is:", j)
			panic("Stopped")
		}
	}
	return "Everything is OK!"
}

func bytesToMb(b uint64) uint64 {
	return b / 1000000
}

func PrintMemoryUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc;%v Mb\n", bytesToMb(m.Alloc))
	fmt.Printf("Sys;%v Mb\n", bytesToMb(m.Sys))
}
