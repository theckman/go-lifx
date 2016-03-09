package lifxpayloads

import (
	"math"
	"time"
)

const (
	maxUint16 float64 = 65535

	colorRangeMax   float64 = 359
	colorRangeMin   float64 = 0
	hueRange                = maxUint16 - colorRangeMin
	colorWheelRange         = colorRangeMax - colorRangeMin

	percRangeMax    float64 = 100
	percRangeMin    float64 = 0
	percValueRange          = maxUint16 - percRangeMin
	percScaledRange         = percRangeMax - percRangeMin
)

func colorRange(value float64) uint16 {
	scaledValue := (((value - colorRangeMin) * colorWheelRange) / hueRange) + colorRangeMin
	return uint16(scaledValue)
}

func percentageRange(value float64) uint8 {
	scaledValue := (((value - percRangeMin) * percScaledRange) / percValueRange) + percRangeMin
	return uint8(scaledValue)
}

// nsecEpochToTime converts a UNIX epoch with nanosecond
// precision in to a time.Time where the Timezone is UTC.
func nsecEpochToTime(nanoseconds uint64) time.Time {
	nanoDur := time.Duration(nanoseconds)

	// convert the value to the UNIX epoch
	// with remaining nanoseconds (npoch)
	epoch := int64(nanoDur / time.Second)
	npoch := int64(nanoDur % time.Second)

	return time.Unix(epoch, npoch).UTC()
}

func durToMs(dur time.Duration) uint32 {
	return uint32(dur / time.Millisecond)
}

func msToDur(ms uint32) time.Duration {
	return time.Duration(ms) * time.Millisecond
}

func round(f float64) float64 {
	if f < 0 {
		return math.Ceil(f - 0.5)
	}

	return math.Floor(f + 0.5)
}
