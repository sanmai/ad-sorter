package bannersorter

import "testing"
import "time"

func TestWrap(t *testing.T) {
	wantBannerID := 42
	wantStartTime := time.Unix(100, 0)
	wantEndTime := time.Unix(200, 0)
	input := Wrap(wantBannerID, wantStartTime, wantEndTime)

	wrap := input.(wrapper)

	if got := wrap.unwrap().(int); got != wantBannerID {
		t.Errorf("got %q, want %q", got, wantBannerID)
	}

	if got := wrap.StartTime(); got != wantStartTime {
		t.Errorf("got %q, want %q", got, wantStartTime)
	}

	if got := wrap.EndTime(); got != wantEndTime {
		t.Errorf("got %q, want %q", got, wantEndTime)
	}
}
