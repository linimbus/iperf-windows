package conversions

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type UnitsFormatter interface {
	Format(units string) string
}

type DefaultFormatter struct {
}

func (df DefaultFormatter) Format(units string) string {
	return fmt.Sprintf("%sbit", units)
}

func StringBitRateToInt(rate string) (value int64, err error) {
	numeric := ""
	units := ""
	valuef := 0.0
	multiplier := 1.0
	for _, c := range rate {
		switch {
		case c >= '0' && c <= '9' || c == '.':
			numeric += string(c)
		default:
			units += string(c)
		}
	}
	valuef, err = strconv.ParseFloat(numeric, 64)
	if err != nil {
		return 0, fmt.Errorf("couldn't parse '%v' into float", rate)
	}
	unitsLower := strings.ToLower(strings.TrimSpace(units))
	switch {
	case strings.HasPrefix(unitsLower, "k"):
		multiplier = 1000
	case strings.HasPrefix(unitsLower, "m"):
		multiplier = 1000 * 1000
	case strings.HasPrefix(unitsLower, "g"):
		multiplier = 1000 * 1000 * 1000
	case unitsLower == "":
		multiplier = 1
	default:
		return 0, fmt.Errorf("invalid units specified '%v'", units)
	}
	value = int64(math.Round(valuef * multiplier))
	return value, err
}

func IntBitRateToString(rate int64) string {
	return IntBitRateToStringFmt(rate, DefaultFormatter{})
}

func IntBitRateToStringFmt(rate int64, formatter UnitsFormatter) string {
	units := "b"
	value := float64(rate)
	if value >= 1000 {
		units = "k"
		value = value / 1000
	}
	if value >= 1000 {
		units = "m"
		value = value / 1000
	}
	if value >= 1000 {
		units = "g"
		value = value / 1000
	}
	if value >= 1000 {
		units = "t"
		value = value / 1000
	}
	suffix := formatter.Format(units)
	srate := fmt.Sprintf("%.2f%s", value, suffix)
	return srate
}
