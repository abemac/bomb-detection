package manager

import (
	"container/heap"
)

const numcols = 16
const numrows = 8

func getRowFromLat(lat float64) int {
	row := (lat + 90.0) / 180.0 * float64(numrows)
	if int(row) == numrows {
		return int(numrows - 1)
	}
	return int(row)
}
func getColFromLong(long float64) int {
	col := (long + 180.0) / 360.0 * float64(numcols)
	if int(col) == numcols {
		return int(numcols - 1)
	}
	return int(col)
}

// An Item is something we manage in a priority queue.
type Block struct {
	lat   float64
	long  float64
	count uint64 // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index       int // The index of the item in the heap.
	lastVisited int64
}

var blocks = [numrows][numcols]*Block{}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Block

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].count < pq[j].count
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Block)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Block, lat float64, long float64, count uint64) {
	item.lat = lat
	item.long = long
	item.count = count
	heap.Fix(pq, item.index)
}
