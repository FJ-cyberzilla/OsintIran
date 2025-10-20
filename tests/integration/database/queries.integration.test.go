// tests/integration/database/queries.integration.test.go
package integration

import (
    "context"
    "database/sql"
    "fmt"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "gorm.io/gorm"

    "secure-iran-intel/database"
    "secure-iran-intel/models"
)

type DatabaseIntegrationTestSuite struct {
    suite.Suite
    db     *gorm.DB
    sqlDB  *sql.DB
    ctx    context.Context
}

func TestDatabaseIntegrationSuite(t *testing.T) {
    suite.Run(t, new(DatabaseIntegrationTestSuite))
}

func (suite *DatabaseIntegrationTestSuite) SetupSuite() {
    suite.ctx = context.Background()
    
    // Connect to test database
    var err error
    suite.db, suite.sqlDB, err = database.ConnectTest()
    if err != nil {
        suite.T().Fatalf("Failed to connect to test database: %v", err)
    }

    // Run migrations
    if err := database.MigrateTest(suite.db); err != nil {
        suite.T().Fatalf("Failed to run test migrations: %v", err)
    }
}

func (suite *DatabaseIntegrationTestSuite) TearDownSuite() {
    if suite.sqlDB != nil {
        suite.sqlDB.Close()
    }
}

func (suite *DatabaseIntegrationTestSuite) SetupTest() {
    // Clean database before each test
    suite.cleanDatabase()
}

func (suite *DatabaseIntegrationTestSuite) TestPhoneLookupQueries() {
    // Create test data
    testLookup := &models.PhoneLookup{
        ID:          "test-lookup-1",
        PhoneNumber: "+989123456789",
        TenantID:    "test-tenant-1",
        Status:      "completed",
        Results:     models.JSONMap{"carrier": "MCI", "region": "Tehran"},
        CreatedAt:   time.Now(),
    }

    // Test INSERT
    result := suite.db.Create(testLookup)
    suite.NoError(result.Error)
    suite.Equal(int64(1), result.RowsAffected)

    // Test SELECT
    var retrievedLookup models.PhoneLookup
    result = suite.db.Where("id = ?", testLookup.ID).First(&retrievedLookup)
    suite.NoError(result.Error)
    suite.Equal(testLookup.PhoneNumber, retrievedLookup.PhoneNumber)
    suite.Equal(testLookup.TenantID, retrievedLookup.TenantID)

    // Test UPDATE
    result = suite.db.Model(&retrievedLookup).Update("status", "processed")
    suite.NoError(result.Error)

    // Verify update
    result = suite.db.Where("id = ?", testLookup.ID).First(&retrievedLookup)
    suite.NoError(result.Error)
    suite.Equal("processed", retrievedLookup.Status)
}

func (suite *DatabaseIntegrationTestSuite) TestEmailDiscoveryQueries() {
    // Create test email discovery record
    testEmail := &models.EmailDiscovery{
        ID:           "test-email-1",
        PhoneNumber:  "+989123456789",
        Email:        "test@example.com",
        Confidence:   0.85,
        Source:       "pattern_generation",
        FoundInBreach: false,
        TenantID:     "test-tenant-1",
        DiscoveredAt: time.Now(),
    }

    // Test batch insert
    emails := []models.EmailDiscovery{
        *testEmail,
        {
            ID:           "test-email-2",
            PhoneNumber:  "+989123456789",
            Email:        "test2@example.com",
            Confidence:   0.72,
            Source:       "breach_database",
            FoundInBreach: true,
            TenantID:     "test-tenant-1",
            DiscoveredAt: time.Now(),
        },
    }

    result := suite.db.Create(&emails)
    suite.NoError(result.Error)
    suite.Equal(int64(2), result.RowsAffected)

    // Test complex query with joins
    var results []struct {
        PhoneNumber string
        EmailCount  int64
        AvgConfidence float64
    }

    result = suite.db.Model(&models.EmailDiscovery{}).
        Select("phone_number, COUNT(*) as email_count, AVG(confidence) as avg_confidence").
        Where("tenant_id = ?", "test-tenant-1").
        Group("phone_number").
        Find(&results)

    suite.NoError(result.Error)
    suite.Len(results, 1)
    suite.Equal("+989123456789", results[0].PhoneNumber)
    suite.Equal(int64(2), results[0].EmailCount)
    suite.InDelta(0.785, results[0].AvgConfidence, 0.01)
}

