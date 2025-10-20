// pkg/normalizer/phone_normalizer.go
package normalizer

import (
    "fmt"
    "regexp"
    "strings"

    "github.com/nyaruka/phonenumbers"
)

type PhoneNormalizer struct {
    defaultRegion string
    strictMode    bool
    geoLocator    *GeoLocator
    carrierDetector *CarrierDetector
}

type NormalizedPhone struct {
    Original       string `json:"original"`
    Normalized     string `json:"normalized"`     // E.164 format
    International  string `json:"international"`  // Human readable
    National       string `json:"national"`       // Local format
    CountryCode    int32  `json:"country_code"`
    Country        string `json:"country"`
    Region         string `json:"region"`
    Carrier        string `json:"carrier"`
    IsValid        bool   `json:"is_valid"`
    IsPossible     bool   `json:"is_possible"`
    Type           string `json:"type"`           // MOBILE, FIXED_LINE, etc.
    Timezone       string `json:"timezone"`
}

func NewPhoneNormalizer(defaultRegion string) *PhoneNormalizer {
    return &PhoneNormalizer{
        defaultRegion: defaultRegion,
        strictMode:    true,
        geoLocator:    NewGeoLocator(),
        carrierDetector: NewCarrierDetector(),
    }
}

// NormalizePhone - Main normalization function using libphonenumber
func (pn *PhoneNormalizer) NormalizePhone(input string, countryHint string) (*NormalizedPhone, error) {
    // Step 1: Clean and preprocess input
    cleaned := pn.cleanPhoneNumber(input)
    
    // Step 2: Detect country if not provided
    countryCode := countryHint
    if countryCode == "" {
        countryCode = pn.detectCountry(cleaned)
    }
    
    // Step 3: Parse with libphonenumber
    num, err := phonenumbers.Parse(cleaned, countryCode)
    if err != nil {
        if pn.strictMode {
            return nil, fmt.Errorf("failed to parse phone number: %w", err)
        }
        // Fallback normalization
        return pn.fallbackNormalize(cleaned, countryCode)
    }
    
    // Step 4: Validate number
    isValid := phonenumbers.IsValidNumber(num)
    isPossible := phonenumbers.IsPossibleNumber(num)
    
    // Step 5: Get additional information
    region := pn.geoLocator.GetRegion(num)
    carrier := pn.carrierDetector.GetCarrier(num)
    numberType := pn.getNumberType(num)
    timezone := pn.geoLocator.GetTimezone(num)
    
    // Step 6: Format in different standards
    normalized := &NormalizedPhone{
        Original:      input,
        Normalized:    phonenumbers.Format(num, phonenumbers.E164),
        International: phonenumbers.Format(num, phonenumbers.INTERNATIONAL),
        National:      phonenumbers.Format(num, phonenumbers.NATIONAL),
        CountryCode:   *num.CountryCode,
        Country:       countryCode,
        Region:        region,
        Carrier:       carrier,
        IsValid:       isValid,
        IsPossible:    isPossible,
        Type:          numberType,
        Timezone:      timezone,
    }
    
    return normalized, nil
}

// Batch normalization for multiple numbers
func (pn *PhoneNormalizer) NormalizeBatch(inputs []string, countryHint string) (map[string]*NormalizedPhone, []error) {
    results := make(map[string]*NormalizedPhone)
    var errors []error
    
    for _, input := range inputs {
        normalized, err := pn.NormalizePhone(input, countryHint)
        if err != nil {
            errors = append(errors, fmt.Errorf("failed to normalize %s: %w", input, err))
            continue
        }
        results[input] = normalized
    }
    
    return results, errors
}

// Clean phone number input
func (pn *PhoneNormalizer) cleanPhoneNumber(input string) string {
    // Remove all non-digit characters except + and spaces
    re := regexp.MustCompile(`[^\d+\s]`)
    cleaned := re.ReplaceAllString(input, "")
    
    // Remove spaces
    cleaned = strings.ReplaceAll(cleaned, " ", "")
    
    // Handle common Iranian number formats
    cleaned = pn.normalizeIranianFormats(cleaned)
    
    return strings.TrimSpace(cleaned)
}

