package main

import (
	"compress/zlib"
	"encoding/csv"
	"fmt"
	"github.com/jvermillard/datacomp/cbor"
	"os"
	"strconv"
)

func main() {

	body := make(map[string]interface{})
	zr, err := zlib.NewReader(os.Stdin)
	if err != nil {
		panic(err)
	}
	decoder := cbor.NewDecoder(zr)
	if err := decoder.Decode(&body); err != nil {
		panic(err)
	}

	w := csv.NewWriter(os.Stdout)

	rawH := body["h"].([]interface{})

	header := make([]string, len(rawH)+1)
	baseValue := make([]float64, len(rawH))

	header[0] = "ts"
	for i, v := range rawH {
		header[i+1] = v.(string)
	}
	w.Write(header)

	var baseTime uint64 = 0

	samples := body["s"].([]interface{})

	for i := 0; i < len(samples); i += len(header) {

		line := make([]string, len(header))
		line[0] = strconv.FormatUint(baseTime+samples[i].(uint64), 10)
		baseTime = baseTime + uint64(samples[i].(uint64))
		for j := 0; j < len(header)-1; j++ {
			delta, err := strconv.ParseFloat(fmt.Sprintf("%v", samples[i+1+j]), 64)
			if err != nil {
				panic(err)
			}
			line[j+1] = fmt.Sprintf("%v", baseValue[j]+delta)
			baseValue[j] += delta
		}
		w.Write(line)
	}
	w.Flush()
}
