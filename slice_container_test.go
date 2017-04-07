package adts

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMakeSliceContainer(t *testing.T) {
	container := MakeSliceContainer()

	if len(container.Backer) != 0 {
		t.Error("Length of empty container should be 0")
	}
	if container.ThreadSafe {
		t.Error("Threadsafe bool should not be set on default make call.")
	}
	if container.Lock == nil {
		t.Error("Lock should not be nil after make call.")
	}

	var l Container
	l = MakeSliceContainer()

	if l.Len() != 0 {
		t.Error("Length of empty container should be 0.")
	}
}

// -------------------------------------------------------
// Test Container Methods
// -------------------------------------------------------

func TestSliceContainerLen(t *testing.T) {
	container := MakeSliceContainer()

	for i := 0; i < 50; i++ {
		container.Backer = append(container.Backer, IntElt(i))

		if container.Len() != i+1 {
			t.Errorf("Container should have length %d, actual length: %d", i+1, container.Len())
			return
		}
	}
}

func TestSliceContainerIsEmpty(t *testing.T) {
	container := MakeSliceContainer()

	if !container.IsEmpty() {
		t.Errorf("Empty container should return true for IsEmpty.\n")
	}

	container.Add(IntElt(0))
	if container.IsEmpty() {
		t.Errorf("Container with 1 element should not return true for IsEmpty.\n")
	}
}

func TestSliceContainerClear(t *testing.T) {
	container := MakeSliceContainer()

	for i := 0; i < 100; i++ {
		container.Add(IntElt(i))
	}

	if container.Len() != 100 {
		t.Errorf("Check Add/Len method, should have length of 100 after 100 adds.\n")
	}

	container.Clear()
	if !container.IsEmpty() {
		t.Errorf("Container should be empty after call to Clear.\n")
	}
}

func TestSliceContainerAdd(t *testing.T) {
	vals := make(map[int][]int) // value -> slice of indices
	r := rand.New(rand.NewSource(99))

	container := MakeSliceContainer()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals[v] = append(vals[v], i)

		if !container.Add(IntElt(v)) {
			t.Errorf("Failed to add %d to container\n", v)
			return
		}

		if len(container.Backer) != i+1 {
			t.Errorf("Container size not correct. Expected: %d, Actual: %d", i+1, len(container.Backer))
			return
		}
		for v, idxs := range vals {
			for _, idx := range idxs {
				if idx >= len(container.Backer) {
					fmt.Printf("Idx: %d, Len: %d\n", idx, len(container.Backer))
					continue
				}
				if !container.Backer[idx].Equals(IntElt(v)) {
					t.Errorf("Add failed to add the value to the proper index | (idx,val) - Expected: (%d, %d), Actual: (%d, %v)",
						idx, v, idx, container.Backer[idx])
				}
			}
		}
	}
}

func TestSliceContainerContains(t *testing.T) {
	vals := make(map[int]bool)
	r := rand.New(rand.NewSource(99))

	container := MakeSliceContainer()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals[v] = true

		container.Add(IntElt(v))

		for v := range vals {
			if !container.Contains(IntElt(v)) {
				t.Errorf("Contains failed to find the value (%d) in the container.", v)
			}
		}
	}
}

func TestSliceContainerRemove(t *testing.T) {
	max := 1000

	container := MakeSliceContainer()
	for i := 0; i < max; i++ {
		container.Add(IntElt(i))
	}

	for i := 0; i < max; i++ {
		idx := rand.Int31n(int32(max - i))
		val := container.Backer[idx]
		if !container.Remove(val) {
			t.Errorf("Failed to remove index %d from container.", idx)
			return
		}

		if container.Len() != max-1-i {
			t.Errorf("Expected length: %d, Actual length: %d", max-1-i, container.Len())
			return
		}

		if container.Contains(val) {
			t.Errorf("Value %v was not actually removed from container.", val)
			return
		}
	}

	// Now we want to check that the resize didn't break the ability to add.
	for i := 0; i < max; i++ {
		container.Add(IntElt(i))
		if container.Len() != i+1 {
			t.Error("Unable to add after removing all elements.")
			return
		}
	}
}
