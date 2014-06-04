package csv

import (
	"errors"
	"github.com/azylman/optimus"
	"github.com/azylman/optimus/sources/csv"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func readFile(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	return lines, nil
}

func TestCSVSink(t *testing.T) {
	source := csv.New("./data.csv")
	err := New(source, "./data_write.csv")
	assert.Nil(t, err)
	expected, err := readFile("./data_write.csv")
	assert.Nil(t, err)
	actual, err := readFile("./data.csv")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

type errorTable struct {
	rows chan optimus.Row
}

func (e errorTable) Err() error {
	return errors.New("failed")
}

func (e errorTable) Rows() <-chan optimus.Row {
	return e.rows
}

func (e errorTable) Stop() {}

func newErrorTable() optimus.Table {
	table := &errorTable{rows: make(chan optimus.Row)}
	close(table.rows)
	return table
}

func TestCSVSinkError(t *testing.T) {
	source := newErrorTable()
	err := New(source, "./data_write.csv")
	assert.Equal(t, err, errors.New("failed"))
}
