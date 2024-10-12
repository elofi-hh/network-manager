package main

import (
	"fmt"
	"os/exec"
)

func kickDevice(mac string) (error) {
  ipsetCmd := exec.Command("ipset", "add", "mac-ban", mac)
  o, err := ipsetCmd.Output()
  if err != nil {
    return err
  }
  fmt.Println(string(o))
  cmd := exec.Command("hostapd_cli", "disassociate", mac)
  o, err = cmd.Output()
  if err != nil {
    return err
  }
  fmt.Println(string(o))

  cmd = exec.Command("hostapd_cli", "deauthenticate", mac)
  o, err = cmd.Output()
  if err != nil {
    return err
  }
  fmt.Println(string(o))
  return nil
}

func main() { 
  err := kickDevice("F2:64:C7:89:4C:42")
  if err != nil {
    fmt.Printf("fuck: %v\n", err)
  }
}
