package speedserver

import(
  //"fmt"
  "requester"
  "net/http"
  "strconv"
  "encoding/json"
)

type badReqError int

func (b badReqError) Error() string {
  return "bad Request"
}

func parseGenerateReq(r *http.Request) (string, string, int, int, error){
  err := r.ParseForm();
  if err != nil {
    return "", "", 0, 0, badReqError(0)
  }
  // method, url, transactions, tps
  trans, err := strconv.Atoi(string(r.Form["numTransactions"][0]))
  if err != nil {
      return "", "", 0, 0, badReqError(0)
  }
  tps, err := strconv.Atoi(string(r.Form["tps"][0]))
  if err != nil {
      return "", "", 0, 0, badReqError(0)
  }
  return string(r.Form["method"][0]), string(r.Form["url"][0]), int(trans), int(tps), nil
}

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
  method, url, requests, rps, err := parseGenerateReq(r);
  if ((err != nil) || ((method=="") || (url=="") || (requests==0) || (rps==0))) {
      w.WriteHeader(http.StatusBadRequest);
  }
  mapassshit, err := requester.MakeRequests(method, url, rps, requests, 1000)
  if err != nil {
    panic(err)
  }
  w.Header().Add("Content-Type", "application/json")
  json.NewEncoder(w).Encode(mapassshit)
}
