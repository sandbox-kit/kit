// Package core assembles an explicitly selected sandbox provider.
//
// Provider implementations are installed separately. Client exposes a shared API;
// adapters will map shared requests and responses to their provider SDKs.
// The caller owns the injected provider and its underlying SDK client.
package core
