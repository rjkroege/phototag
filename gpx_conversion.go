package main

// The HoudahGeo program will directly import GPX files successfully instead of the CSV.
// So make these directly.

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"time"
)

type TrkPoint struct {
	XMLName   xml.Name `xml:"trkpt"`
	Latitude  float64  `xml:"lat,attr"`
	Longitude float64  `xml:"lon,attr"`
	Timestamp string   `xml:"time"`
	Name      string   `xml:"name"`
}

type BoundingRect struct {
	XMLName xml.Name `xml:"bounds"`
	MinLat  float64  `xml:"minlat,attr"`
	MinLong float64  `xml:"minlon,attr"`
	MaxLat  float64  `xml:"maxlat,attr"`
	MaxLong float64  `xml:"maxlon,attr"`
}

type Gpx struct {
	XMLName   xml.Name     `xml:"gpx"`
	Timestamp string       `xml:"time"`
	Version   string       `xml:"version,attr"`
	Creator   string       `xml:"creator,attr"`
	Namespace string       `xml:"xmlns,attr"`
	Bounds    BoundingRect `xml:"bounds"`
	Points    []TrkPoint   `xml:"trk>trkseg>trkpt"`
}

func (br *BoundingRect) update(lat, long float64) {
	br.MinLat = math.Min(br.MinLat, lat)
	br.MinLong = math.Min(br.MinLong, long)
	br.MaxLat = math.Max(br.MaxLat, lat)
	br.MaxLong = math.Max(br.MaxLong, long)
}

func convertToGpx(ofd io.Writer, locs []LocationValue, starttime, endtime time.Time) error {
	if _, err := io.WriteString(ofd, xml.Header); err != nil {
		return fmt.Errorf("can't write header: %v", err)
	}

	gpxout := xml.NewEncoder(ofd)
	gpxout.Indent("", "	")

	points := make([]TrkPoint, 0)
	bounds := BoundingRect{
		MinLat:  math.MaxFloat64,
		MinLong: math.MaxFloat64,
		MaxLat:  -math.MaxFloat64,
		MaxLong: -math.MaxFloat64,
	}

	for i, loc := range locs {
		if loc.timestamp.After(starttime) && loc.timestamp.Before(endtime) {

			// do stuffs
			points = append(points, TrkPoint{
				Latitude:  loc.latitude,
				Longitude: loc.longitude,
				Timestamp: loc.timestamp.UTC().Format(time.RFC3339Nano),
				Name:      fmt.Sprintf("WPT.%d", i),
			})

			bounds.update(loc.latitude, loc.longitude)
		}
	}

	// TODO(rjk): probably do more stuffs

	gpx := Gpx{
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Version:   "1.0",
		Creator:   "phototag",
		Namespace: "http://www.topografix.com/GPX/1/0",
		Bounds:    bounds,
		Points:    points,
	}

	if err := gpxout.Encode(gpx); err != nil {
		return fmt.Errorf("can't encode to XML: %v", err)
	}

	return nil
}

/*

Need structure like this.

<?xml version="1.0" encoding="UTF-8"?>
<gpx version="1.0" creator="GPSBabel - http://www.gpsbabel.org" xmlns="http://www.topografix.com/GPX/1/0">
  <time>2021-08-17T23:50:25.010Z</time>
  <bounds minlat="38.707051000" minlon="-9.147488000" maxlat="38.714079000" maxlon="-9.127050000"/>
  <trk>
    <trkseg>
      <trkpt lat="38.709349000" lon="-9.146423000">
        <time>2018-10-08T09:05:55.956Z</time>
        <name>WPT</name>
      </trkpt>



*/
