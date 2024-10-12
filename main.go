package main

import (
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*
#mac,ip,iface,in,out,total,first_date,last_date
a0:ce:c8:b0:e8:eb,192.168.2.10,br-lan,0,0,0,12-10-2024_06:50:03,12-10-2024_06:50:03
c6:f8:85:ce:62:d6,192.168.2.101,br-lan,0,0,0,12-10-2024_06:50:03,12-10-2024_06:50:03
f2:64:c7:89:4c:42,192.168.2.110,br-lan,468,0,468,12-10-2024_06:50:03,12-10-2024_06:50:59
*/

type DeviceDetailEntry struct {
  MAC string
  IP string
  Interface string
  In int64
  Out int64
  Total int64
  FirstDate time.Time
  LastDate time.Time
}

func parseDetailDate(s string) (*time.Time, error) {
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

  t := time.Date(yearI, time.Month(monthI), dayI, hoursI, minsI, secI, 0, time.Local)
  return &t, nil
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
    FirstDate: *firstDate,
    LastDate: *lastDate,
  }, nil
}

type NetworkManager struct {
  networkDetailDumperGroup sync.WaitGroup
  networkDumpInterval time.Duration

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

    log.Printf("%v", d)
  }

  // write to sqlite

  // update last.db

  return nil
}

func (rm *NetworkManager) listen() {
  for {
    mac := <- rm.kickChan

    err := rm.kickDevice(mac)
    if err != nil {
      log.Printf("error kicking device (%v): %v", mac, err)
    }
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
    networkDumpInterval: time.Second * 5,
    kickChan: make(chan string),
  }
}

var networkManager *NetworkManager

func main() {
  networkManager = NewNetworkManager()
  networkManager.start()
}

