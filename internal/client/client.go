package client

import (
	"fmt"
	"math"
	"text/tabwriter"

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/cli/internal/utils"
)

type Client struct {
	ots pkg.OvertimeService
}

func Init(host string, authHeader string) Client {
	return Client{
		ots: pkg.InitOvertimeClient(host, authHeader),
	}
}

func printHoliday(w *tabwriter.Writer, a *pkg.Holiday) {
	fmt.Fprintf(w, "ID\t: %d\n", a.ID)
	fmt.Fprintf(w, "Description\t: %s\n", a.Description)
	fmt.Fprintf(w, "Type\t: %s\n", a.Type)
	fmt.Fprintf(w, "Start\t: %s\n", utils.FormatTime(a.Start))
	fmt.Fprintf(w, "End\t: %s\n", utils.FormatTime(a.End))
	diff := a.End.Sub(a.Start)
	fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(diff.Minutes())))
}

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
