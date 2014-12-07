package main

// Pool holds bufs.
type BufPool struct {
	pool chan []byte
}

// NewBufPool creates a new pool of Clients.
func NewBufPool(max int) *BufPool {
	return &BufPool{
		pool: make(chan []byte, max),
	}
}

// Borrow a buf from the pool.
func (p *BufPool) Borrow() []byte {
	var b []byte
	select {
	case b = <-p.pool:
	default:
		b = make([]byte, 64000)
	}
	return b
}

// Return returns a buf to the pool.
func (p *BufPool) Return(b []byte) {
	select {
	case p.pool <- b:
	default:
		// let it go, let it go...
	}
}
