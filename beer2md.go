package main

import (
    "bufio"
	"encoding/csv"
	"fmt"
    "io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/gosimple/slug"
    "golang.org/x/text/encoding/charmap"
    "golang.org/x/text/transform"
)

type Beer struct {
	id         string
	brewery_id string
	name       string
	cat_id     string
	style_id   string
	abv        string
	ibu        string
	srm        string
	upc        string
	filepath   string
	descript   string
	last_mod   string
}

type Brewery struct {
	id       string
	name     string
	address1 string
	address2 string
	city     string
	state    string
	code     string
	country  string
	phone    string
	website  string
	filepath string
	descript string
	last_mod string
}

type Geocode struct {
	id         string
	brewery_id string
	latitude   string
	longitude  string
	accuracy   string
}

type Category struct {
	id         string
	cat_name string
	last_mod   string
}

type Style struct {
	id         string
	cat_id string
	style_name   string
	last_mod  string
}

func createIndexFile(brewery Brewery, geocode Geocode, template *template.Template) {
    brewerySlug := slug.MakeLang(brewery.name, "en")
	indexFile := "breweries/" + brewerySlug + "/_index.md"
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		f, err := os.Create(indexFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		data := map[string]interface{}{
			"title":       brewery.name,
			"url":         "/" + brewerySlug + "/",
            "latitude":    geocode.latitude,
            "longitude":   geocode.longitude,
			"address1":    brewery.address1,
			"address2":    brewery.address2,
			"city":        brewery.city,
			"state":       brewery.state,
			"code":        brewery.code,
			"country":     brewery.country,
			"phone":       brewery.phone,
			"website":     brewery.website,
            "description": brewery.descript,

		}
		if err = template.Execute(f, data); err != nil {
			panic(err)
		}
		f.Close()
	}
}

func createElementFile(brewerySlug string, beer Beer, category Category, style Style, template *template.Template) {
    nameSlug := slug.MakeLang(beer.name, "en")
	elementFile := "breweries/" + brewerySlug + "/" + nameSlug + ".md"
	if _, err := os.Stat(elementFile); !os.IsNotExist(err) {
		re := regexp.MustCompile(`\d+$`)
		i := 2
		nameSlug = fmt.Sprintf("%s-%d", nameSlug, i)
		for {
			elementFile = "breweries/" + brewerySlug + "/" + nameSlug + ".md"
			if _, err = os.Stat(elementFile); os.IsNotExist(err) {
				break
			}
			i++
			nameSlug = re.ReplaceAllString(nameSlug, strconv.Itoa(i))
		}
	}
	f, err := os.Create(elementFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := map[string]interface{}{
		"title":       strings.Replace(beer.name, "\"", "", -1),
		"url":         "/" + brewerySlug + "/" + nameSlug + "/",
        "category":    category.cat_name,
        "style":       style.style_name,
        "abv":         beer.abv,
        "ibu":         beer.ibu,
        "srm":         beer.srm,
        "upc":         beer.upc,
        "description": beer.descript,

	}
	if err = template.Execute(f, data); err != nil {
		panic(err)
	}
	f.Close()
}

func createBeerList(data [][]string) []Beer {
    var beerList []Beer
    for i, line := range data {
        if i > 0 { // omit header line
            var rec Beer
            for j, field := range line {
                if j == 0 {
                    rec.id = field
                } else if j == 1 {
                    rec.brewery_id = field
                } else if j == 2 {
                    rec.name = field
                } else if j == 3 {
                    rec.cat_id = field
                } else if j == 4 {
                    rec.style_id = field
                } else if j == 5 {
                    rec.abv = field
                } else if j == 6 {
                    rec.ibu = field
                } else if j == 7 {
                    rec.srm = field
                } else if j == 8 {
                    rec.upc = field
                } else if j == 9 {
                    rec.filepath = field
                } else if j == 10 {
                    rec.descript = field
                } else if j == 11 {
                    rec.last_mod = field
                }
            }
            beerList = append(beerList, rec)
        }
    }
    return beerList
}

func createBreweryList(data [][]string) []Brewery {
    var breweryList []Brewery
    for i, line := range data {
        if i > 0 { // omit header line
            var rec Brewery
            for j, field := range line {
                if j == 0 {
                    rec.id = field
				} else if j == 1 {
					rec.name = field
                } else if j == 2 {
                    rec.address1 = field
                } else if j == 3 {
                    rec.address2 = field
                } else if j == 4 {
                    rec.city = field
                } else if j == 5 {
                    rec.state = field
                } else if j == 6 {
                    rec.code = field
                } else if j == 7 {
                    rec.country = field
                } else if j == 8 {
                    rec.phone = field
                } else if j == 9 {
                    rec.website = field
                } else if j == 10 {
                    rec.filepath = field
                } else if j == 11 {
                    rec.descript = field
                } else if j == 12 {
                    rec.last_mod = field
                }
            }
            breweryList = append(breweryList, rec)
        }
    }
    return breweryList
}

func createGeocodeList(data [][]string) []Geocode {
    var geocodeList []Geocode
    for i, line := range data {
        if i > 0 { // omit header line
            var rec Geocode
            for j, field := range line {
                if j == 0 {
                    rec.id = field
				} else if j == 1 {
					rec.brewery_id = field
                } else if j == 2 {
                    rec.latitude = field
                } else if j == 3 {
                    rec.longitude = field
                } else if j == 4 {
                    rec.accuracy = field
                }
            }
            geocodeList = append(geocodeList, rec)
        }
    }
    return geocodeList
}

func createCategoryList(data [][]string) []Category {
    var categoryList []Category
    for i, line := range data {
        if i > 0 { // omit header line
            var rec Category
            for j, field := range line {
                if j == 0 {
                    rec.id = field
				} else if j == 1 {
					rec.cat_name = field
                } else if j == 2 {
                    rec.last_mod = field
                }
            }
            categoryList = append(categoryList, rec)
        }
    }
    return categoryList
}

func createStyleList(data [][]string) []Style {
    var styleList []Style
    for i, line := range data {
        if i > 0 { // omit header line
            var rec Style
            for j, field := range line {
                if j == 0 {
                    rec.id = field
				} else if j == 1 {
					rec.cat_id = field
                } else if j == 2 {
                    rec.style_name = field
                } else if j == 3 {
                    rec.last_mod = field
                }
            }
            styleList = append(styleList, rec)
        }
    }
    return styleList
}

func ReadLines(filename string, decode bool) [][]string {
    var r io.Reader
    f, err := os.Open(filename)
    if err != nil {
        return nil
    }
    defer f.Close()

    if decode {
        decodingReader := transform.NewReader(f, charmap.Windows1252.NewDecoder())

        var lines string

        scanner := bufio.NewScanner(decodingReader)
        for scanner.Scan() {
            lines += scanner.Text() + "\n"
        }
        r = strings.NewReader(lines)
    } else {
        r = f
    }
    csvReader := csv.NewReader(r)
    data, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal(err)
    }

    return data
}

