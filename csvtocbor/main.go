package main

import (
	"compress/zlib"
	"encoding/csv"
	"fmt"
	"github.com/jvermillard/datacomp/cbor"
	"io"
	"os"
	"strconv"
)

func main() {

	csvReader := csv.NewReader(os.Stdin)

	body := make(map[string]interface{})

	samples := make([]interface{}, 0)

	baseTime := 0

	var baseValue []float64

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
			header := make([]string, len(record)-1)
			baseValue = make([]float64, len(record)-1)
			for j, v := range record {
				if j != 0 {
					header[j-1] = v
				}
			}
			body["h"] = header
			continue
		}

		ts := 0
		for j, v := range record {
			if j == 0 {
				ts, err = strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				samples = append(samples, ts-baseTime)

				baseTime = ts
			} else {
				val, err := strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}

				delta := val - baseValue[j-1]
				if delta == float64(int64(delta)) {
					samples = append(samples, int64(delta))
				} else if delta == float64(float32(delta)) {
					samples = append(samples, float32(delta))
				} else {
					samples = append(samples, delta)
				}
				baseValue[j-1] = val
			}
		}
	}

	fmt.Fprintln(os.Stderr, len(samples))
	body["s"] = samples

	wz := zlib.NewWriter(os.Stdout)

	if err := cbor.Encode(wz, body); err != nil {
		panic(err)
	}

	wz.Flush()

}
