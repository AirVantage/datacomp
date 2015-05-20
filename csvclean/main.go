package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func main() {

	file, err := os.Open("input.csv")

	if err != nil {
		log.Fatal(err)
	}

	csvReader := csv.NewReader(file)

	for i := 0; ; i++ {
		record, err := csvReader.Read()

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}

		if record == nil {
			break
		}

		if i != 0 {
			t, err := time.Parse("2006-01-02T15:04:05 MST", record[0])
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Print(t)
			//fmt.Println(record[1:])

			timeInMs := t.Unix() * 1000

			// find X
			for j := 0; j < 10; j++ {
				fmt.Print((timeInMs + int64(j)*100), ",")
				fmt.Print(record[1+j], ",")
				fmt.Print(record[11+j], ",")
				fmt.Print(record[21+j], ",")
				fmt.Printf("%v,%v,%v,%v\n", record[31], record[32], record[33], record[34])
			}
		}
	}

}
