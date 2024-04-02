package influxdb

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
)

type Writable interface {
	Measurement() string
	Tags() []byte
	Fields() []byte
	Timestamp() int64
}

func (i *InfluxDB) Write(data []Writable) error {
	if len(data) == 0 {
		return nil
	}

	b := make([]byte, 0, len(data)*1024)
	for i := 0; i < len(data); i++ {
		b = append(b, data[i].Measurement()...)
		b = append(b, ',')
		b = append(b, data[i].Tags()...)
		b = append(b, ' ')
		b = append(b, data[i].Fields()...)
		b = append(b, ' ')
		b = strconv.AppendInt(b, data[i].Timestamp(), 10)
		b = append(b, '\n')
	}

	req, err := http.NewRequest(http.MethodPost, i.Conn.writeHost, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "")
	if i.Conn.username != "" {
		req.SetBasicAuth(i.Conn.username, i.Conn.password)
	}

	resp, err := i.Conn.client.Do(req)
	if err != nil {
		return err
	}

	b = nil
	defer resp.Body.Close()

	body := new(bytes.Buffer)
	_, err = io.Copy(body, resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		var err = errors.New(body.String())
		return err
	}

	return nil
}
