package speedclip

type ValueUnit string

const (
	None         ValueUnit = "none"
	Nanoseconds            = "nanoseconds"
	Microseconds           = "microseconds"
	Milliseconds           = "milliseconds"
	Seconds                = "seconds"
	Bytes                  = "bytes"
)

type ProfileSampled struct {
	Type       string
	Name       string
	Unit       ValueUnit
	StartValue float64
	EndValue   float64
	Samples    [][]int
	Weights    [][]int
}

type Frame map[string]string
