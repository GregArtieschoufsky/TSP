package main

import (
	"os"
	"strconv"
	"fmt"
	"math"
	"math/rand"
	"time"
	"os/exec"
)

// Point is a 2D point aligned with an x and y intersection.
type Point struct {
	index int
	x     int
	y     int
}

// String returns a string representation of the Point.
func (point Point) String() string {
	return fmt.Sprintf("[%d] %d,%d", point.index, point.x, point.y)
}

// Pair is two points with a distance between them.
type Pair struct {
	index    int
	pointA   Point
	pointB   Point
	distance float64
}

// String returns a string representation of the Pair.
func (p Pair) String() string {
	return fmt.Sprintf("[%d] %s,%s,%f",
		p.index,
		p.pointA,
		p.pointB,
		p.distance)
}

func main() {

	startTime := time.Now()

	maxPoints, e := strconv.Atoi(os.Args[1])
	minPoints, e := strconv.Atoi(os.Args[2])

	if e != nil {
		fmt.Println(e)
	}

	// Seed the random generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Randomly determine a number of points to generate
	var numberOfPoints = rand.Intn(maxPoints-minPoints) + minPoints
	fmt.Printf("random = %d\n", numberOfPoints)

	pairs := SolveShortestByShortest(numberOfPoints)

	elapsedTime := time.Since(startTime)
	fmt.Printf("Calculation took %s\n", elapsedTime)
	fmt.Printf("Calculation took %s per point\n", 
		(elapsedTime / time.Duration(numberOfPoints)))

	PrintToFile(pairs)

	c := exec.Command("/usr/bin/python", "./plot.py")
	if err := c.Run(); err != nil { 
		fmt.Println("Error: ", err)
	}
}

// CoordsAsString Prepare string describing route coordinates
func CoordsAsString(pairs []Pair) string {

	x := ""
	y := ""
	p := ""

	lastIndex := -1

	for _, pair := range pairs {

		if lastIndex == pair.pointA.index {
			x += strconv.Itoa(pair.pointB.x) + ","
			y += strconv.Itoa(pair.pointB.y) + ","
			p += strconv.Itoa(pair.pointB.index) + ","
			lastIndex = pair.pointB.index
		} else {
			x += strconv.Itoa(pair.pointA.x) + ","
			y += strconv.Itoa(pair.pointA.y) + ","
			p += strconv.Itoa(pair.pointA.index) + ","
			lastIndex = pair.pointA.index
		}
	}

	x += strconv.Itoa(pairs[0].pointA.x)
	y += strconv.Itoa(pairs[0].pointA.y)
	p += strconv.Itoa(pairs[0].pointA.index)

	s := x + "\n" + y + "\n" + p + "\n"
	return s
}

// PrintToFile Prints route to a file
func PrintToFile(pairs []Pair) {

	f, err := os.Create("route.txt")
    if err != nil {
        fmt.Println(err)
        return
	}

	s := CoordsAsString(pairs)
	
    l, err := f.WriteString(s)
    if err != nil {
        fmt.Println(err)
        f.Close()
        return
	}
	
    fmt.Println(l, "bytes written successfully")
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
}

// SolveShortestByShortest calculates and returns a route 
// for TSP problem using shortest path first
func SolveShortestByShortest(numberOfPoints int) []Pair {

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

	return route
}

// GetAvgLength computes the average length of distance between a set
// of paired points.
func GetAvgLength(pairs []Pair) float64 {
	return GetTotalLength(pairs) / float64(len(pairs))
}

// GetTotalLength computes the aggregated length of distance between a
// set of paired points.
func GetTotalLength(pairs []Pair) float64 {
	total := 0.0
	for _, pair := range pairs {
		total += pair.distance
	}
	return total
}

// Distance returns the distance between two points.
func (point Point) Distance(p2 Point) float64 {
	first := math.Pow(float64(p2.x-point.x), 2)
	second := math.Pow(float64(p2.y-point.y), 2)
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
		Pair{
			index:    0,
			pointA:   startingPoint,
			pointB:   pointToConnect,
			distance: startingPairDistance,
		})

	// Redefine the remaining pairs of points to connect yet to be without already paired point
	pairsRefined := startingPoint.RemovePairsWithPoint(pairs)

	// Recursively pair shortest points until unmapped points no longer exist
	routeOfPairs = pointToConnect.ConnectPointToClosest(pairsRefined, routeOfPairs)

	// The last point must be joined with the first point of the route to complete the route
	lastPoint := GetLastPointToConnect(routeOfPairs)
	lastPair := GetPairFromPairs(lastPoint, startingPoint, pairs)
	lastPair.index = len(routeOfPairs)
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
func (point Point) ConnectPointToClosest(pairs []Pair, routeOfPairs []Pair) []Pair {

	if len(pairs) > 0 {

		pairsWithPoint := point.GetPairsContainingPoint(pairs)

		if len(pairsWithPoint) > 0 {

			// Determine the pair of shortest distance
			focusPair := PairWithShortestDistance(pairsWithPoint)

			// Determine the point needing a subsequent connection
			pointToConnect := point.GetOtherPoint(focusPair)

			distance := focusPair.distance

			// Ensure route without duplicate trips to same point by reducing pair set
			pairsRefined := point.RemovePairsWithPoint(pairs)

			// Add the shortest distance pair to the route
			pairToAdd := Pair{index: len(routeOfPairs),
				pointA: point, pointB: pointToConnect, distance: distance}

			routeOfPairs = append(routeOfPairs, pairToAdd)

			// Continue connecting points until exhausted
			routeOfPairs = pointToConnect.ConnectPointToClosest(pairsRefined, routeOfPairs)
		}
	}
	return routeOfPairs
}

// RemovePairsWithPoint removes all pairs containing a given point
// from a set of pairs.
func (point Point) RemovePairsWithPoint(pairs []Pair) []Pair {

	pairsWithoutPoint := make([]Pair, len(pairs))
	copy(pairsWithoutPoint, pairs)

	for _, pair := range pairs {
		if pair.pointA.index == point.index || pair.pointB.index == point.index {
			if pairIndexToDelete, have := pair.GetIndex(pairsWithoutPoint); have {
				pairsWithoutPoint = append(pairsWithoutPoint[:pairIndexToDelete], pairsWithoutPoint[pairIndexToDelete+1:]...)
			}
		}
	}
	return pairsWithoutPoint
}

// GetPairsContainingPoint returns all pairs containing a given point.
func (point Point) GetPairsContainingPoint(pairs []Pair) []Pair {
	var pairsWithPoint []Pair
	for _, pair := range pairs {
		if pair.pointA.index == point.index || pair.pointB.index == point.index {
			pairsWithPoint = append(
				pairsWithPoint,
				Pair{
					index:    pair.index,
					pointA:   pair.pointA,
					pointB:   pair.pointB,
					distance: pair.distance,
				})
		}
	}

	return pairsWithPoint
}

// GetOtherPoint returns the paired point of a given point.
func (point Point) GetOtherPoint(pair Pair) Point {
	if pair.pointA == point {
		return pair.pointB
	}
	return pair.pointA
}

// GetIndex computes the index assigned to a pair from a set of pairs.
//
// If the pair isn't found, returns false.
func (p Pair) GetIndex(pairs []Pair) (int, bool) {
	for i, pair := range pairs {
		if p.index == pair.index {
			return i, true
		}
	}
	return 0, false
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
