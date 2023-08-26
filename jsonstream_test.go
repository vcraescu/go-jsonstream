package jsonstream_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
	"github.com/vcraescu/go-jsonstream"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const jsonFilename = "testdata/large.json"

type Payload struct {
	Result bool    `json:"result"`
	Data   []Datum `json:"data"`
}

type Datum struct {
	Field1  string `json:"field1"`
	Field2  string `json:"field2"`
	Field3  string `json:"field3"`
	Field4  string `json:"field4"`
	Field5  string `json:"field5"`
	Field6  string `json:"field6"`
	Field7  string `json:"field7"`
	Field8  string `json:"field8"`
	Field9  string `json:"field9"`
	Field10 string `json:"field10"`
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()

	b := mustGenerateJSON(5)
	entries, err := jsonstream.Unmarshal[Datum](
		context.Background(),
		bytes.NewReader(b),
		jsonstream.WithStartFrom(5),
		jsonstream.WithBatchSize(1),
	)
	require.NoError(t, err)

	var want Payload
	var got []Datum

	err = json.Unmarshal(b, &want)
	require.NoError(t, err)

	for entry := range entries {
		require.NoError(t, entry.Err)
		got = append(got, entry.Value)
	}

	require.Equal(t, want.Data, got)
}

func BenchmarkStdUnmarshal(b *testing.B) {
	printMemUsage()
	//mustGenerateJSONFile(jsonFilename, 100_000)

	b.ResetTimer()
	b.ReportAllocs()

	f, err := os.Open(jsonFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	js := Payload{}

	if err := json.NewDecoder(f).Decode(&js); err != nil {
		panic(err)
	}

	b.StopTimer()
	printMemUsage()
}

func BenchmarkUnmarshal(b *testing.B) {
	printMemUsage()
	//mustGenerateJSONFile(jsonFilename, 100_000)

	b.ResetTimer()
	b.ReportAllocs()

	f, err := os.Open(jsonFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	entries, err := jsonstream.Unmarshal[Datum](
		context.Background(),
		f,
		jsonstream.WithStartFrom(5),
		jsonstream.WithBatchSize(100),
	)
	if err != nil {
		panic(err)
	}

	for entry := range entries {
		if entry.Err != nil {
			panic(entry.Err)
		}
	}

	b.StopTimer()
	printMemUsage()
}

func mustGenerateJSONFile(filename string, n int) {
	_ = os.MkdirAll(filepath.Dir(filename), 0755)
	b := mustGenerateJSON(n)

	if err := os.WriteFile(filename, b, 0755); err != nil {
		panic(err)
	}
}

func mustGenerateJSON(n int) json.RawMessage {
	entries := make([]Datum, 0, n)

	for i := 0; i < n; i++ {
		entries = append(entries, Datum{
			Field1:  gofakeit.BookTitle(),
			Field2:  gofakeit.BookTitle(),
			Field3:  gofakeit.BookTitle(),
			Field4:  gofakeit.BookTitle(),
			Field5:  gofakeit.BookTitle(),
			Field6:  gofakeit.BookTitle(),
			Field7:  gofakeit.BookTitle(),
			Field8:  gofakeit.BookTitle(),
			Field9:  gofakeit.BookTitle(),
			Field10: gofakeit.BookTitle(),
		})
	}

	f := Payload{
		Result: true,
		Data:   entries,
	}

	b, err := json.Marshal(f)
	if err != nil {
		panic(err)
	}

	return b
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
