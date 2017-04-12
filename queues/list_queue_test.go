package queueadts

import (
	"fmt"
	"math/rand"
	"testing"

	adts "github.com/johnsrd7/go-adts"
)

func TestMakeListQueue(t *testing.T) {
	queue := MakeListQueue()

	if queue.backer == nil {
		t.Error("Backing Container should not be nil after make call.")
	}
	if queue.lock == nil {
		t.Error("Lock for queue should not be nil after make call.")
	}
	if queue.threadSafe {
		t.Error("Threadsafe should be false for call to non-threadsafe make.")
	}

	var s Queue
	s = MakeListQueue()

	if s.Len() != 0 {
		t.Error("Length of empty queue should be 0.")
	}
}

// -------------------------------------------------------
// Test Container Methods
// -------------------------------------------------------

func TestListQueueLen(t *testing.T) {
	queue := MakeListQueue()

	for i := 0; i < 50; i++ {
		queue.Add(adts.IntElt(i))

		if queue.Len() != i+1 {
			t.Errorf("Queue should have length %d, actual length: %d", i+1, queue.Len())
			return
		}
	}
}

func TestListQueueIsEmpty(t *testing.T) {
	queue := MakeListQueue()

	if !queue.IsEmpty() {
		t.Errorf("Empty queue should return true for IsEmpty.\n")
	}

	queue.Add(adts.IntElt(0))
	if queue.IsEmpty() {
		t.Errorf("Queue with 1 element should not return true for IsEmpty.\n")
	}
}

func TestListQueueClear(t *testing.T) {
	queue := MakeListQueue()

	for i := 0; i < 100; i++ {
		queue.Add(adts.IntElt(i))
	}

	if queue.Len() != 100 {
		t.Errorf("Check Add/Len method, should have length of 100 after 100 adds.\n")
	}

	queue.Clear()
	if !queue.IsEmpty() {
		t.Errorf("Queue should be empty after call to Clear.\n")
	}
}

func TestListQueueAdd(t *testing.T) {
	vals := []int{}
	r := rand.New(rand.NewSource(99))

	queue := MakeListQueue()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals = append(vals, v)

		if !queue.Add(adts.IntElt(v)) {
			t.Errorf("Failed to add %d to list\n", v)
			return
		}

		if queue.Len() != i+1 {
			t.Errorf("Queue size not correct. Expected: %d, Actual: %d", i+1, queue.Len())
			return
		}

		iter := queue.backer.Front()
		for idx, v := range vals {
			iterVal, ok := iter.Value.(adts.IntElt)
			if !ok || !iterVal.Equals(adts.IntElt(v)) {
				t.Errorf("Add failed to add the value to the proper index | (idx,val) - Expected: (%d, %d), Actual: (%d, %v)",
					idx, v, idx, iterVal)
			}
			iter = iter.Next()
		}
	}
}

func TestListQueueContains(t *testing.T) {
	vals := make(map[int]bool)
	r := rand.New(rand.NewSource(99))

	queue := MakeListQueue()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals[v] = true

		queue.Add(adts.IntElt(v))

		for v := range vals {
			if !queue.Contains(adts.IntElt(v)) {
				t.Errorf("Contains failed to find the value (%d) in the list.", v)
			}
		}
	}
}

func TestListQueueRemove(t *testing.T) {
	max := 1000

	queue := MakeListQueue()
	for i := 0; i < max; i++ {
		queue.Add(adts.IntElt(i))
	}

	for i := 0; i < max; i++ {
		idx := rand.Int31n(int32(max - i))
		// get to the index
		tmp := queue.backer.Front()
		for j := 0; j < int(idx); j++ {
			tmp = tmp.Next()
		}
		val, _ := tmp.Value.(adts.IntElt)
		if !queue.Remove(val) {
			t.Errorf("Failed to remove index %d from list.", idx)
			return
		}

		if queue.Len() != max-1-i {
			t.Errorf("Expected length: %d, Actual length: %d", max-1-i, queue.Len())
			return
		}

		if queue.Contains(val) {
			t.Errorf("Value %v was not actually removed from list.", val)
			return
		}
	}

	// Now we want to check that the resize didn't break the ability to add.
	for i := 0; i < max; i++ {
		queue.Add(adts.IntElt(i))
		if queue.Len() != i+1 {
			t.Error("Unable to add after removing all elements.")
			return
		}
	}
}

// -------------------------------------------------------
// Test Queue Methods
// -------------------------------------------------------

func TestListQueueEnqueue(t *testing.T) {
	r := rand.New(rand.NewSource(99))

	expected := []int{}
	queue := MakeListQueue()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		expected = append(expected, v)
		if !queue.Enqueue(adts.IntElt(v)) {
			t.Errorf("Enqueue didn't succeed for enqueue #%d\n", i+1)
			return
		}

		if queue.Len() != len(expected) {
			t.Errorf("Enqueue didn't add the element properly to the queue. Expected len: %d, Actual len: %d",
				len(expected), queue.Len())
			return
		}
		iter := queue.backer.Front()
		for idx, exp := range expected {
			act, ok := iter.Value.(adts.IntElt)
			if !ok || !act.Equals(adts.IntElt(exp)) {
				t.Errorf("Queue order was ruined by push #%d. (idx, val) - Expected: (%d, %d), Actual: (%d, %v)",
					i+1, idx, exp, idx, act)
				return
			}

			iter = iter.Next()
		}
	}
}

func TestListQueueDequeue(t *testing.T) {
	queue := MakeListQueue()

	for i := 0; i < 100; i++ {
		queue.Enqueue(adts.IntElt(i))
	}

	for i := 0; i < 100; i++ {
		dequed, ok := queue.Dequeue()
		if !ok {
			t.Errorf("Dequeue didn't succeed for dequeue #%d\n", i+1)
			return
		}

		if !dequed.Equals(adts.IntElt(i)) {
			t.Errorf("Dequeue didn't return the proper element. Expected: %v, Actual: %v\n", i, dequed)
			return
		}
	}

	_, ok := queue.Dequeue()
	if ok {
		t.Errorf("Dequeue should not succeed for empty queue.\n")
	}
}

func listQueueToString(ls *ListQueue) string {
	res := "["
	tmp := ls.backer.Front()
	if tmp != nil {
		res += fmt.Sprintf("%v", tmp.Value)
		tmp = tmp.Next()
	}

	for {
		if tmp == nil {
			break
		}

		res += fmt.Sprintf(",%v", tmp.Value)
		tmp = tmp.Next()
	}

	return res + "]"
}
