package syncdto

import (
	"sync"
)

type SafeList[T any] struct {
	mu    sync.Mutex
	items []T
}

type SafeMap[T any] struct {
	mu     sync.Mutex
	mapSys map[string]T
}

func NewSafeList[T any]() SafeList[T] {
	return SafeList[T]{
		mu:    sync.Mutex{},
		items: make([]T, 0),
	}
}

func NewSafeMap[T any]() SafeMap[T] {
	return SafeMap[T]{
		mu:     sync.Mutex{},
		mapSys: make(map[string]T),
	}
}

func (s *SafeMap[T]) ForEach(f func(key string, channel T)) {
	for key, value := range s.mapSys {
		f(key, value)
	}
}

func (safeList *SafeList[T]) PushItem(item T) {
	safeList.mu.Lock()
	defer safeList.mu.Unlock()
	safeList.items = append(safeList.items, item)
}

func (safeList *SafeList[T]) GetSize() int {
	return len(safeList.items)
}

func (safeList *SafeList[T]) PopItem(index int) T {
	safeList.mu.Lock()
	defer safeList.mu.Unlock()

	poppedItem := safeList.items[index]
	safeList.items = append(safeList.items[:index], safeList.items[index+1:]...)

	return poppedItem
}

func (safeMap *SafeMap[T]) PushMap(key string, value T) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()

	safeMap.mapSys[key] = value
}

func (safeMap *SafeMap[T]) GetMap(key string) (T, bool) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()

	value, exists := safeMap.mapSys[key]
	return value, exists
}

func (safeMap *SafeMap[T]) IsEmpty() bool {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()

	return len(safeMap.mapSys) == 0
}

func (safeMap *SafeMap[T]) ThereIsKey(key string) bool {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()

	_, exists := safeMap.mapSys[key]
	return exists
}