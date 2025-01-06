// Copyright (c) 2024 Andrei Gill. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package peakdetect

import "math/big"

// Verifies that all peak samples are actually peaks,
// and that all the remaining samples, are not peaks,
func isValid[T Number](p Peaks[T]) bool {
	var peakCount int
	var seenPeaks big.Int
	peaks := p.GetPeaks()
	samples := p.GetSamples()
	for _, peakIndex := range peaks {
		if isPeak(peakIndex, samples, true) {
			seenPeaks.SetBit(&seenPeaks, peakIndex, 1)
			peakCount++
			continue
		}
		return false
	}
	var nonPeakCount int
	for i, _ := range samples {
		if seenPeaks.Bit(i) == 0 && !isPeak(i, samples, false) {
			nonPeakCount++
		}
	}
	return peakCount+nonPeakCount == len(samples)
}

// Check whether the sample located at the specified index is a peak.
//
// It is considered a peak, if it is greater than either one, or both
// of its near neighbors. If one or both of its near neighbors are the
// same value as itself, we must extend to either the left or right,
// looking for a sample that is a different value.
func isPeak[T Number](at int, samples []T, initialSampleIsAPeak bool) bool {
	if at < 0 || at >= len(samples) {
		return false
	}
	if len(samples) == 0 {
		return false
	} else if len(samples) == 1 {
		return false
	} else /* len(samples) >= 2 */ {
		if at == 0 {
			return samples[at] > samples[at+1] || peakExtendsRight[T](at, samples, initialSampleIsAPeak)
		} else if at == len(samples)-1 {
			return samples[at-1] < samples[at] || peakExtendsLeft[T](at, samples, initialSampleIsAPeak)
		} else /* at > 0 && < len(p.samples)-1*/ {
			greaterThanRight := samples[at] > samples[at+1]
			smallerThanLeft := samples[at-1] < samples[at]
			return greaterThanRight && smallerThanLeft ||
				greaterThanRight && peakExtendsLeft[T](at, samples, initialSampleIsAPeak) ||
				smallerThanLeft && peakExtendsRight[T](at, samples, initialSampleIsAPeak) ||
				peakExtendsLeft[T](at, samples, initialSampleIsAPeak) &&
					peakExtendsRight[T](at, samples, initialSampleIsAPeak)
		}
	}
}
