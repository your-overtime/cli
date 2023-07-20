package client

import "github.com/your-overtime/api/v2/pkg"

type ExportData struct {
	Activities []pkg.Activity
	Holidays   []pkg.Holiday
	WorkDays   []pkg.WorkDay
}
