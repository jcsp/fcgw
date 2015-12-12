package fcgw

/**
 A raw series stores (in memory) uncompressed time/value points
 */


import (
    "github.com/dgryski/go-tsz"
)

type TimeVal uint32
type SeriesTags map[string]string
type SeriesKey string
type SeriesName string

type RawPoint struct {
 t TimeVal
 v float64
}

type Series struct {
 tags SeriesTags
 name SeriesName

 // Raw data points
 current []RawPoint

 // Archive (compressed hours)
 archive map[TimeVal]tsz.Series
}


func NewSeries(name SeriesName, tags SeriesTags) *Series {
 s := &Series{
  name: name,
  tags: tags,
  current: make([]RawPoint, 720),
  archive: make(map[TimeVal]tsz.Series),
 }

 return s
}

func (s *Series) Key() SeriesKey {
 // Series are identified by their name plus their tags
 // For lookups we generate a concatenation
 var key string
 key = string(s.name)

 key += ","
 for k, v := range(s.tags) {
  // todo be more efficient
  key += k
  key += ":"
  key += v
  key += ","
 }

 return SeriesKey(key)
}


func (s *Series) Append(t TimeVal, v float64) {

 // Have we advanced hour?  Time to compress and archive

}
