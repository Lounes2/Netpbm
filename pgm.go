package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           uint8
}

// Get a PGM image from a file and returns a struct that represents the image.
func GetPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	scanner.Scan()
	magicNumber := scanner.Text()

	// Dodge useless caracters
	for strings.HasPrefix(scanner.Text(), "#") {
		scanner.Scan()
	}

	// Read width, height, and max value
	scanner.Scan()
	size := strings.Fields(scanner.Text())
	width, _ := strconv.Atoi(size[0])
	height, _ := strconv.Atoi(size[1])

	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	// Read image data
	var data [][]uint8
	if magicNumber == "P2" {
		data, _ = readP2(scanner, width, height)
	} else if magicNumber == "P5" {
		data, _ = readP5(file, width, height)
	} else {
		return nil, fmt.Errorf("unsupported PGM format: %s", magicNumber)
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         uint8(max),
	}, nil
}

// Dimensions of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// Returns the value of the pixel(x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Sets the pixel's value at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Saves the image in case there is a problem and return error
func (pgm *PGM) Save(Name string) error {
	file, err := os.Create(Name)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write PGM's dimensions and magic number.
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Write image data
	if pgm.magicNumber == "P2" {
		createP2(writer, pgm.data)
	} else if pgm.magicNumber == "P5" {
		createP5(file, pgm.data)
	} else {
		return fmt.Errorf("ERROR: %s", pgm.magicNumber)
	}

	writer.Flush()
	return nil
}

// Swaps the colors of the image.
func (pgm *PGM) Invert() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pgm.data[y][x] = uint8(pgm.max) - pgm.data[y][x]
		}
	}
}

// Turns the image horizontally.
func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width/2; x++ {
			pgm.data[y][x], pgm.data[y][pgm.width-x-1] = pgm.data[y][pgm.width-x-1], pgm.data[y][x]
		}
	}
}

// Turns the image vertically.
func (pgm *PGM) Flop() {
	for y := 0; y < pgm.height/2; y++ {
		pgm.data[y], pgm.data[pgm.height-y-1] = pgm.data[pgm.height-y-1], pgm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// Sets the max value of the PGM image.
func (pgm *PGM) PGMMAX(maxVal uint8) {
	pgm.max = uint8(maxVal)
}

// Rotate the image (90 degrees)
func (pgm *PGM) ROTA() {
	newData := make([][]uint8, pgm.width)
	for x := 0; x < pgm.width; x++ {
		newData[x] = make([]uint8, pgm.height)
		for y := 0; y < pgm.height; y++ {
			newData[x][y] = pgm.data[pgm.height-y-1][x]
		}
	}
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// Converts the image to PBM from PGM
func (pgm *PGM) PBMConv() *PBM {
	data := make([][]bool, pgm.height)
	for y := 0; y < pgm.height; y++ {
		data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			data[y][x] = pgm.data[y][x] > uint8(pgm.max/2)
		}
	}
	return &PBM{
		data:        data,
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P4",
	}
}

func readP2(scanner *bufio.Scanner, width, height int) ([][]uint8, error) {
	data := make([][]uint8, height)
	for y := 0; y < height; y++ {
		data[y] = make([]uint8, width)
		VALl := strings.Fields(scanner.Text())
		for x := 0; x < width; x++ {
			value, err := strconv.Atoi(VALl[x])
			if err != nil {
				return nil, err
			}
			data[y][x] = uint8(value)
		}
		scanner.Scan()
	}
	return data, nil
}

func readP5(file *os.File, width, height int) ([][]uint8, error) {
	data := make([][]uint8, height)
	for y := 0; y < height; y++ {
		data[y] = make([]uint8, width)
		yes := make([]byte, width)
		_, err := file.Read(yes)
		if err != nil {
			return nil, err
		}
		for x := 0; x < width; x++ {
			data[y][x] = uint8(yes[x])
		}
	}
	return data, nil
}

func createP2(create *bufio.Writer, data [][]uint8) {
	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[y]); x++ {
			fmt.Fprintf(create, "%d ", data[y][x])
		}
		fmt.Fprintln(create)
	}
}

func createP5(file *os.File, data [][]uint8) {
	for y := 0; y < len(data); y++ {
		yes := make([]byte, len(data[y]))
		for x := 0; x < len(data[y]); x++ {
			yes[x] = byte(data[y][x])
		}
		file.Write(yes)
	}
}
