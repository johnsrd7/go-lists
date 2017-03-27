package listadts

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestMakeSliceList(t *testing.T) {
	list := MakeSliceList()

	if len(list.backer) != 0 {
		t.Error("Length of empty list should be 0")
	}
	if list.threadSafe {
		t.Error("Threadsafe bool should not be set on default make call.")
	}
	if list.lock == nil {
		t.Error("Lock should not be nil after make call.")
	}

	var l List
	l = MakeSliceList()

	if l.Len() != 0 {
		t.Error("Length of empty list should be 0.")
	}
}

// -------------------------------------------------------
// Test Container Methods
// -------------------------------------------------------

func TestSliceListLen(t *testing.T) {
	list := MakeSliceList()

	for i := 0; i < 50; i++ {
		list.backer = append(list.backer, intElt(i))

		if list.Len() != i+1 {
			t.Errorf("List should have length %d, actual length: %d", i+1, list.Len())
			return
		}
	}
}

func TestSliceListIsEmpty(t *testing.T) {
	list := MakeSliceList()

	if !list.IsEmpty() {
		t.Errorf("Empty list should return true for IsEmpty.\n")
	}

	list.Add(intElt(0))
	if list.IsEmpty() {
		t.Errorf("List with 1 element should not return true for IsEmpty.\n")
	}
}

func TestSliceListClear(t *testing.T) {
	list := MakeSliceList()

	for i := 0; i < 100; i++ {
		list.Add(intElt(i))
	}

	if list.Len() != 100 {
		t.Errorf("Check Add/Len method, should have length of 100 after 100 adds.\n")
	}

	list.Clear()
	if !list.IsEmpty() {
		t.Errorf("List should be empty after call to Clear.\n")
	}
}

func TestSliceListAdd(t *testing.T) {
	vals := make(map[int][]int) // value -> slice of indices
	r := rand.New(rand.NewSource(99))

	list := MakeSliceList()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals[v] = append(vals[v], i)

		if !list.Add(intElt(v)) {
			t.Errorf("Failed to add %d to list\n", v)
			return
		}

		if len(list.backer) != i+1 {
			t.Errorf("List size not correct. Expected: %d, Actual: %d", i+1, len(list.backer))
			return
		}
		for v, idxs := range vals {
			for _, idx := range idxs {
				if idx >= len(list.backer) {
					fmt.Printf("Idx: %d, Len: %d\n", idx, len(list.backer))
					continue
				}
				if !list.backer[idx].Equals(intElt(v)) {
					t.Errorf("Add failed to add the value to the proper index | (idx,val) - Expected: (%d, %d), Actual: (%d, %v)",
						idx, v, idx, list.backer[idx])
				}
			}
		}
	}
}

func TestSliceListContains(t *testing.T) {
	vals := make(map[int]bool)
	r := rand.New(rand.NewSource(99))

	list := MakeSliceList()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals[v] = true

		list.Add(intElt(v))

		for v := range vals {
			if !list.Contains(intElt(v)) {
				t.Errorf("Contains failed to find the value (%d) in the list.", v)
			}
		}
	}
}

func TestSliceListRemove(t *testing.T) {
	max := 1000

	list := MakeSliceList()
	for i := 0; i < max; i++ {
		list.Add(intElt(i))
	}

	for i := 0; i < max; i++ {
		idx := rand.Int31n(int32(max - i))
		val := list.Get(int(idx))
		if !list.Remove(val) {
			t.Errorf("Failed to remove index %d from list.", idx)
			return
		}

		if list.Len() != max-1-i {
			t.Errorf("Expected length: %d, Actual length: %d", max-1-i, list.Len())
			return
		}

		if list.Contains(val) {
			t.Errorf("Value %v was not actually removed from list.", val)
			return
		}
	}

	// Now we want to check that the resize didn't break the ability to add.
	for i := 0; i < max; i++ {
		list.Add(intElt(i))
		if list.Len() != i+1 {
			t.Error("Unable to add after removing all elements.")
			return
		}
	}
}

// -------------------------------------------------------
// Test List Methods
// -------------------------------------------------------

func TestSliceListGet(t *testing.T) {
	r := rand.New(rand.NewSource(99))

	expected := []int{}
	list := MakeSliceList()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		expected = append(expected, v)
		list.Add(intElt(v))
	}

	for idx, v := range expected {
		if !list.Get(idx).Equals(intElt(v)) {
			t.Errorf("Expected: %d, Actual: %v", v, list.Get(idx))
			return
		}
	}
}

func TestSliceListSet(t *testing.T) {
	list := MakeSliceList()

	for i := 0; i < 100; i++ {
		list.Add(intElt(i))
	}

	for i := 0; i < 100; i++ {
		list.Set(i, intElt(i*2))

		if !list.Get(i).Equals(intElt(i * 2)) {
			t.Errorf("Set did not set element at index %d properly. Expected: %v, Actual: %v\n", i, i*2, list.Get(i))
			return
		}
	}
}
