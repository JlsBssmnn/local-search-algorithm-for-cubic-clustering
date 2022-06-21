package algorithm

import (
	"math"

	"github.com/JlsBssmnn/local-search-algorithm-for-cubic-clustering/src/utils"
)

type GreedyMovingAlgorithm[data any] struct {
	input        *[]data
	calc         CostCalculator[data]
	constraints  *Constraints
	partitioning PartitioningArray
	partitions   map[int]*[]int
	tripleCosts  *TripleCosts
	removeCosts  map[int]float64
	costs        *GreedyMovingCosts
}

type GreedyMovingCosts []*MovingSecondDim

type MovingSecondDim struct {
	minCost  float64
	bestMove int
	moves    []*OneElementMove
}

type OneElementMove struct {
	valid       bool
	cost        float64
	bestMove    int
	doubleMoves *[]DoubleMove
}

type DoubleMove struct {
	valid bool
	cost  float64
}

type TripleCosts [][][]float64

// -------------------------- Methods for the cost data structure and the GreedyJoiningAlgorithm struct

// This function gets the triple cost for the elements i, j and k out of the TripleCost struct
func (costs *TripleCosts) GetTripleCost(i, j, k int) float64 {
	utils.SortInts(&i, &j, &k)
	return (*costs)[i][j-i-1][k-j-1]
}

// This function computes the index for the double move slice for the move of element i to the partition
// of element j together with element k
func getDoubleMoveIndex(i, j, k int) int {
	index := k - i - 1
	if j > i && j < k {
		index--
	}
	return index
}

// This function computes the element that corresponds to the index k in the double move slice for
// the move of element i to the partition of element j (inverse of getDoubleMoveIndex)
func getDoubleMoveElement(i, j, k int) int {
	element := i + 1 + k
	if j > i && j <= element {
		element++
	}
	return element
}

// Gets the cost for moving the given element into a singleton partition
func (algorithm *GreedyMovingAlgorithm[data]) getRemoveCost(element int) float64 {
	if oem := (*algorithm.costs)[element].moves[element]; oem.valid {
		return oem.cost
	} else {
		return 0
	}
}

// Invalidates the cost of moving element1 into the partition of element2
func (algorithm *GreedyMovingAlgorithm[data]) invalidateCost(element1, element2 int) {
	(*algorithm.costs)[element1].moves[element2].valid = false
}

// Makes the cost of moving element1 into the partition of element2 valid if it's not forbidden by
// the constaints
func (algorithm *GreedyMovingAlgorithm[data]) validateCost(element1, element2 int) {
	(*algorithm.costs)[element1].moves[element2].valid = !algorithm.constraints.Get(element1, element2)
}

// returns the cost of moving at index x in the first dimension and index y in the second dimension
func (cost *GreedyMovingCosts) moveCost(x, y int) float64 {
	return (*cost)[x].moves[y].cost
}

// This function will adjust the OneElementMove struct s.t. double moves are considered. This should
// be used if a partition lost an element and now has size 1 or the singleton element was moved into
// a singleton set.
func (algorithm *GreedyMovingAlgorithm[data]) validateDoubleMove(element, singleton int) {
	oem := (*algorithm.costs)[element].moves[singleton]
	if oem.bestMove > -1 {
		return
	}
	oem.bestMove *= -1
	oem.bestMove--
	index := getDoubleMoveIndex(element, singleton, oem.bestMove)
	if oem.doubleMoves == nil {
		oem.valid = false
	} else if index < 0 {
		panic("best move for double move yields invalid index")
	} else {
		oem.cost = (*oem.doubleMoves)[index].cost
	}
}

// This function will adjust the OneElementMove struct s.t. double moves are not considered. This should
// be used if a partition got an element and now has size > 1.
func (algorithm *GreedyMovingAlgorithm[data]) invalidateDoubleMove(element1, element2 int) {
	oem := (*algorithm.costs)[element1].moves[element2]
	if oem.bestMove < 0 {
		return
	}
	oem.bestMove++
	oem.bestMove *= -1
}

// Computes the sum of all triple costs where the first two elements in the triple are the two
// given elements and the third element is in the partition of the given `part` element
func (algorithm *GreedyMovingAlgorithm[data]) tripleCostSum(element1, element2, part int) float64 {
	partition := algorithm.partitions[part]
	sum := 0.0

	for _, thirdElement := range *partition {
		if thirdElement == element1 || thirdElement == element2 {
			continue
		}
		sum += algorithm.tripleCosts.GetTripleCost(element1, element2, thirdElement)
	}
	return sum
}

