package certs

/*
	Sliver Implant Framework
	Copyright (C) 2021  Bishop Fox

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// changed subject to different cities than default main branch
// also change crypto

import (
	"crypto/x509/pkix"
	"fmt"
	insecureRand "math/rand"
	"strings"

	"github.com/bishopfox/sliver/server/codenames"
)

var (
	// Finally found a good use for Github Co-Pilot!
	// Country -> State -> Localities -> Street Addresses
	subjects = map[string]map[string]map[string][]string{
		"US": {
			"": {
				"": {
					"",
				},
			},
			"Arizona": {
				"Chandler":      {""},
				"Tempe":         {""},
				"Oro Valley":    {""},
				"Flagstaff":     {""},
				"Bullhead City": {""},
				"Queen Creek":   {""},
				"Yuma":          {""},
				"Glendale":      {""},
			},
			"California": {
				"Anaheim":          {""},
				"Clovis":           {""},
				"Pomona":           {""},
				"Simi Valley":      {""},
				"San Diego":        {""},
				"Lancaster":        {""},
				"West Covina":      {""},
				"Vista":            {""},
				"Mountain View":    {""},
				"Tracy":            {""},
				"Redwood City":     {""},
				"Chino":            {""},
				"Hesperia":         {""},
				"San Carlos":       {""},
				"San Leandro":      {""},
				"San Marcos":       {""},
				"Merced":           {""},
				"Santa Monica":     {""},
				"Westminister":     {""},
				"East Los Angeles": {""},
			},
			"Colorado": {
				"Erie":           {""},
				"Evans":          {""},
				"Grand Junction": {""},
				"Fort Collins":   {""},
				"Pueblo":         {""},
				"Loveland":       {""},
				"Denver":         {""},
			},
			"Connecticut": {
				"New Haven": {""},
				"Bethel":    {""},
				"Easton":    {""},
				"Hebron":    {""},
			},
			"Washington": {
				"Seattle":   {""},
				"Yakima":    {""},
				"Redmond":   {""},
				"Sammamish": {""},
			},
			"Florida": {
				"Miami":     {""},
				"Hollywood": {""},
				"Palm Bay":  {""},
				"Tampa":     {""},
			},
			"Illinois": {
				"Chicago": {""},
				"Elgin":   {""},
				"Decatur": {""},
				"Skokie":  {""},
			},
			"Indiana": {
				"Indianapolis": {""},
				"Gary":         {""},
				"Fishers":      {""},
				"Kokomo":       {""},
			},
			"Massachusetts": {
				"Alford": {""},
				"Worcester": {
					"",
					"State University",
				},
				"Colrain": {""},
				"Concord": {""},
			},
			"Michigan": {
				"Flint City": {""},
				"Troy City": {
					"",
					"Community College",
				},
				"Canton": {""},
				"Detroit": {
					"",
					"University",
				},
			},
			"Minnesota": {
				"Blaine":        {""},
				"Maple Grove":   {""},
				"Duluth":        {""},
				"Brooklyn Park": {""},
			},
			"New Jersey": {
				"Brick":       {""},
				"Passaic":     {""},
				"East Orange": {""},
				"Mercer": {
					"",
					"Princeton",
				},
			},
			"New York": {
				"New York": {""},
				"Geneva":   {""},
				"Dunkirk": {
					"",
					"Community College",
				},
				"Beacon": {""},
			},
			"North Carolina": {
				"Ferguson": {""},
				"Boone":    {""},
				"Fulton": {
					"",
					"Community College",
				},
				"Plattsburgh": {""},
			},
			"Ohio": {
				"Youngstown": {""},
				"Poland":     {""},
				"Dayton":     {""},
				"Lorain":     {""},
			},
		},
		"CA": {
			"": {
				"": {
					"",
				},
			},
			"Alberta": {
				"Calgary": {""},
				"Edmonton": {
					"",
					"University",
				},
				"Red Deer": {""},
				"Fort McMurray": {
					"",
					"University",
				},
			},
			"British Columbia": {
				"Vancouver": {""},
				"Victoria":  {""},
				"Kelowna":   {""},
				"Richmond":  {""},
			},
			"Manitoba": {
				"Winnipeg": {""},
				"Brandon":  {""},
				"Thompson": {""},
				"Portage la Prairie": {
					"",
					"University",
				},
			},
			"New Brunswick": {
				"Fredericton": {""},
				"Moncton":     {""},
				"Saint John":  {""},
				"Dieppe":      {""},
			},
			"Newfoundland and Labrador": {
				"St. John's": {""},
				"Mount Pearl": {
					"",
					"College",
				},
				"Conception Bay South": {""},
				"Paradise": {
					"",
					"College",
				},
			},
		},
		"JP": {
			"": {
				"": {
					"",
				},
			},
			"Aichi": {
				"Nagoya": {""},
				"Kasugai": {
					"",
					"University",
				},
				"Okazaki": {""},
				"Handa":   {""},
			},
			"Chiba": {
				"Chiba": {""},
				"Kashiwa": {
					"",
					"University",
				},
				"Funabashi": {""},
				"Kimitsu":   {""},
			},
		},
	}
)

func randomSubject(commonName string) *pkix.Name {
	country, province, locale, street := randomProvinceLocalityStreetAddress()
	return &pkix.Name{
		Organization:  randomOrganization(),
		Country:       country,
		Province:      province,
		Locality:      locale,
		StreetAddress: street,
		PostalCode:    randomPostalCode(country),
		CommonName:    commonName,
	}
}

func randomPostalCode(country []string) []string {
	// 1 in `n` will include a postal code
	// From my cursory view of a few TLS certs it seems uncommon to include this
	// in the distinguished name so right now it's set to 1/5
	const postalProbability = 5

	if len(country) == 0 {
		return []string{}
	}
	switch country[0] {

	case "US":
		// American postal codes are 5 digits
		switch insecureRand.Intn(postalProbability) {
		case 0:
			return []string{fmt.Sprintf("%05d", insecureRand.Intn(90000)+1000)}
		default:
			return []string{}
		}

	case "CA":
		// Canadian postal codes are weird and include letter/number combo's
		letters := "ABHLMNKGJPRSTVYX"
		switch insecureRand.Intn(postalProbability) {
		case 0:
			letter1 := string(letters[insecureRand.Intn(len(letters))])
			letter2 := string(letters[insecureRand.Intn(len(letters))])
			if insecureRand.Intn(2) == 0 {
				letter1 = strings.ToLower(letter1)
				letter2 = strings.ToLower(letter2)
			}
			return []string{
				fmt.Sprintf("%s%d%s", letter1, insecureRand.Intn(9), letter2),
			}
		default:
			return []string{}
		}
	}
	return []string{}
}

func randomProvinceLocalityStreetAddress() ([]string, []string, []string, []string) {
	country := randomCountry()
	state := randomState(country)
	locality := randomLocality(country, state)
	streetAddress := randomStreetAddress(country, state, locality)
	return []string{country}, []string{state}, []string{locality}, []string{streetAddress}
}

func randomCountry() string {
	keys := make([]string, 0, len(subjects))
	for k := range subjects {
		keys = append(keys, k)
	}
	return keys[insecureRand.Intn(len(keys))]
}

func randomState(country string) string {
	keys := make([]string, 0, len(subjects[country]))
	for k := range subjects[country] {
		keys = append(keys, k)
	}
	return keys[insecureRand.Intn(len(keys))]
}

func randomLocality(country string, state string) string {
	locales := subjects[country][state]
	keys := make([]string, 0, len(locales))
	for k := range locales {
		keys = append(keys, k)
	}
	return keys[insecureRand.Intn(len(keys))]
}

func randomStreetAddress(country string, state string, locality string) string {
	addresses := subjects[country][state][locality]
	return addresses[insecureRand.Intn(len(addresses))]
}

var (
	orgSuffixes = []string{
		"",
		"",
		"co",
		"llc",
		"inc",
		"corp",
		"ltd",
		"plc",
		"inc.",
		"corp.",
		"ltd.",
		"plc.",
		"co.",
		"llc.",
		"incorporated",
		"limited",
		"corporation",
		"company",
		"incorporated",
		"limited",
		"corporation",
		"company",
	}
)

func randomOrganization() []string {
	adjective, _ := codenames.RandomAdjective()
	noun, _ := codenames.RandomNoun()
	suffix := orgSuffixes[insecureRand.Intn(len(orgSuffixes))]

	var orgName string
	switch insecureRand.Intn(8) {
	case 0:
		orgName = strings.TrimSpace(fmt.Sprintf("%s %s, %s", adjective, noun, suffix))
	case 1:
		orgName = strings.TrimSpace(strings.ToLower(fmt.Sprintf("%s %s, %s", adjective, noun, suffix)))
	case 2:
		orgName = strings.TrimSpace(fmt.Sprintf("%s, %s", noun, suffix))
	case 3:
		orgName = strings.TrimSpace(strings.Title(fmt.Sprintf("%s %s, %s", adjective, noun, suffix)))
	case 4:
		orgName = strings.TrimSpace(strings.Title(fmt.Sprintf("%s %s", adjective, noun)))
	case 5:
		orgName = strings.TrimSpace(strings.ToLower(fmt.Sprintf("%s %s", adjective, noun)))
	case 6:
		orgName = strings.TrimSpace(strings.Title(fmt.Sprintf("%s", noun)))
	case 7:
		noun2, _ := codenames.RandomNoun()
		orgName = strings.TrimSpace(strings.ToLower(fmt.Sprintf("%s-%s", noun, noun2)))
	default:
		orgName = ""
	}

	return []string{orgName}
}
