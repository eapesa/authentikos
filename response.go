package main

import (
  "encoding/json"
)

type Response struct {
  Code string `json:"code,omitempty"`
  Message string `json:"message,omitempty"`
  Params string `json:"params,omitempty"`
}

var (
  OK, _ = json.Marshal(&Response{Code:"OK", Message:"Operation completed successfully"})
  NOT_FOUND, _ = json.Marshal(&Response{Code:"NOT_FOUND", Message:"Does not exist"})
)
