package bp

import (
	"bytes"
	"container/list"
)

// BP includes Get and Give directional channels for retrieving and returning buffers
type BP struct {
	Get   <-chan *bytes.Buffer // Get a buffer from the pool
	Give  chan<- *bytes.Buffer // Give a buffer to be stored in the pool
	Makes int                  // total buffers allocated
	quit  chan struct{}
}

// New initalizes a new buffer pool
func New(bufsz int) *BP {
	get := make(chan *bytes.Buffer)
	give := make(chan *bytes.Buffer)
	np := &BP{
		Get:  get,
		Give: give,
		quit: make(chan struct{}),
	}
	go func() {
		q := new(list.List)
		for {
			if q.Len() == 0 {
				q.PushFront(bytes.NewBuffer(make([]byte, 0, bufsz)))
				np.Makes++
			}

			e := q.Front()

			select {
			case b := <-give:
				q.PushFront(b)
			case get <- e.Value.(*bytes.Buffer):
				q.Remove(e)
			case <-np.quit:
				close(give)
				close(get)
				return
			}
		}
	}()
	return np
}

// Shutdown terminates the pool, closing the Get and Give channels and freeing pool buffers to be garbage collected
func (bp *BP) Shutdown() {
	close(bp.quit)
	<-bp.Get // wait for close
}