func main() {
	// Beers
    data := ReadLines("beers.csv", false)
    beerList := createBeerList(data)

	// Breweries
    data = ReadLines("breweries.csv", true)
    breweryList := createBreweryList(data)

	// Geocode
    data = ReadLines("breweries_geocode.csv", false)
    geocodeList := createGeocodeList(data)

	// Categories
    data = ReadLines("categories.csv", false)
    categoryList := createCategoryList(data)

	// Styles
    data = ReadLines("styles.csv", false)
    styleList := createStyleList(data)

	// Create templates
	indexTmpl := `---
title: {{ .title }}
url: {{ .url }}
latitude: {{ with .latitude }}{{ . }}{{ end }}
longitude: {{ with .longitude }}{{ . }}{{ end }}
address1: {{ with .address1 }}{{ . }}{{ end }}
address2: {{ with .address2 }}{{ . }}{{ end }}
city: {{ with .city }}{{ . }}{{ end }}
state: {{ with .state }}{{ . }}{{ end }}
code: {{ with .code }}{{ . }}{{ end }}
country: {{ with .country }}{{ . }}{{ end }}
phone: {{ with .phone }}{{ . }}{{ end }}
website: {{ with .website }}{{ . }}{{ end }}
---
{{ with .description }}{{ . }}{{ end }}
`
	indexTemplate := template.Must(template.New("index").Parse(indexTmpl))
	mdTmpl := `---
title: "{{ .title }}"
url: {{ .url }}
category: {{ with .category }}{{ . }}{{ end }}
style: {{ with .style }}{{ . }}{{ end }}
abv: {{ with .abv }}{{ . }}{{ end }}
ibu: {{ with .ibu }}{{ . }}{{ end }}
srm: {{ with .srm }}{{ . }}{{ end }}
upc: {{ with .upc }}{{ . }}{{ end }}
---
{{ with .description }}{{ . }}{{ end }}
`
	mdTemplate := template.Must(template.New("markdown").Parse(mdTmpl))

    // Parse data
    for _, beer := range beerList {
        foundGeocode := Geocode{}
        foundBrewery := Brewery{}
        foundCategory := Category{}
        foundStyle := Style{}
        for _, brewery := range breweryList {
            for _, geocode := range geocodeList {
                if brewery.id == geocode.brewery_id {
                    foundGeocode = geocode
                    break;
                }
            
            }
            if beer.brewery_id == brewery.id {
                foundBrewery = brewery
                break
            }
        }
        for _, category := range categoryList {
            if beer.cat_id == category.id {
                foundCategory = category
                break
            }
            
        }
        for _, style := range styleList {
            if beer.style_id == style.id {
                foundStyle = style
                break
            }
        }
        if foundBrewery.name != "" {
            brewerySlug := slug.MakeLang(foundBrewery.name, "en")
            os.MkdirAll("breweries/"+brewerySlug, 0755)
		    createIndexFile(foundBrewery, foundGeocode, indexTemplate)
		    createElementFile(brewerySlug, beer, foundCategory, foundStyle, mdTemplate)
        }
    }
}
