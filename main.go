package main

import (
	"encoding/csv"
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

	log.Println(starttime, endtime)

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

func main() {
	log.Println("hello")

	// Setup options before calling Parse()
	flag.Parse()

	starttime, endtime, err := parseDateRange(*daterange)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(starttime.String(), endtime.String())

	// for each arg
	for _, f := range flag.Args() {
		log.Println(f)
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

		for i := 0; i < 10 && i < len(recs); i++ {
			log.Println(i, recs[i])
		}

		locations := convertFormat(recs, starttime, endtime)

		for i, loc := range locations {
			if loc.timestamp.After(starttime) && loc.timestamp.Before(endtime) {
				log.Printf("[%d] %s %f %f", i, loc.timestamp.String(), loc.latitude, loc.longitude)
			}
		}

		of := f + ".csv"
		ofd, err := os.Create(of)
		if err != nil {
			log.Printf("can't open output %q: %v", of, err)
			continue
		}

		if err := convertToCsv(ofd, locations, starttime, endtime); err != nil {
			log.Println("convertToCsv failed:", err)
			ofd.Close()
			continue
		}
		ofd.Close()
	}

}

func convertToCsv(ofd io.Writer, locs []LocationValue, starttime, endtime time.Time) error {
	csvout := csv.NewWriter(ofd)

	// write me a header line
	if err := csvout.Write([]string{
		"No",
		"Latitude",
		"Longitude",
		"UTC date",
		"UTC time",
	}); err != nil {
		return fmt.Errorf("convertToCsv write: %v", err)
	}

	for i, loc := range locs {
		if loc.timestamp.After(starttime) && loc.timestamp.Before(endtime) {

			if err := csvout.Write([]string{
				fmt.Sprintf("%d", i+1),
				fmt.Sprintf("%f", loc.latitude),
				fmt.Sprintf("%f", loc.longitude),
				loc.timestamp.UTC().Format("2006/01/02"),
				loc.timestamp.UTC().Format("15:04:05.000"),
			}); err != nil {
				return fmt.Errorf("convertToCsv write: %v", err)
			}
		}
	}
	csvout.Flush()
	return csvout.Error()
}

/*

	No // the number of the record (do I need? is easy.)
      lat =      Latitude
      lon =      Longitude

      utc_d =    UTC date
      utc_t =    UTC time

      date =     Date (yyyy/mm/dd)
      time =     Time (hh:mm:ss[.msec])


alt =      Elevation (in meters) of the point. Add "ft" or "feet" for feet.
      arch =     Geocache archived flag
      avail =    Geocache available flag
      bng_e =    British National Grid's easting
      bng =      full coordinate in BNG format (zone easting northing)
      bng_pos =  full coordinate in BNG format (zone easting northing)
      bng_n =    British National Grid's northing
      bng_z =    British National Grid's zone
      caden =    Cadence
      comment =  Notes
      cont =     Geocache container
      cour =     Heading / Course true
      depth =    Depth (in meters).  Add "ft" or "feet" for feet.
      desc =     Description
      diff =     Geocache difficulty
      ele =      Elevation (in meters) of the point. Add "ft" or "feet" for feet.
      e/w =      'e' for eastern hemisphere, 'w' for western
      exported = Geocache export date
      found =    Geocache last found date
      fix =      3d, 2d, etc.
      gcid =     Geocache cache id
      geschw =   Geschwindigkeit (speed)
      hdop =     Horizontal dilution of precision
      head =     Heading / Course true
      heart =    Heartrate
      height =   Elevation (in meters) of the point
      hint =     Geocache cache hint
      icon =     Symbol (icon) name
      name =     Waypoint name ("Shortname")
      n/s =      'n' for northern hemisphere, 's' for southern
      notes =    Notes
      pdop =     Position dilution of precision
      placer =   Geocache placer
      placer_id =Geocache placer id
      power =    Cycling power (in Watts)
      prox =     Proximity (in meters).  Add "ft" or "feet" for feet.
      sat =      Number of sats used for fix
      speed =    Speed, in meters per second. (See below)
      symb =     Symbol (icon) name
      tempf =    Temperature (degrees Fahrenheit)
      temp =     Temperature (degrees Celsius)
      terr =     Geocache terrain
      time =     Time (hh:mm:ss[.msec])
      type =     Geocache cache type
      url =      URL
      utc_d =    UTC date
      utc_t =    UTC time
      utm_c =    UTM zone character
      utm_e =    UTM easting
      utm =      full coordinate in UTM format (zone zone-ch easting northing)
      utm_pos =  full coordinate in UTM format (zone zone-ch easting northing)
      utm_n =    UTM northing
      utm_z =    UTM zone
      vdop =     Vertical dilution of precision
      x =        Longitude
      x_pos =    Longitude
      y =        Latitude
      y_pos =    Latitude
      z =        Altitude (elevation).  See "elevation".

No,Latitude,Longitude,Name,Description,Symbol,Date,Time,Format
1,51.075139,12.463689,"578","578","Waypoint",2005/04/26,16:27:23,"gdb"
2,51.081104,12.465277,"579","579","Waypoint",2005/04/26,16:27:23,"gdb"
3,50.844126,12.408757,"Gosel","Gosel","Exit",2005/02/26,10:10:47,"gpx"
4,50.654763,12.204957,"Greiz",,"Exit",2005/02/26,09:57:04,"gpx"
*/
