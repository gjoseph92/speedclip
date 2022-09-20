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
			return fmt.Errorf("'startValue' field missing or wrong type")
		}
		endValue, ok := profile["endValue"].(float64)
		if !ok {
			return fmt.Errorf("'endValue' field missing or wrong type")
		}

		profileStart := time.Duration(startValue * float64(duration))
		profileEnd := time.Duration(endValue * float64(duration))

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

		origLen := len(samples)
		current := profileStart
		start_i := 0
		new_start := profileStart
		new_end := profileEnd
		haveFirst := false
		end_i := len(weights)

		// silly linear search for first and last sample
		for i, wi := range weights {
			w := wi.(float64)
			current += time.Duration(w * float64(duration))
			if !haveFirst && current > start {
				start_i = i
				new_start = current
				haveFirst = true
			}

			if current > end {
				end_i = i
				new_end = current
				break
			}
		}

		log.Printf("%v Kept %v/%v: %v -> %v\n", profile["name"], end_i-start_i, origLen, start_i, end_i)

		profile["startValue"] = float64(new_start / duration)
		profile["endValue"] = float64(new_end / duration)
		profile["weights"] = weights[start_i:end_i]
		profile["samples"] = samples[start_i:end_i]
	}

	writeSpeedscope(out, file)
	return nil
}
