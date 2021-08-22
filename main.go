package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// App to make takeout .json into .csv

type LocationValue struct {
	longitude float64
	latitude  float64
	timestamp time.Time
}

type LocationSample struct {
	TimestampMs string `json:timestampMs`
	LatitudeE7  int    `json:latitudeE7`
	LongitudeE7 int    `json:longitudeE7`
	Accuracy    int    `json:accuracy`
	Source      string `json:source`
}

type AllLocations struct {
	Locations []LocationSample `json:locations`
}

func parseJsonFile(fd io.Reader) ([]LocationSample, error) {
	recs := AllLocations{}

	decoder := json.NewDecoder(fd)
	if err := decoder.Decode(&recs); err != nil {
		return nil, fmt.Errorf("can't decode file: %v", err)
	}

	return recs.Locations, nil
}

func convertFormat(recs []LocationSample, starttime, endtime time.Time) []LocationValue {
	processed := make([]LocationValue, 0, len(recs))

	for i, ls := range recs {
		// Time
		ts, err := strconv.ParseInt(ls.TimestampMs, 10, 64)
		if err != nil {
			log.Printf("can't parse [%d]%s: %v", i, ls.TimestampMs, err)
			continue
		}
		tm := time.Unix(0, ts*1000*1000)

		if !tm.After(starttime) {
			continue
		}
		if tm.After(endtime) {
			break
		}

		processed = append(processed, LocationValue{
			timestamp: tm,
			latitude:  float64(ls.LatitudeE7) / 10e6,
			longitude: float64(ls.LongitudeE7) / 10e6,
		})
	}

	return processed
}

const daterangeformat = "20060102"

func parseDateRange(dr string) (time.Time, time.Time, error) {
	drs := strings.Split(dr, "-")

	if len(drs) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("can't parse daterange %s", dr)
	}

	starttime, err := time.Parse(daterangeformat, drs[0])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("can't parse daterange %s: %v", dr, err)
	}

	endtime, err := time.Parse(daterangeformat, drs[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("can't parse daterange %s: %v", dr, err)
	}

	return starttime, endtime, nil
}

var daterange = flag.String("d", "19000101-21001231", "[year month date - year month date) e.g. 20170101-20171241")
var verbose = flag.Bool("v", false, "verbose output for debugging")

func main() {
	// Setup options before calling Parse()
	flag.Parse()

	starttime, endtime, err := parseDateRange(*daterange)
	if err != nil {
		log.Fatal(err)
	}
	if *verbose {
		log.Println(starttime.String(), endtime.String())
	}

	// for each arg
	for _, f := range flag.Args() {
		fd, err := os.Open(f)
		if err != nil {
			log.Printf("can't open %q: %v", f, err)
			continue
		}

		recs, err := parseJsonFile(fd)
		if err != nil {
			log.Printf("can't parse %q: %v", f, err)
			fd.Close()
			continue
		}
		fd.Close()

		if *verbose {
			for i := 0; i < 10 && i < len(recs); i++ {
				log.Println(i, recs[i])
			}
		}

		locations := convertFormat(recs, starttime, endtime)

		if *verbose {
			for i, loc := range locations {
				if loc.timestamp.After(starttime) && loc.timestamp.Before(endtime) {
					log.Printf("[%d] %s %f %f", i, loc.timestamp.String(), loc.latitude, loc.longitude)
				}
			}
		}

//		of := f + ".csv"
		of := f + ".gpx"
		ofd, err := os.Create(of)
		if err != nil {
			log.Printf("can't open output %q: %v", of, err)
			continue
		}

		// TODO(rjk): control via CLI
/*
		if err := convertToCsv(ofd, locations, starttime, endtime); err != nil {
			log.Println("convertToCsv failed:", err)
			ofd.Close()
			continue
		}
*/
		
		if err := convertToGpx(ofd, locations, starttime, endtime); err != nil {
			log.Println("convertToGpx failed:", err)
			ofd.Close()
			continue
		}
		ofd.Close()
	}

}

