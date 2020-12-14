package ezBufferPool

type Pool struct {
	pool       chan []byte
	bufferSize int
}

func NewPool(reuseCount, bufferSize int) *Pool {
	return &Pool{make(chan []byte, reuseCount), bufferSize}
}
func (p *Pool) Get() []byte {
	select {
	case buffer := <-p.pool:
		return buffer
	default:
		return make([]byte, p.bufferSize)
	}
}
func (p *Pool) Put(buffer []byte) {
	select {
	case p.pool <- buffer:
	default:
	}
}
func (p *Pool) Free() {
	close(p.pool)
}
