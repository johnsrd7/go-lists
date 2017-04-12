package queueadts

import (
	"container/list"
	"sync"

	adts "github.com/johnsrd7/go-adts"
)

// ListQueue is a simple type that implements the Stack interface (both threadsafe and not).
type ListQueue struct {
	backer     *list.List
	lock       *sync.Mutex
	threadSafe bool
}

// MakeListQueue creates a non-threadsafe ListQueue
func MakeListQueue() *ListQueue {
	return &ListQueue{list.New(), &sync.Mutex{}, false}
}

// MakeListQueueThreadSafe creates a threadsafe ListQueue
func MakeListQueueThreadSafe() *ListQueue {
	return &ListQueue{list.New(), &sync.Mutex{}, true}
}

// -------------------------------------------------------
// Container Methods
// -------------------------------------------------------

// Len returns the number of elements in the queue.
func (lq *ListQueue) Len() int {
	if lq.threadSafe {
		lq.lock.Lock()
		defer lq.lock.Unlock()
		return lq.backer.Len()
	}

	return lq.backer.Len()
}

// IsEmpty returns if the queue is empty or not.
func (lq *ListQueue) IsEmpty() bool {
	return lq.Len() == 0
}

// Clear removes all elements from the queue.
func (lq *ListQueue) Clear() {
	if lq.threadSafe {
		lq.lock.Lock()
		defer lq.lock.Unlock()
		lq.backer.Init()
		return
	}

	lq.backer.Init()
}

// Contains returns true if the given item is in the queue.
func (lq *ListQueue) Contains(item adts.ContainerElement) bool {
	if lq.threadSafe {
		lq.lock.Lock()
		defer lq.lock.Unlock()
		return lq.containsHelper(item)
	}

	return lq.containsHelper(item)
}

// containsHelper checks the list for the given element (in a non-threadsafe way).
func (lq *ListQueue) containsHelper(item adts.ContainerElement) bool {
	for tmp := lq.backer.Front(); tmp != nil; tmp = tmp.Next() {
		if v, ok := tmp.Value.(adts.ContainerElement); ok {
			if v.Equals(item) {
				return true
			}
		}
	}

	return false
}

// Add returns true if the given element was added to the top of the queue.
func (lq *ListQueue) Add(item adts.ContainerElement) bool {
	if lq.threadSafe {
		lq.lock.Lock()
		defer lq.lock.Unlock()
		return lq.backer.PushBack(item) != nil
	}

	return lq.backer.PushBack(item) != nil
}

// Remove returns true if the given element was removed.
func (lq *ListQueue) Remove(item adts.ContainerElement) bool {
	if lq.threadSafe {
		lq.lock.Lock()
		defer lq.lock.Unlock()
		return lq.removeHelper(item)
	}

	return lq.removeHelper(item)
}

func (lq *ListQueue) removeHelper(item adts.ContainerElement) bool {
	for tmp := lq.backer.Front(); tmp != nil; tmp = tmp.Next() {
		if v, ok := tmp.Value.(adts.ContainerElement); ok {
			if v.Equals(item) {
				removed, ok := lq.backer.Remove(tmp).(adts.ContainerElement)
				return ok && removed.Equals(item)
			}
		}
	}

	return false
}

// -------------------------------------------------------
// Queue Methods
// -------------------------------------------------------

// Enqueue pushes the given element onto the back of the queue.
func (lq *ListQueue) Enqueue(item adts.ContainerElement) bool {
	return lq.Add(item)
}

// Dequeue removes the element from the front of the queue and returns the element.
func (lq *ListQueue) Dequeue() (adts.ContainerElement, bool) {
	if lq.threadSafe {
		lq.lock.Lock()
		defer lq.lock.Unlock()
		return lq.dequeueHelper()
	}

	return lq.dequeueHelper()
}

// dequeueHelper removes the element from the front of the queue and returns the element.
func (lq *ListQueue) dequeueHelper() (adts.ContainerElement, bool) {
	// Have to use list.List's Len function since we might already
	// be inside a lock and we don't want to deadlock.
	if lq.backer.Len() == 0 {
		return adts.EmptyContainerElement{}, false
	}

	if lastElt, ok := lq.backer.Front().Value.(adts.ContainerElement); ok {
		lq.backer.Remove(lq.backer.Front())
		return lastElt, true
	}

	return adts.EmptyContainerElement{}, false
}
