// Copyright (c) 2024 Andrei Gill. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package peakdetect

import (
	"fmt"
	"github.com/mowshon/iterium"
	"math"
	"math/bits"
	"os"
	"reflect"
	"slices"
)

func expect(expected []int, result PrimaryPeaks[int]) {
	if !reflect.DeepEqual(expected, result.peaks) {
		fmt.Println(fmt.Sprintf("expected %v, got %v", expected, result))
		os.Exit(1)
	}
}

func TestCountBinary(bitWidth int) {
	// Counts in a binary space of the specified number of bits,
	// wherein the generated peaks are also all binary. This was
	// an early test, designed to thoroughly test the merge logic.
	maxCount := int(math.Pow(2, float64(bitWidth)))
	for i := 0; i < maxCount; i++ {
		testEach(i, bitWidth)
	}
}

func TestCountDecimal() {
	// Counts in this fixed radix system.
	//
	// |3|3|3|3|3|3|3|3|
	// |2|2|2|2|2|2|2|2|
	// |1|1|1|1|1|1|1|1|
	// |0|0|0|0|0|0|0|0|
	//
	// In this case there are 8 places, each of which
	// can be in only one of 4 possible values, leading
	// to the total number of possibilities being
	// 4^8=65536
	//
	// For example, one input may be something like:
	//
	// |0|1|2|3|4|3|2|1|   with the only peak located at index: 4
	//
	// or this input
	//
	// |0|3|3|2|1|1|2|1|   with peaks at 1, 2 and 6
	//

	//
	// If we run it with the alphabet being: {0, 1, 2, 3, 4, 5, 6, 7}, then given
	// 8 places as above, we would have a total of: 8^8=16777216 possibilities.
	//
	// I ran it, and it produced the correct output, but took about a minute or so
	// to run. Therefore, I left it at a much faster, but still representative test,
	// using only 4 unique symbols instead of 8.
	//
	// p := iterium.Product([]int{0, 1, 2, 3, 4, 5, 6, 7}, numberOfPlaces)
	//
	// When I ran it with the above, this was the output (trailing part of it):
	//
	// [7 7 7 7 7 7 7 1]
	// [0 1 2 3 4 5 6]
	// [7 7 7 7 7 7 7 2]
	// [0 1 2 3 4 5 6]
	// ...
	// [7 7 7 7 7 7 7 6]
	// [0 1 2 3 4 5 6]
	// [7 7 7 7 7 7 7 7]
	// []
	// Total: 16777216
	//

	maxNumberOfPlaces := 8
	for numberOfPlaces := 1; numberOfPlaces <= maxNumberOfPlaces; numberOfPlaces++ {
		p := iterium.Product([]int{0, 1, 2, 3}, numberOfPlaces)
		s, _ := p.Slice()
		for _, samples := range s {
			peaks := DetectPeaks(samples)
			if isValid[int](&peaks) {
				fmt.Println(peaks.samples)
				fmt.Println(peaks.peaks)
			} else {
				fmt.Println(" FAILURE ")
				fmt.Println(peaks.samples)
				fmt.Println(peaks.peaks)
				os.Exit(1)
			}
		}
		fmt.Println("Total:", p.Count())
	}
}

func testEach(value int, bitWidth int) {
	samples := make([]int, bitWidth)
	samplePeaks := []int{} // TODO: if I change this to 'nil slice declaration' it will fail these tests
	unsignedValue := uint(value)
	allBitsSet := bits.OnesCount(unsignedValue) == bitWidth

	for j := 0; j < bitWidth; j++ {
		if value&(1<<j) == 0 {
			samples[bitWidth-j-1] = 0
		} else {
			if !allBitsSet {
				samplePeaks = append(samplePeaks, bitWidth-j-1)
			}
			samples[bitWidth-j-1] = 1
		}
	}

	if !allBitsSet {
		slices.Reverse(samplePeaks)
	} else {
		samplePeaks = []int{}
	}

	peaks := DetectPeaks(samples)

	fmt.Println(value, " ", peaks.samples)
	fmt.Println(peaks.peaks)

	if false {
		fmt.Println(fmt.Sprintf("expected %v, got %v", samplePeaks, peaks.peaks))
	}

	expect(samplePeaks, peaks)
}

