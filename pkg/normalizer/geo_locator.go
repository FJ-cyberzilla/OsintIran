// pkg/normalizer/geo_locator.go
package normalizer

import (
    "github.com/nyaruka/phonenumbers"
)

type GeoLocator struct {
    countryMapping map[int32]string
    regionMapping  map[string]string
}

func NewGeoLocator() *GeoLocator {
    return &GeoLocator{
        countryMapping: map[int32]string{
            98: "IR", // Iran
            1:  "US", // United States
            44: "GB", // United Kingdom
            49: "DE", // Germany
            33: "FR", // France
            // Add all country codes...
        },
        regionMapping: loadRegionMapping(),
    }
}

func (gl *GeoLocator) GetRegion(num *phonenumbers.PhoneNumber) string {
    countryCode := *num.CountryCode
    
    // For Iran, detect specific regions based on area codes
    if countryCode == 98 {
        return gl.getIranianRegion(num)
    }
    
    // Generic region detection for other countries
    return gl.countryMapping[countryCode]
}

func (gl *GeoLocator) GetTimezone(num *phonenumbers.PhoneNumber) string {
    countryCode := *num.CountryCode
    
    // Map country codes to timezones
    timezoneMap := map[int32]string{
        98: "Asia/Tehran",    // Iran
        1:  "America/New_York", // US
        44: "Europe/London",  // UK
        49: "Europe/Berlin",  // Germany
        33: "Europe/Paris",   // France
    }
    
    if tz, exists := timezoneMap[countryCode]; exists {
        return tz
    }
    
    return "Unknown"
}

func (gl *GeoLocator) getIranianRegion(num *phonenumbers.PhoneNumber) string {
    nationalNumber := phonenumbers.GetNationalSignificantNumber(num)
    
    if len(nationalNumber) < 3 {
        return "Unknown"
    }
    
    areaCode := nationalNumber[:3]
    
    // Iranian mobile operator regions
    iranRegionMap := map[string]string{
        "091": "Tehran",      // MCI
        "0910": "Tehran",     // MCI
        "0911": "Tehran",     // MCI
        "0912": "Tehran",     // MCI
        "0913": "Tehran",     // MCI
        "0914": "Tehran",     // MCI
        "0915": "Tehran",     // MCI
        "0916": "Tehran",     // MCI
        "0917": "Tehran",     // MCI
        "0918": "Tehran",     // MCI
        "0919": "Tehran",     // MCI
        "093": "Tehran",      // MTN Irancell
        "0930": "Tehran",     // MTN Irancell
        "0931": "Tehran",     // MTN Irancell
        "0932": "Tehran",     // MTN Irancell
        "0933": "Tehran",     // MTN Irancell
        "0934": "Tehran",     // MTN Irancell
        "0935": "Tehran",     // MTN Irancell
        "0936": "Tehran",     // MTN Irancell
        "0937": "Tehran",     // MTN Irancell
        "0938": "Tehran",     // MTN Irancell
        "0939": "Tehran",     // MTN Irancell
        "092": "Tehran",      // Rightel
        "0920": "Tehran",     // Rightel
        "0921": "Tehran",     // Rightel
        "0922": "Tehran",     // Rightel
        "099": "Tehran",      // Shatel
        "0990": "Tehran",     // Shatel
        "0991": "Tehran",     // Shatel
        // Add more specific regional mappings
    }
    
    if region, exists := iranRegionMap[areaCode]; exists {
        return region
    }
    
    return "Iran"
}

func loadRegionMapping() map[string]string {
    // Load from configuration file
    return map[string]string{
        "091": "MCI",
        "093": "MTN Irancell", 
        "092": "Rightel",
        "099": "Shatel",
        "090": "Irancell",
        "094": "Irancell",
    }
}
