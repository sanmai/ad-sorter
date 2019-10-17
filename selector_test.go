package bannersorter

import "testing"
import "time"

func TestSelectorTrue(t *testing.T) {
	want := 42
	wrap := Wrap(want, time.Now(), time.Now())
	selector := selector{[]Input{wrap}, func(in Input) bool {
		return true
	}}

	if got := (*selector.earliestExpiringInput()).(wrapper).unwrap(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSelectorFalse(t *testing.T) {
	wrap := Wrap(42, time.Now(), time.Now())
	selector := selector{[]Input{wrap}, func(in Input) bool {
		return false
	}}

	if got := selector.earliestExpiringInput(); got != nil {
		t.Errorf("got %T, want nil", got)
	}
}

func TestSameEndTime(t *testing.T) {
	now := time.Now()

	want := 1

	first := Wrap(want, now, now)
	second := Wrap(want+1, now, now)

	selector := selector{[]Input{first, second}, func(in Input) bool {
		return true
	}}

	if got := (*selector.earliestExpiringInput()).(wrapper).unwrap(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestEarliestEnding(t *testing.T) {
	want := 2

	first := Wrap(want-1, time.Now(), timeWithAdded("3h"))
	second := Wrap(want, time.Now(), timeWithAdded("20m"))
	third := Wrap(want+1, time.Now(), timeWithAdded("40m"))

	selector := selector{[]Input{first, second, third}, func(in Input) bool {
		return true
	}}

	if got := (*selector.earliestExpiringInput()).(wrapper).unwrap(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func timeWithAdded(duration string) time.Time {
	interval, err := time.ParseDuration(duration)
	if err != nil {
		panic(err)
	}

	return time.Now().Add(interval)
}
