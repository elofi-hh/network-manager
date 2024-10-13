package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"manager/store"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gammazero/deque"
	"modernc.org/ql"
)

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
  log.Printf("done -1")
  ipsetCreateCmd := exec.Command("ipset", "create", "mac-ban", "hash:mac", "-!")
  err := ipsetCreateCmd.Run()
  if err != nil {
    return err
  }

  log.Printf("done 0")

  ipsetCmd := exec.Command("ipset", "add", "mac-ban", mac)
  err = ipsetCmd.Run()
  if err != nil {
    return err
  }

  log.Printf("done 1")

  cmd := exec.Command("hostapd_cli", "disassociate", mac)
  err = cmd.Run()
  if err != nil {
    return err
  }

  log.Printf("done 2")

  cmd = exec.Command("hostapd_cli", "deauthenticate", mac)
  err = cmd.Run()
  if err != nil {
    return err
  }

  log.Printf("done 4")
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
  log.Printf("not chillin")
  rm.networkDumpChan <- entries
  log.Printf("chillin")
  return nil
}

func (rm *NetworkManager) listen() {
  for {
    mac := <- rm.kickChan
    networkManager.kickChan <- mac
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
  file := "bruh1.db"

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
  CREATE TABLE IF NOT EXISTS prompt (
    id int,
    prompt string,
    threshold int64,
    window int64,
  );
  CREATE UNIQUE INDEX IF NOT EXISTS xprompt ON prompt (id);
  COMMIT;
  `

  db, err := ql.OpenFile(file, &ql.Options{
    CanCreate: true,
  })
  if err != nil {
    log.Printf("bruh")
    return nil, err
  }

  ctx := ql.NewRWCtx()

  if _, _, err := db.Run(ctx, tableCreateCmd); err != nil {
    log.Printf("fuck")
    return nil, err
  } 
  return &Database{
    db: db,
    ctx: ctx,
    dataFrameSize: frameSize,
  }, nil
}

func (db *Database) GetOnboarded() (*bool, error) {
  r, _, err := db.db.Run(db.ctx, "SELECT * FROM prompt")
  if err != nil {
    return nil, err
  }

  rows, err := r[0].Rows(-1, 0)
  if err != nil {
    return nil, err
  }

  a := len(rows) == 1

  fmt.Printf("prompt: %v", rows)
  return &a, nil
}

type PromptDetails struct {
  prompt string
  threshold int64
  window int64
}

func (db *Database) GetPrompt() (*PromptDetails, error) {
  r, _, err := db.db.Run(db.ctx, `SELECT * FROM prompt WHERE id = 1`)
  if err != nil {
    return nil, err
  }

  rows, err := r[0].Rows(-1, 0)
  if err != nil {
    return nil, err
  }

  if len(rows) == 0 {
    return nil, errors.New("no prompt was set")
  }

  prompt := ""

  var threshold int64 = 30000
  var window int64 = 10

  if v, ok := rows[0][1].(string); ok {
    prompt = v
  }

  if v, ok := rows[0][2].(int); ok {
     threshold = int64(v)
  } else if v, ok := rows[0][2].(int64); ok {
    threshold = v
  }

  if v, ok := rows[0][3].(int); ok {
    window = int64(v)
  } else if v, ok := rows[0][3].(int64); ok {
    window = v
  }

  p := &PromptDetails{
    prompt: prompt,
    threshold: threshold,
    window: window,
  }
  return p, nil
}

func (db *Database) SetPrompt(prompt string, threshold int64, window int64) (error) {
  // the reason we cant do a simple ON CONFLICT here is because we arent using regular SQL. its some weird sql-like language so we can avoid using C compilation

  r, _, err := db.db.Run(db.ctx, `SELECT * FROM prompt WHERE id = 1`)
  if err != nil {
    return err
  }

  rows, err := r[0].Rows(-1, 0)
  if err != nil {
    return err
  }

  a := len(rows) == 1

  var buffer bytes.Buffer
  buffer.WriteString(`BEGIN TRANSACTION;`)

  if a {
    buffer.WriteString(`DELETE FROM prompt WHERE id = 1;`)
  }

  buffer.WriteString(fmt.Sprintf(`INSERT INTO prompt (id, prompt, threshold, window) VALUES (1, "%v", %v, %v);`, prompt, threshold, window))

  buffer.WriteString(`COMMIT;`)
  _, _, err = db.db.Run(db.ctx, buffer.String())
  return err
}

func (db *Database) ResetOnboarded() (error) {
  _, _, err := db.db.Run(db.ctx, `
    BEGIN TRANSACTION;
      DELETE FROM prompt WHERE id = 1;
    COMMIT;
  `)
  return err
}

func (db *Database) GetData() (*map[int64][]DeviceDetailEntry, error) {
  r, _, err := db.db.Run(db.ctx, "SELECT * FROM network_data")
  if err != nil {
    return nil, err
  }

  rows, err := r[0].Rows(-1, 0)
  if err != nil {
    return nil, err
  } 
  data := make(map[int64][]DeviceDetailEntry)

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
    
    data[id] = append(data[id], *entry)
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

  //log.Printf("[database] successfully inserted id: %v", id)

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
  j, err := json.Marshal(*d)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Content-Length", fmt.Sprintf("%d", len(j)))
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.WriteHeader(http.StatusOK)
  w.Write(j)
}

func (b *Backend) CheckOnboardedHandler(w http.ResponseWriter, r *http.Request) {
  d, err := database.GetOnboarded()
  if err != nil {
    log.Printf("error getting data: %v", err)
  }
  j, err := json.Marshal(*d)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.Header().Set("Content-Length", fmt.Sprintf("%d", len(j)))
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.WriteHeader(http.StatusOK)
  w.Write(j)
}

type OnboardRequest struct {
  Prompt string `json:"prompt"`
  Threshold int64 `json:"threshold"`
  Window int64 `json:"window"`
}

func (b *Backend) HandleOnboard(w http.ResponseWriter, r *http.Request) {
  if r.Method == "DELETE" {
    err := database.ResetOnboarded()
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    w.WriteHeader(http.StatusOK)
    return
  }
  by, err := io.ReadAll(r.Body)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  body := string(by)
  if body == "" || r.Method != "POST" {
    w.WriteHeader(http.StatusBadRequest)
    return 
  }

  var req OnboardRequest

  err = json.Unmarshal(by, &req)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  err = database.SetPrompt(req.Prompt, req.Threshold, req.Window)
  if err != nil {
    log.Printf("error setting prompt: %v", err)
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  log.Printf("set prompt to: %v", body)
  w.WriteHeader(http.StatusOK)
}

func (b *Backend) start() (error) {
  r := http.NewServeMux()
  r.HandleFunc("/", b.handler)
  r.HandleFunc("/check_onboarded", b.CheckOnboardedHandler)
  r.HandleFunc("/onboard", b.HandleOnboard)

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

  blockchain *store.Store
  opts *bind.TransactOpts
}

func NewReporter() (*Reporter, error) {
  bc, err := ethclient.Dial("http://192.168.2.145:7545")
  if err != nil {
    return nil, err
  } 

  privateKey, err := crypto.HexToECDSA("c65c8be273db3bee1a55ae8f267bba4c193ea3435e9e9610d812c62a34592bc1")
  if err != nil {
    return nil, err
  }

  publicKey := privateKey.Public()
  publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
  if !ok {
    return nil, errors.New("error casting public key to ECDSA")
  }

  fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

  _, err = bc.PendingNonceAt(context.Background(), fromAddress)
  if err != nil {
    log.Fatal(err)
  }

  gasPrice, err := bc.SuggestGasPrice(context.Background())
  if err != nil {
    log.Fatal(err)
  }

  auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
  if err != nil {
    return nil, err
  }

  auth.Nonce = nil
  auth.Value = big.NewInt(0)     // in wei
  auth.GasLimit = uint64(5000000) // in units
  auth.GasPrice = gasPrice

  contractAddr := common.HexToAddress("0x8c4C0114C20bbf16B3bD61CA0E0179E80E4c3454")

  instance, err := store.NewStore(contractAddr, bc)
  if err != nil {
    return nil, err
  }

  r := &Reporter{
    threshold: 0.5, // this should be determined by the prompt as well
    blockchain: instance,
    opts: auth,
  }

  go r.ConsumeTrafficAnalysisResults()

  return r, nil
}

func (r *Reporter) updateRanking(mac string, networkEloRanking uint16, isAbuser bool) (error) {
  _, err := r.blockchain.UpdateRanking(r.opts, mac, networkEloRanking, isAbuser)
  if err != nil {
    return err
  }

  i, err := r.blockchain.GetRanking(nil, mac)
  if err != nil {
    return err
  }

  if i < 800 {
    networkManager.kickChan <- mac
  }

  log.Printf("[reporter] rank: %v", i)
  return nil
}

func (r *Reporter) ConsumeTrafficAnalysisResults() {
  for {
    log.Printf("waiting for analysis results")
    results := <- trafficAnalyzer.analysisResultsChan

    log.Printf("start abuse check")
    isAbuser := func() (bool) {
      if results.abuserRating > r.threshold {
        return true
      }

      return false
    }()

    log.Printf("end abuse check")

    // report to the blockchain here
    err := r.updateRanking(results.MAC, 1000, isAbuser)
    if err != nil {
      log.Printf("[reporter] failed to report to blockchain: %v", err)
    }

    log.Printf("end update ranking")

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

  pastTrafficBuffer deque.Deque[[]DeviceDetailEntry]
}

func NewTrafficAnalyzer() (*TrafficAnalyzer, error) {
  ta := &TrafficAnalyzer{
    analysisResultsChan: make(chan TrafficAnalysisResults),
    pastTrafficBuffer: deque.Deque[[]DeviceDetailEntry]{},
  }

  go ta.ConsumeNetworkDump()

  return ta, nil
}

func (ta *TrafficAnalyzer) ConsumeNetworkDump() {
  for {
    log.Printf("waiting for dump")
    entries := <- networkManager.networkDumpChan

    log.Printf("[traffic-analyzer] received network snapshot %v", entries)

    ta.pastTrafficBuffer.PushFront(entries)

    prompt, err := database.GetPrompt()
    if err == nil {
      for {
        if int64(ta.pastTrafficBuffer.Len()) <= prompt.window {
          log.Printf("this 1")
          break
        }
        ta.pastTrafficBuffer.PopBack()
      }

      m := ta.pastTrafficBuffer.Len()
      i := 1

      abuserTotal := map[string]int64{}

      for {
        if i >= m {
          log.Printf("this 2")
          break
        }
        // mapping every mac address to its total data use
        curr := map[string]int64{}
        for _, v := range ta.pastTrafficBuffer.At(i) {
          curr[v.MAC] = v.Total
        }

        last := map[string]int64{}
        for _, v := range ta.pastTrafficBuffer.At(i - 1) {
          last[v.MAC] = v.Total
        }

        for mac, tot := range curr {
          lastTot, ok := last[mac]
          // if this is the first moment we have seen this device we can skip
          if !ok {
            log.Printf("this 2")
            continue
          }

          if lastTot - tot > prompt.threshold {
            abuserTotal[mac] += 1
          }
        }

        i += 1
      }

      log.Printf("window: %v", prompt.window)
      log.Printf("abuserTotals: %v", abuserTotal) 

      for mac, total := range abuserTotal {
        log.Printf("bruh ting %v", ta.analysisResultsChan)

        // this is fuck
        ta.analysisResultsChan <- TrafficAnalysisResults{
          MAC: mac,
          abuserRating: float64(total) / float64(prompt.window),
        }
        log.Printf("bruh idk")
      }
    } else {
      log.Printf("error getting prompt: %v", err)
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
