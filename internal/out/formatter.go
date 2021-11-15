package out

import (
	"encoding/json"
	"fmt"
	"math"
)

func formatMinutes(t int64) string {
	ds, hs1 := math.Modf(float64(t) / (24 * 60))
	hs2, mf := math.Modf(hs1 * 24)
	if ds == 0 {
		return fmt.Sprintf("%02dh:%02dm", int(hs2), int(mf*60))
	}
	return fmt.Sprintf("%02dd:%02dh:%02dm", int(ds), int(hs2), int(mf*60))
}

func formatMinutesToHoursAndMinutes(t int64) string {
	hs2, mf := math.Modf(float64(t) / 60)
	return fmt.Sprintf("%02dh:%02dm", int(hs2), int(mf*60))
}

func PrintJson(v interface{}) error {
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", data)
	return nil
}
