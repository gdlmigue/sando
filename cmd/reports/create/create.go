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
		err := createClient(params.startDate, params.endDate)
		if err != nil {
			return err
		}
		return nil
	}()

	cmdutil.ExitIfError(err)

	cmdutil.Success("Report created \n")
}

func createClient(startDate, endDate string) error {
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
		req, err := http.NewRequest(http.MethodGet, "/reporting-api/v1/reports/bytes-by-cpcode/versions/1/report-data", nil)
		if err != nil {
			return err
		}
		q := req.URL.Query()
		q.Add("start", startDate)
		q.Add("end", endDate)
		q.Add("objectIds", rec[0])
		req.URL.RawQuery = q.Encode()

		edgerc.SignRequest(req)

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Wrong credentials for fetching data or akamai error, either way save it for later
		if resp.StatusCode == http.StatusForbidden {
			cpcode := rec[0]
			row := []string{cpcode}
			missingData = append(missingData, row)
			continue
		}

		var report Report
		err = json.NewDecoder(resp.Body).Decode(&report)
		if err != nil {
			return err
		}

		row := []string{report.Data[0].Cpcode, report.Data[0].OriginBytes, report.Data[0].EdgeBytes, report.Data[0].MidgressBytes, report.Data[0].BytesOffload}
		data = append(data, row)
	}

	fileReport, err := os.Create("report.csv")
	if err != nil {
		return err
	}
	defer fileReport.Close()
	w := csv.NewWriter(fileReport)
	defer w.Flush()
	w.WriteAll(data)

	missingReport, err := os.Create("report_missing.csv")
	if err != nil {
		return err
	}
	defer missingReport.Close()
	wm := csv.NewWriter(missingReport)
	defer wm.Flush()

	wm.WriteAll(missingData)

	return nil
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
