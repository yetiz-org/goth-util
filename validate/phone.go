package validate

import (
	"strconv"

	"github.com/biter777/countries"
)

var ExistCountryCodeMap = map[string]struct{}{"441624": struct{}{}, "441534": struct{}{}, "441481": struct{}{}, "1876": struct{}{}, "1869": struct{}{}, "1868": struct{}{}, "1809": struct{}{}, "1829": struct{}{}, "1849": struct{}{}, "1787": struct{}{}, "1939": struct{}{}, "1784": struct{}{}, "1767": struct{}{}, "1758": struct{}{}, "1721": struct{}{}, "1684": struct{}{}, "1671": struct{}{}, "1670": struct{}{}, "1664": struct{}{}, "1649": struct{}{}, "1473": struct{}{}, "1441": struct{}{}, "1345": struct{}{}, "1340": struct{}{}, "1284": struct{}{}, "1268": struct{}{}, "1264": struct{}{}, "1246": struct{}{}, "1242": struct{}{}, "998": struct{}{}, "996": struct{}{}, "995": struct{}{}, "994": struct{}{}, "993": struct{}{}, "992": struct{}{}, "977": struct{}{}, "976": struct{}{}, "975": struct{}{}, "974": struct{}{}, "973": struct{}{}, "972": struct{}{}, "971": struct{}{}, "970": struct{}{}, "968": struct{}{}, "967": struct{}{}, "966": struct{}{}, "965": struct{}{}, "964": struct{}{}, "963": struct{}{}, "962": struct{}{}, "961": struct{}{}, "960": struct{}{}, "886": struct{}{}, "880": struct{}{}, "856": struct{}{}, "855": struct{}{}, "853": struct{}{}, "852": struct{}{}, "850": struct{}{}, "692": struct{}{}, "691": struct{}{}, "690": struct{}{}, "689": struct{}{}, "688": struct{}{}, "687": struct{}{}, "686": struct{}{}, "685": struct{}{}, "683": struct{}{}, "682": struct{}{}, "681": struct{}{}, "680": struct{}{}, "679": struct{}{}, "678": struct{}{}, "677": struct{}{}, "676": struct{}{}, "675": struct{}{}, "674": struct{}{}, "673": struct{}{}, "672": struct{}{}, "670": struct{}{}, "599": struct{}{}, "598": struct{}{}, "597": struct{}{}, "595": struct{}{}, "593": struct{}{}, "592": struct{}{}, "591": struct{}{}, "590": struct{}{}, "509": struct{}{}, "508": struct{}{}, "507": struct{}{}, "506": struct{}{}, "505": struct{}{}, "504": struct{}{}, "503": struct{}{}, "502": struct{}{}, "501": struct{}{}, "500": struct{}{}, "423": struct{}{}, "421": struct{}{}, "420": struct{}{}, "389": struct{}{}, "387": struct{}{}, "386": struct{}{}, "385": struct{}{}, "383": struct{}{}, "382": struct{}{}, "381": struct{}{}, "380": struct{}{}, "379": struct{}{}, "378": struct{}{}, "377": struct{}{}, "376": struct{}{}, "375": struct{}{}, "374": struct{}{}, "373": struct{}{}, "372": struct{}{}, "371": struct{}{}, "370": struct{}{}, "359": struct{}{}, "358": struct{}{}, "357": struct{}{}, "356": struct{}{}, "355": struct{}{}, "354": struct{}{}, "353": struct{}{}, "352": struct{}{}, "351": struct{}{}, "350": struct{}{}, "299": struct{}{}, "298": struct{}{}, "297": struct{}{}, "291": struct{}{}, "290": struct{}{}, "269": struct{}{}, "268": struct{}{}, "267": struct{}{}, "266": struct{}{}, "265": struct{}{}, "264": struct{}{}, "263": struct{}{}, "262": struct{}{}, "261": struct{}{}, "260": struct{}{}, "258": struct{}{}, "257": struct{}{}, "256": struct{}{}, "255": struct{}{}, "254": struct{}{}, "253": struct{}{}, "252": struct{}{}, "251": struct{}{}, "250": struct{}{}, "249": struct{}{}, "248": struct{}{}, "246": struct{}{}, "245": struct{}{}, "244": struct{}{}, "243": struct{}{}, "242": struct{}{}, "241": struct{}{}, "240": struct{}{}, "239": struct{}{}, "238": struct{}{}, "237": struct{}{}, "236": struct{}{}, "235": struct{}{}, "234": struct{}{}, "233": struct{}{}, "232": struct{}{}, "231": struct{}{}, "230": struct{}{}, "229": struct{}{}, "228": struct{}{}, "227": struct{}{}, "226": struct{}{}, "225": struct{}{}, "224": struct{}{}, "223": struct{}{}, "222": struct{}{}, "221": struct{}{}, "220": struct{}{}, "218": struct{}{}, "216": struct{}{}, "213": struct{}{}, "212": struct{}{}, "211": struct{}{}, "98": struct{}{}, "95": struct{}{}, "94": struct{}{}, "93": struct{}{}, "92": struct{}{}, "91": struct{}{}, "90": struct{}{}, "86": struct{}{}, "84": struct{}{}, "82": struct{}{}, "81": struct{}{}, "66": struct{}{}, "65": struct{}{}, "64": struct{}{}, "63": struct{}{}, "62": struct{}{}, "61": struct{}{}, "60": struct{}{}, "58": struct{}{}, "57": struct{}{}, "56": struct{}{}, "55": struct{}{}, "54": struct{}{}, "53": struct{}{}, "52": struct{}{}, "51": struct{}{}, "49": struct{}{}, "48": struct{}{}, "47": struct{}{}, "46": struct{}{}, "45": struct{}{}, "44": struct{}{}, "43": struct{}{}, "41": struct{}{}, "40": struct{}{}, "39": struct{}{}, "36": struct{}{}, "34": struct{}{}, "33": struct{}{}, "32": struct{}{}, "31": struct{}{}, "30": struct{}{}, "27": struct{}{}, "20": struct{}{}, "7": struct{}{}, "1": struct{}{}}

func IsValidPhoneCountryCode(code string) bool {
	_, f := ExistCountryCodeMap[code]
	return f
}

func IsValidPhoneNumber(number string) bool {
	if len(number) < 2 {
		return false
	}

	if number[0] == '+' {
		number = number[1:]
	}

	return IsDigits(number)
}

func IsValidPhoneTerritoryCode(code string) bool {
	return len(code) == 2 && countries.ByName(code) != countries.Unknown
}

func IsCodeRegionEqual(countryCode string, territoryCode string) bool {
	if countryCode == "" || territoryCode == "" {
		return false
	}

	code, err := strconv.Atoi(countryCode)
	if err != nil {
		return false
	}

	countryListFromCountryCode := countries.CallCode(code).Countries()
	if countryListFromCountryCode[0] == countries.Unknown {
		return false
	}

	countryFromTerritoryCode := countries.ByName(territoryCode)
	for _, country := range countryListFromCountryCode {
		if country == countryFromTerritoryCode {
			return true
		}
	}

	return false
}
