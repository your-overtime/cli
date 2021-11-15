package out

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/your-overtime/api/pkg"
	"github.com/your-overtime/cli/internal/utils"
)

func PrintActivity(a *pkg.Activity) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	printActivityWithWriter(w, a)
}

func PrintActivities(activities []pkg.Activity) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, '.', tabwriter.TabIndent)
	mins := 0
	now := time.Now()
	for i, a := range activities {
		fmt.Fprintf(w, "No\t: %d\n", i+1)
		printActivityWithWriter(w, &a)
		fmt.Fprintln(w)
		if a.End != nil {
			mins += int(a.End.Sub(*a.Start).Minutes())
		} else {
			mins += int(now.Sub(*a.Start).Minutes())
		}
	}
	fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(mins)))
	w.Flush()
}

func printActivityWithWriter(w *tabwriter.Writer, a *pkg.Activity) {
	fmt.Fprintf(w, "ID\t: %d\n", a.ID)
	fmt.Fprintf(w, "Description\t: %s\n", a.Description)
	fmt.Fprintf(w, "Start\t: %s\n", utils.FormatTime(*a.Start))
	if a.End != nil {
		fmt.Fprintf(w, "End\t: %s\n", utils.FormatTime(*a.End))
		diff := a.End.Sub(*a.Start)
		fmt.Fprintf(w, "Duration\t: %s\n", formatMinutes(int64(diff.Minutes())))
	}
	w.Flush()
}
