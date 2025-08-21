package access

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockReloadable is a test implementation of the Reloadable interface
type MockReloadable struct {
	reloadCount int
	wasReloaded bool
}

// Reload implements the Reloadable interface
func (m *MockReloadable) Reload() {
	m.reloadCount++
	m.wasReloaded = true
}

// GetReloadCount returns the number of times Reload was called
func (m *MockReloadable) GetReloadCount() int {
	return m.reloadCount
}

// WasReloaded returns whether Reload was called at least once
func (m *MockReloadable) WasReloaded() bool {
	return m.wasReloaded
}

// Reset resets the mock to its initial state
func (m *MockReloadable) Reset() {
	m.reloadCount = 0
	m.wasReloaded = false
}

func TestReloadableInterface(t *testing.T) {
	// Test that MockReloadable implements Reloadable interface
	var reloadable Reloadable = &MockReloadable{}
	assert.NotNil(t, reloadable)

	// Test that we can call Reload method through interface
	mock := &MockReloadable{}
	var iface Reloadable = mock

	// Initial state
	assert.False(t, mock.WasReloaded())
	assert.Equal(t, 0, mock.GetReloadCount())

	// Call Reload once
	iface.Reload()
	assert.True(t, mock.WasReloaded())
	assert.Equal(t, 1, mock.GetReloadCount())

	// Call Reload multiple times
	iface.Reload()
	iface.Reload()
	assert.True(t, mock.WasReloaded())
	assert.Equal(t, 3, mock.GetReloadCount())
}

func TestMockReloadable(t *testing.T) {
	mock := &MockReloadable{}

	// Test initial state
	assert.False(t, mock.WasReloaded())
	assert.Equal(t, 0, mock.GetReloadCount())

	// Test single reload
	mock.Reload()
	assert.True(t, mock.WasReloaded())
	assert.Equal(t, 1, mock.GetReloadCount())

	// Test multiple reloads
	mock.Reload()
	mock.Reload()
	assert.True(t, mock.WasReloaded())
	assert.Equal(t, 3, mock.GetReloadCount())

	// Test reset functionality
	mock.Reset()
	assert.False(t, mock.WasReloaded())
	assert.Equal(t, 0, mock.GetReloadCount())

	// Test reload after reset
	mock.Reload()
	assert.True(t, mock.WasReloaded())
	assert.Equal(t, 1, mock.GetReloadCount())
}

// ConfigManager is an example implementation of Reloadable for testing
type ConfigManager struct {
	config map[string]string
	loaded bool
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		config: make(map[string]string),
		loaded: false,
	}
}

func (c *ConfigManager) Reload() {
	// Simulate loading configuration
	c.config = map[string]string{
		"database_url": "localhost:5432",
		"cache_size":   "1000",
		"debug_mode":   "true",
	}
	c.loaded = true
}

func (c *ConfigManager) Get(key string) string {
	return c.config[key]
}

func (c *ConfigManager) IsLoaded() bool {
	return c.loaded
}

func TestConfigManagerReloadable(t *testing.T) {
	config := NewConfigManager()

	// Test initial state
	assert.False(t, config.IsLoaded())
	assert.Equal(t, "", config.Get("database_url"))

	// Test that it implements Reloadable
	var reloadable Reloadable = config
	assert.NotNil(t, reloadable)

	// Test reload functionality
	reloadable.Reload()
	assert.True(t, config.IsLoaded())
	assert.Equal(t, "localhost:5432", config.Get("database_url"))
	assert.Equal(t, "1000", config.Get("cache_size"))
	assert.Equal(t, "true", config.Get("debug_mode"))

	// Test multiple reloads
	reloadable.Reload()
	assert.True(t, config.IsLoaded())
	assert.Equal(t, "localhost:5432", config.Get("database_url"))
}

// CacheManager is another example implementation for testing
type CacheManager struct {
	cache     map[string]interface{}
	hitCount  int
	missCount int
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		cache: make(map[string]interface{}),
	}
}

func (c *CacheManager) Reload() {
	// Simulate cache reload - clear cache and reset counters
	c.cache = make(map[string]interface{})
	c.hitCount = 0
	c.missCount = 0
}

func (c *CacheManager) Set(key string, value interface{}) {
	c.cache[key] = value
}

func (c *CacheManager) Get(key string) (interface{}, bool) {
	value, exists := c.cache[key]
	if exists {
		c.hitCount++
	} else {
		c.missCount++
	}
	return value, exists
}

func (c *CacheManager) GetStats() (hits, misses int) {
	return c.hitCount, c.missCount
}

func TestCacheManagerReloadable(t *testing.T) {
	cache := NewCacheManager()

	// Test that it implements Reloadable
	var reloadable Reloadable = cache
	assert.NotNil(t, reloadable)

	// Add some data and generate stats
	cache.Set("key1", "value1")
	cache.Set("key2", 42)

	value, exists := cache.Get("key1")
	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	_, exists = cache.Get("nonexistent")
	assert.False(t, exists)

	hits, misses := cache.GetStats()
	assert.Equal(t, 1, hits)
	assert.Equal(t, 1, misses)

	// Test reload - should clear everything
	reloadable.Reload()

	// Check that cache was cleared
	_, exists = cache.Get("key1")
	assert.False(t, exists)

	// Check that stats were reset
	hits, misses = cache.GetStats()
	assert.Equal(t, 0, hits)
	assert.Equal(t, 1, misses) // This miss from checking key1
}
