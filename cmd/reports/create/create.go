package create

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"os"
	"sando/internal/cmdcommon"
	"sando/internal/cmdutil"
	"sando/internal/query"

	ec "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgegrid"
	"github.com/spf13/cobra"
)

const (
	helpText = `Create a bandwidth report for akamai`
	examples = `$ sando reports create
	#Create a bandwidth report for akamai
	$ sando reports create -s 2023-09-01T00:00:00Z -e 2023-09-02T00:00:00Z`
)

func NewCmdCreate() *cobra.Command {
	return &cobra.Command{
		Use:     "create",
		Short:   "Create an akamai bw report",
		Long:    helpText,
		Example: examples,
		Run:     create,
	}
}

func SetFlags(cmd *cobra.Command) {
	cmdcommon.SetCreateReportFlags(cmd)
}

func create(cmd *cobra.Command, _ []string) {
	params := parseFlags(cmd.Flags())
	err := func() error {
		s := cmdutil.Info("Creating report...")
		defer s.Stop()
		err := generateReport(params.startDate, params.endDate)
		if err != nil {
			return err
		}
		return nil
	}()

	cmdutil.ExitIfError(err)

	cmdutil.Success("Report created \n")
}

func generateReport(startDate, endDate string) error {
	filePath := "/Users/sjimenez/Documents/sando/test.csv"
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}

	edgerc := ec.Must(ec.New(ec.WithEnv(true)))

	client := http.Client{}

	var data [][]string
	var missingData [][]string

	for _, rec := range lines {
		report, err := fetchReport(startDate, endDate, rec[0], "", edgerc, client)
		if err != nil {
			return err
		}

		// Means that the cpcode could be under outside the main house contract
		if report == nil {
			// The accountSwitchKeys let us check in restricted contracts
			for _, acc := range accountSwitchKeys {
				report, err := fetchReport(startDate, endDate, rec[0], acc, edgerc, client)
				if err != nil {
					return err
				}
				if report == nil {
					continue
				}
				row := []string{report.Data[0].Cpcode, report.Data[0].OriginBytes, report.Data[0].EdgeBytes, report.Data[0].MidgressBytes, report.Data[0].BytesOffload}
				data = append(data, row)
				break
			}
		} else {
			row := []string{report.Data[0].Cpcode, report.Data[0].OriginBytes, report.Data[0].EdgeBytes, report.Data[0].MidgressBytes, report.Data[0].BytesOffload}
			data = append(data, row)
		}
	}

	// Write reports to CSV files
	if err := writeReportToCSV("report3.csv", data); err != nil {
		return err
	}

	if err := writeReportToCSV("report_missing3.csv", missingData); err != nil {
		return err
	}

	return nil

}

func fetchReport(startDate, endDate, objectId, accountSwitchKey string, edgerc *ec.Config, client http.Client) (*Report, error) {
	// Create HTTP request with the specified account switch key
	req, err := createRequest(startDate, endDate, objectId, accountSwitchKey, edgerc)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return nil, nil
	}

	var report Report
	err = json.NewDecoder(resp.Body).Decode(&report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func createRequest(startDate, endDate, objectId, accountSwitchKey string, edgerc *ec.Config) (*http.Request, error) {
	// Create HTTP request with the specified account switch key
	req, err := http.NewRequest(http.MethodGet, "/reporting-api/v1/reports/bytes-by-cpcode/versions/1/report-data", nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("start", startDate)
	q.Add("end", endDate)
	q.Add("objectIds", objectId)
	if accountSwitchKey != "" {
		q.Add("accountSwitchKey", accountSwitchKey)
	}
	req.URL.RawQuery = q.Encode()
	edgerc.SignRequest(req)
	return req, nil
}

func writeReportToCSV(filePath string, data [][]string) error {
	// Write data to CSV file
	fileReport, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fileReport.Close()

	w := csv.NewWriter(fileReport)
	defer w.Flush()

	if err := w.WriteAll(data); err != nil {
		return err
	}

	return nil
}

var accountSwitchKeys = []string{
	"B-3-QCCVOP:1-9OGH",
	"F-AC-835737:1-5G3LB",
	"AANA-73OX09:1-5G3LB",
	"F-AC-2444579:1-5G3LB",
	"1-10M0I:1-5G3LB",
	"149-2KXM:1-5G3LB",
	"1-36NFAF:1-2233J1",
	"1-36NFAF:1-9OGH",
	"F-AC-725166:1-5G3LB",
	"AANA-2WTYU4:1-5G3LB",
	"B-3-QCCVP5:1-9OGH",
	"B-3-QCCVK7:1-9OGH",
	"AANA-2X23BL:1-5G3LB",
	"F-AC-4902988:1-5G3LB",
	"1-KFKPC3:1-5G3LB",
	"B-3-QCCVNZ:1-9OGH",
}

type createParams struct {
	startDate string
	endDate   string
}

func parseFlags(flags query.FlagParser) *createParams {
	startDate, err := flags.GetString("startDate")
	cmdutil.ExitIfError(err)

	endDate, err := flags.GetString("endDate")
	cmdutil.ExitIfError(err)

	return &createParams{
		startDate: startDate,
		endDate:   endDate,
	}
}

type Report struct {
	Metadata          Metadata          `json:"metadata"`
	Data              []Datum           `json:"data"`
	SummaryStatistics SummaryStatistics `json:"summaryStatistics"`
}

type Datum struct {
	Cpcode        string `json:"cpcode"`
	BytesOffload  string `json:"bytesOffload"`
	EdgeBytes     string `json:"edgeBytes"`
	MidgressBytes string `json:"midgressBytes"`
	OriginBytes   string `json:"originBytes"`
}

type Metadata struct {
	Name               string        `json:"name"`
	Version            string        `json:"version"`
	OutputType         string        `json:"outputType"`
	GroupBy            []string      `json:"groupBy"`
	Interval           string        `json:"interval"`
	Start              string        `json:"start"`
	End                string        `json:"end"`
	AvailableDataEnds  string        `json:"availableDataEnds"`
	SuggestedRetryTime interface{}   `json:"suggestedRetryTime"`
	RowCount           int64         `json:"rowCount"`
	Filters            []interface{} `json:"filters"`
	Columns            []Column      `json:"columns"`
	ObjectType         string        `json:"objectType"`
	ObjectIDS          []string      `json:"objectIds"`
}

type Column struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type SummaryStatistics struct {
}
