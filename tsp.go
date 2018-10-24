package main

import (
	"os"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

//A point in 2D space, aligned with an x and y intersection
type Point struct {
	index int
	x int
	y int
}

//A pair of points with a distance between them
type Pair struct {
	index int
	pointA Point
	pointB Point
	distance float64
}

//Entry function to generate points and solve the TSP problem
func main() {

	startTime := time.Now()

	maxPoints, e := strconv.Atoi(os.Args[1])
	minPoints, e := strconv.Atoi(os.Args[2])

	if e != nil {
        fmt.Println(e)
    }

	//Seed the random generator with the current time
	rand.Seed(startTime.UnixNano())

	//Randomly determine a number of points to generate
	var numberOfPoints = rand.Intn(maxPoints-minPoints) + minPoints
	fmt.Println("random = " + strconv.Itoa(numberOfPoints))

	//Build up a set of a number of random points
	points := BuildPoints(numberOfPoints)
	PrintPoints(points)

	//For each point, pair it with another point
	pairs := PairPoints(points)
	PrintPairs(pairs)

	//Calculate a route between points with shortest distance between an two points first
	fmt.Println("Calculating shortest path now")
	route := RouteShortestFirst(pairs)
	PrintPairs(route)

	//Print to console a description of the route determined
	fmt.Println("Total Distance = " + strconv.FormatFloat(GetTotalLength(route), 'f', 2, 64))
	fmt.Println("Avg Distance Of All = " + strconv.FormatFloat(GetAvgLength(pairs), 'f', 2, 64))
	fmt.Println("Avg Distance Of Route = " + strconv.FormatFloat(GetAvgLength(route), 'f', 2, 64))

	elapsedTime := time.Since(startTime)
	fmt.Printf("Calculation took %s\n", elapsedTime)
	fmt.Printf("Calculation took %s per point\n", 
		(elapsedTime / time.Duration(numberOfPoints)))
}

//Get the average length of distance between a set of paired points
func GetAvgLength(pairs []Pair) float64 {
	
	return (GetTotalLength(pairs) / float64(len(pairs)))
}

//Get the aggregated length of distance between a set of paired points
func GetTotalLength(pairs []Pair) float64 {
	total := 0.0
	for _, pair := range pairs {

		total = (total + pair.distance)
	}
	return total
}

//Calculate the distance between two points
func (p Point) Distance(p2 Point) float64 {
	first := math.Pow(float64(p2.x-p.x), 2)
	second := math.Pow(float64(p2.y-p.y), 2)
	return math.Sqrt(first + second)
}

//Use the fmt import to itteratively print description of each point in a set
func PrintPoints(points []Point) {
	fmt.Println("\nPoints:")
	for _, point := range points {
		fmt.Println(DescribePoint(point))
	}
}

//Retrun a concatenated string describing a point in x, y space
func DescribePoint(point Point) string {
	pointDescription := 
		"[" + strconv.Itoa(point.index) + "] " +
		strconv.Itoa(point.x) + "," + strconv.Itoa(point.y)

	return pointDescription
}

//Use the fmt import to itteratively print a description of each pair in a set
func PrintPairs(pairs []Pair) {
	fmt.Println("\nPairs of length (" + strconv.Itoa(len(pairs)) + ")")
	for _, pair := range pairs {
		fmt.Println(
			"[",
			strconv.Itoa(pair.index),
			"] ",
			DescribePoint(pair.pointA),
			",",
			DescribePoint(pair.pointB),
			",",
			strconv.FormatFloat(pair.distance, 'f', 2, 64))
	}
}

//Return a concatenated string describing a pair of points in x, y space
func DescribePair(pair Pair) string {
	pairDescription := 
		"[" + strconv.Itoa(pair.index) + "] " + 
		DescribePoint(pair.pointA) + "," + DescribePoint(pair.pointB) + "," + 
		strconv.FormatFloat(pair.distance, 'f', 2, 64)

	return pairDescription
}

//Return a number of randomly generated points
func BuildPoints(numberOfPoints int) []Point {
	points := make([]Point, numberOfPoints)
	for i:=0; i<numberOfPoints; i++ {
		var randomX = rand.Intn(100)
		var randomY = rand.Intn(100)
		var index = i
		points[i] = Point{index:index, x:randomX, y:randomY}
	}
	return points
}

//Create a pair between each given point in a set
func PairPoints(points []Point) []Pair {
	var count = 0
	var skip = false
	pairs := make([]Pair, 0)
	for _, pointA := range points {
		for _, pointB := range points {

			//Avoid duplicates
			for _, pair := range pairs {
				if(pair.pointA.index==pointA.index && pair.pointB.index==pointB.index) { 
					skip = true 
				}
				if(pair.pointA.index==pointB.index && pair.pointB.index==pointA.index) { 
					skip = true 
				}
			}

			//Add new pair to set of pairs among points
			if(pointA.index!=pointB.index && !skip) {
				count++
				var distance = pointA.Distance(pointB)
				pairs = append(
					pairs, 
					Pair{index:count, pointA:pointA, pointB:pointB, distance:distance})
				skip = false
			} else {
				skip = false
			}
		}
	}
	return pairs
}

//Begin routine to define a path between all points beginning with the shortest path
func RouteShortestFirst(pairs []Pair) []Pair {

	//Begin with a point inside of a pair of the shortest distance
	startingPair := PairWithShortestDistance(pairs)
	startingPoint := startingPair.pointA
	fmt.Println("Starting Point = " + DescribePoint(startingPoint))

	//The complimenting point from the starting pair must be identified and connected next
	pointToConnect := startingPair.pointB

	//Construct an empty set of pairs to define as the shortest route across all points
	routeOfPairs := make([]Pair, 0)
	startingPairDistance := startingPair.distance
	
	//Add the first, shortest pair to the route
	routeOfPairs = append(
		routeOfPairs, 
		Pair{index:0, pointA:startingPoint, pointB:pointToConnect, 
			distance:startingPairDistance})

	//Redefine the remaining pairs of points to connect yet to be without already paired point
	pairsRefined := RemovePairsWithPoint(startingPoint, pairs)
	
	//Recursively pair shortest points until unmapped points no longer exist
	routeOfPairs = ConnectPointToClosest(pointToConnect, pairsRefined, routeOfPairs)

	//The last point must be joined with the first point of the route to complete the route
	lastPoint := GetLastPointToConnect(routeOfPairs)
	lastPair := GetPairFromPairs(lastPoint, startingPoint, pairs)
	lastPair.index = len(routeOfPairs)
	routeOfPairs = append(routeOfPairs, lastPair)
	
	return routeOfPairs
}

//Return a pair from a set that includes two given points
func GetPairFromPairs(pointA Point, pointB Point, pairs []Pair) Pair {
	for _, pair := range pairs {
		
		if pair.pointA.index==pointA.index && pair.pointB.index==pointB.index {

			return pair
		}

		if pair.pointA.index==pointB.index && pair.pointB.index==pointA.index {

			return pair
		}
	}

	return pairs[len(pairs)-1]
}

//Determine which point from a route has to be connected yet
func GetLastPointToConnect(pairs []Pair) Point {
	lastPair := pairs[len(pairs)-1]
	secondToLastPair := pairs[len(pairs)-2]

	if lastPair.pointA.index==secondToLastPair.pointA.index || 
		lastPair.pointA.index==secondToLastPair.pointB.index {
			return lastPair.pointB
	}

	return lastPair.pointA
}

//Connect to a given point a pair of points with a short distance
func ConnectPointToClosest(pointIn Point, pairs []Pair, routeOfPairs []Pair) []Pair {

	if len(pairs)>0 {

		pairsWithPoint := GetPairsContainingPoint(pointIn, pairs)

		if len(pairsWithPoint)>0 {

			//Determine the pair of shortest distance
			focusPair := PairWithShortestDistance(pairsWithPoint)

			//Determine the point needing a subsequent connection
			pointToConnect := GetOtherPointInPair(pointIn, focusPair);

			distance := focusPair.distance

			//Ensure route without duplicate trips to same point by reducing pair set
			pairsRefined := RemovePairsWithPoint(pointIn, pairs)

			//Add the shortest distance pair to the route
			pairToAdd := Pair{index:len(routeOfPairs), 
				pointA:pointIn, pointB:pointToConnect, distance:distance}

			routeOfPairs = append(routeOfPairs, pairToAdd)
			
			//Continue connecting points until exhausted
			routeOfPairs = ConnectPointToClosest(pointToConnect, pairsRefined, routeOfPairs);
		}
	}
	return routeOfPairs
}

//Remove all pairs containing a given point from a set of pairs
func RemovePairsWithPoint(point Point, pairs []Pair) []Pair {
	
	pairsWithoutPoint := make([]Pair, len(pairs))
	copy(pairsWithoutPoint, pairs)

	for _, pair := range pairs {
		
		if pair.pointA.index==point.index || pair.pointB.index==point.index {
			pairIndexToDelete := GetIndexOfPair(pair, pairsWithoutPoint)
			
			if pairIndexToDelete >= 0 {
				pairsWithoutPoint = append(pairsWithoutPoint[:pairIndexToDelete], 
					pairsWithoutPoint[pairIndexToDelete+1:]...)
			}
		}
	}
	return pairsWithoutPoint
}

//Return all pairs containing a given point
func GetPairsContainingPoint(point Point, pairs []Pair) []Pair {
	var pairsWithPoint []Pair
	for _, pair := range pairs {
		if pair.pointA.index==point.index || pair.pointB.index==point.index {
			pairsWithPoint = append(
				pairsWithPoint, 
				Pair{index:pair.index, pointA:pair.pointA, pointB:pair.pointB, 
					distance:pair.distance})	
		}
	}
		
	return pairsWithPoint
}

//From a given pair, return the paired point of a given point
func GetOtherPointInPair(pointIn Point, pair Pair) Point {
	if pair.pointA==pointIn { return pair.pointB }
	return pair.pointA
}

//Get the index assigned to a pair from a set of pairs
func GetIndexOfPair(pairIn Pair, pairs []Pair) int {
	i := 0
	for _, pair := range pairs {
		if pairIn.index==pair.index {
			return i
		}
		i++
	}
	return -1
}

//Return the pair with the shortest distance
func PairWithShortestDistance(pairs []Pair) Pair {
	var distance float64
	var shortestDistancePair Pair
	for _, pair := range pairs {
		if(distance==0) {
			shortestDistancePair = pair
			distance = pair.distance
		}

		if(pair.distance<distance) {
			shortestDistancePair = pair
			distance = pair.distance
		}
	}
	return shortestDistancePair
}