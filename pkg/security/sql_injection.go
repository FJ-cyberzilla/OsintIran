// pkg/security/sql_injection.go
package security

import (
    "database/sql"
    "fmt"
    "regexp"
    "strings"
    "unicode"
)

type SQLInjectionProtector struct {
    blacklistPatterns []*regexp.Regexp
    whitelistPatterns []*regexp.Regexp
    maxQueryLength    int
}

func NewSQLInjectionProtector() *SQLInjectionProtector {
    return &SQLInjectionProtector{
        maxQueryLength: 10000,
        blacklistPatterns: []*regexp.Regexp{
            regexp.MustCompile(`(?i)(\bUNION\b.*\bSELECT\b)`),
            regexp.MustCompile(`(?i)(\bDROP\b|\bDELETE\b|\bINSERT\b|\bUPDATE\b|\bALTER\b)`),
            regexp.MustCompile(`(?i)(\bOR\b.*=.*\bOR\b)`),
            regexp.MustCompile(`(?i)(\bEXEC\b|\bEXECUTE\b|\bEXECSP\b)`),
            regexp.MustCompile(`(?i)(\bWAITFOR\b.*\bDELAY\b)`),
            regexp.MustCompile(`--|\/\*|\*\/|;`),
            regexp.MustCompile(`(\bUNION\b.*\bALL\b.*\bSELECT\b)`),
            regexp.MustCompile(`(\bLOAD_FILE\b|\bINTO\b.*\bOUTFILE\b|\bINTO\b.*\bDUMPFILE\b)`),
        },
    }
}

// Multi-layer SQL injection protection
func (sip *SQLInjectionProtector) SecureQuery(query string, args ...interface{}) (string, []interface{}, error) {
    // Layer 1: Length validation
    if len(query) > sip.maxQueryLength {
        return "", nil, fmt.Errorf("query too long")
    }

    // Layer 2: Blacklist patterns
    for _, pattern := range sip.blacklistPatterns {
        if pattern.MatchString(query) {
            return "", nil, fmt.Errorf("potential SQL injection detected")
        }
    }

    // Layer 3: Parameterized query enforcement
    secureQuery, secureArgs, err := sip.parameterizeQuery(query, args)
    if err != nil {
        return "", nil, err
    }

    // Layer 4: Input sanitization
    sanitizedArgs := make([]interface{}, len(secureArgs))
    for i, arg := range secureArgs {
        sanitizedArgs[i] = sip.sanitizeInput(arg)
    }

    return secureQuery, sanitizedArgs, nil
}

func (sip *SQLInjectionProtector) parameterizeQuery(query string, args []interface{}) (string, []interface{}, error) {
    // Ensure all user inputs are parameterized
    if strings.Contains(query, "%s") || strings.Contains(query, "%v") {
        return "", nil, fmt.Errorf("raw string formatting not allowed - use parameterized queries")
    }

    // Count expected parameters
    expectedParams := strings.Count(query, "?")
    if expectedParams != len(args) {
        return "", nil, fmt.Errorf("parameter count mismatch")
    }

    return query, args, nil
}

func (sip *SQLInjectionProtector) sanitizeInput(input interface{}) interface{} {
    switch v := input.(type) {
    case string:
        // Remove control characters and excessive whitespace
        cleaned := strings.Map(func(r rune) rune {
            if unicode.IsControl(r) {
                return -1
            }
            return r
        }, v)
        return strings.TrimSpace(cleaned)
    default:
        return input
    }
}

// Secure database wrapper
type SecureDB struct {
    db *sql.DB
    sip *SQLInjectionProtector
}

func (sdb *SecureDB) SecureExec(query string, args ...interface{}) (sql.Result, error) {
    secureQuery, secureArgs, err := sdb.sip.SecureQuery(query, args...)
    if err != nil {
        return nil, fmt.Errorf("sql injection protection: %w", err)
    }
    return sdb.db.Exec(secureQuery, secureArgs...)
}

func (sdb *SecureDB) SecureQuery(query string, args ...interface{}) (*sql.Rows, error) {
    secureQuery, secureArgs, err := sdb.sip.SecureQuery(query, args...)
    if err != nil {
        return nil, fmt.Errorf("sql injection protection: %w", err)
    }
    return sdb.db.Query(secureQuery, secureArgs...)
}
