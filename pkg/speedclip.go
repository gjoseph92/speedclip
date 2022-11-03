package speedclip

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

func readSpeedscope(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var file map[string]interface{}
	json.Unmarshal(data, &file)
	return file, nil
}

// function writeSpeedscope writes JSON data to a file, or to stdout if `"-"` is given as the path
func writeSpeedscope(path string, contents interface{}) error {
	data, err := json.Marshal(contents)
	if err != nil {
		return err
	}

	var file *os.File
	if path == "-" {
		file = os.Stdout
	} else {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	_, err = file.Write(data)
	return err
}

// function Clip crops a speedscope profile to the given start and end points
func Clip(path string, out string, start time.Duration, end time.Duration) error {
	file, err := readSpeedscope(path)
	if err != nil {
		return err
	}

	profiles, ok := file["profiles"].([]interface{})
	if !ok {
		return fmt.Errorf("'profiles' field missing or wrong type")
	}

	for _, p := range profiles {
		profile := p.(map[string]interface{})

		unit, ok := profile["unit"].(string)
		if !ok {
			return fmt.Errorf("'unit' field missing or wrong type from profile %v", profile["name"])
		}
		var duration time.Duration
		switch unit {
		case "nanoseconds":
			duration = time.Nanosecond
		case "microseconds":
			duration = time.Microsecond
		case "milliseconds":
			duration = time.Millisecond
		case "seconds":
			duration = time.Second
		default:
			// return fmt.Errorf("unsupported duration %v", unit)
			duration = time.Second
		}

		startValue, ok := profile["startValue"].(float64)
		if !ok {
			// TODO we don't really need this, if empty could assume 0
			return fmt.Errorf("'startValue' field missing or wrong type")
		}

		profileStart := time.Duration(startValue * float64(duration))

		samples, ok := profile["samples"].([]interface{})
		if !ok {
			return fmt.Errorf("'samples' field missing or wrong type")
		}

		weights, ok := profile["weights"].([]interface{})
		if !ok {
			return fmt.Errorf("'weights' field missing or wrong type")
		}

		if len(samples) != len(weights) {
			return fmt.Errorf("'samples' and 'weights' have different lengths: %v, %v", len(samples), len(weights))
		}

		if start < 0 || end < 0 {
			// sum up total duration to calculate offset from end
			totalDuration := time.Duration(0)
			for _, wi := range weights {
				w := wi.(float64)
				totalDuration += time.Duration(w * float64(duration))
			}

			if start < 0 {
				start = totalDuration + start
			}
			if end < 0 {
				end = totalDuration + end
			}
		}

		if end != 0 && end < start {
			return fmt.Errorf("end %v < start %v", end, start)
		}

		log.Printf("%v %v", start, end)

		origLen := len(samples)
		current := time.Duration(0)
		start_i := 0
		new_start := time.Duration(0)
		new_end := time.Duration(0)
		end_i := len(weights)

		// silly linear search for first and last sample
		if start > 0 {
			for i := 0; i < len(weights); i++ {
				w := weights[i].(float64)
				current += time.Duration(w * float64(duration))
				if current >= start {
					start_i = i
					new_start = current + profileStart
					break
				}
			}
		}

		if end > 0 {
			for i := start_i; i < len(weights); i++ {
				w := weights[i].(float64)
				current += time.Duration(w * float64(duration))
				if current > end {
					end_i = i
					new_end = current + profileStart
					break
				}
			}
		}

		log.Printf("%v Kept %v/%v: %v -> %v\n", profile["name"], end_i-start_i, origLen, start_i, end_i)

		if start > 0 {
			profile["startValue"] = float64(new_start / duration)
		}
		if end > 0 {
			profile["endValue"] = float64(new_end / duration)
		}
		profile["weights"] = weights[start_i:end_i]
		profile["samples"] = samples[start_i:end_i]
	}

	writeSpeedscope(out, file)
	return nil
}
