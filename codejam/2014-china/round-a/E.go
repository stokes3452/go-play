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

type pair map[string] bool

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

func main() {
	testsRemaining := readIntLine()

	for testNumber := 1; testNumber <= testsRemaining; testNumber++ {
		fmt.Fprintf(outputFile, "Case #%d:\n", testNumber)
		Debug("Test number %d\n", testNumber)
		roomCount := readIntLine()
		roomColors := make(map[string]int)
		equivalentRooms := make(map[int]int)
		//we'll go hideous DFS implementation for now, for speedz.  Structs if I open this file again, ever, I promise
		finalRooms := make(map[int]*room)
		for i := 1; i <= roomCount; i++ {
			roomColor := string(readLine())
			oldRoom, exists := roomColors[roomColor]
			if exists {
				equivalentRooms[i] = oldRoom
			} else {
				roomColors[roomColor] = i
				//why not go all in?
				newRoom := &room{}
				newRoom.teleporters = make(map[int]int)
				finalRooms[i] = newRoom
			}
		}

		teleporterCount := readIntLine()
		for i := 0; i < teleporterCount; i++ {
			teleporterValues := getIntsFromLine()
			teleporterSource := getRoomEquivalency(equivalentRooms, teleporterValues[0])
			teleporterDestination := getRoomEquivalency(equivalentRooms, teleporterValues[1])
			Debug("Initializing a teleporter from %d to %d at a cost of %d\n", teleporterSource, teleporterDestination, teleporterValues[2])
			if teleporterDestination == teleporterSource {
				//no point in adding a teleporter we're never going to use
				continue
			}
			teleporterRoom := finalRooms[teleporterSource]
			teleporterRoom.setTeleporter(teleporterDestination, teleporterValues[2])
		}
		
		soldierCount := readIntLine()
		for i := 0; i < soldierCount; i++ {
			soldierValues := getIntsFromLine()
			soldierSource := getRoomEquivalency(equivalentRooms, soldierValues[0])
			soldierDestination := getRoomEquivalency(equivalentRooms, soldierValues[1])
			fmt.Fprintf(outputFile, "%d\n", getPathLength(soldierSource, soldierDestination, finalRooms))
		}

		fmt.Printf("Processed %d of %d tests\n", testNumber, testsRemaining)
	}
}

func getPathLength(source, destination int, rooms map[int]*room) (int) {
	Debug("Exploring path from %d to %d\n", source, destination)
	if source == destination {
		Debug("Short-circuiting - we're on the same equivalent room\n")
		return 0
	}
	for _, curRoom := range rooms {
		curRoom.pathLength = -1
	}
	
	//M turbolifts, M channels max.. M is 3k, that shouldn't really hurt our memory eh?
	roomQueue := make(chan int, 9000)
	sourceRoom := rooms[source]
	roomQueue <- source
	sourceRoom.pathLength = 0
	destinationRoom := rooms[destination]
	var currentRoom *room
	current := 0
	for {
		select {
			//go until we run out of stuff to go on
			case current = <- roomQueue:
				currentRoom = rooms[current]
				for target, cost := range currentRoom.teleporters {
					Debug("Examining teleporter to %d at cost %d\n", target, cost)
					targetRoom := rooms[target]
					newCost := currentRoom.pathLength + cost
					//If at any point we're looking further ahead than our destination room, stop
					if destinationRoom.pathLength != -1 && newCost > destinationRoom.pathLength {
						continue;
					}
					if targetRoom.pathLength == -1 || targetRoom.pathLength > newCost {
						Debug("Explored to %d from %s and replaced %d with %d\n", target, currentRoom, cost, newCost)
						targetRoom.pathLength = newCost
						roomQueue <- target
					}
				}
			default:
				return destinationRoom.pathLength
		}
	}	
	
	return -1
}

type room struct {
	teleporters map[int]int
	pathLength int
}

func (myRoom room) setTeleporter(destination, cost int) {
	oldTime, exists := myRoom.teleporters[destination]
	if !exists || oldTime > cost {
		myRoom.teleporters[destination] = cost
	}
	return
}

func getRoomEquivalency(equivalentRooms map[int]int, sourceRoom int) (int) {
	alternateRoom, exists := equivalentRooms[sourceRoom]
	if exists {
		return alternateRoom
	} else {
		return sourceRoom
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