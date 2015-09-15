package main

import (
	"compress/zlib"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/jvermillard/datacomp/cbor"
)

func main() {

	csvReader := csv.NewReader(os.Stdin)

	body := make(map[string]interface{})

	samples := make([]interface{}, 0)

	baseTime := 0.0

	var baseValue []float64

	var factors []float64

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
			factors = make([]float64, len(record))
			body["f"] = factors
			continue
		}

		// fill factors, I know row 0 is TS (sampling at 100ms)
		factors[0] = 0.01
		factors[1] = 1
		factors[2] = 1
		factors[3] = 1
		factors[4] = 1000000
		factors[5] = 1000000
		factors[6] = 1
		factors[7] = 100

		ts := 0.0
		for j, v := range record {
			if j == 0 {
				tsint, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				ts = float64(tsint) * factors[0]

				appendCompactType(ts-baseTime, &samples)

				baseTime = ts
			} else {
				val, err := strconv.ParseFloat(v, 64)
				if err != nil {
					panic(err)
				}
				val = val * factors[j]
				delta := val - baseValue[j-1]
				appendCompactType(delta, &samples)
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
func appendCompactType(value float64, samples *[]interface{}) {
	if value == float64(int64(value)) {
		*samples = append(*samples, int64(value))
	} else if value == float64(float32(value)) {
		*samples = append(*samples, float32(value))
	} else {
		*samples = append(*samples, value)
	}
}
