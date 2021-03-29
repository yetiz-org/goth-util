package validate

import (
	"strconv"

	"github.com/biter777/countries"
)

var ExistCountryCodeMap = map[string]interface{}{"441624": nil, "441534": nil, "441481": nil, "1876": nil, "1869": nil, "1868": nil, "1809": nil, "1829": nil, "1849": nil, "1787": nil, "1939": nil, "1784": nil, "1767": nil, "1758": nil, "1721": nil, "1684": nil, "1671": nil, "1670": nil, "1664": nil, "1649": nil, "1473": nil, "1441": nil, "1345": nil, "1340": nil, "1284": nil, "1268": nil, "1264": nil, "1246": nil, "1242": nil, "998": nil, "996": nil, "995": nil, "994": nil, "993": nil, "992": nil, "977": nil, "976": nil, "975": nil, "974": nil, "973": nil, "972": nil, "971": nil, "970": nil, "968": nil, "967": nil, "966": nil, "965": nil, "964": nil, "963": nil, "962": nil, "961": nil, "960": nil, "886": nil, "880": nil, "856": nil, "855": nil, "853": nil, "852": nil, "850": nil, "692": nil, "691": nil, "690": nil, "689": nil, "688": nil, "687": nil, "686": nil, "685": nil, "683": nil, "682": nil, "681": nil, "680": nil, "679": nil, "678": nil, "677": nil, "676": nil, "675": nil, "674": nil, "673": nil, "672": nil, "670": nil, "599": nil, "598": nil, "597": nil, "595": nil, "593": nil, "592": nil, "591": nil, "590": nil, "509": nil, "508": nil, "507": nil, "506": nil, "505": nil, "504": nil, "503": nil, "502": nil, "501": nil, "500": nil, "423": nil, "421": nil, "420": nil, "389": nil, "387": nil, "386": nil, "385": nil, "383": nil, "382": nil, "381": nil, "380": nil, "379": nil, "378": nil, "377": nil, "376": nil, "375": nil, "374": nil, "373": nil, "372": nil, "371": nil, "370": nil, "359": nil, "358": nil, "357": nil, "356": nil, "355": nil, "354": nil, "353": nil, "352": nil, "351": nil, "350": nil, "299": nil, "298": nil, "297": nil, "291": nil, "290": nil, "269": nil, "268": nil, "267": nil, "266": nil, "265": nil, "264": nil, "263": nil, "262": nil, "261": nil, "260": nil, "258": nil, "257": nil, "256": nil, "255": nil, "254": nil, "253": nil, "252": nil, "251": nil, "250": nil, "249": nil, "248": nil, "246": nil, "245": nil, "244": nil, "243": nil, "242": nil, "241": nil, "240": nil, "239": nil, "238": nil, "237": nil, "236": nil, "235": nil, "234": nil, "233": nil, "232": nil, "231": nil, "230": nil, "229": nil, "228": nil, "227": nil, "226": nil, "225": nil, "224": nil, "223": nil, "222": nil, "221": nil, "220": nil, "218": nil, "216": nil, "213": nil, "212": nil, "211": nil, "98": nil, "95": nil, "94": nil, "93": nil, "92": nil, "91": nil, "90": nil, "86": nil, "84": nil, "82": nil, "81": nil, "66": nil, "65": nil, "64": nil, "63": nil, "62": nil, "61": nil, "60": nil, "58": nil, "57": nil, "56": nil, "55": nil, "54": nil, "53": nil, "52": nil, "51": nil, "49": nil, "48": nil, "47": nil, "46": nil, "45": nil, "44": nil, "43": nil, "41": nil, "40": nil, "39": nil, "36": nil, "34": nil, "33": nil, "32": nil, "31": nil, "30": nil, "27": nil, "20": nil, "7": nil, "1": nil}

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
