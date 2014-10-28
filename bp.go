package bp

import (
	"bytes"
	"container/list"
)

type BP struct {
	Get   <-chan *bytes.Buffer
	Give  chan<- *bytes.Buffer
	Makes int
	quit  chan struct{}
}

func newBufferPool(bufsz int) *BP {
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

func (bp *BP) Shutdown() {
	close(bp.quit)
	<-bp.Get // wait for close
}