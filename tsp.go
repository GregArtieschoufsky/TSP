package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// Point is a 2D point aligned with an x and y intersection.
type Point struct {
	index int
	x     int
	y     int
}

// String returns a string representation of the Point.
func (p Point) String() string {
	return fmt.Sprintf("[%d] %d,%d", p.index, p.x, p.y)
}

// Pair is two points with a distance between them.
type Pair struct {
	index    int
	pointA   Point
	pointB   Point
	distance float64
}

func main() {

	// Seed the random generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Randomly determine a number of points to generate
	var numberOfPoints = rand.Intn(20) + 4
	fmt.Printf("random = %d\n", numberOfPoints)

	// Build up a set of a number of random points
	points := BuildPoints(numberOfPoints)
	PrintPoints(points)

	// For each point, pair it with another point
	pairs := PairPoints(points)
	PrintPairs(pairs)

	// Calculate a route between points with shortest distance between an two points first
	fmt.Println("Calculating shortest path now")
	route := RouteShortestFirst(pairs)
	PrintPairs(route)

	// Print to console a description of the route determined
	fmt.Printf("Total Distance = %f\n", GetTotalLength(route))
	fmt.Printf("Avg Distance Of All = %f\n", GetAvgLength(pairs))
	fmt.Printf("Avg Distance Of Route = %f\n", GetAvgLength(route))
}

// GetAvgLength computes the average length of distance between a set
// of paired points.
func GetAvgLength(pairs []Pair) float64 {

	return (GetTotalLength(pairs) / float64(len(pairs)))
}

// GetTotalLength computes the aggregated length of distance between a
// set of paired points.
func GetTotalLength(pairs []Pair) float64 {
	total := 0.0
	for _, pair := range pairs {

		total = (total + pair.distance)
	}
	return total
}

// Distance returns the distance between two points.
func (p Point) Distance(p2 Point) float64 {
	first := math.Pow(float64(p2.x-p.x), 2)
	second := math.Pow(float64(p2.y-p.y), 2)
	return math.Sqrt(first + second)
}

// PrintPoints prints the list of points.
func PrintPoints(points []Point) {
	fmt.Println("\nPoints:")
	for _, point := range points {
		fmt.Println(point)
	}
}

// PrintPairs prints a list of pair descriptions.
func PrintPairs(pairs []Pair) {
	fmt.Printf("\nPairs of length (%d)\n", len(pairs))
	for _, pair := range pairs {
		fmt.Printf("[%d] %s, %s, %f\n",
			pair.index,
			pair.pointA,
			pair.pointB,
			pair.distance)
	}
}

// String returns a string representation of the Pair.
func (pair *Pair) String() string {
	return fmt.Sprintf("[%d] %s,%s,%f",
		pair.index,
		pair.pointA,
		pair.pointB,
		pair.distance)
}

// BuildPoints returns a number of randomly generated points.
func BuildPoints(numberOfPoints int) []Point {
	points := make([]Point, numberOfPoints)
	for i := 0; i < numberOfPoints; i++ {
		var randomX = rand.Intn(100)
		var randomY = rand.Intn(100)
		var index = i
		points[i] = Point{index: index, x: randomX, y: randomY}
	}
	return points
}

// PairPoints creates a pair between each given point in a set.
func PairPoints(points []Point) []Pair {
	var count = 0
	var skip = false
	pairs := make([]Pair, 0)
	for _, pointA := range points {
		for _, pointB := range points {

			// Avoid duplicates
			for _, pair := range pairs {
				if pair.pointA.index == pointA.index && pair.pointB.index == pointB.index {
					skip = true
				}
				if pair.pointA.index == pointB.index && pair.pointB.index == pointA.index {
					skip = true
				}
			}

			// Add new pair to set of pairs among points
			if pointA.index != pointB.index && !skip {
				count++
				var distance = pointA.Distance(pointB)
				pairs = append(
					pairs,
					Pair{index: count, pointA: pointA, pointB: pointB, distance: distance})
				skip = false
			} else {
				skip = false
			}
		}
	}
	return pairs
}

// RouteShortestFirst starts a routine to define a path between all
// points beginning with the shortest path.
func RouteShortestFirst(pairs []Pair) []Pair {

	// Begin with a point inside of a pair of the shortest distance
	startingPair := PairWithShortestDistance(pairs)
	startingPoint := startingPair.pointA
	fmt.Printf("Starting Point = %s\n", startingPoint)

	// The complimenting point from the starting pair must be identified and connected next
	pointToConnect := startingPair.pointB

	// Construct an empty set of pairs to define as the shortest route across all points
	routeOfPairs := make([]Pair, 0)
	startingPairDistance := startingPair.distance

	// Add the first, shortest pair to the route
	routeOfPairs = append(
		routeOfPairs,
		Pair{index: 0, pointA: startingPoint, pointB: pointToConnect,
			distance: startingPairDistance})

	// Redefine the remaining pairs of points to connect yet to be without already paired point
	pairsRefined := RemovePairsWithPoint(startingPoint, pairs)

	// Recursively pair shortest points until unmapped points no longer exist
	routeOfPairs = ConnectPointToClosest(pointToConnect, pairsRefined, routeOfPairs)

	// The last point must be joined with the first point of the route to complete the route
	lastPoint := GetLastPointToConnect(routeOfPairs)
	lastPair := GetPairFromPairs(lastPoint, startingPoint, pairs)
	routeOfPairs = append(routeOfPairs, lastPair)

	return routeOfPairs
}

