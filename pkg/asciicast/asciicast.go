package asciicast

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Asciicast struct {
	w              io.Writer
	timestampStart time.Time
}

type header struct {
	Version   int   `json:"version"`
	Width     int   `json:"width"`
	Height    int   `json:"height"`
	Timestamp int64 `json:"timestamp"`
}

func NewAsciicastEncoder(w io.Writer, width, height int) *Asciicast {
	headerData := header{
		Version:   2,
		Width:     width,
		Height:    height,
		Timestamp: time.Now().Unix(),
	}
	headerJSON, _ := json.Marshal(headerData)
	w.Write(headerJSON)
	return &Asciicast{
		w:              w,
		timestampStart: time.Now(),
	}
}

func (a *Asciicast) Write(data []byte) (int, error) {
	return a.write("o", data)
}

func (a *Asciicast) WriteSize(width, height int) (int, error) {
	return a.write("r", []byte(fmt.Sprintf("%dx%d", width, height)))
}

func (a *Asciicast) write(code string, data []byte) (int, error) {
	ts := time.Since(a.timestampStart).Seconds()
	escapedData, _ := json.Marshal(string(data))
	_, err := fmt.Fprintf(a.w, "\n"+`[%.6f, "%s", %s]`, ts, code, escapedData)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}
