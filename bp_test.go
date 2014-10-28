package bp

import "testing"

func TestBP(t *testing.T) {

	bp := newBufferPool(10)
	b := <-bp.Get
	if cap(b.Bytes()) != int(10) {
		t.Errorf("Expected buffer capacity: %d. Actual: %d", 10, b.Len())
	}
	bp.Give <- b
	if bp.Makes != 2 {
		t.Errorf("Expected makes: %d. Actual: %d", 2, bp.Makes)
	}

	b = <-bp.Get
	bp.Give <- b
	b = <-bp.Get
	if bp.Makes != 2 {
		t.Errorf("Expected makes: %d. Actual: %d", 2, bp.Makes)
	}
}

func TestShutdown(t *testing.T) {
	bp := newBufferPool(1)
	b := <-bp.Get
	bp.Give <- b
	bp.Shutdown()
	for _ = range bp.Get {
		t.Fatal("bp.Get should be closed")
	}

}