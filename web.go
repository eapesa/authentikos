package main

import(
  "fmt"
  "net/http"
  "github.com/pquerna/otp/totp"
  "bytes"
  "image/png"
  "github.com/go-redis/redis"
)

var CacheClient *redis.Client

func initializeHelperClients() {
  CacheClient = redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
		Password: "",
		DB:       0,
  })
  CacheClient.FlushDB()
}

func initializeServer() {
  port := ":8000"

  fmt.Printf("Initializing server at %s\n", port)
  viewServer   := http.FileServer(http.Dir("priv/views"))
  assetsServer := http.FileServer(http.Dir("priv/assets"))
  http.Handle("/", http.StripPrefix("/", viewServer))
  http.Handle("/assets/", http.StripPrefix("/assets/", assetsServer))
  http.HandleFunc("/otp/generate", generateHandler)
  http.HandleFunc("/otp/verify", verifyHandler)

  http.ListenAndServe(port, nil)
}

func main() {
  // NOTE: Initialize everything here
  initializeHelperClients()
  initializeServer()
}

func generateHandler(res http.ResponseWriter, req *http.Request) {
  i := 1
  if req.Method != "GET" {
    fmt.Printf("[generateHandler:%d] Encountered error Request method not supported\n", i)
    res.Write(NOT_FOUND)
    return
  }

  account_name := req.URL.Query().Get("account_name")
  i++
  if (account_name == "") {
    fmt.Printf("[generateHandler:%d] Validation error.\n", i)
    res.Write(INVALID_INPUT)
    return
  }

  qrCode, err1 := generateQrCode(account_name)
  i++
  if err1 != nil {
    fmt.Printf("[generateHandler:%d] Server Error: %v\n", i, err1)
    res.Write(SERVER_ERROR)
    return
  }

  res.Header().Set("Content-Type", "image/jpeg")
  _, err2 := res.Write(qrCode)
  i++
  if err2 != nil {
    fmt.Printf("[generateHandler:%d] Server Error: %v\n", i, err2)
    return
  }

  return
}

func verifyHandler(res http.ResponseWriter, req *http.Request) {
  i := 1
  if req.Method != "GET" {
    fmt.Printf("[verifyHandler:%d] Encountered error Request method not supported\n", i)
    res.Write(NOT_FOUND)
    return
  }

  passcode     := req.URL.Query().Get("passcode")
  account_name := req.URL.Query().Get("account_name")
  i++
  if passcode == "" || account_name == "" {
    fmt.Printf("[verifyHandler:%d] Validation error.\n", i)
    res.Write(INVALID_INPUT)
    return
  }

  ok := validatePasscode(passcode, account_name)
  i++
  if ok != true {
    fmt.Printf("[verifyHandler:%d] Incorrect passcode: %v\n", i, ok)
    res.Write(UNAUTHORIZED)
    return
  }

  res.Header().Set("Content-Type", "application/json")
  _, err2 := res.Write(OK)
  if err2 != nil {
    fmt.Printf("[verifyHandler:%d] Server Error: %v\n", i, err2)
    return
  }

  return
}

func generateQrCode(account_name string) ([]byte, error) {
  var opts = totp.GenerateOpts{
    Issuer: "Authentikos",
    AccountName: account_name,
  }

  key, err1 := totp.Generate(opts)
  if err1 != nil {
    return nil, err1
  }

  _, err2 := storeTotpKey(account_name, key.Secret())
  if err2 != nil {
    return nil, err2
  }

  var imageBuf bytes.Buffer
  image, err3 := key.Image(200, 200)
  if err3 != nil {
    return nil, err3
  }

  png.Encode(&imageBuf, image)
  return imageBuf.Bytes(), nil
}

func validatePasscode(passcode, account_name string) bool {
  secret, _ := getTotpKey(account_name)
  return totp.Validate(passcode, secret)
}

func storeTotpKey(cacheKey, totpKey string) (string, error) {
  return CacheClient.Set(cacheKey, totpKey, 0).Result()
}

func getTotpKey(cacheKey string) (string, error) {
  return CacheClient.Get(cacheKey).Result()
}