/*
Exhaustive test of all possible peaks in a cluster of 3 samples, wherein
each sample can only be one of 3 possible unique values.

We have a total of three samples, each of which can be an arbitrary integer.

Let us look at the problem by first limiting each sample to a set of
the following possible unique values: {0, 1, 2}

Because we have a total of 3 samples, and each of which can only take on
the above possible values, the resulting space of possibilities is 3^3=27.

In other words, there are 27 ways in which we can express 3 samples, such
that we can encode all their possible differences. For example, the first
sample being greater than the second or less than it, or equal to it. The
same of the second and third samples, and the third and first.

Suppose for example, that our alphabet was limited only to the symbols: {0, 1}

In this case, we would not be able to express all possible differences in
the cluster of 3 samples, because we could never express a full ramp.
In order to be able to express a full ramp, we must have as many unique
symbols in our alphabet, as there are samples.

Hence, since we have 3 samples, our alphabet needs 3 unique symbols: {0, 1, 2}

We can now express a full ramp:

	                     |2|
		           |1|1|
	   |0|1|2|       |0|0|0|

Next, we need to enumerate all possible differences between every two
samples, out of the 3 total. This is the standard 'N choose K' problem,
wherein we have to choose 2 values out of 3, without regard for order
of the 2 chosen items.

Out of a set of 3 samples, 'a', 'b' and 'c', we have the following 3 possibilities:

				a b c       Let us encode our difference relations as such:
				-----
				a b         The choices 'ab', 'bc' and 'cd', allow
				  b c       for the following possibilities:
				a   c
	                            (A > b)  Ab       (B > c)  Bc      (C > a)  Ca
			                    (a < B)  aB       (b < C)  bC      (c < A)  aC
			                    (a == b) ab       (b == c) bc      (c == a) ac

We can now see what the logic would look like, when trying to express the above
state space of 27 possible differences between two 2 samples, out of 3 total:

if a < b { // aB

	  if b < c { // bC
	    if c < a { // cA
	      <leaf>
	    } else if c > a { // Ca
	      <leaf>
	    } else { // c == a, ca
	      <leaf>
	    }
	  } else if b > c { Bc
	      <branch>
	  } else { // b == c, bc
	      <branch>
	  }
	} else if a > b { // Ab

	  <branch>
	} else { // a == b, ab

	  <branch>
	}

Below, we feed all 27 possible clusters of 3 samples, wherein the alphabet
is limited to 3 unique symbols, leading to 3^3=27 total possibilities.
*/
func TestPeakDetectTriple() {
	expect([]int{}, peakDetectTriple(1, 1, 1))  /*      (0)  |1|1|1|            [2]            [3]    */
	expect([]int{2}, peakDetectTriple(1, 1, 2)) /*                     (1)  |1|1|1|            |2|    */
	expect([]int{2}, peakDetectTriple(1, 1, 3)) /*                                    (2)  |1|1|1|    */

	expect([]int{1}, peakDetectTriple(1, 2, 1))    /*   (3)    [2]            [2|2]            [3]    */
	expect([]int{1, 2}, peakDetectTriple(1, 2, 2)) /*        |1|1|1|   (4)  |1|1|1|          |2|2|    */
	expect([]int{2}, peakDetectTriple(1, 2, 3))    /*                                 (5)  |1|1|1|    */

	expect([]int{1}, peakDetectTriple(1, 3, 1))    /*   (6)    [3]            [3]            [3|3]    */
	expect([]int{1}, peakDetectTriple(1, 3, 2))    /*          |2|     (7)    |2|2|          |2|2|    */
	expect([]int{1, 2}, peakDetectTriple(1, 3, 3)) /*        |1|1|1|        |1|1|1|   (8)  |1|1|1|    */

	expect([]int{0}, peakDetectTriple(2, 1, 1))    /*   (9)  [2]                               [3]   */
	expect([]int{0, 2}, peakDetectTriple(2, 1, 2)) /*        |1|1|1|   (10) [2] [2]        [2] |2|   */
	expect([]int{0, 2}, peakDetectTriple(2, 1, 3)) /*                       |1|1|1|   (11) |1|1|1|   */

	expect([]int{0, 1}, peakDetectTriple(2, 2, 1)) /*   (12) [2|2]                             [3]   */
	expect([]int{}, peakDetectTriple(2, 2, 2))     /*        |1|1|1|   (13) |2|2|2|        |2|2|2|   */
	expect([]int{2}, peakDetectTriple(2, 2, 3))    /*                       |1|1|1|   (14) |1|1|1|   */

	expect([]int{1}, peakDetectTriple(2, 3, 1))    /*   (15)   [3]            [3]            [3|3]   */
	expect([]int{1}, peakDetectTriple(2, 3, 2))    /*        |2|2|     (16) |2|2|2|        |2|2|2|   */
	expect([]int{1, 2}, peakDetectTriple(2, 3, 3)) /*        |1|1|1|        |1|1|1|   (17) |1|1|1|   */

	expect([]int{0}, peakDetectTriple(3, 1, 1))    /*   (18) [3]            [3]            [3] [3]   */
	expect([]int{0, 2}, peakDetectTriple(3, 1, 2)) /*        |2|       (19) |2| [2]        |2| |2|   */
	expect([]int{0, 2}, peakDetectTriple(3, 1, 3)) /*        |1|1|1|        |1|1|1|   (20) |1|1|1|   */

	expect([]int{0}, peakDetectTriple(3, 2, 1))    /*   (21) [3]            [3]            [3] [3]   */
	expect([]int{0}, peakDetectTriple(3, 2, 2))    /*        |2|2|     (22) |2|2|2|        |2|2|2|   */
	expect([]int{0, 2}, peakDetectTriple(3, 2, 3)) /*        |1|1|1|        |1|1|1|   (23) |1|1|1|   */

	expect([]int{0, 1}, peakDetectTriple(3, 3, 1)) /*   (24) [3|3]          [3|3]          |3|3|3|   */
	expect([]int{0, 1}, peakDetectTriple(3, 3, 2)) /*        |2|2|     (25) |2|2|2|        |2|2|2|   */
	expect([]int{}, peakDetectTriple(3, 3, 3))     /*        |1|1|1|        |1|1|1|   (26) |1|1|1|   */
}

