package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

func main() {

	csvReader := csv.NewReader(os.Stdin)

	senml := new(SenML)
	senml.Ver = 1
	senml.Bn = "myapp."
	senml.E = make([]*Element, 0)
	var header []string

	var baseTime int

	for i := 0; ; i++ {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if record == nil {
			break
		}

		if i == 0 {
			// read the header
			header = make([]string, len(record)-1)
			for j, v := range record {
				if j != 0 {
					header[j-1] = v
				}
			}
			continue
		}

		ts := 0
		for j, v := range record {
			if j == 0 {
				ts, err = strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				if i == 1 {
					baseTime = ts
				}
			} else {
				element := new(Element)

				element.Dt = ts - baseTime
				baseTime = ts
				val, err := strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				element.N = header[j-1]
				element.V = val
				senml.E = append(senml.E, element)
			}
		}
	}

	bin, err := json.MarshalIndent(senml, "", "")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bin))

	/*var buffTest bytes.Buffer
	encoder := cbor.NewEncoder(&buffTest)
	ok, err := encoder.Marshal(senml)

	if err != nil || !ok {
		panic(err)
	}*/

	//fmt.Printf("%v", buffTest.Bytes())

}

type SenML struct {
	Ver int        `json:"ver"`
	Bt  int64      `json:"bt"`
	Bn  string     `json:"bn"`
	E   []*Element `json:"e"`
}

type Element struct {
	N  string  `json:"n"`
	Dt int     `json:"dt"`
	V  float64 `json:"v"`
}
