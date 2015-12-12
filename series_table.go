package fcgw

/**
 The series table
 */

import (
 "log"
 "os"
)


type SeriesTable struct {
 by_key map[SeriesKey]*Series
 Logger *log.Logger
}

func NewSeriesTable() *SeriesTable {
 return &SeriesTable{
  by_key: make(map[SeriesKey]*Series),
  Logger: log.New(os.Stderr, "[SeriesTable] ", log.LstdFlags),
 }
}

func (st *SeriesTable) Series(key SeriesKey) *Series {
 val, ok := st.by_key[key]
 if ok == false {
  return val
 } else {
  return nil
 }
}

func (st *SeriesTable) Insert(series *Series) error {
 st.by_key[series.Key()] = series
 return nil
}
