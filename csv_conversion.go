package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

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
