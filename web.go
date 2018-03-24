package main

import(
  "fmt"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "github.com/pquerna/otp/totp"
  "bytes"
  "image/png"
)

type GenerateBody struct {
  Issuer      string `json:"issuer,omitempty"`
  AccountName string `json:"account_name,omitempty"`
}

func main() {
  port := ":8000"
  fmt.Printf("Initializing server at %s\n", port)
  viewServer  := http.FileServer(http.Dir("priv/views"))
  assetsServer := http.FileServer(http.Dir("priv/assets"))
  http.Handle("/", http.StripPrefix("/", viewServer))
  http.Handle("/assets/", http.StripPrefix("/assets/", assetsServer))
  http.HandleFunc("/otp/generate", generateHandler)
  http.HandleFunc("/otp/verify", verifyHandler)

  http.ListenAndServe(port, nil)
}

func generateHandler(res http.ResponseWriter, req *http.Request) {
  bodyBytes, err := ioutil.ReadAll(req.Body)
  if err != nil {
    fmt.Printf("[API-Generate] 1:Encountered error: %v\n", err)
  }
  var bodyJson GenerateBody
  err2 := json.Unmarshal(bodyBytes, &bodyJson)
  if err2 != nil {
    fmt.Printf("[API-Generate] 2:Encountered error: %v\n", err2)
  }
  qrCode, err3 := generateQrCode(bodyJson)
  if err3 != nil {
    fmt.Printf("[API-Generate] 3:Encountered error: %v\n", err3)
  }
  res.Header().Set("Content-Type", "image/jpeg")
  _, err4 := res.Write(qrCode)
  if err4 != nil {
    fmt.Printf("[API-Generate] 4:Encountered error: %v\n", err4)
  }
}

func verifyHandler(res http.ResponseWriter, req *http.Request) {
  res.Write(OK)
}

func generateQrCode(body GenerateBody) ([]byte, error) {
  key, err := totp.Generate(totp.GenerateOpts{
      Issuer:       "gmail.com",//body.Issuer,
      AccountName:  "eapesa@gmail.com",//body.AccountName,
  })
  if err != nil {
    fmt.Printf("[API-Generate:generateQrCode] 1:Encountered error: %v\n", err)
    return nil, err
  }
  var imageBuf bytes.Buffer
  image, err2 := key.Image(200, 200)
  if err2 != nil {
    fmt.Printf("[API-Generate:generateQrCode] 2:Encountered error: %v\n", err)
    return nil, err
  }
  png.Encode(&imageBuf, image)
  return imageBuf.Bytes(), nil
}
