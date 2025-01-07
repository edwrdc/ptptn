package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

const UJRAH_RATE = 0.01 // 1% per annum

type LoanSummary struct {
	Principal              float64
	RepaymentPeriodInMonth int
	MonthlyPayment         float64
	TotalPayment           float64
	TotalInterest          float64
}

func main() {
	app := tview.NewApplication()
	principalInput := tview.NewInputField().
		SetLabel("Your total principal amount. (Jumlah hutang): RM").
		SetFieldWidth(10).
		SetAcceptanceFunc(tview.InputFieldFloat).
		SetDoneFunc(func(key tcell.Key) {
			app.Stop()
		})

	output := tview.NewTextView().SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	principalInput.SetDoneFunc(func(key tcell.Key) {
		//output.SetText(principalInput.GetText())
		principalAmount, err := strconv.ParseFloat(principalInput.GetText(), 64)
		if err != nil {
			panic(err)
		}
		loan := &LoanSummary{
			Principal: principalAmount,
		}
		output.SetText(loan.Summary())
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(principalInput, 0, 1, true).
		AddItem(output, 0, 1, false)
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func (l *LoanSummary) calculateRepaymentPeriod() {
	// https://www.ptptn.gov.my/en/bayaran-balik/
	switch {
	case l.Principal <= 10000:
		l.RepaymentPeriodInMonth = 60
	case l.Principal <= 22000:
		l.RepaymentPeriodInMonth = 120
	case l.Principal <= 50000:
		l.RepaymentPeriodInMonth = 180
	default:
		l.RepaymentPeriodInMonth = 240
	}
}

func (l *LoanSummary) calculatePayment() {
	l.calculateRepaymentPeriod()
	l.TotalInterest = l.Principal * UJRAH_RATE * (float64(l.RepaymentPeriodInMonth) / 12)
	l.TotalPayment = l.Principal + l.TotalInterest
	l.MonthlyPayment = l.TotalPayment / float64(l.RepaymentPeriodInMonth)
}

func (l *LoanSummary) Summary() string {
	//TODO: this would be the summary
	l.calculatePayment()

	return fmt.Sprintf(
		"PTPTN Loan Summary\n"+
			"─────────────────\n"+
			"Principal: RM %.2f\n"+
			"Period: %d months (%.0f years)\n"+
			"Monthly: RM %.2f\n"+
			"Interest: RM %.2f\n"+
			"Total: RM %.2f",
		l.Principal,
		l.RepaymentPeriodInMonth, float64(l.RepaymentPeriodInMonth)/12,
		l.MonthlyPayment,
		l.TotalInterest,
		l.TotalPayment)
}