func TestMergeOfTriples() {
	// both are peaks, samples are equal
	expect([]int{2, 3}, merge(peakDetectThreeSamples([]int{1, 2, 3}), peakDetectThreeSamples([]int{3, 2, 1})))

	// both are peaks, left sample is greater
	expect([]int{2}, merge(peakDetectThreeSamples([]int{1, 2, 3}), peakDetectThreeSamples([]int{2, 2, 1})))

	// both are peaks, right sample is greater
	expect([]int{3}, merge(peakDetectThreeSamples([]int{1, 2, 2}), peakDetectThreeSamples([]int{3, 2, 1})))

	// neither is a peak, samples are equal
	expect([]int{4}, merge(peakDetectThreeSamples([]int{1, 2, 2}), peakDetectThreeSamples([]int{2, 3, 1})))

	// neither is a peak, left sample is greater
	expect([]int{1, 4}, merge(peakDetectThreeSamples([]int{1, 3, 2}), peakDetectThreeSamples([]int{1, 3, 1})))

	// neither is a peak, right sample is greater
	expect([]int{1, 4}, merge(peakDetectThreeSamples([]int{1, 3, 1}), peakDetectThreeSamples([]int{2, 3, 1})))

	// left side is a peak, samples are equal
	expect([]int{4}, merge(peakDetectThreeSamples([]int{1, 2, 3}), peakDetectThreeSamples([]int{3, 4, 1})))

	// left side is a peak, left sample is greater
	expect([]int{2, 4}, merge(peakDetectThreeSamples([]int{1, 2, 4}), peakDetectThreeSamples([]int{3, 4, 1})))

	// left side is a peak, right sample is greater
	expect([]int{4}, merge(peakDetectThreeSamples([]int{1, 2, 3}), peakDetectThreeSamples([]int{4, 5, 1})))

	// right side is a peak, samples are equal
	expect([]int{1}, merge(peakDetectThreeSamples([]int{1, 4, 3}), peakDetectThreeSamples([]int{3, 2, 1})))

	// right side is a peak, left sample is greater
	expect([]int{1}, merge(peakDetectThreeSamples([]int{1, 5, 4}), peakDetectThreeSamples([]int{3, 2, 1})))

	// right side is a peak, right sample is greater
	expect([]int{1, 3}, merge(peakDetectThreeSamples([]int{1, 4, 3}), peakDetectThreeSamples([]int{4, 2, 1})))
}
