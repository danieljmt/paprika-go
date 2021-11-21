package paprika

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

func (r *Recipe) compress() ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed marshaling recipe `%s`: %w", r.Name, err)
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close()
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decompress(reader io.Reader) (*Recipe, error) {
	r := &Recipe{}
	gr, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gr.Close()

	if err := json.NewDecoder(gr).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
