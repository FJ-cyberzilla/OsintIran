// backend/platforms/iranian/iranian_platforms.go
package iranian

var IranianPlatforms = map[string]PlatformConfig{
    "rubika": {
        Name:         "روبیکا",
        BaseURL:      "https://rubika.ir",
        SearchMethod: "phone",
        Priority:     1,
        Requires:     []string{"mobile_app_simulation"},
    },
    "eitaa": {
        Name:         "ایتا",
        BaseURL:      "https://eitaa.com",
        SearchMethod: "phone",
        Priority:     1,
    },
    "soroush": {
        Name:         "سروش",
        BaseURL:      "https://splus.ir",
        SearchMethod: "phone",
        Priority:     1,
    },
    "bale": {
        Name:         "بله",
        BaseURL:      "https://bale.ai",
        SearchMethod: "phone",
        Priority:     1,
    },
    "gap": {
        Name:         "گپ",
        BaseURL:      "https://gap.im",
        SearchMethod: "phone",
        Priority:     2,
    },
    "shad": {
        Name:         "شاد",
        BaseURL:      "https://shad.ir",
        SearchMethod: "phone",
        Priority:     2,
    },
    "ita": {
        Name:         "ایتا",
        BaseURL:      "https://itaa.ir",
        SearchMethod: "phone",
        Priority:     2,
    },
    "bisphone": {
        Name:         "بیسفون",
        BaseURL:      "https://bisphone.com",
        SearchMethod: "phone",
        Priority:     2,
    },
    // ... 22 more Iranian platforms
}

// Iranian Mobile Operators
var IranianOperators = map[string]OperatorConfig{
    "mci": {
        Name:       "همراه اول",
        Code:       "091",
        APISupport: true,
        Endpoints:  []string{"https://mci.ir/api"},
    },
    "mtn": {
        Name:       "ایرانسل",
        Code:       "093",
        APISupport: true,
        Endpoints:  []string{"https://mtn.ir/api"},
    },
    "rightel": {
        Name:       "رایتل",
        Code:       "092",
        APISupport: true,
        Endpoints:  []string{"https://rightel.ir/api"},
    },
    "shatel": {
        Name:       "شاتل",
        Code:       "099",
        APISupport: true,
        Endpoints:  []string{"https://shatel.ir/api"},
    },
    "tci": {
        Name:       "ثابت ایران",
        Code:       "021",
        APISupport: true,
        Endpoints:  []string{"https://tci.ir/api"},
    },
}
