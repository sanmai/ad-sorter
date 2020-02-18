package bannersorter_test

import (
	"testing"
	"time"

	bs "github.com/sanmai/adsorter"
)

func TestDateTimeMath(t *testing.T) {
	jst, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+09:00")
	cet, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+02:00")

	if cet.After(jst) == false {
		t.Errorf("CET %q shall be before JST %q", cet.Local(), jst.Local())
	}
}

func TestProdExample(t *testing.T) {
	want := 42
	jst, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+09:00")
	cet, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+02:00")
	now, _ := time.Parse(time.RFC3339, "2019-10-11T12:30:00+09:00")

	finder := bs.NewBannerFinderEx([]bs.Input{bs.Wrap(want, jst, cet)}, now, false)

	if got, err := finder.CurrentBanner(); got != want || err != nil {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdExampleNegativePast(t *testing.T) {
	jst, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+09:00")
	cet, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+02:00")
	now, _ := time.Parse(time.RFC3339, "2019-10-09T12:30:00+09:00")

	finder := bs.NewBannerFinderEx([]bs.Input{bs.Wrap(42, jst, cet)}, now, false)

	if got, err := finder.CurrentBanner(); err == nil {
		t.Errorf("got %q, expected nil", got)
	}
}

func TestProdExampleNegativeFuture(t *testing.T) {
	jst, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+09:00")
	cet, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+02:00")

	finder := bs.NewBannerFinder([]bs.Input{bs.Wrap(42, jst, cet)})

	if got, _ := finder.CurrentBanner(); got != nil {
		t.Errorf("got %q, expected nil", got)
	}
}

func TestProdEmpty(t *testing.T) {
	finder := bs.NewBannerFinder([]bs.Input{})

	if got, _ := finder.CurrentBanner(); got != nil {
		t.Errorf("got %q, expected nil", got)
	}
}

func makeInput(value int, starts string, ends string) bs.Input {
	return makeInputFixedTime(value, starts, ends, time.Now())
}

func makeInputFixedTime(value int, starts string, ends string, now time.Time) bs.Input {
	startDuration, err := time.ParseDuration(starts)
	if err != nil {
		panic(err)
	}

	endDuration, err := time.ParseDuration(ends)
	if err != nil {
		panic(err)
	}

	return bs.Wrap(value, now.Add(startDuration), now.Add(endDuration))
}

func TestMakeInput(t *testing.T) {
	now := time.Now()
	input := makeInput(1, "-1h", "+1h")

	if !input.StartTime().Before(input.EndTime()) {
		t.Errorf("%q not before %q", input.StartTime(), input.EndTime())
	}

	if !input.StartTime().Before(now) {
		t.Errorf("%q not before %q", input.StartTime(), now)
	}

	if !input.EndTime().After(now) {
		t.Errorf("%q not after %q", input.EndTime(), now)
	}
}

func TestProdFindsSingle(t *testing.T) {
	want := 42

	finder := bs.NewBannerFinder([]bs.Input{makeInput(want, "-1h", "+1h")})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsNothingInGap(t *testing.T) {
	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(2, "+1h", "+2h"),
	})

	if got, err := finder.CurrentBanner(); got != nil || err == nil {
		t.Errorf("got %q, want nil", got)
	}
}

func TestStagingFindsNextBannerInGap(t *testing.T) {
	want := 2

	finder := bs.NewStagingBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(want, "+1h", "+2h"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsNothingInFuture(t *testing.T) {
	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(2, "-1h", "-10m"),
	})

	if got, err := finder.CurrentBanner(); got != nil || err == nil {
		t.Errorf("got %q, want nil", got)
	}
}

func TestStagingFindsNothingInFuture(t *testing.T) {
	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(2, "-1h", "-10m"),
	})

	if got, _ := finder.CurrentBanner(); got != nil {
		t.Errorf("got %q, want nil", got)
	}
}

func TestProdFindsCurrentBanner(t *testing.T) {
	want := 2

	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(want, "-1h", "+1h"),
		makeInput(3, "+1h", "+2h"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStagingFindsCurrentBanner(t *testing.T) {
	want := 2

	finder := bs.NewStagingBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(want, "-1h", "+1h"),
		makeInput(3, "+1h", "+2h"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsCurrentBannerOverlaidByNext(t *testing.T) {
	want := 2

	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(want, "-1h", "+1h"),
		makeInput(3, "+10m", "+20m"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStagingFindsNextBannerOverlaidOverCurrent(t *testing.T) {
	want := 3

	finder := bs.NewStagingBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(2, "-1h", "+1h"),
		makeInput(want, "+10m", "+20m"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsCurrentBannerOverlaidAnotherCurrent(t *testing.T) {
	want := 2

	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInput(want, "-1h", "+1h"),
		makeInput(3, "-10m", "+2h"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsCurrentBannerUnordered(t *testing.T) {
	want := 2

	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(3, "+1h", "+2h"),
		makeInput(want, "-1h", "+1h"),
		makeInput(1, "-2h", "-1h"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestStagingFindsCurrentBannerUnordered(t *testing.T) {
	want := 2

	finder := bs.NewStagingBannerFinder([]bs.Input{
		makeInput(3, "+1h", "+2h"),
		makeInput(want, "-1h", "+1h"),
		makeInput(1, "-2h", "-1h"),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsFirstWhenEndTimeIsSame(t *testing.T) {
	want := 2
	now := time.Now()

	finder := bs.NewBannerFinder([]bs.Input{
		makeInput(1, "-2h", "-1h"),
		makeInputFixedTime(want, "-1h", "+1h", now),
		makeInputFixedTime(3, "-10m", "+1h", now),
	})

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsExpiringNow(t *testing.T) {
	want := 1
	now := time.Now()

	finder := bs.NewBannerFinderEx([]bs.Input{
		makeInputFixedTime(want, "-1h", "+0h", now),
	}, now, false)

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestProdFindsStartedNow(t *testing.T) {
	want := 1
	now := time.Now()

	finder := bs.NewBannerFinderEx([]bs.Input{
		makeInputFixedTime(want, "-0h", "+1h", now),
	}, now, false)

	if got, _ := finder.CurrentBanner(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

type testBanner struct {
	name               string
	startTime, endTime time.Time
}

func (tb testBanner) StartTime() time.Time {
	return tb.startTime
}

func (tb testBanner) EndTime() time.Time {
	return tb.endTime
}

func TestProdFindsNonWrapped(t *testing.T) {
	prot := makeInput(0, "-1h", "+1h")

	want := "example"

	finder := bs.NewBannerFinder([]bs.Input{
		testBanner{want, prot.StartTime(), prot.EndTime()},
	})

	got, err := finder.CurrentBanner()
	tb := got.(testBanner)

	if err != nil {
		t.Errorf("Expected no error, found %q", err)
	}

	if tb.name != want {
		t.Errorf("got %q, want %q", tb.name, want)
	}
}
