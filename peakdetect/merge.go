// Copyright (c) 2024 Andrei Gill. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package peakdetect

import (
	"log"
)

func merge[T Number](left, right PrimaryPeaks[T]) PrimaryPeaks[T] {
	if left.samples != nil && right.samples == nil {
		return left
	} else if left.samples == nil && right.samples != nil {
		return right
	} else if left.samples == nil && right.samples == nil {
		log.Fatal("unexpected state")
	}
	leftSample := left.getLastSample()
	rightSample := right.getFirstSample()
	if left.isLastSamplePeak() {
		if right.isFirstSamplePeak() {
			if leftSample == rightSample {
				// The last sample of the 'left' peaks cluster, as well as
				// the first sample of the 'right' peaks cluster, are both
				// peaks, and both samples are equal to the same value.
				//
				// In this case, the peak does not propagate neither left
				// nor right. Therefore, we simply merge the two sample
				// arrays into one, and then merge the peaks, while fixing
				// up the offsets.
				return mergeSamples[T](left, right)
			} else {
				// In this case, the last sample of the 'left' peaks cluster,
				// as well as the first sample of the 'right' peaks cluster
				// are both peaks, but the sample values are not equals.
				//
				// This means that only the higher of the two values is a peak,
				// and we must recalculate, and propagate the peak property.
				if leftSample > rightSample {
					// The left sample is greater, therefore the peak is in the 'left'
					// peaks cluster. We want to go through the 'right' side, and remove
					// the peak property from all samples that are equal to the first
					// sample, and are contiguous with it.
					return mergeSamples[T](left, removeContiguousPeaksFromRightSide[T](right, rightSample))
				} else {
					// The right sample is greater, therefore the peak is in the
					// 'right' peaks cluster. We want to go through the 'left' side,
					// and remove the peak property from all samples that are equal
					// to the last sample, and are contiguous with it.
					return mergeSamples[T](removeContiguousPeaksFromLeftSide[T](left, leftSample), right)
				}
			}
		} else {
			// The last sample of the 'left' peaks cluster is a peak,
			// but the first sample of the 'right' peaks cluster is not.

			// First, we need to check if the neighboring samples on the
			// merge boundary are equal. If they're equal, but only one
			// of them is a peak, it means that we must propagate this
			// property through all contiguous samples. In this case,
			// the peak is on the 'left'. This means that we must propagate
			// the fact that it is not actually a peak, from 'right' to 'left'.
			//
			// This is because while the 'left' sample is a peak, the 'right'
			// sample is not. However, both are the same value. This indicates
			// that on the 'right' side, there is a condition which disqualifies
			// the first sample from being a peak. For example, the second
			// sample in the 'right' cluster may be a peak, or all the samples
			// may be equal to the same value, and therefore there is no peak
			// at all. In either case, because the first sample on the 'right'
			// is not a peak, but it is in fact equal to the last sample on
			// the 'left', we must ensure that all contiguous samples on the
			// 'left' are not marked as peaks.
			//
			// We iterate on the 'left' side from right to left, and remove
			// the peak property from all contiguous samples that are equal
			// to the value of the first sample on the 'right', or the last
			// sample on the 'left', since the two are equal in this case.
			//
			// Important point here, is that if there are peaks on the 'right'
			// side, then we remove the contiguous peaks from the 'left'. If,
			// however, there are no peaks on the 'right' side, e.g., due to
			// all samples being equal, then we must instead add all the peaks
			// on the 'right' side.
			if leftSample == rightSample {
				if len(right.peaks) == 0 {
					return mergeSamples[T](left, addContiguousPeaksToRightSide[T](right))
				} else {
					return mergeSamples[T](removeContiguousPeaksFromLeftSide[T](left, leftSample), right)
				}
			} else {
				// The last sample of the 'left' peaks cluster is a peak,
				// but the first sample of the 'right' peaks cluster is not.
				//
				// Additionally, the last sample of the 'left' peaks cluster,
				// is not equal to the first sample of the 'right' peaks cluster,
				// but is greater than it.
				//
				// This means that only the higher of the two values is a peak.
				//
				// In this case we needn't do any recalculation, because the last
				// sample of the 'left' peaks cluster, is already greater than the
				// first sample of the 'right' peaks cluster.
				//
				// Therefore, we simply merge the two sample arrays into one, and
				// then merge the peaks, while fixing up the offsets.
				if leftSample > rightSample {
					// The left sample is greater, therefore the peak is in the
					// 'left' peaks cluster. We want to go through the 'right' side,
					// and remove the peak property from all samples that are contiguous
					// to the 'rightSample', and are equal to it.
					return mergeSamples[T](left, right)
				} else {
					// The 'right' sample is greater, but it is the 'left' sample that is a peak.
					// Therefore, we must remove the peak property from all contiguous trailing
					// samples on the 'left' side.
					//
					// However, we do not need to make the first sample of the 'right' peaks
					// cluster a peak. Just because it is greater than its neighbor, i.e.,
					// the last sample of the 'left' peaks cluster, does not mean that it is
					// in fact a peak in the 'right' peaks cluster. Whether this sample is
					// a peak in the 'right' peaks cluster has already been determined when
					// we computed the peaks in the right peaks cluster.
					return mergeSamples[T](removeContiguousPeaksFromLeftSide[T](left, leftSample),
						addContiguousPeaksToRightSide[T](right))
				}
			}
		}
	} else {
		if !right.isFirstSamplePeak() {
			// Neither the last sample of the 'left' peaks cluster, nor
			// the first sample of the 'right' peaks cluster, are peaks.
			if leftSample == rightSample {
				// Because the two samples are equal to the same value, one
				// cannot be a peak relative to the other, by definition.
				//
				// Additionally, since there are no peaks on the boundary where
				// the two peak clusters are concatenated, there is no peak to
				// propagate neither 'left' nor 'right'. Because there are no peaks
				// to propagate, we simply merge the two sample arrays into one,
				// and then merge the peaks, and then fix up the peak offsets.
				return mergeSamples[T](left, right)
			} else {
				if leftSample > rightSample {
					// In this case, the last sample of the 'left' peaks cluster is
					// greater than the first sample of the 'right' peaks cluster.
					//
					// However, neither is a peak in its corresponding peak cluster.
					return mergeSamples[T](addContiguousPeaksToLeftSide[T](left), right)
				} else {
					// In this case, the last sample of the 'left' peaks cluster
					// is less than the first sample of the 'right' peaks cluster.
					//
					// This case is similar to the one above, i.e., where the last
					// sample of the 'left' peaks cluster is greater than the first
					// sample of the 'right' peaks cluster, but neither one is a peak.
					//
					// We have to assign the peak property to the first sample of the
					// 'right' peaks cluster, and then extend this property to the right
					// for all samples that are equal to the first sample.
					return mergeSamples[T](left, addContiguousPeaksToRightSide[T](right))
				}
			}
		} else {
			// The last sample of the 'left' peaks cluster is not a peak,
			// but the first sample of the 'right' peaks cluster is a peak.

			// If the two samples are equal
			//
			// We iterate on the right side from left to right, and remove
			// the peak property from all contiguous samples that are equal
			// to the value of the last sample on the left, or the first
			// sample on the right, since the two are equal in this case.
			if leftSample == rightSample {
				if len(left.peaks) == 0 {
					return mergeSamples[T](addContiguousPeaksToLeftSide[T](left), right)
				} else {
					return mergeSamples[T](left, removeContiguousPeaksFromRightSide[T](right, rightSample))
				}
			} else {
				// Additionally, the first sample of the 'right' peaks cluster,
				// is not equal to the last sample of the 'left' peaks cluster,
				// but is less than it.
				//
				// This means that only the higher of the two values is a peak.
				//
				// In this case we needn't do any recalculation, because the first
				// sample of the 'right' peaks cluster, is already less than the
				// last sample of the 'left' peaks cluster.
				//
				// Therefore, we simply merge the two sample arrays into one, and
				// then merge the peaks, while fixing up the offsets.
				if leftSample > rightSample {
					// The 'left' sample is greater, therefore the peak is in the 'left' peaks
					// cluster. We want to go through the 'right' side, and remove the peak
					// property from all samples that are contiguous to the first sample,
					// and are equal to it.
					return mergeSamples[T](addContiguousPeaksToLeftSide[T](left), removeContiguousPeaksFromRightSide[T](right, rightSample))
				} else {
					// The 'right' sample is greater than the 'left' sample, and is a peak.
					//
					// In this case, we simply merge the 'left' and the 'right' side without
					// any additional changes to the peaks of either side.
					return mergeSamples[T](left, right)
				}
			}
		}
	}
}

