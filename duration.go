// Package iso8601duration provides a partial implementation of ISO8601 durations
package iso8601duration

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"text/template"
	"time"
)

var (
	// ErrBadFormat is returned when parsing fails
	ErrBadFormat = errors.New("bad format string")

	tmpl = template.Must(template.New("duration").
		Parse(
			`P{{if .Years}}{{.Years}}Y{{end}}` +
				`{{if .Months}}{{.Months}}M{{end}}` +
				`{{if .Weeks}}{{.Weeks}}W{{end}}` +
				`{{if .Days}}{{.Days}}D{{end}}` +
				`{{if .HasTimePart}}T{{end}}` +
				`{{if .Hours}}{{.Hours}}H{{end}}` +
				`{{if .Minutes}}{{.Minutes}}M{{end}}` +
				`{{if .Seconds}}{{.Seconds}}S{{end}}`,
		),
	)

	full = regexp.MustCompile(`P((?P<year>\d+)Y)?((?P<month>\d+)M)?((?P<week>\d+)W)?((?P<day>\d+)D)?(T((?P<hour>\d+)H)?((?P<minute>\d+)M)?((?P<second>\d+)S)?)?`)
)

type Duration struct {
	Years   int
	Months  int
	Weeks   int
	Days    int
	Hours   int
	Minutes int
	Seconds int
}

func FromString(dur string) (*Duration, error) {
	var (
		match []string
		re    *regexp.Regexp
	)

	if full.MatchString(dur) {
		match = full.FindStringSubmatch(dur)
		re = full
	} else {
		return nil, ErrBadFormat
	}

	d := &Duration{}

	for i, name := range re.SubexpNames() {
		part := match[i]
		if i == 0 || name == "" || part == "" {
			continue
		}

		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		switch name {
		case "year":
			d.Years = val
		case "month":
			d.Months = val
		case "week":
			d.Weeks = val
		case "day":
			d.Days = val
		case "hour":
			d.Hours = val
		case "minute":
			d.Minutes = val
		case "second":
			d.Seconds = val
		default:
			return nil, errors.New(fmt.Sprintf("unknown field %s", name))
		}
	}

	return d, nil
}

// String prints out the value passed in. It's not strictly according to the
// ISO spec, but it's pretty close. In particular, to completely conform it
// would need to round up to the next largest unit. 61 seconds to 1 minute 1
// second, for example. It would also need to disallow weeks mingling with
// other units.
func (d *Duration) String() string {
	var s bytes.Buffer

	err := tmpl.Execute(&s, d)
	if err != nil {
		panic(err)
	}

	return s.String()
}

func (d *Duration) HasTimePart() bool {
	return d.Hours != 0 || d.Minutes != 0 || d.Seconds != 0
}

// ToEstimatedDuration returns an inaccurate duration that
// is independent of when counting is started
func (d *Duration) ToEstimatedDuration() time.Duration {
	day := time.Hour * 24
	month := day * 30
	year := day * 365

	tot := time.Duration(0)

	tot += year * time.Duration(d.Years)
	tot += month * time.Duration(d.Months)
	tot += day * 7 * time.Duration(d.Weeks)
	tot += day * time.Duration(d.Days)
	tot += time.Hour * time.Duration(d.Hours)
	tot += time.Minute * time.Duration(d.Minutes)
	tot += time.Second * time.Duration(d.Seconds)

	return tot
}

// ToDuration returns an accurate duration based on the current
// date in the calendar. As months and years have variable durations
// it's difficult to guess when exactly the duration will be passed.
// This method aims to return a duration that will exactly hit the
// expected time and date.
func (d *Duration) ToDuration(from time.Time) time.Duration {
	targetTime := from.
		AddDate(d.Years, d.Months, 0).
		AddDate(0, 0, 7*d.Weeks).
		AddDate(0, 0, d.Days).
		Add(time.Duration(d.Hours) * time.Hour).
		Add(time.Duration(d.Minutes) * time.Minute).
		Add(time.Duration(d.Seconds) * time.Second)
	return targetTime.Sub(from)
}
