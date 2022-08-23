package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "log"
  "log/syslog"
  "net/http"
  "net/url"
  "os"
  "strconv"
  "strings"
  "time"
  "encoding/json"
  "os/user"

  b64 "encoding/base64"
  
  "golang.org/x/crypto/ssh/terminal"
)

var (
  server     = flag.String("s", "", "JWT Server")
  lock       = flag.String("lock", "jwt-agent.pid", "Process ID")
  userId     = flag.String("l", "", "User Name")
  uid        string
  passphrase string
  dir        string
)

func init() {
  var oldPid string

  cuser, err := user.Current()
  if err != nil {
     log.Fatalln(err)
     panic(err)
  }

  uid = cuser.Uid
  
  dir = "/tmp/jwt_user_u" + uid

  if _, err := os.Stat(dir); os.IsNotExist(err) {
    os.Mkdir(dir, 0755)
  }

  lockFile := dir + "/" + *lock;
  
  logger, err := syslog.New(syslog.LOG_INFO, "jwt-agent")
  if err != nil {
    panic(err)
  }
  log.SetOutput(logger)

  fp, err := os.Open(lockFile)
  if err == nil {
    defer fp.Close()
    scanner := bufio.NewScanner(fp)
    for scanner.Scan() {
      oldPid = scanner.Text()
    }

    targetPid, err := strconv.Atoi(oldPid)
    process, err := os.FindProcess(targetPid)
    if err == nil {
      err = process.Kill()
      if err != nil {
         log.Println(err)
      }
    }
  }

  pid := os.Getpid()

  file, err := os.Create(lockFile)
  if err != nil {
    log.Fatalln(err)  
    panic(err)
  }
  defer file.Close()

  _, err = file.WriteString(strconv.Itoa(pid))
  if err != nil {
    log.Fatalln(err)
    panic(err)
  }
}

func getToken(userId string, passphrase string, initial bool) (string, error) {
  endpoint := fmt.Sprintf("https://%s/jwt-server/jwt", *server)

  values := url.Values{}
  values.Set("user", userId)
  values.Add("pass", passphrase)

  req, err := http.NewRequest(
    "POST",
    endpoint,
    strings.NewReader(values.Encode()),
  )
  if err != nil {
    return "", err
  }

  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

  client := &http.Client{}
  resp, err := client.Do(req)

  if initial && err != nil {
    return "", err    
  } else if err != nil {
    sec := 1

    for err != nil {
      log.Printf("retry after %d seconds\n", sec)
      time.Sleep(time.Duration(sec) * time.Second)
    
      resp, err = client.Do(req)

      if sec >= 64 {
        sec = 64
      } else {
        sec *= 2
      }
    }    
  }
  
  defer resp.Body.Close()

  if resp.StatusCode != 200 {
    return "", fmt.Errorf("bad response status code %d", resp.StatusCode)
  }

  body, _ := io.ReadAll(resp.Body)
  token := string(body)

  if token == "" {
    return "", fmt.Errorf("authentication error")
  }

  filename := dir + "/token.jwt"
  file, err := os.Create(filename)
  if err != nil {
      return "", err
  }
  defer file.Close()
  _, err = file.WriteString(token)
  if err != nil {
      return "", err
  }

  log.Println("get token...")
  
  return token, nil
}

func parseToken(tokenString string) (int64, error) {
  header := strings.Split (tokenString,".")
  str_payload := strings.Replace(header[1],"-","+",-1)
  str_sig := strings.Replace(str_payload,"_","/",-1)

  llx := len(str_sig)
  nnx := ((4 - llx % 4) % 4)
  ssx := strings.Repeat("=" , nnx)
  str := strings.Join([]string{str_sig, ssx}, "")
  bytes, err :=  b64.StdEncoding.DecodeString(str)
  if err != nil {
    return 0, err
  }

  uEnc := b64.URLEncoding.EncodeToString([]byte(bytes))
  decode, _ := b64.URLEncoding.DecodeString(uEnc)

  var decode_data interface{}
  
  if err = json.Unmarshal(decode, &decode_data); err != nil {
   return 0, err
  }

  data := decode_data.(map[string]interface{})
  exp := data["exp"].(float64)
  
  now := time.Now().Unix()

  limit := (exp - float64(now)) * 0.8

  return int64(limit), nil
}

func main() {
  flag.Parse()

  if (*server == "" || *userId == "") {
    fmt.Println("Usage: jwt-agent -s {server} -l {user}")
    return
  }

  fmt.Print("passphrase:")
  passphrase, err := terminal.ReadPassword(0)
  if err != nil {
    log.Fatalln(err)
    panic(err)
  }
  fmt.Println()

  initial := true
  
  for {
    token, err := getToken(*userId, string(passphrase), initial)
    if err != nil {
      log.Fatalln(err)
      panic(err)
    }

    limit, err := parseToken(token)
    if err != nil {
      log.Fatalln(err)
      panic(err)
    }

    initial = false
    time.Sleep(time.Duration(limit) * time.Second)  
  }
}
