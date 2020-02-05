package util

import "os"

func ArgIndex(arg string) (int, error) {
	i := -1
	for index, element := range os.Args {
		if element == arg {
			i = index
			break
		}
	}
	return i, nil
}
