// backend/platforms/international/international_platforms.go
package international

var InternationalPlatforms = map[string]PlatformConfig{
    "instagram": {
        Name:         "Instagram",
        BaseURL:      "https://instagram.com",
        SearchMethod: "phone_email_username",
        Priority:     1,
        Categories:   []string{"social", "photos", "videos"},
    },
    "telegram": {
        Name:         "Telegram",
        BaseURL:      "https://telegram.org",
        SearchMethod: "phone",
        Priority:     1,
        Categories:   []string{"messaging", "social"},
    },
    "whatsapp": {
        Name:         "WhatsApp",
        BaseURL:      "https://whatsapp.com",
        SearchMethod: "phone",
        Priority:     1,
        Categories:   []string{"messaging"},
    },
    "facebook": {
        Name:         "Facebook",
        BaseURL:      "https://facebook.com",
        SearchMethod: "phone_email",
        Priority:     1,
        Categories:   []string{"social", "networking"},
    },
    "twitter": {
        Name:         "Twitter",
        BaseURL:      "https://twitter.com",
        SearchMethod: "phone_email",
        Priority:     1,
        Categories:   []string{"social", "microblogging"},
    },
    "linkedin": {
        Name:         "LinkedIn",
        BaseURL:      "https://linkedin.com",
        SearchMethod: "email",
        Priority:     1,
        Categories:   []string{"professional", "networking"},
    },
    "snapchat": {
        Name:         "Snapchat",
        BaseURL:      "https://snapchat.com",
        SearchMethod: "phone",
        Priority:     2,
        Categories:   []string{"social", "messaging"},
    },
    "tiktok": {
        Name:         "TikTok",
        BaseURL:      "https://tiktok.com",
        SearchMethod: "phone_email",
        Priority:     1,
        Categories:   []string{"social", "videos"},
    },
    "discord": {
        Name:         "Discord",
        BaseURL:      "https://discord.com",
        SearchMethod: "email",
        Priority:     2,
        Categories:   []string{"gaming", "messaging"},
    },
    // ... 61 more international platforms
}
