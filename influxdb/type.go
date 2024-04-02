package influxdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

var tzMap = map[string]*time.Location{}

type Time time.Time

func (t *Time) UnmarshalCSV(data []byte) error {
	if bytes.Equal(data, []byte("0")) {
		*t = Time(time.Now())

		return nil
	}

	msec, err := strconv.Atoi(string(data[0:10]))
	if err != nil {
		return err
	}

	*t = Time(time.Unix(int64(msec), 0))

	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", t.ToRFC3339String())

	return []byte(formatted), nil
}

func (t *Time) ToString(timezone ...string) string {
	if len(timezone) > 0 && timezone[0] != "" {
		loc, ok := tzMap[timezone[0]]
		if !ok {
			loc, _ = time.LoadLocation(timezone[0])
			tzMap[timezone[0]] = loc
		}

		return time.Time(*t).In(loc).Format("2006-01-02 15:04:05")
	}

	return time.Time(*t).Format("2006-01-02 15:04:05")
}

func (t *Time) ToRFC3339String(timezone ...string) string {
	if len(timezone) > 0 && timezone[0] != "" {
		loc, ok := tzMap[timezone[0]]
		if !ok {
			loc, _ = time.LoadLocation(timezone[0])
			tzMap[timezone[0]] = loc
		}

		return time.Time(*t).In(loc).Format(time.RFC3339)
	}

	return time.Time(*t).Format(time.RFC3339)
}

type Tags struct {
	InvSN   string `json:"invsn,omitempty" excel:"invsn" redis:"Gatewaysn"`
	Gateway string `json:"gatewaysn,omitempty" excel:"gateway" redis:"Invsn"`
	BatSN   string `json:"sn,omitempty" excel:"sn"`
}

var (
	TagGateway  = []byte("gatewaysn")
	TagInverter = []byte("invsn")
	TagBattery  = []byte("sn")
)

func (t *Tags) UnmarshalCSV(data []byte) error {
	res := bytes.Split(data, []byte(","))
	for i := 0; i < len(res); i++ {
		switch {
		case bytes.HasPrefix(res[i], TagGateway):
			t.Gateway = string(res[i][10:])
		case bytes.HasPrefix(res[i], TagInverter):
			t.InvSN = string(res[i][6:])
		case bytes.HasPrefix(res[i], TagBattery):
			t.BatSN = string(res[i][3:])
		}
	}

	return nil
}

type Int int

func (i *Int) UnmarshalCSV(data []byte) error {
	if string(data) == "" {
		*i = 0
	} else {
		n, err := strconv.ParseInt(string(data), 10, 64)
		if err != nil {
			return err
		}
		*i = Int(n)
	}

	return nil
}

type Uint8 uint8

func (i *Uint8) UnmarshalCSV(data []byte) error {
	if string(data) == "" {
		*i = 0
	} else {
		n, err := strconv.ParseUint(string(data), 10, 8)
		if err != nil {
			return err
		}
		*i = Uint8(n)
	}

	return nil
}

type Float64 float64

func (f *Float64) UnmarshalCSV(data []byte) error {
	if string(data) == "" {
		*f = 0
	} else {
		m, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			return err
		}
		*f = Float64(m)
	}

	return nil
}

func (f Float64) MarshalJSON() ([]byte, error) {
	if float64(f) == math.NaN() {
		return json.Marshal(0)
	}

	value, err := strconv.ParseFloat(strconv.FormatFloat(float64(f), 'f', 2, 64), 64)
	if err != nil {
		return nil, err
	}

	return json.Marshal(value)
}

func (f *Float64) UnmarshalJSON(value []byte) error {
	if len(value) == 0 {
		*f = 0

		return nil
	}

	m, err := strconv.ParseFloat(string(value), 64)
	if err != nil {
		return err
	}

	*f = Float64(m)

	return nil
}