// Updates remove cost for an element in the destination partition
func (algorithm *GreedyMovingAlgorithm[data]) updateRemoveCostDest(element, rElement int) {
	algorithm.validateCost(element, element)
	diff := -algorithm.tripleCostSum(element, rElement, element)
	(*algorithm.costs)[element].moves[element].cost += diff

	algorithm.removeCosts[element] = diff
}

// Updates the cost for removing the given element from it's partition when the given rElement is not
// in this partition anymore. This function should only be called if element and rElement are in the
// same partition.
func (algorithm *GreedyMovingAlgorithm[data]) updateRemoveCostSource(element, rElement int) {
	oem := (*algorithm.costs)[element].moves[element]

	if len(*algorithm.partitions[element]) <= 2 {
		algorithm.invalidateCost(element, element)
		algorithm.removeCosts[element] = -oem.cost
		oem.cost = 0
	} else {
		diff := algorithm.tripleCostSum(element, rElement, rElement)
		algorithm.removeCosts[element] = diff
		oem.cost += diff
	}
}

// This function applies a change in remove costs to all elements in the second dimension of the element
// except for the ignore elements
func (algorithm *GreedyMovingAlgorithm[data]) updateCostSecondDim(element int, ignore ...int) {
	diff := algorithm.removeCosts[element]
	msd := (*algorithm.costs)[element]
	for j := 0; j < len(msd.moves); j++ {
		oem := msd.moves[j]
		if j == element || !oem.valid || utils.Contains(ignore, j) {
			continue
		}
		msd.moves[j].cost += diff
	}
}

// Gets the costs that arrise when moving the given element into the partition of the given destination element.
// Only the costs that occur in the destination partition are considered.
func (algorithm *GreedyMovingAlgorithm[data]) getRealMoveCost(element, destination int) float64 {
	oem := (*algorithm.costs)[element].moves[destination]
	if oem.bestMove < 0 {
		return oem.cost
	} else {
		return 0
	}
}

// This function creates a OneElementMove struct for the case of moving the given element
// into a singleton partition which just contains the given element `singleton`
func (algorithm *GreedyMovingAlgorithm[data]) createOneElementMove(element, singleton int) *OneElementMove {
	oem := OneElementMove{}
	invalid := algorithm.constraints.Get(element, singleton)
	if invalid || element == singleton {
		oem.valid = false
		return &oem
	}
	n := len(*algorithm.input)
	minCost3D := math.Inf(1)
	bestMove3D := -1

	doubleMovesLen := n - 1 - element
	if singleton > element {
		doubleMovesLen--
	}
	if doubleMovesLen < 1 {
		oem.valid = false
		return &oem
	}
	doubleMoves := make([]DoubleMove, doubleMovesLen)
	offset := element + 1 // the offset from the k to the index in the input data

	for k := 0; k < doubleMovesLen; k++ {
		if k+offset == singleton {
			offset++
		}
		kElement := k + offset
		if algorithm.constraints.Get(element, kElement) || algorithm.constraints.Get(singleton, kElement) {
			doubleMoves[k] = DoubleMove{valid: false}
			continue
		}
		cost := algorithm.tripleCosts.GetTripleCost(element, singleton, kElement)
		doubleMoves[k] = DoubleMove{valid: true, cost: cost}

		if cost < minCost3D {
			minCost3D = cost
			bestMove3D = kElement
		}
	}
	oem.valid = true
	oem.cost = minCost3D
	oem.bestMove = bestMove3D
	oem.doubleMoves = &doubleMoves

	return &oem
}

// Updates the cost of moving the given element into the previous partition of rElement
// when removing rElement from this partition. Umin is the smallest element in the partition
// of rElement without rElement itself. (Umin is -1 if rElement was in a singleton partition)
func (algorithm *GreedyMovingAlgorithm[data]) updateMoveCostSource(element, rElement, Umin int) {
	if Umin == -1 {
		algorithm.invalidateCost(element, rElement)
		return
	}
	partition := algorithm.partitions[rElement]

	if len(*partition) == 2 {
		algorithm.validateCost(element, Umin)
		algorithm.validateDoubleMove(element, Umin)
		return
	}

	diff := algorithm.tripleCostSum(element, rElement, rElement)
	formerRepresentative := utils.Min([]int{rElement, Umin})
	removeCostDiff := algorithm.removeCosts[element]

	oem := (*algorithm.costs)[element].moves[Umin]
	oem.cost = (*algorithm.costs)[element].moves[formerRepresentative].cost - diff + removeCostDiff
	algorithm.validateCost(element, Umin)
	algorithm.invalidateDoubleMove(element, Umin)
}

