package iso8601duration

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	t.Parallel()

	// test with bad format
	_, err := FromString("asdf")
	assert.Equal(t, err, ErrBadFormat)

	// test with good full string
	dur, err := FromString("P1Y2M3W4DT3H4M5S")
	assert.Nil(t, err)
	assert.Equal(t, 1, dur.Years)
	assert.Equal(t, 2, dur.Months)
	assert.Equal(t, 3, dur.Weeks)
	assert.Equal(t, 4, dur.Days)
	assert.Equal(t, 3, dur.Hours)
	assert.Equal(t, 4, dur.Minutes)
	assert.Equal(t, 5, dur.Seconds)

	// test with good week string
	dur, err = FromString("P1W")
	assert.Nil(t, err)
	assert.Equal(t, 1, dur.Weeks)
}

func TestString(t *testing.T) {
	t.Parallel()

	// test empty
	d := Duration{}
	assert.Equal(t, d.String(), "P")

	// test only larger-than-day
	d = Duration{Years: 1, Days: 2}
	assert.Equal(t, d.String(), "P1Y2D")

	// test only smaller-than-day
	d = Duration{Hours: 1, Minutes: 2, Seconds: 3}
	assert.Equal(t, d.String(), "PT1H2M3S")

	// test month
	d = Duration{Hours: 1, Months: 2}
	assert.Equal(t, d.String(), "P2MT1H")

	// test full format
	d = Duration{Years: 1, Days: 2, Hours: 3, Minutes: 4, Seconds: 5}
	assert.Equal(t, d.String(), "P1Y2DT3H4M5S")

	// test week format
	d = Duration{Weeks: 1}
	assert.Equal(t, d.String(), "P1W")
}

func TestToEstimatedDuration(t *testing.T) {
	t.Parallel()

	d := Duration{Years: 1}
	assert.Equal(t, d.ToEstimatedDuration(), time.Hour*24*365)

	d = Duration{Months: 1}
	assert.Equal(t, d.ToEstimatedDuration(), time.Hour*24*30)

	d = Duration{Weeks: 1}
	assert.Equal(t, d.ToEstimatedDuration(), time.Hour*24*7)

	d = Duration{Days: 1}
	assert.Equal(t, d.ToEstimatedDuration(), time.Hour*24)

	d = Duration{Hours: 1}
	assert.Equal(t, d.ToEstimatedDuration(), time.Hour)

	d = Duration{Minutes: 1}
	assert.Equal(t, d.ToEstimatedDuration(), time.Minute)

	d = Duration{Seconds: 1}
	assert.Equal(t, d.ToEstimatedDuration(), time.Second)
}

func TestDuration(t *testing.T) {
	now := time.Now()

	d := Duration{
		Years:   2,
		Months:  3,
		Weeks:   1,
		Days:    1,
		Hours:   2,
		Minutes: 51,
		Seconds: 12,
	}

	dur := d.ToDuration(now)
	inaccurateDur := d.ToEstimatedDuration()

	future :=
		now.
			AddDate(d.Years, d.Months, 0).
			AddDate(0, 0, 7*d.Weeks).
			AddDate(0, 0, d.Days).
			Add(time.Duration(d.Hours) * time.Hour).
			Add(time.Duration(d.Minutes) * time.Minute).
			Add(time.Duration(d.Seconds) * time.Second)

	stdDur := future.Sub(now)

	log.Println("now", now)
	log.Println("future", future)
	log.Println(dur, "vs", stdDur)
	log.Println("inaccurate dur:", inaccurateDur, "time:", now.Add(inaccurateDur), "diff:", inaccurateDur-dur)

	assert.Equal(t, stdDur, dur)
}
