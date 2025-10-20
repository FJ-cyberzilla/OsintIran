// backend/internal/services/social_scraper.go
package services

type SocialScraper struct {
    proxyManager    *ProxyManager
    platformClients map[string]PlatformClient
}

func NewSocialScraper() *SocialScraper {
    return &SocialScraper{
        proxyManager: NewProxyManager(),
        platformClients: map[string]PlatformClient{
            "instagram": NewInstagramClient(),
            "telegram":  NewTelegramClient(),
            "rubika":    NewRubikaClient(),
            "soroush":   NewSoroushClient(),
            "eitaa":     NewEitaaClient(),
        },
    }
}

func (ss *SocialScraper) FindSocialProfiles(phoneNumber string) []SocialProfile {
    var profiles []SocialProfile
    var mu sync.Mutex
    var wg sync.WaitGroup

    for platform, client := range ss.platformClients {
        wg.Add(1)
        go func(platform string, client PlatformClient) {
            defer wg.Done()
            
            profile, err := client.SearchByPhone(phoneNumber)
            if err == nil && profile != nil {
                mu.Lock()
                profiles = append(profiles, *profile)
                mu.Unlock()
            }
        }(platform, client)
    }

    wg.Wait()
    return profiles
}
