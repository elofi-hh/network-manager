package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"modernc.org/ql"
)

/*
#mac,ip,iface,in,out,total,first_date,last_date
a0:ce:c8:b0:e8:eb,192.168.2.10,br-lan,0,0,0,12-10-2024_06:50:03,12-10-2024_06:50:03
c6:f8:85:ce:62:d6,192.168.2.101,br-lan,0,0,0,12-10-2024_06:50:03,12-10-2024_06:50:03
f2:64:c7:89:4c:42,192.168.2.110,br-lan,468,0,468,12-10-2024_06:50:03,12-10-2024_06:50:59
*/

type DeviceDetailEntry struct {
  MAC string `json:"mac"`
  IP string `json:"ip"`
  Interface string `json:"interface"`
  In int64 `json:"data_in"`
  Out int64 `json:"data_out"`
  Total int64 `json:"data_total"`
  FirstDateUnix int64 `json:"first_date_unix"`
  LastDateUnix int64 `json:"last_date_unix"`
}

func parseDetailDate(s string) (*int64, error) {
  if !strings.Contains(s, "_") {
    t, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("invalid time format: %v", s))
    }

    return &t, nil
  }

  sp := strings.Split(s, "_")
  date := sp[0]
  timeg := sp[1]

  dateS := strings.Split(date, "-")
  day := dateS[0]
  dayI, err := strconv.Atoi(string(day))
  if err != nil {
    return nil, err
  }

  month := dateS[1]
  monthI, err := strconv.Atoi(string(month))
  if err != nil {
    return nil, err
  }

  year := dateS[2]
  yearI, err := strconv.Atoi(string(year))
  if err != nil {
    return nil, err
  }

  timeS := strings.Split(timeg, ":")
  hours := timeS[0]
  hoursI, err := strconv.Atoi(string(hours))
  if err != nil {
    return nil, err
  }

  mins := timeS[1]
  minsI, err := strconv.Atoi(string(mins))
  if err != nil {
    return nil, err
  }

  sec := timeS[2]
  secI, err := strconv.Atoi(string(sec))
  if err != nil {
    return nil, err
  }

  t := time.Date(yearI, time.Month(monthI), dayI, hoursI, minsI, secI, 0, time.Local).Unix()
  return &t, nil
}

func DeviceDetailEntryFromRowInterface(row []interface{}) (*DeviceDetailEntry, error) {
  if len(row) != 8 {
    return nil, errors.New("row is not of length 8")
  }

  var stringRow [8]string

  for i, v := range row { 
    stringRow[i] = fmt.Sprintf("%v", v)
  }

  return DeviceDetailEntryFromRow(stringRow[:])
}

func DeviceDetailEntryFromRow(row []string) (*DeviceDetailEntry, error) {
  if len(row) != 8 {
    return nil, errors.New("row is not of length 8")
  }

  in, err := strconv.ParseInt(row[3], 10, 64)
  if err != nil {
    return nil, err
  }

  out, err := strconv.ParseInt(row[4], 10, 64)
  if err != nil {
    return nil, err
  }

  tot, err := strconv.ParseInt(row[5], 10, 64)
  if err != nil {
    return nil, err
  }

  firstDate, err := parseDetailDate(row[6])
  if err != nil {
    return nil, err
  }

  lastDate, err := parseDetailDate(row[7])
  if err != nil {
    return nil, err
  }

  return &DeviceDetailEntry{
    MAC: row[0],
    IP: row[1],
    Interface: row[2],
    In: in,
    Out: out,
    Total: tot,
    FirstDateUnix: *firstDate,
    LastDateUnix: *lastDate,
  }, nil
}

type NetworkManager struct {
  networkDetailDumperGroup sync.WaitGroup
  networkDumpInterval time.Duration

  networkDumpChan chan []DeviceDetailEntry

  kickChan chan string
  kickMutex sync.Mutex
}

func (rm *NetworkManager) kickDevice(mac string) (error) {
  // not 100% sure if we need to lock this but just incase for now
  rm.kickMutex.Lock()
  ipsetCmd := exec.Command("ipset", "add", "mac-ban", mac)
  err := ipsetCmd.Run()
  if err != nil {
    return err
  }
  cmd := exec.Command("hostapd_cli", "disassociate", mac)
  err = cmd.Run()
  if err != nil {
    return err
  }

  cmd = exec.Command("hostapd_cli", "deauthenticate", mac)
  err = cmd.Run()
  if err != nil {
    return err
  }
  rm.kickMutex.Unlock()
  return nil
}