// Special handling for Iranian phone numbers
func (pn *PhoneNormalizer) normalizeIranianFormats(phone string) string {
    // Convert Persian digits to English
    phone = pn.convertPersianDigits(phone)
    
    // Handle common Iranian formats
    patterns := map[string]*regexp.Regexp{
        "with_zero":   regexp.MustCompile(`^0098(\d{10})$`),
        "with_plus":   regexp.MustCompile(`^\+98(\d{10})$`),
        "without_country": regexp.MustCompile(`^0?9(\d{9})$`),
    }
    
    for format, pattern := range patterns {
        if matches := pattern.FindStringSubmatch(phone); matches != nil {
            switch format {
            case "with_zero":
                return "+98" + matches[1]
            case "with_plus":
                return phone // Already correct
            case "without_country":
                return "+98" + matches[1]
            }
        }
    }
    
    return phone
}

// Convert Persian/Arabic digits to English
func (pn *PhoneNormalizer) convertPersianDigits(input string) string {
    persianToEnglish := map[rune]rune{
        '۰': '0', '۱': '1', '۲': '2', '۳': '3', '۴': '4',
        '۵': '5', '۶': '6', '۷': '7', '۸': '8', '۹': '9',
        '٠': '0', '١': '1', '٢': '2', '٣': '3', '٤': '4',
        '٥': '5', '٦': '6', '٧': '7', '٨': '8', '٩': '9',
    }
    
    var result strings.Builder
    for _, char := range input {
        if replacement, exists := persianToEnglish[char]; exists {
            result.WriteRune(replacement)
        } else {
            result.WriteRune(char)
        }
    }
    
    return result.String()
}

// Fallback normalization when libphonenumber fails
func (pn *PhoneNormalizer) fallbackNormalize(phone, country string) (*NormalizedPhone, error) {
    // Simple regex-based normalization as fallback
    re := regexp.MustCompile(`^(?:\+?(\d{1,3})?[-. ]?)?\(?(\d{3})\)?[-. ]?(\d{3})[-. ]?(\d{4})$`)
    
    if matches := re.FindStringSubmatch(phone); matches != nil {
        countryCode := matches[1]
        if countryCode == "" {
            countryCode = "98" // Default to Iran
        }
        
        normalized := &NormalizedPhone{
            Original:   phone,
            Normalized: fmt.Sprintf("+%s%s%s%s", countryCode, matches[2], matches[3], matches[4]),
            Country:    country,
            IsValid:    false, // Mark as invalid since fallback was used
            IsPossible: true,
        }
        
        return normalized, nil
    }
    
    return nil, fmt.Errorf("could not normalize phone number: %s", phone)
}

func (pn *PhoneNormalizer) getNumberType(num *phonenumbers.PhoneNumber) string {
    switch phonenumbers.GetNumberType(num) {
    case phonenumbers.MOBILE:
        return "MOBILE"
    case phonenumbers.FIXED_LINE:
        return "FIXED_LINE"
    case phonenumbers.FIXED_LINE_OR_MOBILE:
        return "FIXED_LINE_OR_MOBILE"
    case phonenumbers.TOLL_FREE:
        return "TOLL_FREE"
    case phonenumbers.PREMIUM_RATE:
        return "PREMIUM_RATE"
    case phonenumbers.SHARED_COST:
        return "SHARED_COST"
    case phonenumbers.VOIP:
        return "VOIP"
    case phonenumbers.PERSONAL_NUMBER:
        return "PERSONAL_NUMBER"
    case phonenumbers.PAGER:
        return "PAGER"
    case phonenumbers.UAN:
        return "UAN"
    case phonenumbers.VOICEMAIL:
        return "VOICEMAIL"
    default:
        return "UNKNOWN"
    }
}
