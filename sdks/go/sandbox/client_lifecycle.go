package sandbox

import (
	"context"
	"sync"
)

// clientLifecycle prevents provider cleanup from racing an in-flight operation.
// It is runtime-only and emitted as an opaque field on Client.
type clientLifecycle struct {
	mu     sync.RWMutex
	closed bool
}

func (c *Client) beginOperation(operation string) (func(), error) {
	provider := ""
	if c != nil && c.backend != nil {
		provider = c.backend.Name()
	}
	if c == nil || c.lifecycle == nil {
		return nil, NewError(ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: provider, Operation: operation, Field: "client", Message: "sandbox-kit: client is not initialized"}, nil)
	}
	c.lifecycle.mu.RLock()
	if c.lifecycle.closed {
		c.lifecycle.mu.RUnlock()
		return nil, NewError(ErrorInfo{Kind: ErrorKindFailedPrecondition, Provider: provider, Operation: operation, Field: "client", Message: "sandbox-kit: client is closed"}, nil)
	}
	return c.lifecycle.mu.RUnlock, nil
}

func (c *Client) closeBackend(ctx context.Context) error {
	if c == nil || c.backend == nil || c.lifecycle == nil {
		return nil
	}
	c.lifecycle.mu.Lock()
	defer c.lifecycle.mu.Unlock()
	if c.lifecycle.closed {
		return nil
	}
	c.lifecycle.closed = true
	return c.backend.Close(ctx)
}

func (c *Client) isClosed() bool {
	if c == nil || c.lifecycle == nil {
		return true
	}
	c.lifecycle.mu.RLock()
	defer c.lifecycle.mu.RUnlock()
	return c.lifecycle.closed
}
