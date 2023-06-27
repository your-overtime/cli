package client

import (
	"fmt"
	"math"
	"strings"
	"text/tabwriter"

	"github.com/your-overtime/api/v2/pkg"
	"github.com/your-overtime/cli/internal/utils"
)

type Client struct {
	ots     pkg.OvertimeService
	APIHost string
}

func Init(host string, authHeader string) Client {
	if !strings.HasSuffix(host, "/") {
		host = host + "/"
	}
	return Client{
		ots:     pkg.InitOvertimeClient(host, authHeader),
		APIHost: host,
	}
}

func printHoliday(w *tabwriter.Writer, a *pkg.Holiday) {
	fmt.Fprintf(w, "ID\t: %d\n", a.ID)
	fmt.Fprintf(w, "Description\t: %s\n", a.Description)
	fmt.Fprintf(w, "Type\t: %s\n", a.Type)
	fmt.Fprintf(w, "Start\t: %s\n", utils.FormatDay(a.Start))
	fmt.Fprintf(w, "End\t: %s\n", utils.FormatDay(a.End))
	diff := a.End.Sub(a.Start)
	fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(diff.Minutes())))
}

func formatMinutes(t int64) string {

	hs2, mf := math.Modf(float64(t) / 60)

	return fmt.Sprintf("%02dh:%02dm", int(hs2), int(mf*60))

}
func formatMinutesToHoursAndMinutes(t int64) string {
	hs2, mf := math.Modf(float64(t) / 60)
	return fmt.Sprintf("%02dh:%02dm", int(hs2), int(mf*60))
}
