package manager

import (
	"container/heap"
	"sync"
)

const numBlockCols = 20
const numBlockRows = 10

// PriorityQueue implements heap.Interface and holds Blocks.
type PriorityQueue []*Block

var blocks = [numBlockRows][numBlockCols]*Block{}
var pq PriorityQueue
var blocksMutex sync.RWMutex

func init() {
	pq = make(PriorityQueue, numBlockRows*numBlockCols)
	for i := 0; i < numBlockRows; i++ {
		for j := 0; j < numBlockCols; j++ {
			blocks[i][j] = &Block{
				row:   i,
				col:   j,
				index: i*numBlockCols + j,
				count: 0}
			pq[i*numBlockCols+j] = blocks[i][j]
		}
	}
	heap.Init(&pq)
}

type Block struct {
	count       int64 // The priority of the item in the queue.
	index       int   // The index of the item in the heap.
	lastVisited int64
	row         int
	col         int
}

func getRowFromLat(lat float64) int {
	row := (lat + 90.0) / 180.0 * float64(numBlockRows)
	if int(row) == numBlockRows {
		return int(numBlockRows - 1)
	}
	return int(row)
}
func getColFromLong(long float64) int {
	col := (long + 180.0) / 360.0 * float64(numBlockCols)
	if int(col) == numBlockCols {
		return int(numBlockCols - 1)
	}
	return int(col)
}
func getBlockLat(row int) float64 {
	latsPerBlock := 180.0 / numBlockRows
	return latsPerBlock*float64(row) + latsPerBlock/2.0 - 90.0
}
func getBlockLong(col int) float64 {
	longsPerBlock := 360.0 / numBlockCols
	return longsPerBlock*float64(col) + longsPerBlock/2.0 - 180
}

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool { //Less = "is more important"
	timediff := pq[j].lastVisited - pq[i].lastVisited
	countdiff := pq[j].count - pq[i].count
	score := timediff + countdiff
	return score > 0

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

func (pq *PriorityQueue) updateCount(b *Block, count int64) {
	b.count = count
	heap.Fix(pq, b.index)
}
func (pq *PriorityQueue) updateVisitedTime(b *Block, time int64) {
	b.lastVisited = time
	heap.Fix(pq, b.index)
}
