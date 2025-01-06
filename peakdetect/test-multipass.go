// Copyright (c) 2024 Andrei Gill. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package peakdetect

import (
	"fmt"
	"github.com/mowshon/iterium"
	"log"
	"math"
	"math/bits"
	"os"
	"slices"
)

func TestMultipass() {
	TestMultipass0()
	//TestCountBinary0(16)
}

func TestMultipass0() {
	//testOneCase()
	//iteratePeakDetect0()
	iteratePeakDetect1(10)

	//testRamp()

	//oneRamp(10)
	//oneRampLevel3(10)
}

func testOneCase() {
	//samples := []int{3, 3, 3, 3, 3, 3, 0, 1, 0, 1}
	//samples := []int{3, 3, 3, 3, 3, 3, 0, 1, 0, 2}

	//samples := []int{0, 1, 5, 2, 6, 8, 7, 1, 8}

	//samples := []int{0, 1, 1, 1, 1}
	//samples := []int{1, 0, 2, 2, 2}

	//samples := []int{2, 2, 2, 0, 1}

	// 7 peaks [1 4 8 12 15 17 19]
	samples := []int{1, 8, 3, 5, 7, 2, 3, 6, 9, 0, 1, 3, 12, 1, 4, 8, 2, 4, 1, 9}
	//
	// [1 8 3 5 7 2 3 6 9 0 1 3 12 1 4 8 2 4 1 9]
	//    ^     ^       ^        ^     ^       ^
	// [1 4 8 12 15 17 19]
	//
	//  1  4  8     12  15  17    19
	// [8][7][9]   [12] [8] [4]   [9]
	//  ^     ^      ^             ^
	//
	// After being merged, [9] that is to the left of [12],
	// is no longer a peak, and therefore the peak property
	// is removed from it.
	//
	// [8][7][9][12][8][4]   [9]
	//  ^         ^
	//
	//  0        3         6
	// [8][7][9][12][8][4][9]
	//  ^         ^        ^
	//
	//    1                     12            19
	// [1 8 3 5 7 2 3 6 9 0 1 3 12 1 4 8 2 4 1 9]
	//    ^                      ^             ^
	//

	//samples := []int{3, 3, 2, 2, 2, 3, 3, 2, 0, 1}
	// [3 3 2 2 2 3 3 2 0 1]
	//  ^ ^       ^ ^     ^
	// [0 1 5 6]
	//

	//peaks2 := DetectPeaks(samples)
	//peaks1 := DetectPeaksInPrimary(peaks2)
	//peaks := peaks1.InflateWithCount(len(samples), PrimaryValuesOnly[int](&peaks1))

	peaks3 := DetectPeaks(samples)
	fmt.Println("peaks3 (primary samples and peaks)")
	fmt.Println(peaks3.GetSamples())
	fmt.Println(peaks3.GetPeaks())

	peaks2 := DetectPeaksInPrimary(peaks3)
	fmt.Println("peaks2 (secondary samples and peaks)")
	fmt.Println(peaks2.GetSamples())
	fmt.Println(peaks2.GetPeaks())

	peaks1 := DetectPeaksInSecondary(peaks2)
	fmt.Println("peaks1 (secondary samples and peaks)")
	fmt.Println(peaks1.GetSamples())
	fmt.Println(peaks1.GetPeaks())

	//peaks := peaks3.InflateWithCount(len(samples) /*PrimaryValuesOnly[int](&peaks1)*/, &peaks3)
	//peaks := peaks2.InflateWithCount(len(samples), PrimaryValuesOnly[int](&peaks2))
	peaks := peaks1.InflateWithCount(len(samples), PrimaryValuesOnly[int](&peaks1))
	fmt.Println("peaks1 (inflated primary peaks)")
	fmt.Println(peaks)
	fmt.Println(peaks1.GetPrimaryPeaks())
}

func iteratePeakDetect0() {
	samples := []int{1, 8, 3, 5, 7, 2, 3, 6, 9, 0, 1, 3, 12, 1, 4, 8, 2, 4, 1, 9}

	iteratePeakDetect(samples)
}

func iteratePeakDetect1(numberOfPlaces int) {
	p := iterium.Product([]int{0, 1, 2, 3}, numberOfPlaces)
	s, _ := p.Slice()
	for _, samples := range s {
		iteratePeakDetect(samples)
	}
	fmt.Println("Total:", p.Count())
}