// Updates the cost of moving the given element into the partition of Umin when rElement was also
// moved to this partition. Umin is the smallest element in the destination partition without rElement
// (-1 if the element is moved to a singleton).
func (algorithm *GreedyMovingAlgorithm[data]) updateMoveCostTarget(element, rElement, Umin int) {
	if Umin == -1 {
		algorithm.validateCost(element, rElement)
		algorithm.validateDoubleMove(element, rElement)
		return
	}
	var newRepresentative int
	var otherElement int
	if rElement < Umin {
		newRepresentative = rElement
		otherElement = Umin
	} else {
		newRepresentative = Umin
		otherElement = rElement
	}

	oem := (*algorithm.costs)[element].moves[newRepresentative]

	if len(*algorithm.partitions[Umin]) == 1 {
		var removeCost float64
		if oem := (*algorithm.costs)[element].moves[element]; oem.valid {
			removeCost = oem.cost
		}
		oem.cost = algorithm.tripleCosts.GetTripleCost(element, rElement, Umin) + removeCost
	} else {
		diff := algorithm.tripleCostSum(element, rElement, Umin)
		removeCostDiff := algorithm.removeCosts[element]
		oem.cost = (*algorithm.costs)[element].moves[Umin].cost + diff + removeCostDiff
	}
	algorithm.validateCost(element, newRepresentative)
	algorithm.invalidateDoubleMove(element, newRepresentative)
	algorithm.invalidateCost(element, otherElement)
}

// This function updates the partitioning array and partitions map of the algorithm
// struct to the move of the given element to the given partition of the element UminDest.
// UminDest has to be the smallest element in it's partition or -1 if the element is moved
// into a singleton. UminSource is the smallest element of the previous partition of the
// element without the element itself or -1 if the element was in a singleton.
func (algorithm *GreedyMovingAlgorithm[data]) updatePartitioning(UminSource, UminDest, element int) {
	// update partitioning array
	var formerSourceRepresentative int
	var newDestRepresentative int
	if element < UminSource {
		formerSourceRepresentative = element
	} else {
		formerSourceRepresentative = UminSource
	}
	if element < UminDest || UminDest == -1 {
		newDestRepresentative = element
	} else {
		newDestRepresentative = UminDest
	}
	for i := 0; i < len(*algorithm.input); i++ {
		if i == element {
			algorithm.partitioning[i] = newDestRepresentative
		} else if algorithm.partitioning[i] == formerSourceRepresentative {
			algorithm.partitioning[i] = UminSource
		} else if algorithm.partitioning[i] == UminDest {
			algorithm.partitioning[i] = newDestRepresentative
		}
	}

	// update partitions map
	if UminSource != -1 {
		utils.DeleteByElement(algorithm.partitions[UminSource], element)
	}
	if UminDest != -1 {
		utils.InsertInOrder(algorithm.partitions[UminDest], element)
		algorithm.partitions[element] = algorithm.partitions[UminDest]
	} else {
		algorithm.partitions[element] = &[]int{element}
	}
}

// --------------------------

// Initializes the TripleCost data structure which will store every combination of triple costs
func (algorithm *GreedyMovingAlgorithm[data]) InitializeTripleCosts() {
	n := len(*algorithm.input)
	firstDim := make(TripleCosts, n-2)

	for i := 0; i < n-2; i++ {
		secondDim := make([][]float64, n-i-2)

		for j := i + 1; j < n-1; j++ {
			thirdDim := make([]float64, n-j-1)

			for k := j + 1; k < n; k++ {
				thirdDim[k-j-1] = algorithm.calc.TripleCost((*algorithm.input)[i], (*algorithm.input)[j], (*algorithm.input)[k])
			}
			secondDim[j-i-1] = thirdDim
		}
		firstDim[i] = secondDim
	}
	algorithm.tripleCosts = &firstDim
}

