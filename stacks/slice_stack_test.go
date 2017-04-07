package stackadts

import (
	"math/rand"
	"testing"

	adts "github.com/johnsrd7/go-adts"
)

func TestMakeSliceStack(t *testing.T) {
	stack := MakeSliceStack()

	if len(stack.backer) != 0 {
		t.Error("Length of empty stack should be 0")
	}
	if stack.threadSafe {
		t.Error("Threadsafe bool should not be set on default make call.")
	}
	if stack.lock == nil {
		t.Error("Lock should not be nil after make call.")
	}

	var s Stack
	s = MakeSliceStack()

	if s.Len() != 0 {
		t.Error("Length of empty stack should be 0.")
	}
}

// -------------------------------------------------------
// Test Container Methods
// -------------------------------------------------------

func TestSliceStackLen(t *testing.T) {
	stack := MakeSliceStack()

	for i := 0; i < 50; i++ {
		stack.backer = append(stack.backer, adts.IntElt(i))

		if stack.Len() != i+1 {
			t.Errorf("Stack should have length %d, actual length: %d", i+1, stack.Len())
			return
		}
	}
}

func TestSliceStackIsEmpty(t *testing.T) {
	stack := MakeSliceStack()

	if !stack.IsEmpty() {
		t.Errorf("Empty stack should return true for IsEmpty.\n")
	}

	stack.Add(adts.IntElt(0))
	if stack.IsEmpty() {
		t.Errorf("Stack with 1 element should not return true for IsEmpty.\n")
	}
}

func TestSliceStackClear(t *testing.T) {
	stack := MakeSliceStack()

	for i := 0; i < 100; i++ {
		stack.Add(adts.IntElt(i))
	}

	if stack.Len() != 100 {
		t.Errorf("Check Add/Len method, should have length of 100 after 100 adds.\n")
	}

	stack.Clear()
	if !stack.IsEmpty() {
		t.Errorf("Stack should be empty after call to Clear.\n")
	}
}

func TestSliceStackAdd(t *testing.T) {
	vals := []int{}
	r := rand.New(rand.NewSource(99))

	stack := MakeSliceStack()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals = append(vals, v)

		if !stack.Add(adts.IntElt(v)) {
			t.Errorf("Failed to add %d to list\n", v)
			return
		}

		if len(stack.backer) != i+1 {
			t.Errorf("Stack size not correct. Expected: %d, Actual: %d", i+1, len(stack.backer))
			return
		}
		for idx, v := range vals {
			if !stack.backer[idx].Equals(adts.IntElt(v)) {
				t.Errorf("Add failed to add the value to the proper index | (idx,val) - Expected: (%d, %d), Actual: (%d, %v)",
					idx, v, idx, stack.backer[idx])
			}

		}
	}
}

func TestSliceStackContains(t *testing.T) {
	vals := make(map[int]bool)
	r := rand.New(rand.NewSource(99))

	stack := MakeSliceStack()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals[v] = true

		stack.Add(adts.IntElt(v))

		for v := range vals {
			if !stack.Contains(adts.IntElt(v)) {
				t.Errorf("Contains failed to find the value (%d) in the list.", v)
			}
		}
	}
}

func TestSliceStackRemove(t *testing.T) {
	max := 1000

	stack := MakeSliceStack()
	for i := 0; i < max; i++ {
		stack.Add(adts.IntElt(i))
	}

	for i := 0; i < max; i++ {
		idx := rand.Int31n(int32(max - i))
		val := stack.backer[idx]
		if !stack.Remove(val) {
			t.Errorf("Failed to remove index %d from list.", idx)
			return
		}

		if stack.Len() != max-1-i {
			t.Errorf("Expected length: %d, Actual length: %d", max-1-i, stack.Len())
			return
		}

		if stack.Contains(val) {
			t.Errorf("Value %v was not actually removed from list.", val)
			return
		}
	}

	// Now we want to check that the resize didn't break the ability to add.
	for i := 0; i < max; i++ {
		stack.Add(adts.IntElt(i))
		if stack.Len() != i+1 {
			t.Error("Unable to add after removing all elements.")
			return
		}
	}
}

// -------------------------------------------------------
// Test Stack Methods
// -------------------------------------------------------

func TestSliceStackPush(t *testing.T) {
	r := rand.New(rand.NewSource(99))

	expected := []int{}
	stack := MakeSliceStack()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		expected = append(expected, v)
		stack.Push(adts.IntElt(v))

		if stack.Len() != len(expected) {
			t.Errorf("Push didn't add the element properly to the stack. Expected len: %d, Actual len: %d",
				len(expected), stack.Len())
			return
		}
		for idx, exp := range expected {
			act := stack.backer[idx]
			if !act.Equals(adts.IntElt(exp)) {
				t.Errorf("Stack order was ruined by push. (idx, val) - Expected: (%d, %d), Actual: (%d, %v)",
					idx, exp, idx, act)
				return
			}

		}
	}
}

func TestSliceStackPop(t *testing.T) {
	stack := MakeSliceStack()

	for i := 0; i < 100; i++ {
		stack.Push(adts.IntElt(i))
	}

	for i := 0; i < 100; i++ {
		popped := stack.Pop()

		if !popped.Equals(adts.IntElt(100 - i - 1)) {
			t.Errorf("Pop didn't return the proper element. Expected: %v, Actual: %v\n", 100-i-1, popped)
			return
		}
	}
}
