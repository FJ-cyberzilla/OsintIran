// backend/internal/security/license_validator.go
package security

type HardwareLicense struct {
    MachineID    string    `json:"machine_id"`
    CPUID        string    `json:"cpu_id"`
    MACAddress   string    `json:"mac_address"`
    DiskSerial   string    `json:"disk_serial"`
    ExpiryDate   time.Time `json:"expiry_date"`
    IranRegion   bool      `json:"iran_region"`
}

type LicenseValidator struct {
    encryptionKey []byte
    allowedRegions []string
}

func NewLicenseValidator() *LicenseValidator {
    return &LicenseValidator{
        encryptionKey:  loadEncryptionKey(),
        allowedRegions: []string{"IR", "+98", "0098"},
    }
}

func (lv *LicenseValidator) ValidateLicense(licenseData []byte) (*HardwareLicense, error) {
    // Decrypt license
    decrypted, err := lv.decryptLicense(licenseData)
    if err != nil {
        return nil, fmt.Errorf("license decryption failed")
    }

    var license HardwareLicense
    if err := json.Unmarshal(decrypted, &license); err != nil {
        return nil, fmt.Errorf("invalid license format")
    }

    // Check expiry
    if time.Now().After(license.ExpiryDate) {
        return nil, fmt.Errorf("license expired")
    }

    // Validate Iran region
    if !license.IranRegion {
        return nil, fmt.Errorf("license not valid for Iran region")
    }

    // Verify hardware binding
    if !lv.verifyHardwareBinding(license) {
        return nil, fmt.Errorf("hardware binding mismatch")
    }

    return &license, nil
}

func (lv *LicenseValidator) verifyHardwareBinding(license HardwareLicense) bool {
    currentMachineID := generateMachineID()
    return license.MachineID == currentMachineID
}
