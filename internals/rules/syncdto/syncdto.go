package syncdto

import (
	"sync"
)

type SafeList[T any] struct {
	mu    sync.Mutex
	items []T
}

type SafeMap[T any] struct {
	mu		 sync.Mutex
	mapSys map[string]T
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

func PushMap[T any](safeMap *SafeMap[T], key string, value T) {
	safeMap.mu.Lock()
	defer safeMap.mu.Unlock()

	safeMap.mapSys[key] = value
}