func iteratePeakDetect(samples []int) {
	primary := DetectPeaks(samples)
	fmt.Println("primary (primary samples and peaks)")
	fmt.Println(primary.GetSamples())
	fmt.Println(primary.GetPeaks())

	secondary := DetectPeaksInPrimary(primary)
	fmt.Println("secondary (secondary samples and peaks)")
	fmt.Println(secondary.GetSamples())
	fmt.Println(secondary.GetPeaks())

	for len(secondary.GetPeaks()) > 0 {
		secondary = DetectPeaksInSecondary(secondary)
		fmt.Println("secondary (loop) (secondary samples and peaks)")
		fmt.Println(secondary.GetSamples())
		fmt.Println(secondary.GetPeaks())
	}
}

func testRamp() {
	maxNumberOfPlaces := 12
	for numberOfPlaces := 1; numberOfPlaces <= maxNumberOfPlaces; numberOfPlaces++ {
		oneRamp(numberOfPlaces)
	}
}

func oneRamp(numberOfPlaces int) {
	p := iterium.Product([]int{0, 1, 2, 3}, numberOfPlaces)
	s, _ := p.Slice()
	for _, samples := range s {
		peaks1 := DetectPeaks(samples)
		peaks := DetectPeaksInPrimary(peaks1)

		if isValid[int](&peaks) {
			fmt.Println(peaks1.GetSamples())
			fmt.Println(peaks1.GetPeaks())
			fmt.Println(peaks.GetSamples())
			fmt.Println(peaks.GetPeaks())
			fmt.Println("------------")
		} else {
			fmt.Println(" FAILURE ")
			fmt.Println(peaks1.GetSamples())
			fmt.Println(peaks1.GetPeaks())
			fmt.Println(peaks.GetSamples())
			fmt.Println(peaks.GetPeaks())
			os.Exit(1)
		}
	}
	fmt.Println("Total:", p.Count())
}

func oneRampLevel3(numberOfPlaces int) {
	p := iterium.Product([]int{0, 1, 2, 3}, numberOfPlaces)
	s, _ := p.Slice()
	for _, samples := range s {
		peaks3 := DetectPeaks(samples)
		fmt.Println("peaks3 (primary)")
		fmt.Println(peaks3.GetSamples())
		fmt.Println(peaks3.GetPeaks())

		peaks2 := DetectPeaksInPrimary(peaks3)
		peaks1 := DetectPeaksInSecondary(peaks2)
		peaks := peaks1.InflateWithCount(len(samples), PrimaryValuesOnly[int](&peaks1))

		if isValid[int](&peaks1) {
			fmt.Println("peaks1 (tertiary)")
			fmt.Println(peaks1.GetSamples())
			fmt.Println(peaks1.GetPeaks())
			fmt.Println("peaks1 (inflated primary peaks)")
			fmt.Println(peaks)
			fmt.Println(peaks1.GetPrimaryPeaks())
			fmt.Println("------------")
		} else {
			fmt.Println(" FAILURE ")
			fmt.Println(peaks1.GetSamples())
			fmt.Println(peaks1.GetPeaks())
			os.Exit(1)
		}
	}
	fmt.Println("Total:", p.Count())
}

func TestCountBinary0(bitWidth int) {
	// Counts in a binary space of the specified number of bits,
	// wherein the generated peaks are also all binary. This was
	// an early test, designed to thoroughly test the merge logic.
	maxCount := int(math.Pow(2, float64(bitWidth)))
	for i := 0; i < maxCount; i++ {
		testEach0(i, bitWidth)
	}
}

func testEach0(value int, bitWidth int) {
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
	expect(samplePeaks, peaks)

	secondaryPeaks := DetectPeaksInPrimary(peaks)
	if len(secondaryPeaks.GetPeaks()) > 0 {
		log.Fatal("For binary peaks, there should be zero secondary peaks, as they should all merge into one" +
			"contiguous series of non-peak samples.")
	}
	//if len(secondaryPeaks.primaryPeaks) > 0 {
	//		log.Fatal("For binary peaks, there should be zero secondary peaks, and therefore zero offsets, since" +
	//			"peaks and offsets must always be a one-to-one mapping.")
	//	}

	fmt.Println("Primary peaks:")
	fmt.Println(value, " ", peaks.GetSamples())
	fmt.Println(peaks.GetPeaks())

	fmt.Println("Secondary peaks:")
	fmt.Println(secondaryPeaks.GetSamples())
	fmt.Println(secondaryPeaks.GetPeaks())
}
