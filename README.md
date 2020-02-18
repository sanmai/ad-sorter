[![Build Status](https://travis-ci.com/sanmai/adsorter.svg?branch=master)](https://travis-ci.com/sanmai/adsorter)

# Ad Sorter

The task is to make a reusable library that will take a list of some current and pending advertisements (AKA banners), then keeping only one of them according to certain rules. Each banner has a time limit set by start and end times, specified with a timezone.

There are two rulesets.

- Normal rules are to keep only the banner that is supposed to be displayed at this moment, with a preference to banners ending the earlies.
- Staging rules are to show any future banner if there isn't a banner to be displayed currently, in addition to the normal requirements.

## Overview

This library is built upon the basic principles of the Go standard library: given two [time.Time](https://golang.org/pkg/time/#Time) objects, they can be compared regardless of the timezones each of them were initialized with. Hence, given a `time.Time` object, it can be compared against the current time just as well.

Consider this example:

	jst, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+09:00")
	cet, _ := time.Parse(time.RFC3339, "2019-10-11T11:20:00+02:00")

	fmt.Println(cet.After(jst)) // prints "true"

Here we can tell that 11:20 in Europe is a very late evening in Japan, and that's what the `time.After` will [tell us right away](https://play.golang.org/p/d_B4Z0wCrTd), without any need for us to know about time zones, meanwhile Go being perfectly aware of different timezones for the object.

Provided the time and date for the banner stored in the database with corresponding timezone, no additional effort required from the user. If a timezone is stored independently, there are means in the standard library [to initialize a `Time` object with a timezone/location](https://golang.org/pkg/time/#Date).

To run tests use the standard `go test` command.

## Inputs

The library takes a slice of objects corresponding to the `Input` interface.
```go
type Input interface {
	StartTime() time.Time
	EndTime() time.Time
}
```
If it is unfeasible to implement this interface in the database layer, there's a handy adapter-wrapper, which takes a banner object, and two `Time` objects.

```go
func Wrap(filling interface{}, startTime, endTime time.Time) Input
```

## Outputs

The main object of the library comes in two main flavors. For production environment:

```go
func NewBannerFinder(inputs []Input) BannerFinder
```
And for the staging environment:
```go
func NewStagingBannerFinder(inputs []Input) BannerFinder
```

Both implement the following interface:
```go
type BannerFinder interface {
	CurrentBanner() (interface{}, error)
}
```
The usage boils down to calling a corresponding method for production or staging environments with a slice of inputs, later calling `CurrentBanner()` on the returned object. This method will return either an original input, or, if a wrapper was used, an embedded object. It will return an error with description if there's nothing to display.

Additionally there's `NewBannerFinderEx` to make a custom finder with either a staging flag toggled on, or a current time set to other time. It is primarily used for testing purposes, and can be used for additional integration testing.

## Example

```go
banner := loadOneFromDatabase()

start, _ := time.Parse(time.RFC3339, banner.start)
end, _ := time.Parse(time.RFC3339, banner.end)

finder := bs.NewBannerFinder([]bs.Input{bs.Wrap(banner, start, end)})

if pBanner, err := finder.CurrentBanner(); err == nil {
	bannerToShow := pBanner.(bannerType)
	// ... respond to the API request with banner details
}
```

## Notes

- `go test -cover` reports 100% coverage at time of writing.

- It should not be a problem to implement a library following the very same approach in PHP.


