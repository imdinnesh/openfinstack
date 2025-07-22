package kafka

import "sync"

type Dispatcher struct {
    handlers map[string]MessageHandler
    mu       sync.RWMutex
}

func NewDispatcher() *Dispatcher {
    return &Dispatcher{
        handlers: make(map[string]MessageHandler),
    }
}

func (d *Dispatcher) RegisterHandler(topic string, handler MessageHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.handlers[topic] = handler
}

func (d *Dispatcher) GetHandler(topic string) (MessageHandler, bool) {
    d.mu.RLock()
    defer d.mu.RUnlock()
    h, ok := d.handlers[topic]
    return h, ok
}
