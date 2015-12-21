package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"runtime"
)

var alfabeto = []rune{
	' ', ' ', '!', '"', '#', '$', '%', '&', '\'', '(',
	')', '*', '+', ',', '-', '.', '/', ':', ';', '<',
	'=', '>', '?', '@', '[', '\\', ']', '^', '_', '{',
	'|', '}', '~',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

var hash string
var data []string

func mustParseArgs() {
	if len(os.Args) <= 3 {
		print("Usage: ", os.Args[0], " hash known_data [[known data] ...]\n")
		os.Exit(1)
	}

	hash = os.Args[1]
	data = os.Args[2:]
}

func main() {
	mustParseArgs()

	completed := make(chan bool)
	combi := make(chan string)
	done := make(chan bool)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer func() { done <- true }()

			digest := sha1.New()
			for c := range combi {
				digest.Reset()
				digest.Write([]byte(c))
				result := hex.EncodeToString(digest.Sum(nil))

				if result != hash {
					continue
				}

				print("Found matched input string. string:", c, " hash:", result, "\n")
				completed <- true
			}
			print("Completing goroutine.\n")
		}()
	}

	indexes := make([]int, len(data)+1)

	print("Generating combinations.\n")
loop:
	for {
		select {
		case <-completed:
			close(combi)
			break loop
		default:
			combi <- combine(indexes)

			indexes[len(data)]++
			for i := len(data); i > 0; i-- {
				if indexes[i] < len(alfabeto) {
					break
				}
				indexes[i] = 0
				indexes[i-1]++
			}

			if indexes[0] >= len(alfabeto) {
				close(combi)
				break loop
			}
		}
	}

	print("Generator completed. Waiting for processors.\n")
	for i := 0; i < runtime.NumCPU(); i++ {
		<-done
	}
}

func combine(indexes []int) string {
	buffer := bytes.Buffer{}
	for i, ai := range indexes {
		if ai > 0 {
			buffer.WriteRune(alfabeto[ai])
		}
		if i < len(data) {
			buffer.WriteString(data[i])
		}
	}
	return buffer.String()
}