func (algorithm *GreedyMovingAlgorithm[data]) InitializeCosts() ([3]int, float64) {
	U, a, b := -1, -1, -1
	n := len(*algorithm.input)
	bestCostOverall := math.Inf(1)

	greedyMovingCosts := make(GreedyMovingCosts, n)
	for i := 0; i < n; i++ {
		minCost := math.Inf(1)
		bestMove := -1
		moves := make([]*OneElementMove, n)
		for j := 0; j < n; j++ {
			valid := i != j
			if !valid {
				moves[j] = &OneElementMove{valid: false}
				continue
			}
			oem := algorithm.createOneElementMove(i, j)
			moves[j] = oem
			if oem.cost < minCost {
				minCost = oem.cost
				bestMove = j
			}
		}
		movingSecondDim := MovingSecondDim{minCost: minCost, bestMove: bestMove, moves: moves}
		greedyMovingCosts[i] = &movingSecondDim

		if minCost < bestCostOverall {
			bestCostOverall = minCost
			a = i
			U = bestMove
			b = movingSecondDim.moves[bestMove].bestMove
		}
	}

	algorithm.costs = &greedyMovingCosts
	return [3]int{U, a, b}, bestCostOverall
}

// Updates the cost data structure when moving the given element into the partition of the
// element `partition`. If the element is moved into a singleton, `partition` should be -1.
//
// This function also returns the best next move in a element-array in the form [U, a, b].
// If U is -1 then element a should be moved into a singleton. If b is -1, then only a should
// be moved, otherwise a and b should be moved into U.
func (algorithm *GreedyMovingAlgorithm[data]) Move(partition, element int) ([3]int, float64) {
	U, a, b := -1, -1, -1
	n := len(*algorithm.input)
	bestCostOverall := math.Inf(1)
	algorithm.removeCosts = make(map[int]float64)

	// The smallest element in the partition of the removed element without the element itself
	UminSource := -1
	for _, e := range *algorithm.partitions[element] {
		if e != element && (e < UminSource || UminSource == -1) {
			UminSource = e
		}
	}

	UminDest := -1
	if partition != -1 {
		for _, e := range *algorithm.partitions[partition] {
			if e < UminDest || UminDest == -1 {
				UminDest = e
			}
		}
	}

	if UminDest == -1 && UminSource == -1 {
		panic("Both UminDest and UminSource were -1")
	}

	// The partition where the element was in
	ePart := algorithm.partitions[element]
	indexEPart := 0

	// The destination partition
	destPart := algorithm.partitions[partition]
	indexDPart := 0

	for i := 0; i < n; i++ {
		if i == element {
			// The element that is moved is considered
			indexEPart++
			oem := (*algorithm.costs)[i].moves[i]

			if UminSource == -1 {
				// The partition where the moved element is comming from is a singleton
				algorithm.invalidateCost(i, UminDest)
				newRemoveCost := -(algorithm.getRealMoveCost(i, UminDest) - oem.cost)

				algorithm.removeCosts[i] = -oem.cost + newRemoveCost
				oem.cost = newRemoveCost
				algorithm.validateCost(i, i)
				algorithm.updateCostSecondDim(i)
			} else if UminDest == -1 {
				// The element is moved into a singleton set
				algorithm.removeCosts[i] = -oem.cost
				oem.cost = 0
				algorithm.invalidateCost(i, i)
				(*algorithm.costs)[i].moves[UminSource].cost = -oem.cost
				algorithm.validateCost(i, UminSource)
				algorithm.updateCostSecondDim(i, UminSource)
				if len(*algorithm.partitions[UminSource]) > 2 {
					algorithm.invalidateDoubleMove(i, UminSource)
				}
			} else {
				(*algorithm.costs)[i].moves[UminSource].cost = -(*algorithm.costs)[i].moves[UminDest].cost
				newRemoveCost := -(algorithm.getRealMoveCost(i, UminDest) - oem.cost)
				algorithm.removeCosts[i] = -oem.cost + newRemoveCost
				algorithm.invalidateCost(i, UminDest)
				oem.cost = newRemoveCost
				algorithm.validateCost(i, UminSource)
				algorithm.updateCostSecondDim(i, UminSource)
				if len(*algorithm.partitions[UminSource]) > 2 {
					algorithm.invalidateDoubleMove(i, UminSource)
				}
			}
		} else if indexEPart < len(*ePart) && i == (*ePart)[indexEPart] {
			// An element in the partition of the moved element is considered
			indexEPart++
			algorithm.updateRemoveCostSource(i, element)
			algorithm.updateMoveCostTarget(i, element, UminDest)
			var newRepresentative int
			if element < UminDest {
				newRepresentative = element
			} else {
				newRepresentative = UminDest
			}
			algorithm.updateCostSecondDim(i, newRepresentative)
			// algorithm.updateCostSecondDim(i)
		} else if destPart != nil && indexDPart < len(*destPart) && i == (*destPart)[indexDPart] {
			// An element in the partition that the moved element is moved to is considered
			indexDPart++
			algorithm.updateRemoveCostDest(i, element)
			algorithm.updateMoveCostSource(i, element, UminSource)
			algorithm.invalidateCost(i, element)
			algorithm.updateCostSecondDim(i, UminSource)
		} else {
			algorithm.updateMoveCostSource(i, element, UminSource)
			algorithm.updateMoveCostTarget(i, element, UminDest)
		}
	}

	algorithm.updatePartitioning(UminSource, UminDest, element)

	// Update new bestCost for element i, this must be done in a new loop because it
	// uses the adjusted costs of other elements
	for i := 0; i < n; i++ {
		minCostOneMove := math.Inf(1)
		bestMoveOneMove := -1
		msd := (*algorithm.costs)[i]

		for j, oem := range msd.moves {
			if j == i {
				continue
			} else if oem.doubleMoves == nil {
				if oem.valid && oem.cost < minCostOneMove {
					minCostOneMove = oem.cost
					bestMoveOneMove = j
				}
				continue
			}

			// recompute all double moves
			invalid := oem.bestMove < 0
			minCostDoubleMove := math.Inf(1)
			bestMoveDoubleMove := -1
			for k := 0; k < len(*oem.doubleMoves); k++ {
				kElement := getDoubleMoveElement(i, j, k)
				(*oem.doubleMoves)[k].cost = algorithm.tripleCosts.GetTripleCost(i, j, kElement) + algorithm.getRemoveCost(i) + algorithm.getRemoveCost(kElement)

				// if i and k are in the same partition, some costs were considered twice
				// so they have to be substracted again
				if utils.Contains(*algorithm.partitions[i], kElement) {
					(*oem.doubleMoves)[k].cost += algorithm.tripleCostSum(i, kElement, i)
				}
				if (*oem.doubleMoves)[k].valid && (*oem.doubleMoves)[k].cost < minCostDoubleMove {
					minCostDoubleMove = (*oem.doubleMoves)[k].cost
					bestMoveDoubleMove = k
				}
			}

			if invalid {
				oem.bestMove = -(getDoubleMoveElement(i, j, bestMoveDoubleMove) + 1)
			} else {
				oem.bestMove = getDoubleMoveElement(i, j, bestMoveDoubleMove)
				oem.cost = minCostDoubleMove
			}

			if !oem.valid {
				continue
			} else if oem.cost < minCostOneMove {
				minCostOneMove = oem.cost
				bestMoveOneMove = j
			}
		}

		// The cost for moving element i into a singleton was not considered yet
		if msd.moves[i].valid && msd.moves[i].cost < minCostOneMove {
			minCostOneMove = msd.moves[i].cost
			bestMoveOneMove = i
		}

		msd.minCost = minCostOneMove
		msd.bestMove = bestMoveOneMove

		// check if the bestMove for i is better overall
		if msd.minCost < bestCostOverall {
			bestCostOverall = msd.minCost
			a = i
			if msd.bestMove == i {
				U = -1
			} else {
				U = msd.bestMove
			}
			if msd.moves[msd.bestMove].bestMove < 0 {
				b = -1
			} else {
				b = msd.moves[msd.bestMove].bestMove
			}
		}
	}

	return [3]int{U, a, b}, bestCostOverall
}

