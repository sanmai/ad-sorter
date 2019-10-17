package bannersorter

import (
	"errors"
	"time"
)

// BannerFinder defined a method to get a banner object to display.
type BannerFinder interface {
	CurrentBanner() (interface{}, error)
}

type bannerFinder struct {
	inputs []Input
	now    time.Time
	qa     bool
}

// NewBannerFinder makes new production banner finder.
func NewBannerFinder(inputs []Input) BannerFinder {
	return NewBannerFinderEx(inputs, time.Now(), false)
}

// NewStagingBannerFinder makes new Staging banner finder.
func NewStagingBannerFinder(inputs []Input) BannerFinder {
	return NewBannerFinderEx(inputs, time.Now(), true)
}

// NewBannerFinderEx makes new banner finder for E2E testing purposes.
func NewBannerFinderEx(inputs []Input, now time.Time, qa bool) BannerFinder {
	return bannerFinder{inputs, now, qa}
}

func (bf bannerFinder) CurrentBanner() (interface{}, error) {
	selector := selector{bf.inputs, func(input Input) bool {
		// end time in the past - bail
		if input.EndTime().Before(bf.now) {
			return false
		}

		// Staging assumes future banners are shown always
		if bf.qa {
			return true
		}

		// production, and start time in the future - bail
		if input.StartTime().After(bf.now) {
			return false
		}

		return true
	}}

	result := selector.earliestExpiringInput()

	if result == nil {
		return nil, errors.New("Nothing to display")
	}

	switch value := (*result).(type) {
	case wrapper:
		return value.unwrap(), nil
	default:
		return value, nil
	}
}
