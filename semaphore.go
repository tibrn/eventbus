package eventbus

type PingPong struct{}

var (
	defaultPingPong = PingPong{}
)

type Semaphore struct {
	ch chan PingPong
}

func NewSemaphore(resources int) *Semaphore {

	if resources < 1 {
		resources = 1
	}

	return &Semaphore{
		ch: make(chan PingPong, resources),
	}
}

func (s *Semaphore) Acquire() {
	s.ch <- defaultPingPong
}

func (s *Semaphore) Release() {
	<-s.ch
}
