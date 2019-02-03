package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/major1201/geo-cli"
	"github.com/major1201/goutils"
	"github.com/major1201/goutils/logging"
	"github.com/oschwald/geoip2-golang"
	"github.com/urfave/cli"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var logger = logging.New("GEO")

// AppVer means the project's version
const AppVer = "0.2.0"

func pipe(db *geoip2.Reader, language string) {
	ipReg, err := regexp.Compile(goutils.RegIPv4)
	if err != nil {
		logger.Fatal(err)
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		lineByte, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := string(lineByte)
		ip := ipReg.FindString(line)
		if ip == "" {
			fmt.Println(line)
		} else {
			data := geo.Query(ip, db, language)[0]
			if !data.Error {
				res := goutils.FilterBlankString([]string{data.Geo.CountryName, data.Geo.SubdivisionName, data.Geo.CityName})
				if len(res) == 0 {
					fmt.Println(line)
				} else {
					fmt.Printf("%v\t%v\n", line, strings.Join(res, " - "))
				}
			}
		}
	}
}

func display(dataArray []*geo.RetData, detail bool, language string) {
	for _, data := range dataArray {
		if data.Error {
			fmt.Printf("%v - %v\n", data.Host, data.Message)
			continue
		}

		var hostString string
		if data.IP.String() == data.Host {
			hostString = data.Host
		} else {
			hostString = fmt.Sprintf("%v(%v)", data.Host, data.IP.String())
		}

		if detail {
			fmt.Println(hostString)

			t := reflect.TypeOf(geo.Geo{})
			for i := 0; i < t.NumField(); i++ {
				fmt.Printf("  %v: %v\n", t.Field(i).Tag.Get(language), reflect.Indirect(reflect.ValueOf(data.Geo)).Field(i))
			}
			fmt.Println()
		} else {
			result := goutils.FilterBlankString([]string{data.Geo.ContinentName, data.Geo.CountryName, data.Geo.SubdivisionName, data.Geo.CityName})
			if len(result) == 0 {
				fmt.Printf("%v - not found\n", hostString)
			} else {
				fmt.Printf("%v - %v\n", hostString, strings.Join(result, " - "))
			}
		}
	}
}

func runApp(c *cli.Context) {
	db, err := geo.Open(c.String("mmdb-file"))
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	if c.NArg() == 0 {
		// read from stdin
		pipe(db, c.String("language"))
		return
	}

	dataArray := make([]*geo.RetData, 0)
	for _, host := range c.Args() {
		dataArray = append(dataArray, geo.Query(host, db, c.String("language"))...)
	}

	if c.Bool("json") {
		js, _ := json.Marshal(dataArray)
		fmt.Print(string(js))
		return
	}

	display(dataArray, c.Bool("detail"), c.String("language"))
}

func main() {
	// init logging
	logging.AddStdout(0)

	// parse flags
	app := getApp()
	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}