func (rm *NetworkManager) dump() (error) {
  // dump raw.db
  cmd := exec.Command("wrtbwmon", "update", "raw.db")
  err := cmd.Run()
  if err != nil {
    return err
  }

  // diff with last.db
  s, err := os.ReadFile("raw.db")
  if err != nil {
    return err
  }

  r := csv.NewReader(strings.NewReader(string(s)))

  _, err = r.Read()
  if err != nil {
    return err
  }

  entries := []DeviceDetailEntry{}

  for {
    row, err := r.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      return err
    }

    d, err := DeviceDetailEntryFromRow(row)
    if err != nil {
      return err
    }

    entries = append(entries, *d) 
  }

  log.Printf("[network] took snapshot: %v", entries)

  // write to db
  err = database.InsertDeviceNetworkData(entries)
  if err != nil {
    return err
  }

  // publish
  rm.networkDumpChan <- entries

  return nil
}

func (rm *NetworkManager) listen() {
  for {
    mac := <- rm.kickChan

    err := rm.kickDevice(mac)
    if err != nil {
      log.Printf("error kicking device (%v): %v", mac, err)
    }

    log.Printf("[network] kicked %v", mac)
  }
}

func (nm *NetworkManager) start() {
  // start dumping details
  // start receiving on channl to kick people off
  go nm.listen()

  for {
    err := nm.dump()
    if err != nil {
      log.Printf("error dumping: %v", err)
    }

    time.Sleep(nm.networkDumpInterval)
  }
}

func NewNetworkManager() (*NetworkManager) {
  return &NetworkManager{
    networkDumpInterval: time.Second * 3,
    networkDumpChan: make(chan []DeviceDetailEntry),
    kickChan: make(chan string),
  }
}

var networkManager *NetworkManager

type Database struct {
  db *ql.DB
  ctx *ql.TCtx

  // how many dumps do we want to keep in the DB at once (anything older will be purged to save space)
  dataFrameSize int
}

func NewDatabase() (*Database, error) {
  // read in the last entry id from db
  frameSize := 120
  file := "db.db"

  tableCreateCmd := `
  BEGIN TRANSACTION;
  CREATE TABLE IF NOT EXISTS network_data (
    entry_id int64,
    mac string,
    ip string,
    interface string,
    data_int int64,
    data_out int64,
    data_total int64,
    first_date_unix int64,
    last_date_unix int64,
  );
  CREATE UNIQUE INDEX IF NOT EXISTS xnetwork_data ON network_data (entry_id, mac);
  COMMIT;
  `

  db, err := ql.OpenFile(file, &ql.Options{
    CanCreate: true,
  })
  if err != nil {
    return nil, err
  }

  ctx := ql.NewRWCtx()

  if _, _, err := db.Run(ctx, tableCreateCmd); err != nil {
    return nil, err
  } 
  return &Database{
    db: db,
    ctx: ctx,
    dataFrameSize: frameSize,
  }, nil
}

func (db *Database) GetData() (*[][]DeviceDetailEntry, error) {
  r, _, err := db.db.Run(db.ctx, "SELECT * FROM network_data")
  if err != nil {
    return nil, err
  }

  rows, err := r[0].Rows(-1, 0)
  if err != nil {
    return nil, err
  }

  var size int64 = 0
  var minimum int64 = 0
  for _, v := range rows {
    if val, ok  := v[0].(int); ok {
      size = max(size, int64(val))
      minimum = min(minimum, int64(val))
    } else if val, ok := v[0].(int64); ok {
      size = max(size, val)
      minimum = min(minimum, val)
    } else {
      return nil, errors.New(fmt.Sprintf("entry_id is not an int or int64: %v", v[0]))
    }
  }

  data := make([][]DeviceDetailEntry, size - minimum + 1)

  for _, v := range rows {
    entry, err := DeviceDetailEntryFromRowInterface(v[1:])
    if err != nil {
      return nil, err
    }

    var id int64
    if val, ok  := v[0].(int); ok {
      id = int64(val)
    } else if val, ok := v[0].(int64); ok {
      id = val
    } else {
      return nil, errors.New(fmt.Sprintf("2 entry_id is not an int or int64: %v", v[0]))
    }
    
    data[id - minimum] = append(data[id - minimum], *entry)
  } 

  return &data, nil
}