func (algorithm *GreedyMovingAlgorithm[data]) Initialize() ([3]int, float64) {
	algorithm.InitializeTripleCosts()
	n := len(*algorithm.input)
	algorithm.partitioning = make(PartitioningArray, n)
	algorithm.partitions = make(map[int]*[]int, n)

	for i := 0; i < n; i++ {
		algorithm.partitioning[i] = i
		algorithm.partitions[i] = &[]int{i}
	}
	return algorithm.InitializeCosts()
}

// The greedy moving algorithm with following properties:
// 	- it will only evaluate moves of 2 elements if the destination partition has 1 element
// 	- it will move one element if the destination partition has more than 1 element,
//		otherwise it will move 2 elements
// 	- if there is only one partition left the algorithm terminates
func ImprovedGreedyMoving[data any](input *[]data, calc CostCalculator[data]) PartitioningArray {
	algorithm := GreedyMovingAlgorithm[data]{input: input, calc: calc}
	nextMove, cost := algorithm.Initialize()
	i := 0

	for cost < 0 && nextMove[1] != -1 {
		i++
		newNextMove, newCost := algorithm.Move(nextMove[0], nextMove[1])
		if nextMove[2] != -1 {
			nextMove, cost = algorithm.Move(nextMove[0], nextMove[2])
		} else {
			nextMove = newNextMove
			cost = newCost
		}
	}
	return algorithm.partitioning
}
