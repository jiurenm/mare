package influxdb

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
)

type FormatType int8

const (
	JSON FormatType = iota
	CSV
)

type InfluxDB struct {
	Conn    *InfluxClient
	limiter *Limiter
}

var ErrNoData = fmt.Errorf("no data found")

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

func NewInfluxDB(cfg Config) (*InfluxDB, func(), error) {
	u, err := url.Parse(fmt.Sprintf("%s:%d/influxdb/v1", cfg.Host, cfg.Port))
	if err != nil {
		return nil, nil, err
	}

	u.Path = path.Join(u.Path, "write")
	params := u.Query()
	params.Set("db", "hypon")
	params.Set("rp", "")
	params.Set("precision", "s")
	params.Set("consistency", "")
	u.RawQuery = params.Encode()

	i := &InfluxDB{Conn: &InfluxClient{
		readHost:  fmt.Sprintf("%s:%d/rest/sql/%s", cfg.Host, cfg.Port, cfg.Database),
		writeHost: u.String(),
		auth:      "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.Username, cfg.Password))),
		username:  cfg.Username,
		password:  cfg.Password,
		client: &http.Client{
			Timeout: time.Minute,
		},
	}, limiter: NewLimiter(150, 1)}

	return i, i.Close, nil
}

// Close 关闭连接。
func (i *InfluxDB) Close() {
	i.Conn.client.CloseIdleConnections()
}

type InfluxClient struct {
	client    *http.Client
	readHost  string
	writeHost string
	baseURL   string
	auth      string
	username  string
	password  string
}

func (i *InfluxDB) Query(ctx context.Context, query string, dst interface{}, format ...FormatType) error {
	q := url.QueryEscape(query)
	u := i.Conn.readHost + q

	ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := i.limiter.Wait(ctx1)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx1, http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	if len(format) > 0 && format[0] == CSV {
		req.Header.Set("Accept", "application/csv")
	}

	resp, err := i.Conn.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(format) > 0 && format[0] == CSV {
		err = csvutil.Unmarshal(body, dst)

		if err != nil {
			return err
		}

		return nil
	}

	var res Response

	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}

	if len(res.Results) != 1 {
		return fmt.Errorf("unexpected response format")
	}

	results := res.Results[0]
	if len(results.Series) == 0 {
		return fmt.Errorf("no series found")
	}

	return handleSeries(results.Series, dst)
}

func handleSeries(series []Series, dst interface{}) error {
	var jsonStr []byte

	var err error

	var m []map[string]interface{}

	for i := 0; i < len(series); i++ {
		for j := 0; j < len(series[i].Values); j++ {
			mm := make(map[string]interface{}, len(series[i].Columns)+len(series[i].Tags))
			for k, v := range series[i].Tags {
				mm[k] = v
			}

			for k := 0; k < len(series[i].Columns); k++ {
				if series[i].Values[j][k] == nil {
					series[i].Values[j][k] = float64(0)
				}

				mm[series[i].Columns[k]] = series[i].Values[j][k]
			}

			m = append(m, mm)
		}
	}

	jsonStr, err = json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonStr, &dst)
}

func (i *InfluxDB) Delete(query string) error {
	uri := i.Conn.baseURL

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBufferString(query))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", i.Conn.auth)

	resp, err := i.Conn.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (i *InfluxDB) Query2(ctx context.Context, query string, dst interface{}, tz ...string) error {
	ctx1, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := i.limiter.Wait(ctx1)
	if err != nil {
		return err
	}

	uri := i.Conn.baseURL
	if len(tz) > 0 && tz[0] != "" {
		uri = uri + "?tz=" + tz[0]
	}

	method := http.MethodPost
	payload := strings.NewReader(query)

	req, err := http.NewRequest(method, uri, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", i.Conn.auth)

	resp, err := i.Conn.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	_, err = io.Copy(body, resp.Body)
	if err != nil {
		return err
	}

	var res TDResponse
	err = json.Unmarshal(body.Bytes(), &res)
	if err != nil {
		return err
	}

	if res.Code != 0 {
		if strings.Contains(res.Desc, "Table does not exist") {
			return ErrNoData
		}

		return fmt.Errorf(res.Desc)
	}

	if res.Rows == 0 {
		if strings.Contains(res.Desc, "Table does not exist") {
			return ErrNoData
		}

		return ErrNoData
	}

	return handleSeries2(res, dst)
}

type Response struct {
	Results []Result `json:"results"`
}

type Result struct {
	Series      []Series `json:"series"`
	StatementID int      `json:"statement_id"`
}

type Series struct {
	Name    string            `json:"name"`
	Tags    map[string]string `json:"tags"`
	Columns []string          `json:"columns"`
	Values  [][]interface{}   `json:"values"`
}

type TDResponse struct {
	Code       int              `json:"code"`
	Desc       string           `json:"desc"`
	ColumnMeta [][3]interface{} `json:"column_meta"`
	Data       [][]interface{}  `json:"data"`
	Rows       int              `json:"rows"`
}

func handleSeries2(series TDResponse, dst any) error {
	head := make([]string, len(series.ColumnMeta))

	for i := 0; i < len(series.ColumnMeta); i++ {
		head[i] = series.ColumnMeta[i][0].(string)
	}

	res := make([]map[string]any, 0, len(series.Data))

	for i := 0; i < len(series.Data); i++ {
		data := series.Data[i]
		val := make(map[string]any, len(head))

		for j := 0; j < len(head); j++ {
			if data[j] == nil {
				val[head[j]] = 0
			} else {
				val[head[j]] = data[j]
			}
		}

		res = append(res, val)
	}

	jsonStr, err := json.Marshal(res)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonStr, &dst)
}
