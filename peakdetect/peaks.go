// Copyright (c) 2024 Andrei Gill. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package peakdetect

import (
	"golang.org/x/exp/constraints"
	"log"
)

type Number interface {
	constraints.Integer | constraints.Float
}

type Peaks[T Number] interface {
	GetSampleCount() int
	GetPeakCount() int
	GetSamples() []T
	GetPeaks() []int
}

type PrimaryPeaks[T Number] struct {
	samples []T
	peaks   []int
}

type SecondaryPeaks[T Number] struct {
	PrimaryPeaks[T]
	primarySamples []T
	primaryPeaks   []int
	originalPeaks  []int
}

func CreatePeaks[T Number](a, b, c T, peaks []int) PrimaryPeaks[T] {
	return CreatePeaksWith([]T{a, b, c}, peaks)
}

func CreatePeaksWith[T Number](samples []T, peaks []int) PrimaryPeaks[T] {
	return PrimaryPeaks[T]{samples, peaks}
}

func CreateSecondaryPeaksWith[T Number](p PrimaryPeaks[T], primaryPeaks []int, originalPeaks []int) SecondaryPeaks[T] {
	return SecondaryPeaks[T]{p, []T{}, primaryPeaks, originalPeaks}
}

func (p *PrimaryPeaks[T]) GetSampleCount() int {
	return len(p.samples)
}

func (p *PrimaryPeaks[T]) GetPeakCount() int {
	return len(p.peaks)
}

func (p *PrimaryPeaks[T]) GetSamples() []T {
	return p.samples
}

func (p *PrimaryPeaks[_]) GetPeaks() []int {
	return p.peaks
}

func (p *SecondaryPeaks[T]) GetPrimarySamples() []T {
	return p.primarySamples
}

func (p *SecondaryPeaks[_]) GetPrimaryPeaks() []int {
	return p.primaryPeaks
}

func (p *PrimaryPeaks[T]) getFirstSample() T {
	return p.samples[0]
}

func (p *PrimaryPeaks[T]) getLastSample() T {
	return p.samples[len(p.samples)-1]
}

func (p *PrimaryPeaks[_]) isFirstSamplePeak() bool {
	return p.isSampleFoundInPeaks(0)
}

func (p *PrimaryPeaks[_]) isLastSamplePeak() bool {
	return p.isSampleFoundInPeaks(len(p.samples) - 1)
}

func (p *PrimaryPeaks[_]) isSampleFoundInPeaks(sampleIndex int) bool {
	for i := 0; i < len(p.peaks); i++ {
		if p.peaks[i] == sampleIndex {
			return true
		}
	}
	return false
}

func (p *PrimaryPeaks[_]) createAllPeaks() []int {
	peaks := make([]int, len(p.samples))
	for i := 0; i < len(peaks); i++ {
		peaks[i] = i
	}
	return peaks
}

// Inflate
// returns an array wherein all samples that were originally
// not peaks, are set to zero, and wherein all peaks are set
// to their original peak values
func (p *PrimaryPeaks[T]) Inflate() []T {
	return p.InflateWithCount(len(p.samples), p)
}

func (p *PrimaryPeaks[T]) InflateWithCount(samples int, from Peaks[T]) []T {
	if samples < 0 {
		log.Fatal("Number of samples must be greater than zero")
	}
	peaks := make([]T, samples)
	p.InflateInto(peaks, from)
	return peaks
}

func (p *PrimaryPeaks[T]) InflateInto(peaks []T, from Peaks[T]) {
	if len(p.peaks) > len(peaks) {
		log.Fatal("Invalid state: more peaks than space for samples in the array")
	}
	fromPeaks := from.GetPeaks()
	fromSamples := from.GetSamples()
	for i := 0; i < len(p.peaks); i++ {
		at := fromPeaks[i]
		peaks[at] = fromSamples[at]
	}
}

func AlignPeaksToSamplePositions(sampleCount int, peaks []int) []int {
	result := make([]int, sampleCount)
	for i := 0; i < len(peaks); i++ {
		result[peaks[i]] = peaks[i]
	}
	return result
}

type secondaryAsPrimary[T Number] struct {
	*SecondaryPeaks[T]
}

func PrimaryValuesOnly[T Number](p *SecondaryPeaks[T]) Peaks[T] {
	return &secondaryAsPrimary[T]{p}
}

func (p *secondaryAsPrimary[T]) GetSamples() []T {
	return p.GetPrimarySamples()
}

func (p *secondaryAsPrimary[_]) GetPeaks() []int {
	return p.GetPrimaryPeaks()
}