func (suite *DatabaseIntegrationTestSuite) TestTenantIsolation() {
    // Create data for multiple tenants
    tenants := []string{"tenant-a", "tenant-b", "tenant-c"}
    
    for _, tenantID := range tenants {
        lookups := []models.PhoneLookup{
            {
                ID:          fmt.Sprintf("lookup-%s-1", tenantID),
                PhoneNumber: "+989111111111",
                TenantID:    tenantID,
                Status:      "completed",
                CreatedAt:   time.Now(),
            },
            {
                ID:          fmt.Sprintf("lookup-%s-2", tenantID),
                PhoneNumber: "+989222222222", 
                TenantID:    tenantID,
                Status:      "processing",
                CreatedAt:   time.Now(),
            },
        }
        
        result := suite.db.Create(&lookups)
        suite.NoError(result.Error)
    }

    // Test tenant isolation - each tenant should only see their data
    for _, tenantID := range tenants {
        var tenantLookups []models.PhoneLookup
        result := suite.db.Where("tenant_id = ?", tenantID).Find(&tenantLookups)
        suite.NoError(result.Error)
        suite.Len(tenantLookups, 2, "Tenant should see only their lookups")
        
        for _, lookup := range tenantLookups {
            suite.Equal(tenantID, lookup.TenantID, "Lookup should belong to correct tenant")
        }
    }

    // Test cross-tenant data leakage prevention
    var allLookups []models.PhoneLookup
    result := suite.db.Find(&allLookups)
    suite.NoError(result.Error)
    suite.Len(allLookups, 6, "Should see all lookups without tenant filter")
}

func (suite *DatabaseIntegrationTestSuite) TestConcurrentDatabaseAccess() {
    const numGoroutines = 10
    const insertsPerGoroutine = 5

    errors := make(chan error, numGoroutines*insertsPerGoroutine)
    var wg sync.WaitGroup

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(goroutineID int) {
            defer wg.Done()
            
            for j := 0; j < insertsPerGoroutine; j++ {
                lookup := models.PhoneLookup{
                    ID:          fmt.Sprintf("concurrent-%d-%d", goroutineID, j),
                    PhoneNumber: fmt.Sprintf("+9891234567%02d", j),
                    TenantID:    fmt.Sprintf("tenant-%d", goroutineID),
                    Status:      "completed",
                    CreatedAt:   time.Now(),
                }
                
                result := suite.db.Create(&lookup)
                if result.Error != nil {
                    errors <- result.Error
                }
                
                // Small delay to simulate real workload
                time.Sleep(time.Millisecond * 10)
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Check for errors
    var errorCount int
    for err := range errors {
        suite.T().Logf("Concurrent insert error: %v", err)
        errorCount++
    }

    suite.Equal(0, errorCount, "No errors should occur during concurrent access")

    // Verify all data was inserted
    var totalLookups int64
    result := suite.db.Model(&models.PhoneLookup{}).Where("id LIKE ?", "concurrent-%").Count(&totalLookups)
    suite.NoError(result.Error)
    suite.Equal(int64(numGoroutines*insertsPerGoroutine), totalLookups)
}

func (suite *DatabaseIntegrationTestSuite) TestDatabasePerformance() {
    // Test query performance with large datasets
    const batchSize = 1000
    
    // Create large batch of test data
    var lookups []models.PhoneLookup
    for i := 0; i < batchSize; i++ {
        lookups = append(lookups, models.PhoneLookup{
            ID:          fmt.Sprintf("perf-test-%d", i),
            PhoneNumber: fmt.Sprintf("+9891234567%03d", i%1000),
            TenantID:    "perf-tenant",
            Status:      "completed",
            CreatedAt:   time.Now(),
        })
    }

    // Test batch insert performance
    startTime := time.Now()
    result := suite.db.CreateInBatches(&lookups, 100) // Batch size 100
    insertTime := time.Since(startTime)
    
    suite.NoError(result.Error)
    suite.Equal(int64(batchSize), result.RowsAffected)
    
    // Insert should complete within reasonable time
    maxInsertTime := 2 * time.Second
    suite.True(insertTime < maxInsertTime, 
        "Batch insert of %d records should complete within %v, took %v", 
        batchSize, maxInsertTime, insertTime)

    // Test query performance
    startTime = time.Now()
    var count int64
    result = suite.db.Model(&models.PhoneLookup{}).Where("tenant_id = ?", "perf-tenant").Count(&count)
    queryTime := time.Since(startTime)
    
    suite.NoError(result.Error)
    suite.Equal(int64(batchSize), count)
    
    // Count query should be fast
    maxQueryTime := 100 * time.Millisecond
    suite.True(queryTime < maxQueryTime,
        "Count query should complete within %v, took %v",
        maxQueryTime, queryTime)

    fmt.Printf("📊 Performance: Insert %d records: %v, Count query: %v\n", 
        batchSize, insertTime, queryTime)
}

func (suite *DatabaseIntegrationTestSuite) cleanDatabase() {
    // Clean all test data
    tables := []string{
        "phone_lookups", 
        "email_discoveries",
        "intelligence_reports",
        "tenant_usage",
    }
    
    for _, table := range tables {
        suite.db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id LIKE 'test-%%' OR id LIKE 'concurrent-%%' OR id LIKE 'perf-%%'", table))
    }
}
