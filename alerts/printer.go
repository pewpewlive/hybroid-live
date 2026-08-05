package alerts

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Printer struct {
	alertsByFile map[string][]Alert
}

func NewPrinter() Printer {
	return Printer{
		alertsByFile: make(map[string][]Alert),
	}
}

func (p *Printer) GetAlerts(sourcePath string) []Alert {
	return p.alertsByFile[sourcePath]
}

func (p *Printer) AllAlerts() map[string][]Alert {
	return p.alertsByFile
}

func (p *Printer) StageAlerts(sourcePath string, alerts []Alert) {
	p.alertsByFile[sourcePath] = append(p.alertsByFile[sourcePath], alerts...)
}

func (p *Printer) PrintAlerts() error {
	warningsCount, errorsCount := 0, 0
	for sourcePath, fileAlerts := range p.alertsByFile {
		sort.Slice(fileAlerts, func(i, j int) bool {
			return fileAlerts[i].Span().Line < fileAlerts[j].Span().Line
		})

		alertMsg := strings.Builder{}
		for _, alert := range fileAlerts {
			switch alert.Type() {
			case Error:
				errorsCount++
				alertMsg.WriteString(fmt.Sprintf("[light_red][bold]error[%s]: [reset]", alert.ID()))
			case Warning:
				warningsCount++
				alertMsg.WriteString(fmt.Sprintf("[light_yellow][bold]warning[%s]: [default]", alert.ID()))
			}

			alertMsg.WriteString(fmt.Sprintf("[bold]%s[reset]\n", alert.Message()))

			span := alert.Span()
			alertMsg.WriteString(fmt.Sprintf("[light_gray]  --- %s:%d:%d ---[default]\n", sourcePath, span.Line, span.Column))

			if note := alert.Note(); note != "" {
				alertMsg.WriteString(fmt.Sprintf("[cyan]  = note:[default] %s\n", note))
			}

			fmt.Print(alertMsg.String() + "\n")
			alertMsg.Reset()
		}
	}

	extra := ""
	if warningsCount != 0 {
		extra = "[light_yellow]" + strconv.Itoa(warningsCount) + " warning(s)[white]"
	}
	if errorsCount != 0 {
		if warningsCount != 0 {
			extra += " and [light_red]" + strconv.Itoa(errorsCount) + " error(s)"
		} else {
			extra += "[light_red]" + strconv.Itoa(errorsCount) + " error(s)"
		}
	}
	if warningsCount == 0 && errorsCount == 0 {
		fmt.Printf("Compilation finished\n")
	} else {
		fmt.Printf("Compilation finished with %s\n", extra)
	}

	return nil
}