func mergeSamples[T Number](left, right PrimaryPeaks[T]) PrimaryPeaks[T] {
	samples := append(left.samples, right.samples...)
	peaks := append(left.peaks, right.peaks...)
	offset := len(left.samples)
	for i := 0; i < len(right.peaks); i++ {
		peaks[len(left.peaks)+i] = offset + right.peaks[i]
	}
	return CreatePeaksWith[T](samples, peaks)
}

// invariant: the rightmost sample of the 'left' peaks cluster must not be a peak
//
// The only case in which we need to add peaks to the 'left' side, is in the event
// the 'left' side does not have any peaks at all. This could happen for only one
// reason, when all samples have the same value. If there is a peak on the 'left'
// side, then we do not need to do anything. If there is not a peak on the 'left'
// side, then it means that all the samples have the same value, and in that case
// we just add the 3 peaks, one for each sample.
func addContiguousPeaksToLeftSide[T Number](left PrimaryPeaks[T]) PrimaryPeaks[T] {
	if len(left.peaks) == 0 {
		// No peaks on the 'left' side, means that all the samples have the same value.
		//
		// In this case, upon being concatenated with the 'right' side, the 'left'
		// side samples now have the property of being peaks, because the last
		// sample in 'left' peaks cluster, is greater than the first sample in the
		// 'right' peaks cluster.
		return CreatePeaksWith[T](left.samples, left.createAllPeaks())
	} else {
		return left
	}
}

