package geo

import (
	"github.com/oschwald/geoip2-golang"
	"net"
)

// Geo is the parsed record geo data
type Geo struct {
	ContinentCode      string  `de:"Continent Code" en:"Continent Code" es:"Código Continente" fr:"Code Continent" ja:"大陸コード" pt-BR:"Código do Continente" ru:"Континентный код" zh-CN:"洲代码"`
	ContinentName      string  `de:"Continent naam" en:"Continent Name" es:"Nombre del continente" fr:"Nom du continent" ja:"大陸名" pt-BR:"Nome do Continente" ru:"Название континента" zh-CN:"洲名"`
	CountryIsoCode     string  `de:"ISO-Ländercode" en:"Country ISO code" es:"código de país ISO" fr:"Code ISO du pays" ja:"国のISOコード" pt-BR:"código de país ISO" ru:"Код страны ISO" zh-CN:"国家ISO代码"`
	CountryName        string  `de:"Ländername" en:"Country name" es:"Nombre del país" fr:"Nom du pays" ja:"国名" pt-BR:"Nome do país" ru:"Имя страны" zh-CN:"国家"`
	SubdivisionIsoCode string  `de:"Vorort ISO-Code" en:"Subdivision ISO code" es:"Barrio código ISO" fr:"Barrio código ISO" ja:"区ISOコード" pt-BR:"código ISO subdivisão" ru:"Код подразделения ISO" zh-CN:"省级ISO代码"`
	SubdivisionName    string  `de:"Vorort Name" en:"Subdivision name" es:"nombre de la subdivisión" fr:"Nom de lotissement" ja:"区名" pt-BR:"nome do Bairro" ru:"название Район" zh-CN:"省级名称"`
	CityName           string  `de:"Stadtname" en:"City name" es:"Nombre de la ciudad" fr:"Nom de Ville" ja:"市名" pt-BR:"Nome da Cidade" ru:"Название города" zh-CN:"城市名称"`
	PostalCode         string  `de:"Postleitzahl" en:"Postal code" es:"Código postal" fr:"Code postal" ja:"郵便番号" pt-BR:"Código postal" ru:"Почтовый индекс" zh-CN:"邮政编码"`
	Latitude           float64 `de:"Breite" en:"Latitude" es:"Latitud" fr:"Latitude" ja:"緯度" pt-BR:"Latitude" ru:"широта" zh-CN:"纬度"`
	Longitude          float64 `de:"Länge" en:"Longitude" es:"Longitud" fr:"Longitude" ja:"経度" pt-BR:"Longitude" ru:"долгота" zh-CN:"经度"`
	TimeZone           string  `de:"Tijdzone" en:"Time Zone" es:"Horzono" fr:"Fuseau horaire" ja:"タイムゾーン" pt-BR:"Fuso horário" ru:"Часовой пояс" zh-CN:"时区"`
}

// RetData is the struct to be returned by Query
type RetData struct {
	Host    string
	IP      net.IP
	Error   bool
	Message string
	Geo     *Geo
}

// Open an mmdb file by MaxMind
func Open(file string) (*geoip2.Reader, error) {
	return geoip2.Open(file)
}

func parseRecord(record *geoip2.City, language string) *Geo {
	result := Geo{
		ContinentCode:  record.Continent.Code,
		ContinentName:  record.Continent.Names[language],
		CountryIsoCode: record.Country.IsoCode,
		CountryName:    record.Country.Names[language],
		CityName:       record.City.Names[language],
		PostalCode:     record.Postal.Code,
		Latitude:       record.Location.Latitude,
		Longitude:      record.Location.Longitude,
		TimeZone:       record.Location.TimeZone,
	}
	if len(record.Subdivisions) > 0 {
		result.SubdivisionIsoCode = record.Subdivisions[0].IsoCode
		result.SubdivisionName = record.Subdivisions[0].Names[language]
	}
	return &result
}

// Query a host for record data
func Query(host string, db *geoip2.Reader, language string) []*RetData {
	ips, err := net.LookupIP(host)
	if err != nil {
		return []*RetData{
			{
				Host:    host,
				Error:   true,
				Message: "can't resolve host",
			},
		}
	}

	rets := make([]*RetData, 0)

MainLoop:
	for i, ip := range ips {
		// distinct results
		for j := 0; j < i; j++ {
			if ip.Equal(ips[j]) {
				continue MainLoop
			}
		}

		// query geo data
		ret := &RetData{
			Host:  host,
			IP:    ip,
			Error: false,
		}
		record, err := db.City(ip)
		if err != nil {
			ret.Error = true
			ret.Message = err.Error()
		} else {
			ret.Geo = parseRecord(record, language)
		}

		rets = append(rets, ret)
	}

	return rets
}
