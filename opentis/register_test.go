package opentis

import "testing"

// TestSimpleRegister tests the simpleregister component
func TestSimpleRegister(t *testing.T) {

	t.Run("SimpleRegister has a default value of 0", func(t *testing.T) {
		r := newRegister()
		if r.Read() != 0 {
			t.Errorf("SimpleRegister should have a default value of 0")
		}
	})

	t.Run("SimpleRegister.Write writes to internal", func(t *testing.T) {
		r := newRegister()
		val := 18
		r.Write(val)

		out := r.Read()

		if out != val {
			t.Errorf("Expected register value to be %d, got %d instead", val, out)
		}

	})
}

func TestReadOnlyRegister(t *testing.T) {

	t.Run("ReadOnlyRegister outputs values from input channel", func(t *testing.T) {
		in := make(chan int)

		val := 18

		go func() {
			in <- val
		}()

		r := newReadOnlyRegister(in)

		out := r.Read()

		if out != val {
			t.Errorf("Expected register value to be %d, got %d instead", val, out)
		}
		close(in)
	})

	t.Run("ReadOnlyRegister.Write() does nothing", func(t *testing.T) {
		in := make(chan int)

		real := 18
		go func() {
			in <- 18
		}()

		r := newReadOnlyRegister(in)

		val := 3

		r.Write(val)

		out := r.Read()

		if out == val {
			t.Errorf("Expected read value to be %d, got %d instead", real, out)
		}
		close(in)
	})
}

func TestWriteOnlyRegister(t *testing.T) {
	t.Run("WriteOnlyRegister outputs values to input channel", func(t *testing.T) {
		out := make(chan int)

		r := newWriteOnlyRegister(out)

		val := 13

		go func() {
			r.Write(val)
		}()

		expect := <-out

		if expect != val {
			t.Errorf("expected read value to be %d, got %d instead", expect, val)
		}
		close(out)
	})

	t.Run("WriteOnlyRegister.Read() outputs -1", func(t *testing.T) {
		out := make(chan int)

		r := newWriteOnlyRegister(out)

		expect := r.Read()

		if expect != -1 {
			t.Errorf("WriteOnlyRegister.Read() should return -1, got %d", expect)
		}
	})
}

func TestVirtualRegister(t *testing.T) {
	t.Run("VirtualRegister.Write()", func(t *testing.T) {
		reg := newVirtualRegister()

		val := 15

		go func() {
			reg.Write(val)
		}()

		out := reg.Read()

		if out != val {
			t.Errorf("Expected read value to be %d, got %d instead", out, val)
		}
	})
}