func (db *Database) InsertDeviceNetworkData(entries []DeviceDetailEntry) (error) {
  // insert new data
  r, _, err := db.db.Run(db.ctx, "SELECT max(entry_id) FROM network_data")
  if err != nil {
    return err
  }

  fr, err := r[0].FirstRow()
  if err != nil {
    return err
  }

  val := fr[0]

  var id int64
  if val == nil {
    id = 0
  } else if v, ok := val.(int); ok {
    id = int64(v) + 1
  } else if v, ok := val.(int64); ok {
    id = v + 1
  } else {
    return errors.New(fmt.Sprintf("invalid response for entry_id: %v", fr))
  }

  log.Printf("[database] inserting id: %v", id)

  var buffer bytes.Buffer
  buffer.WriteString("BEGIN TRANSACTION;\n")
  buffer.WriteString(`INSERT INTO network_data VALUES`)
  for i, entry := range entries {
    if i != 0 {
      buffer.WriteString(`,`)
    }

    buffer.WriteString(fmt.Sprintf(` (%v, "%v", "%v", "%v", %v, %v, %v, %v, %v)`, id, entry.MAC, entry.IP, entry.Interface, entry.In, entry.Out, entry.Total, entry.FirstDateUnix, entry.LastDateUnix))
  }
  buffer.WriteString(";\n")
  buffer.WriteString(`COMMIT;`)
  _, _, err = db.db.Run(db.ctx, buffer.String())
  if err != nil {
    return err
  }

  //log.Printf("successfully inserted %v", entries)

  old := id - int64(db.dataFrameSize)

  // drop data older than
  _, _, err = db.db.Run(db.ctx, fmt.Sprintf(`
    BEGIN TRANSACTION;
      DELETE FROM network_data WHERE entry_id < %v;
    COMMIT;
    `, old))
  if err != nil {
    return err
  }

  //log.Printf("successfully purged old data (before %v)", old)
  return nil
}

var database *Database

type Backend struct {}

func (b *Backend) handler(w http.ResponseWriter, r *http.Request) {
  d, err := database.GetData()
  if err != nil {
    log.Printf("error getting data: %v", err)
  }
  json, err := json.Marshal(*d)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Write(json)
}

func (b *Backend) start() (error) {
  r := http.NewServeMux()
  r.HandleFunc("/", b.handler)

  err := http.ListenAndServe(":8080", r)
  if err != nil {
    return err
  }

  return nil
}

func NewBackend() (*Backend, error) {
  return &Backend{}, nil
}

var backend *Backend

type Reporter struct {
  threshold float64
}

func NewReporter() (*Reporter, error) {
  r := &Reporter{
    threshold: 0.5,
  }

  go r.ConsumeTrafficAnalysisResults()

  return r, nil
}

func (r *Reporter) ConsumeTrafficAnalysisResults() {
  for {
    results := <- trafficAnalyzer.analysisResultsChan

    // report to the blockchain here

    // report to network manager here

    log.Printf("[reporter] %v", results)
  } 
}

var reporter *Reporter

type TrafficAnalysisResults struct {
  MAC string
  abuserRating float64
}

type TrafficAnalyzer struct {
  analysisResultsChan chan TrafficAnalysisResults
}

func NewTrafficAnalyzer() (*TrafficAnalyzer, error) {
  ta := &TrafficAnalyzer{
    analysisResultsChan: make(chan TrafficAnalysisResults),
  }

  go ta.ConsumeNetworkDump()

  return ta, nil
}

func (ta *TrafficAnalyzer) ConsumeNetworkDump() {
  for {
    entries := <- networkManager.networkDumpChan

    log.Printf("[traffic-analyzer] received network snapshot %v", entries)

    // run analysis here to calc abuser rating

    for _, v := range entries {
      ta.analysisResultsChan <- TrafficAnalysisResults{
        MAC: v.MAC,
        abuserRating: 0.1,
      }
    } 
  }
}

var trafficAnalyzer *TrafficAnalyzer

func main() {
  var err error
  database, err = NewDatabase()
  if err != nil {
    log.Panicf("error initializing db: %v", err)
    return
  }

  networkManager = NewNetworkManager()
  go networkManager.start() 

  trafficAnalyzer, err = NewTrafficAnalyzer()
  if err != nil {
    log.Printf("failed to start traffic analyzer: %v", err)
    return
  }

  reporter, err = NewReporter()
  if err != nil {
    log.Printf("failed to start reporter: %v", err)
    return
  }

  backend, err = NewBackend()
  err = backend.start()
  if err != nil {
    log.Printf("failed to start backend: %v", err)
    return
  }
}

func handleKeyboardInterrupt() {
  // nicely close db here
}
