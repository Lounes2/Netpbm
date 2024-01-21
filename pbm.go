package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	var dimension string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	line = strings.TrimSpace(line)
	if line != "P1" && line != "P4" {
		return nil, fmt.Errorf("error %s", line)
	}
	magicNumber := line

	// Lecture des dimensions
	for scanner.Scan() {
		if scanner.Text()[0] == '#' {
			continue
		}
		break

	}

	dimension = scanner.Text()
	res := strings.Split(dimension, " ")
	height, _ := strconv.Atoi(res[0])
	width, _ := strconv.Atoi(res[1])

	// Lecture des données binaires
	var pbm *PBM

	if magicNumber == "P1" {
		data := make([][]bool, height)
		for i := range data {
			data[i] = make([]bool, width)
		}

		for i := 0; i < height; i++ {
			scanner.Scan()
			line := scanner.Text()
			hori := strings.Fields(line)
			for j := 0; j < width; j++ {
				verti, _ := strconv.Atoi(hori[j])
				if verti == 1 {
					data[i][j] = true
				}
			}
		}

		pbm = &PBM{
			data:        data,
			width:       width,
			height:      height,
			magicNumber: magicNumber,
		}
		fmt.Printf("%+v\n", PBM{data, width, height, magicNumber})
	}
	return pbm, nil
}

func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At renvoie la valeur du pixel en (x, y).
func (pbm *PBM) At(x, y int) bool {
	if len(pbm.data) == 0 || x < 0 || y < 0 || x >= pbm.width || y >= pbm.height {
		return false
	}

	return pbm.data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

func (pbm *PBM) Save(filename string) error {
	fileSave, error := os.Create(filename)
	if error != nil {
		return error
	}
	defer fileSave.Close()

	fmt.Fprintf(fileSave, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	if pbm.magicNumber == "P1" {
		for _, i := range pbm.data {
			for _, j := range i {
				if j {
					fmt.Fprint(fileSave, "1 ")
				} else {
					fmt.Fprint(fileSave, "0 ")
				}
			}
			fmt.Fprintln(fileSave)
		}
	}
	return nil
}

func (pbm *PBM) Invert() {
	for x := 0; x < pbm.height; x++ {
		for y := 0; y < pbm.width; y++ {
			pbm.data[x][y] = !pbm.data[x][y]
		}
	}
}

func (pbm *PBM) Flip() {
	for x := 0; x < pbm.height; x++ {
		for y := 0; y < pbm.width/2; y++ {
			pbm.data[x][y], pbm.data[x][pbm.width-y-1] = pbm.data[x][pbm.width-y-1], pbm.data[x][y]
		}
	}
}

func (pbm *PBM) Flop() {
	for x := 0; x < pbm.height/2; x++ {
		pbm.data[x], pbm.data[pbm.height-x-1] = pbm.data[pbm.height-x-1], pbm.data[x]
	}
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
