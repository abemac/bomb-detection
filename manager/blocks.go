package manager

import (
	"container/heap"
	"sync"
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

func getBlockLat(row int) float64 {
	latsPerBlock := 180.0 / numrows
	return latsPerBlock*float64(row) + latsPerBlock/2.0 - 90.0
}
func getBlockLong(col int) float64 {
	longsPerBlock := 360.0 / numcols
	return longsPerBlock*float64(col) + longsPerBlock/2.0 - 180
}

// An Item is something we manage in a priority queue.
type Block struct {
	count       int64 // The priority of the item in the queue.
	index       int   // The index of the item in the heap.
	lastVisited int64
	row         int
	col         int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Block

var blocks = [numrows][numcols]*Block{}
var pq PriorityQueue
var pqMutex sync.RWMutex
var blockMutex sync.RWMutex

func init() {
	pq = make(PriorityQueue, numrows*numcols)
	for i := 0; i < numrows; i++ {
		for j := 0; j < numcols; j++ {
			blocks[i][j] = &Block{row: i,
				col:   j,
				index: i*numcols + j,
				count: 0}
			pq[i*numcols+j] = blocks[i][j]
		}
	}
	heap.Init(&pq)
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].lastVisited < pq[j].lastVisited
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
func (pq *PriorityQueue) updateCount(b *Block, count int64) {
	b.count = count
	heap.Fix(pq, b.index)
}
func (pq *PriorityQueue) updateVisitedTime(b *Block, time int64) {
	b.lastVisited = time
	heap.Fix(pq, b.index)
}
