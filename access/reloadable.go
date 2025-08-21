// Package access provides interfaces for components that support runtime reloading.
package access

// Reloadable defines an interface for components that can be reloaded at runtime.
// This is commonly used for configuration objects, caches, or other stateful
// components that need to refresh their state without restarting the application.
type Reloadable interface {
	// Reload refreshes the component's internal state.
	// Implementations should handle errors gracefully and maintain consistency.
	Reload()
}