// invariant: the rightmost sample of the 'left' peaks cluster must be a peak
func removeContiguousPeaksFromLeftSide[T Number](left PrimaryPeaks[T], sampleValue T) PrimaryPeaks[T] {
	lastPeak := len(left.peaks) - 1
	contiguousEqualSamplesCount := 1
	for i := len(left.peaks) - 2; i >= 0; i-- {
		if left.samples[left.peaks[i]] == sampleValue && left.peaks[lastPeak]-left.peaks[i] == 1 {
			contiguousEqualSamplesCount++
			lastPeak = i
		} else {
			break
		}
	}
	return CreatePeaksWith[T](left.samples, left.peaks[:len(left.peaks)-contiguousEqualSamplesCount])
}

// invariant: the leftmost sample of the 'right' peaks cluster must not be a peak
//
// The only case in which we need to add peaks to the 'right' side, is in the event
// the 'right' side does not have any peaks at all. This could happen for only one
// reason, when all samples have the same value. If there is a peak on the 'right'
// side, then we do not need to do anything. If there is not a peak on the 'right'
// side, then it means that all the samples have the same value, and in that case
// we just add the 3 peaks, one for each sample.
func addContiguousPeaksToRightSide[T Number](right PrimaryPeaks[T]) PrimaryPeaks[T] {
	if len(right.peaks) == 0 {
		// No peaks on the 'right' side, means that all the samples have the same value.
		//
		// In this case, upon being concatenated with the 'left' side, the 'right'
		// side samples now have the property of being peaks, because the first
		// sample in 'right' peaks cluster, is greater than the last sample in the
		// 'left' peaks cluster.
		return CreatePeaksWith[T](right.samples, right.createAllPeaks())
	} else {
		return right
	}
}

// invariant: the leftmost sample of the 'right' peaks cluster must be a peak
func removeContiguousPeaksFromRightSide[T Number](right PrimaryPeaks[T], sampleValue T) PrimaryPeaks[T] {
	lastPeak := 0
	contiguousEqualSamplesCount := 1
	for i := 1; i < len(right.peaks); i++ {
		if right.samples[right.peaks[i]] == sampleValue && right.peaks[i]-right.peaks[lastPeak] == 1 {
			contiguousEqualSamplesCount++
			lastPeak = i
		} else {
			break
		}
	}
	return CreatePeaksWith[T](right.samples, right.peaks[contiguousEqualSamplesCount:])
}
