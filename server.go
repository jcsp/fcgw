package fcgw

/**
 The server wraps the influxdb httpd server, and instead of
 passing input data on to a database, passes it through to
 SeriesTable.
 */

import (
	"github.com/influxdb/influxdb/services/httpd"
	"github.com/influxdb/influxdb/cluster"
	"github.com/influxdb/influxdb/meta"
	"time"
)

type InfluxServer struct {
	service *httpd.Service
	series_table *SeriesTable
}

func NewInfluxServer(st *SeriesTable) *InfluxServer {
	is := &InfluxServer{
		series_table: st,
	}

	c := httpd.Config{
		Enabled: true,
		BindAddress: "0.0.0.0:8086",
		AuthEnabled: false,
		LogEnabled: true,
		WriteTracing: false,
		PprofEnabled: false,
		HTTPSEnabled: false,
		HTTPSCertificate: "",
	}
	is.service = httpd.NewService(c)
	is.service.Handler.PointsWriter = NewPointsConverter(st)
	is.service.Handler.Version = "embedded"
	is.service.Handler.MetaStore = MetaStore{}
	//is.service.Handle.QueryExecutor =

	return is
}

type MetaStore struct {}

func (ms MetaStore) WaitForLeader(timeout time.Duration) error {
	return nil
}

func (ms MetaStore) Database(name string) (*meta.DatabaseInfo, error) {
	di := &meta.DatabaseInfo{
		Name: "dummy",
	}
	return di, nil
}

func (ms MetaStore) Authenticate(username, password string) (ui *meta.UserInfo, err error) {
	return nil, nil
}

func (ms MetaStore) Users() ([]meta.UserInfo, error) {
	return nil, nil
}

// Define a PointsWriter class that takes a WritePointsRequest from
// the influxdb server code, and sends it in our simplfied form
// up to SeriesTable

type PointsConverter struct {
 st *SeriesTable
}

func NewPointsConverter(st *SeriesTable) *PointsConverter {
 pc := &PointsConverter{
  st: st,
 }

 return pc
}

func (pc *PointsConverter) WritePoints(p *cluster.WritePointsRequest) error {
	pc.st.Logger.Println("WritePoints")
	for _, point := range(p.Points) {
		for field_name, field := range(point.Fields()) {
			// Whereas influx has time series
			// with multiple fields, we define
			// a time series as a single value series.
			// Append field to name to make the conversion
			name := point.Name()
			name += "." + field_name

//			if field.(type) == float64 {
//
//			}

			var v float64
			switch field.(type) {
				case float64:
					v, _ = field.(float64)
				case uint32:
					vtmp, _ := field.(uint32)
					v = float64(vtmp)
				case uint64:
					vtmp, _ := field.(uint64)
					v = float64(vtmp)
				case float32:
					vtmp, _ := field.(float32)
					v = float64(vtmp)
				default:
					continue
			}

			// Work out a series key
			sk := NewSeries(SeriesName(name), SeriesTags(point.Tags())  )
			key := sk.Key()
			// Look up the series object in seriestable
			series := pc.st.Series(key)
			if series == nil {
				pc.st.Insert(sk)
				series = sk
			}
			// append the data
			series.Append(TimeVal(point.Time().Unix()), v)
		}
	}

	return nil
}

func (is *InfluxServer) Open() error {
	return is.service.Open()
}


func (is *InfluxServer) Close() error {
	return is.service.Close()
}
