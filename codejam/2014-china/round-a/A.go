package main

import (
	"flag"
	"fmt"
	"bufio"
	"bytes"
	"strconv"
	"os"
	"sort"
)

var fileIn = flag.String("fileIn", "in.txt", "The file to inport")
var fileOut = flag.String("fileOut", "out.txt", "The output file")
var debug = flag.Bool("debug", false, "Debug on or off?")
var reader *bufio.Reader
var outputFile *os.File

func Debug(format string, a ...interface{}) {
	if (*debug) {
		fmt.Printf(format, a...)
	}
}

func init() {
	flag.Parse()
	Debug("Debug is %t\n", *debug)
	Debug("File in is: %s\n", *fileIn)
	Debug("File out is: %s\n", *fileOut)
	
	file, err := os.Open(*fileIn)
	if err != nil {
		panic(fmt.Sprintf("File read failed with err: %s", err))
	}
	
	reader = bufio.NewReaderSize(file, 16000)
	outputFile, err = os.Create(*fileOut)
	if err != nil {
		panic(fmt.Sprintf("File creation failed with err: %s", err))
	}
}

var doubles map[int]string
var digits map[byte]string

func init() {
	doubles = make(map[int]string)
	digits = make(map[byte]string)
	//english hard, I don't want to think about this
	doubles[2] = "double"
	doubles[3] = "triple"
	doubles[4] = "quadruple"
	doubles[5] = "quintuple"
	doubles[6] = "sextuple"
	doubles[7] = "septuple"
	doubles[8] = "octuple"
	doubles[9] = "nonuple"
	doubles[10] = "decuple"
	digits['0'] = "zero"
	digits['1'] = "one"
	digits['2'] = "two"
	digits['3'] = "three"
	digits['4'] = "four"
	digits['5'] = "five"
	digits['6'] = "six"
	digits['7'] = "seven"
	digits['8'] = "eight"
	digits['9'] = "nine"
}

func main() {
	testsRemaining := readIntLine()

	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		Debug("Test number %d\n", testNumber)
		line := readLine()
		values := bytes.Split(line, []byte{' '})
		if len(values) != 2 {
			panic(fmt.Sprint("Invalid line: %s", line))
		}
		
		phoneNumber := values[0]
		format := values[1]
		formats := []int{}
		response := [][]byte{}
		//get some ints up in here..
		formatParts := bytes.Split(format, []byte{'-'})
		for _, part := range formatParts {
			currentValue, err := strconv.Atoi(string(part))
			if err != nil {
				panic(fmt.Sprintf("bad int provided: %s", part))
			}
			formats = append(formats, currentValue)
		}
		
		phoneCounter := 0
		for _, formatLength := range formats {
			digit := phoneNumber[phoneCounter]
			digitCounter := 1
			phoneCounter++
			for i := 1; i < formatLength; i++ {
				if phoneNumber[phoneCounter] == digit {
					digitCounter++
				} else {
					_, ok := doubles[digitCounter]
					if ok {
						response = append(response, []byte(doubles[digitCounter]))
					} else {
						for ;digitCounter > 1; digitCounter-- {
							response = append(response, []byte(digits[digit]))
						}
					}
					response = append(response, []byte(digits[digit]))
					//print remaining digits
					digit = phoneNumber[phoneCounter]
					digitCounter = 1
				}
				phoneCounter++
			}
			_, ok := doubles[digitCounter]
			if ok {
				response = append(response, []byte(doubles[digitCounter]))
			} else {
				for ;digitCounter > 1; digitCounter-- {
					response = append(response, []byte(digits[digit]))
				}
			}
			response = append(response, []byte(digits[digit]))
		}


		fmt.Fprintf(outputFile, "Case #%d: %s\n", testNumber, bytes.Join(response, []byte{' '}))
		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
}

func readLine() (line []byte) {
	line, prefix, err := reader.ReadLine()
	if err != nil {
		panic("Error reading line for test")
	}
	if prefix {
		panic("Did not complete read.  Increase buffer size and try try again")
	}
	return
}

func getFloatsFromLine() (lineValues sort.Float64Slice) {
	line := readLine()
	parts := bytes.Split(line, []byte{' '})
	for _, part := range parts {
		currentValue, err := strconv.ParseFloat(string(part), 64)
		if err != nil {
			panic(fmt.Sprintf("bad float provided: %s", part))
		}
		lineValues = append(lineValues, currentValue)
	}
	return
}

func getIntsFromLine() (lineValues sort.IntSlice) {
	line := readLine()
	parts := bytes.Split(line, []byte{' '})
	for _, part := range parts {
		currentValue, err := strconv.Atoi(string(part))
		if err != nil {
			panic(fmt.Sprintf("bad int provided: %s", part))
		}
		lineValues = append(lineValues, currentValue)
	}
	return
}

func readIntLine() (response int) {
	line := readLine()
	response, err := strconv.Atoi(string(line))
	if err != nil {
		panic(fmt.Sprintf("Err reading int value from line: %d", line))
	}
	return
}