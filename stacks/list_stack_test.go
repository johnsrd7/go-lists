package stackadts

import (
	"fmt"
	"math/rand"
	"testing"

	adts "github.com/johnsrd7/go-adts"
)

func TestMakeListStack(t *testing.T) {
	stack := MakeListStack()

	if stack.backer == nil {
		t.Error("Backing Container should not be nil after make call.")
	}
	if stack.lock == nil {
		t.Error("Lock for stack should not be nil after make call.")
	}
	if stack.threadSafe {
		t.Error("Threadsafe should be false for call to non-threadsafe make.")
	}

	var s Stack
	s = MakeListStack()

	if s.Len() != 0 {
		t.Error("Length of empty stack should be 0.")
	}
}

// -------------------------------------------------------
// Test Container Methods
// -------------------------------------------------------

func TestListStackLen(t *testing.T) {
	stack := MakeListStack()

	for i := 0; i < 50; i++ {
		stack.Add(adts.IntElt(i))

		if stack.Len() != i+1 {
			t.Errorf("Stack should have length %d, actual length: %d", i+1, stack.Len())
			return
		}
	}
}

func TestListStackIsEmpty(t *testing.T) {
	stack := MakeListStack()

	if !stack.IsEmpty() {
		t.Errorf("Empty stack should return true for IsEmpty.\n")
	}

	stack.Add(adts.IntElt(0))
	if stack.IsEmpty() {
		t.Errorf("Stack with 1 element should not return true for IsEmpty.\n")
	}
}

func TestListStackClear(t *testing.T) {
	stack := MakeListStack()

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

func TestListStackAdd(t *testing.T) {
	vals := []int{}
	r := rand.New(rand.NewSource(99))

	stack := MakeListStack()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		vals = append(vals, v)

		if !stack.Add(adts.IntElt(v)) {
			t.Errorf("Failed to add %d to list\n", v)
			return
		}

		if stack.Len() != i+1 {
			t.Errorf("Stack size not correct. Expected: %d, Actual: %d", i+1, stack.Len())
			return
		}

		iter := stack.backer.Front()
		for j := 0; j < len(vals); j++ {
			idx := len(vals) - 1 - j
			v = vals[idx]
			iterVal, ok := iter.Value.(adts.IntElt)
			if !ok || !iterVal.Equals(adts.IntElt(v)) {
				t.Errorf("Add failed to add the value to the proper index | (idx,val) - Expected: (%d, %d), Actual: (%d, %v)",
					idx, v, j, iterVal)
			}
			iter = iter.Next()
		}
	}
}

func TestListStackContains(t *testing.T) {
	vals := make(map[int]bool)
	r := rand.New(rand.NewSource(99))

	stack := MakeListStack()

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

func TestListStackRemove(t *testing.T) {
	max := 1000

	stack := MakeListStack()
	for i := 0; i < max; i++ {
		stack.Add(adts.IntElt(i))
	}

	for i := 0; i < max; i++ {
		idx := rand.Int31n(int32(max - i))
		// get to the index
		tmp := stack.backer.Front()
		for j := 0; j < int(idx); j++ {
			tmp = tmp.Next()
		}
		val, _ := tmp.Value.(adts.IntElt)
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

func TestListStackPush(t *testing.T) {
	r := rand.New(rand.NewSource(99))

	expected := []int{}
	stack := MakeListStack()

	for i := 0; i < 1000; i++ {
		v := r.Int()

		expected = append(expected, v)
		if !stack.Push(adts.IntElt(v)) {
			t.Errorf("Push didn't succeed for push #%d\n", i+1)
			return
		}

		if stack.Len() != len(expected) {
			t.Errorf("Push didn't add the element properly to the stack. Expected len: %d, Actual len: %d",
				len(expected), stack.Len())
			return
		}
		iter := stack.backer.Front()
		for j := 0; j < len(expected); j++ {
			idx := len(expected) - 1 - j
			exp := expected[idx]
			act, ok := iter.Value.(adts.IntElt)
			if !ok || !act.Equals(adts.IntElt(exp)) {
				t.Errorf("Stack order was ruined by push #%d. (idx, val) - Expected: (%d, %d), Actual: (%d, %v)",
					i+1, idx, exp, j, act)
				return
			}

			iter = iter.Next()
		}
	}
}

func TestListStackPop(t *testing.T) {
	stack := MakeListStack()

	for i := 0; i < 100; i++ {
		stack.Push(adts.IntElt(i))
	}

	for i := 0; i < 100; i++ {
		popped, ok := stack.Pop()
		if !ok {
			t.Errorf("Pop didn't succeed for pop #%d\n", i+1)
			return
		}

		if !popped.Equals(adts.IntElt(100 - i - 1)) {
			t.Errorf("Pop didn't return the proper element. Expected: %v, Actual: %v\n", 100-i-1, popped)
			return
		}
	}

	_, ok := stack.Pop()
	if ok {
		t.Errorf("Pop should not succeed for empty stack.\n")
	}
}

func listStackToString(ls *ListStack) string {
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
