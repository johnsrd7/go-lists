package listadts

import (
	"fmt"
	"math/rand"
	"testing"

	adts "github.com/johnsrd7/go-adts"
)

func printSinglyLinkedList(ll *SinglyLinkedList) {
	for tmp := ll.head; tmp != nil; tmp = tmp.next {
		fmt.Printf("%v -> ", tmp.elt)
	}

	fmt.Println()
}

func TestMakeSinglyLinkedList(t *testing.T) {
	list := MakeSinglyLinkedList()

	if list.len != 0 {
		t.Error("Length of empty list should be 0")
	}
	if list.threadSafe {
		t.Error("Threadsafe bool should not be set on default make call.")
	}
	if list.lock == nil {
		t.Error("Lock should not be nil after make call.")
	}

	var l List
	l = MakeSinglyLinkedList()

	if l.Len() != 0 {
		t.Error("Length of empty list should be 0.")
	}
}

// -------------------------------------------------------
// Test Container Methods
// -------------------------------------------------------

func TestSinglyLinkedListLen(t *testing.T) {
	list := MakeSinglyLinkedList()

	for i := 0; i < 50; i++ {
		list.Add(adts.IntElt(i))

		if list.Len() != i+1 {
			t.Errorf("List should have length %d, actual length: %d", i+1, list.Len())
			return
		}
	}
}

func TestSinglyLinkedListIsEmpty(t *testing.T) {
	list := MakeSinglyLinkedList()

	if !list.IsEmpty() {
		t.Errorf("Empty list should return true for IsEmpty.\n")
	}

	list.Add(adts.IntElt(0))
	if list.IsEmpty() {
		t.Errorf("List with 1 element should not return true for IsEmpty.\n")
	}
}

func TestSinglyLinkedListClear(t *testing.T) {
	list := MakeSinglyLinkedList()

	for i := 0; i < 100; i++ {
		list.Add(adts.IntElt(i))
	}

	if list.Len() != 100 {
		t.Errorf("Check Add/Len method, should have length of 100 after 100 adds.\n")
	}

	list.Clear()
	if !list.IsEmpty() {
		t.Errorf("List should be empty after call to Clear.\n")
	}
}

func TestSinglyLinkedListAdd(t *testing.T) {
	expected := []adts.ContainerElement{}
	r := rand.New(rand.NewSource(99))

	list := MakeSinglyLinkedList()

	for i := 0; i < 1000; i++ {
		v := adts.IntElt(r.Int())

		expected = append(expected, v)

		if !list.Add(v) {
			t.Errorf("Failed to add %d to list\n", v)
			return
		}

		if list.Len() != i+1 {
			t.Errorf("List size not correct. Expected: %d, Actual: %d", i+1, list.len)
			return
		}

		tmp := list.head
		for idx, val := range expected {
			if !tmp.elt.Equals(val) {
				t.Errorf("Add failed to add the value to the proper index (%d) | (idx,val) - Expected: (%d, %v), Actual: (%d, %v)\n",
					i, idx, val, idx, tmp.elt)
				return
			}

			tmp = tmp.next
		}
	}
}

func TestSinglyLinkedListContains(t *testing.T) {
	vals := make(map[int]bool)
	r := rand.New(rand.NewSource(99))

	list := MakeSinglyLinkedList()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals[v] = true

		list.Add(adts.IntElt(v))

		for v := range vals {
			if !list.Contains(adts.IntElt(v)) {
				t.Errorf("Contains failed to find the value (%d) in the list.", v)
			}
		}
	}
}

func TestSinglyLinkedListRemove(t *testing.T) {
	max := 1000

	list := MakeSinglyLinkedList()
	for i := 0; i < max; i++ {
		list.Add(adts.IntElt(i))
	}

	for i := 0; i < max; i++ {
		idx := rand.Int31n(int32(max - i))
		val := list.Get(int(idx))
		if !list.Remove(val) {
			t.Errorf("Failed to remove index %d from list.", idx)
			return
		}

		if list.Len() != max-1-i {
			t.Errorf("(%d) Expected length: %d, Actual length: %d", i, max-1-i, list.Len())
			return
		}

		if list.Contains(val) {
			t.Errorf("Value %v was not actually removed from list.", val)
			return
		}
	}

	if list.Len() != 0 {
		t.Error("Length of list should be 0 after all elements are removed.\n")
	}
}

// -------------------------------------------------------
// Test List Methods
// -------------------------------------------------------

func TestGet(t *testing.T) {
	r := rand.New(rand.NewSource(99))

	expected := []int{}
	list := MakeSinglyLinkedList()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		expected = append(expected, v)
		list.Add(adts.IntElt(v))
	}

	for idx, v := range expected {
		if !list.Get(idx).Equals(adts.IntElt(v)) {
			t.Errorf("Expected: %d, Actual: %v", v, list.Get(idx))
			return
		}
	}
}

func TestSet(t *testing.T) {
	list := MakeSinglyLinkedList()

	for i := 0; i < 100; i++ {
		list.Add(adts.IntElt(i))
	}

	for i := 0; i < 100; i++ {
		list.Set(i, adts.IntElt(i*2))

		if !list.Get(i).Equals(adts.IntElt(i * 2)) {
			t.Errorf("Set did not set element at index %d properly. Expected: %v, Actual: %v\n", i, i*2, list.Get(i))
			return
		}
	}
}
