// Copyright (c) 2024 Andrei Gill. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package peakdetect

import (
	"errors"
	"fmt"
	"os"
)

func peakDetectSample[T Number](a T) PrimaryPeaks[T] {
	return CreatePeaksWith[T]([]T{a}, []int{})
}

func peakDetectPair[T Number](a, b T) PrimaryPeaks[T] {
	if err, result := peakDetectPair0[T](a, b); err == nil {
		return result
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	return PrimaryPeaks[T]{}
}

func peakDetectTriple[T Number](a, b, c T) PrimaryPeaks[T] {
	if err, result := peakDetectTriple0[T](a, b, c); err == nil {
		return result
	} else {
		fmt.Println(err)
		os.Exit(1)
	}
	return PrimaryPeaks[T]{}
}

func peakDetectThreeSamples[T Number](samples []T) PrimaryPeaks[T] {
	return peakDetectTriple[T](samples[0], samples[1], samples[2])
}

func DetectPeaks[T Number](samples []T) PrimaryPeaks[T] {
	at := 0
	stride := 3
	left := PrimaryPeaks[T]{}

	for i := 0; at+stride <= len(samples); i++ {
		if at > 0 {
			right := peakDetectTriple[T](samples[at], samples[at+1], samples[at+2])
			left = merge[T](left, right)
		} else {
			left = peakDetectTriple[T](samples[i], samples[i+1], samples[i+2])
		}
		at += stride
	}

	if len(samples)-at == 1 {
		left = merge[T](left, peakDetectSample[T](samples[at]))
	} else if len(samples)-at == 2 {
		left = merge[T](left, peakDetectPair[T](samples[at], samples[at+1]))
	}

	return left
}

func DetectPeaksInPrimary[T Number](p PrimaryPeaks[T]) SecondaryPeaks[T] {
	at := 0
	stride := 3
	left := SecondaryPeaks[T]{}

	for i := 0; at+stride <= len(p.peaks); i++ {
		if at > 0 {
			a := p.peaks[at]
			b := p.peaks[at+1]
			c := p.peaks[at+2]
			originalPeaks := []int{a, b, c}
			rightTriple := peakDetectTriple[T](p.samples[a], p.samples[b], p.samples[c])
			right := CreateSecondaryPeaksWith[T](rightTriple, alignPrimaryPeaks[T](rightTriple, originalPeaks), originalPeaks)
			left = mergeSecondary[T](left, right)
		} else {
			a := p.peaks[i]
			b := p.peaks[i+1]
			c := p.peaks[i+2]
			originalPeaks := []int{a, b, c}
			leftTriple := peakDetectTriple[T](p.samples[a], p.samples[b], p.samples[c])
			left = CreateSecondaryPeaksWith[T](leftTriple, alignPrimaryPeaks[T](leftTriple, originalPeaks), originalPeaks)
		}
		at += stride
	}

	if len(p.peaks)-at == 1 {
		a := p.peaks[at]
		originalPeaks := []int{a}
		rightSample := peakDetectSample[T](p.samples[a])
		right := CreateSecondaryPeaksWith[T](rightSample, alignPrimaryPeaks[T](rightSample, originalPeaks), originalPeaks)
		left = mergeSecondary[T](left, right)
	} else if len(p.peaks)-at == 2 {
		a := p.peaks[at]
		b := p.peaks[at+1]
		originalPeaks := []int{a, b}
		rightPair := peakDetectPair[T](p.samples[a], p.samples[b])
		right := CreateSecondaryPeaksWith[T](rightPair, alignPrimaryPeaks[T](rightPair, originalPeaks), originalPeaks)
		left = mergeSecondary[T](left, right)
	}

	left.primarySamples = p.samples
	return left
}

func DetectPeaksInSecondary[T Number](p SecondaryPeaks[T]) SecondaryPeaks[T] {
	at := 0
	stride := 3
	left := SecondaryPeaks[T]{}

	for i := 0; at+stride <= len(p.peaks); i++ {
		if at > 0 {
			a := p.peaks[at]
			b := p.peaks[at+1]
			c := p.peaks[at+2]
			rightTriple := peakDetectTriple[T](p.samples[a], p.samples[b], p.samples[c])
			right := CreateSecondaryPeaksWith[T](rightTriple, alignPrimaryPeaks[T](rightTriple, p.primaryPeaks[at:at+stride]), p.originalPeaks[at:at+stride])
			left = mergeSecondary[T](left, right)
		} else {
			a := p.peaks[i]
			b := p.peaks[i+1]
			c := p.peaks[i+2]
			leftTriple := peakDetectTriple[T](p.samples[a], p.samples[b], p.samples[c])
			left = CreateSecondaryPeaksWith[T](leftTriple, alignPrimaryPeaks[T](leftTriple, p.primaryPeaks[:stride]), p.originalPeaks[:stride])
		}
		at += stride
	}

	if len(p.peaks)-at == 1 {
		a := p.peaks[at]
		rightSample := peakDetectSample[T](p.samples[a])
		right := CreateSecondaryPeaksWith[T](rightSample, alignPrimaryPeaks[T](rightSample, p.primaryPeaks[len(p.primaryPeaks)-1:]), p.originalPeaks[len(p.originalPeaks)-1:])
		left = mergeSecondary[T](left, right)
	} else if len(p.peaks)-at == 2 {
		a := p.peaks[at]
		b := p.peaks[at+1]
		rightPair := peakDetectPair[T](p.samples[a], p.samples[b])
		right := CreateSecondaryPeaksWith[T](rightPair, alignPrimaryPeaks[T](rightPair, p.primaryPeaks[len(p.primaryPeaks)-2:]), p.originalPeaks[len(p.originalPeaks)-2:])
		left = mergeSecondary[T](left, right)
	}

	left.primarySamples = p.primarySamples
	return left
}

func IteratePeakDetectToCompletion[T Number](samples []T) SecondaryPeaks[T] {
	primary := DetectPeaks[T](samples)
	secondary := DetectPeaksInPrimary[T](primary)
	for secondary.GetPeakCount() > 0 {
		secondary = DetectPeaksInSecondary[T](secondary)
	}
	return secondary
}

func IteratePeakDetect[T Number](iterations uint, samples []T) (SecondaryPeaks[T], bool) {
	if iterations == 0 {
		return SecondaryPeaks[T]{}, false
	}
	primary := DetectPeaks[T](samples)
	secondary := DetectPeaksInPrimary[T](primary)
	if iterations == 1 {
		return secondary, true
	} else {
		iterations--
	}
	for iterations > 1 && secondary.GetPeakCount() > 0 {
		secondary = DetectPeaksInSecondary[T](secondary)
		iterations--
	}
	return secondary, true
}

func alignPrimaryPeaks[T Number](p PrimaryPeaks[T], originalPeaks []int) []int {
	var primaryPeaks []int
	peaksDetected := len(p.peaks)
	if peaksDetected > 0 {
		primaryPeaks = make([]int, peaksDetected)
		for j, peakIndex := range p.peaks {
			primaryPeaks[j] = originalPeaks[peakIndex]
		}
	} else {
		primaryPeaks = []int{}
	}
	return primaryPeaks
}

func peakDetectPair0[T Number](a, b T) (error, PrimaryPeaks[T]) {
	if a > b {
		return nil, CreatePeaksWith[T]([]T{a, b}, []int{0})
	} else if a < b {
		return nil, CreatePeaksWith[T]([]T{a, b}, []int{1})
	} else /* a == b */ {
		return nil, CreatePeaksWith[T]([]T{a, b}, []int{})
	}
}

func peakDetectTriple0[T Number](a, b, c T) (error, PrimaryPeaks[T]) {
	if a < b {
		if b < c {
			if c < a {
				// (0) impossible state
				//
				// aB   bC   cA
				//
				//   |2|
				// |1|1|?|
				//
				// In this case 'c' is less than 'a', and 'a' is less than
				// 'b', therefore, transitively, 'c' is also less than 'b'.
				// However, in this case, it is also said that 'b' is less
				// than 'c',  which is a contradiction.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else if c > a {
				// (1) Peak at: {2}
				//
				// aB   bC   Ca
				//
				//     |3|
				//   |2|2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{2})
			} else /* c == a */ {
				// (2) impossible state
				//
				// aB   bC   ca
				//
				//   |2|
				// |1|1|?|
				//
				// In this case, 'a' is less than 'b', and 'b' is less than 'c'.
				// However, since 'a' is also equal to 'c', therefore 'b' is less
				// than 'a'. But in this case, 'b' is said to be greater than 'a',
				// which is a contradiction.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			}
		} else if b > c {
			if c < a {
				// (3) Peak at: {1}
				//
				// aB   Bc   cA
				//
				//   |3|
				// |2|2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{1})
			} else if c > a {
				// (4) Peak at: {1}
				//
				// aB   Bc   Ca
				//
				//   |3|
				//   |2|2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{1})
			} else /* c == a */ {
				// (5) Peak at: {1}
				//
				// aB   Bc   ca
				//
				//                 |3|
				//   |2|           |2|
				// |1|1|1|       |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{1})
			}
		} else /* b == c */ {
			if c < a {
				// (6) impossible state
				//
				// aB   bc   cA
				//
				// In this case, 'a' is less than 'b', and 'b' is equal to 'c'.
				// Therefore, this implies that 'a' must also be less than 'c',
				// since, well, 'b' and 'c' are equal, by definition.
				//
				// However, in this case, 'a' is defined to be less than 'b',
				// and therefore cannot equal to 'b', by definition.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else if c > a {
				// (7) Peaks at: {1, 2}
				//
				// aB   bc   Ca
				//
				//   |2|2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{1, 2})
			} else /* c == a */ {
				// (8) impossible state
				//
				// aB   bc   ca
				//
				// In this case, 'b' is equal to 'c', and 'c' is equal to 'a',
				// and therefore, since 'b' and 'c' are equal, and 'c' and 'a'
				// are also equal, then 'a' and 'b' must also be equal.
				// However, in this case, 'a' must be less than 'b'.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			}
		}
	} else if a > b {
		if b < c {
			if c < a {
				// (9) Peak at: {0, 2}
				//
				// Ab   bC   cA
				//
				// |3|
				// |2| |2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{0, 2})
			} else if c > a {
				// (10) Peak at: {0, 2}
				//
				// Ab   bC   Ca
				//
				//     |3|
				// |2| |2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{0, 2})
			} else /* c == a */ {
				// (11) Peaks at: {0, 2}
				//
				// Ab   bC   ca
				//
				//               |3| |3|
				// |2| |2|       |2| |2|
				// |1|1|1|       |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{0, 2})
			}
		} else if b > c {
			if c < a {
				// (12) Peak at: {0}
				//
				// Ab   Bc   cA
				//
				// |3|
				// |2|2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{0})
			} else if c > a {
				// (13) impossible state
				//
				// Ab   Bc   Ca
				//
				// |3|
				// |2|2|
				// |1|1|?|
				//
				// In this case 'b' is bigger than 'c'. For its part,
				// 'c' is supposed to be bigger than 'a'.  But 'a' is
				// bigger than 'b', and 'b' is bigger than 'c'. Therefore,
				// 'a' must also be bigger than 'c', since it is bigger
				// than 'b', which is bigger than 'c'.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else /* c == a */ {
				// (14) impossible state
				//
				// Ab   Bc   ca
				//
				// |2| |2|
				// |1|?|1|
				//
				// In this state, 'c' is equal to 'a', but 'a' is said to be
				// bigger than 'b'. However, it is also said that 'b' is bigger
				// than 'c'. It cannot be so, because 'c' and 'a' are equal
				// by definition, and since 'a' is bigger than 'b', then so
				// must 'c' be also, by definition, be bigger than 'b'.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			}
		} else /* b == c */ {
			if c < a {
				// (15) Peak at: {0}
				//
				// Ab   bc   cA
				//
				// |3|
				// |2|2|2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{0})
			} else if c > a {
				// (16) impossible state
				//
				// Ab   bc   Ca
				//
				// |3|
				// |2|2|?|
				// |1|1|?|
				//
				// Because 'b' is equal to 'c', and 'a' is bigger than 'b',
				// therefore, 'a' must also be bigger than 'c'. However, in
				// this state, 'c' is said to be bigger than 'a', which is a
				// contradiction.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else /* c == a */ {
				// (17) impossible state
				//
				// Ab   bc   ca
				//
				// |3| |?|
				// |2|2|?|
				// |1|1|?|
				//
				// In this state, 'b' is said to be equal to 'c', and 'c'
				// is said to be equal to 'a'. However, in this state 'a'
				// is also said to be bigger than 'b', which is a contradiction,
				// because 'b' is equal to 'c', which is equal to 'a', and therefore
				// cannot be bigger than neither 'b' nor 'c', since again, the two
				// are equal.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			}
		}
	} else /* a == b */ {
		if b < c {
			if c < a {
				// (18) impossible state
				//
				// ab   bC   cA
				//
				//     |?|
				// |2|2|?|
				// |1|1|?|
				//
				// In this case, 'a' is equal to 'b', and 'b' is less than 'c'.
				// In this case also, 'c' is said to be less than 'a', which is a contradiction,
				// because 'a' is said to also be less than 'c', by way of being equal to 'b',
				// which is said to be less than 'c', by definition.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else if c > a {
				// (19) Peak at: {2}
				//
				// ab   bC   Ca
				//
				//                   |3|           |3|
				//     |2|           |2|       |2|2|2|
				// |1|1|1|       |1|1|1|       |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{2})
			} else /* c == a */ {
				// (20) impossible state
				//
				// ab   bC   ca
				//
				//     |?|
				// |1|1|?|
				//
				// In this case, 'a' is equal to 'b', and 'b' is less than 'c'.
				// However, 'c' is also said to be equal to 'a', and therefore since
				// 'a' is equal to be', 'c' must then also be equal to 'b'. However,
				// this case states that 'b' is less than 'c', which is a contradiction.
				// Therefore, this state is impossible.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			}
		} else if b > c {
			if c < a {
				// (21) Peaks at: {0, 1}
				//
				// ab   Bc   cA
				//
				// |2|2|
				// |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{0, 1})
			} else if c > a {
				// (22) impossible state
				//
				// ab   Bc   Ca
				//
				//     |?|
				// |2|2|?|
				// |1|1|?|
				//
				// In this case, 'a' is equal to 'b', and 'b' is greater than 'c'.
				// However, in this case 'c' is also said to be greater than 'a',
				// but since 'a' and 'b' are equal, then this says that 'c' is also
				// greater than 'b', which is a contradiction.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else /* c == a */ {
				// (23) impossible state
				//
				// ab   Bc   ca
				//
				// |2|2|
				// |1|1|?|
				//
				// In this case, 'a' is equal to 'b' and 'c' is equal to 'a'.
				// Therefore, transitively, 'b' is equal to 'c'. However, in this
				// case, 'b' is said to be greater than 'c', which is a contradiction.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			}
		} else /* b == c */ {
			if c < a {
				// (24) impossible state
				//
				// ab   bc   cA
				//
				// |1|1|?|
				//
				// In this case, 'a' is equal to 'b', and 'b' is equal to 'c'.
				// However, 'c' is also said to be less than 'a', which is a contradiction.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else if c > a {
				// (25) impossible state
				//
				// ab   bc   Ca
				//
				// |1|1|?|
				//
				// In this case, 'a' is equal to 'b', and 'b' is equal to 'c'.
				// However, 'c' is also said to be greater than 'a', which is a contradiction.
				// Therefore, this is an impossible state.
				return errors.New("impossible state"), PrimaryPeaks[T]{}
			} else /* c == a */ {
				// (26) Peaks: {}
				//
				// ab   bc   ca
				//
				//                             |3|3|3|
				//               |2|2|2|       |2|2|2|
				// |1|1|1|       |1|1|1|       |1|1|1|
				return nil, CreatePeaks[T](a, b, c, []int{})
			}
		}
	}
}

// We continue through all equal samples, terminating only upon encountering a sample
// that is either greater, or less. We return 'false' if we find a greater sample, we
// return 'true' if we find a smaller sample, to the right of the sample at index 'at'.
func peakExtendsRight[T Number](at int, samples []T, initialSampleIsAPeak bool) bool {
	for i := at; i < len(samples); i++ {
		if i+1 < len(samples)-1 {
			if samples[i] > samples[i+1] {
				return true
			} else if samples[i] < samples[i+1] {
				return false
			}
		}
	}
	// If all samples to the right are equal, and we ran out of samples,
	// whether the peak extends to the right will depend on if the first
	// sample was a peak or not.
	return initialSampleIsAPeak
}

// We continue through all equal samples, terminating only upon encountering a sample
// that is either greater, or less. We return 'false' if we find a greater sample, we
// return 'true' if we find a smaller sample, to the left of the sample at index 'at'.
func peakExtendsLeft[T Number](at int, samples []T, initialSampleIsAPeak bool) bool {
	for i := at; i >= 0; i-- {
		if i-1 >= 0 {
			if samples[i] > samples[i-1] {
				return true
			} else if samples[i] < samples[i-1] {
				return false
			}
		}
	}
	// If all samples to the left are equal, and we ran out of samples,
	// whether the peak extends to the left will depend on if the first
	// sample was a peak or not.
	return initialSampleIsAPeak
}