// GetPairFromPairs returns a pair from a set that includes two given
// points.
func GetPairFromPairs(pointA Point, pointB Point, pairs []Pair) Pair {
	for _, pair := range pairs {

		if pair.pointA.index == pointA.index && pair.pointB.index == pointB.index {

			return pair
		}

		if pair.pointA.index == pointB.index && pair.pointB.index == pointA.index {

			return pair
		}
	}

	return pairs[len(pairs)-1]
}

// GetLastPointToConnect determines which point from a route has to be
// connected yet.
func GetLastPointToConnect(pairs []Pair) Point {
	lastPair := pairs[len(pairs)-1]
	secondToLastPair := pairs[len(pairs)-2]

	if lastPair.pointA.index == secondToLastPair.pointA.index ||
		lastPair.pointA.index == secondToLastPair.pointB.index {
		return lastPair.pointB
	}

	return lastPair.pointA
}

// ConnectPointToClosest connects pointIn to a given point a pair of
// points with a short distance.
func ConnectPointToClosest(pointIn Point, pairs []Pair, routeOfPairs []Pair) []Pair {

	if len(pairs) > 0 {

		pairsWithPoint := GetPairsContainingPoint(pointIn, pairs)

		if len(pairsWithPoint) > 0 {

			// Determine the pair of shortest distance
			focusPair := PairWithShortestDistance(pairsWithPoint)

			// Determine the point needing a subsequent connection
			pointToConnect := GetOtherPointInPair(pointIn, focusPair)

			distance := focusPair.distance

			// Ensure route without duplicate trips to same point by reducing pair set
			pairsRefined := RemovePairsWithPoint(pointIn, pairs)

			// Add the shortest distance pair to the route
			pairToAdd := Pair{index: len(routeOfPairs),
				pointA: pointIn, pointB: pointToConnect, distance: distance}

			routeOfPairs = append(routeOfPairs, pairToAdd)

			// Continue connecting points until exhausted
			routeOfPairs = ConnectPointToClosest(pointToConnect, pairsRefined, routeOfPairs)
		}
	}
	return routeOfPairs
}

// RemovePairsWithPoint removes all pairs containing a given point
// from a set of pairs.
func RemovePairsWithPoint(point Point, pairs []Pair) []Pair {

	pairsWithoutPoint := make([]Pair, len(pairs))
	copy(pairsWithoutPoint, pairs)

	for _, pair := range pairs {

		if pair.pointA.index == point.index || pair.pointB.index == point.index {
			pairIndexToDelete := GetIndexOfPair(pair, pairsWithoutPoint)

			if pairIndexToDelete >= 0 {
				pairsWithoutPoint = append(pairsWithoutPoint[:pairIndexToDelete], pairsWithoutPoint[pairIndexToDelete+1:]...)
			}
		}
	}
	return pairsWithoutPoint
}

// GetPairsContainingPoint returns all pairs containing a given point.
func GetPairsContainingPoint(point Point, pairs []Pair) []Pair {
	var pairsWithPoint []Pair
	for _, pair := range pairs {
		if pair.pointA.index == point.index || pair.pointB.index == point.index {
			pairsWithPoint = append(
				pairsWithPoint,
				Pair{index: pair.index, pointA: pair.pointA, pointB: pair.pointB,
					distance: pair.distance})
		}
	}

	return pairsWithPoint
}

// GetOtherPointInPair returns the paired point of a given point.
func GetOtherPointInPair(pointIn Point, pair Pair) Point {
	if pair.pointA == pointIn {
		return pair.pointB
	}
	return pair.pointA
}

// GetIndexOfPair computes the index assigned to a pair from a set of pairs.
func GetIndexOfPair(pairIn Pair, pairs []Pair) int {
	i := 0
	for _, pair := range pairs {
		if pairIn.index == pair.index {
			return i
		}
		i++
	}
	return -1
}

// PairWithShortestDistance returns the pair with the shortest
// distance.
func PairWithShortestDistance(pairs []Pair) Pair {
	var distance float64
	var shortestDistancePair Pair
	for _, pair := range pairs {
		if distance == 0 {
			shortestDistancePair = pair
			distance = pair.distance
		}

		if pair.distance < distance {
			shortestDistancePair = pair
			distance = pair.distance
		}
	}
	return shortestDistancePair
}
