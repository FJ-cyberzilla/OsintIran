// pkg/normalizer/carrier_detector.go
package normalizer

import (
    "github.com/nyaruka/phonenumbers"
)

type CarrierDetector struct {
    carrierMapping map[string]string
}

func NewCarrierDetector() *CarrierDetector {
    return &CarrierDetector{
        carrierMapping: loadCarrierMapping(),
    }
}

func (cd *CarrierDetector) GetCarrier(num *phonenumbers.PhoneNumber) string {
    countryCode := *num.CountryCode
    nationalNumber := phonenumbers.GetNationalSignificantNumber(num)
    
    // Special handling for Iranian carriers
    if countryCode == 98 {
        return cd.getIranianCarrier(nationalNumber)
    }
    
    // Generic carrier detection for other countries
    // This could be enhanced with external carrier lookup APIs
    return "Unknown"
}

func (cd *CarrierDetector) getIranianCarrier(nationalNumber string) string {
    if len(nationalNumber) < 3 {
        return "Unknown"
    }
    
    prefix := nationalNumber[:3]
    
    iranCarrierMap := map[string]string{
        "091": "MCI (Hamrah Aval)",
        "0910": "MCI (Hamrah Aval)",
        "0911": "MCI (Hamrah Aval)", 
        "0912": "MCI (Hamrah Aval)",
        "0913": "MCI (Hamrah Aval)",
        "0914": "MCI (Hamrah Aval)",
        "0915": "MCI (Hamrah Aval)",
        "0916": "MCI (Hamrah Aval)",
        "0917": "MCI (Hamrah Aval)",
        "0918": "MCI (Hamrah Aval)",
        "0919": "MCI (Hamrah Aval)",
        "093": "MTN Irancell",
        "0930": "MTN Irancell",
        "0931": "MTN Irancell",
        "0932": "MTN Irancell",
        "0933": "MTN Irancell",
        "0934": "MTN Irancell",
        "0935": "MTN Irancell",
        "0936": "MTN Irancell",
        "0937": "MTN Irancell",
        "0938": "MTN Irancell",
        "0939": "MTN Irancell",
        "092": "Rightel",
        "0920": "Rightel",
        "0921": "Rightel",
        "0922": "Rightel",
        "099": "Shatel",
        "0990": "Shatel",
        "0991": "Shatel",
        "090": "Irancell",
        "094": "Irancell",
        "0900": "Irancell",
        "0901": "Irancell",
        "0902": "Irancell",
        "0903": "Irancell",
        "0904": "Irancell",
        "0905": "Irancell",
    }
    
    if carrier, exists := iranCarrierMap[prefix]; exists {
        return carrier
    }
    
    return "Unknown Iranian Carrier"
}

func loadCarrierMapping() map[string]string {
    // Load comprehensive carrier mapping from config
    return map[string]string{
        // Global carrier mappings can be added here
    }
}
