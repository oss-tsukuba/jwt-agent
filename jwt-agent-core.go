package main

import (
  "bufio"
  "flag"
  "fmt"
  "io"
  "io/ioutil"
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
  "path/filepath"

  b64 "encoding/base64"

)

var (
  server     = flag.String("s", "", "JWT Server URL")
  lock       = flag.String("lock", "jwt-agent.pid", "Process ID")
  userId     = flag.String("l", "", "User Name")
  uid        string
  dir        string
  servers     []string
  basename   = "token.jwt"
)

func init() {
  var oldPid string

  cuser, err := user.Current()
  if err != nil {
     log.Fatalln(err)
     panic(err)
  }

  uid = cuser.Uid

  path := os.Getenv("JWT_USER_PATH")

  if (path != "") {
    basename = filepath.Base(path)
    dir = filepath.Dir(path)
  } else {
    dir = "/tmp/jwt_user_u" + uid
  }

  if _, err := os.Stat(dir); os.IsNotExist(err) {
    os.Mkdir(dir, 0755)
  }

  lockFile := dir + "/" + *lock;

  logger, err := syslog.New(syslog.LOG_INFO, "jwt-agent")
  if err == nil {
    log.SetOutput(logger)
  } else {
    fmt.Fprintln(os.Stderr, "syslog disabled, error messages will be abandoned")
    log.SetOutput(ioutil.Discard)
  }

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

  file, err := os.OpenFile(lockFile, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
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

func new_servers(ss []string, pos int) []string {
  sliceSize := len(ss)

  if pos == 0 || sliceSize <= pos {
    return ss
  }
  org1 := ss[0:pos]
  org2 := ss[pos:sliceSize]
  return append(org2, org1...)
}

func getToken(userId string, passphrase string, initial bool) (string, error) {
  var new_ss []string
  values := url.Values{}
  values.Set("user", userId)
  values.Add("pass", passphrase)

  client := &http.Client{
    Timeout: 10 * time.Second,
  }

  var resp *http.Response
  sec := 1

  for {
    all_err := true

    for i:= 0; i < len(servers); i++ {
      endpoint := fmt.Sprintf("%sjwt", servers[i])
      req, err := http.NewRequest(
        "POST",
        endpoint,
        strings.NewReader(values.Encode()),
      )
      if err != nil {
        return "", err
      }

      req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

      resp, err = client.Do(req)

      if err != nil {
        log.Println(err)
        return "", err
      } else {
        defer resp.Body.Close()
      }

      if err == nil && resp.StatusCode == 200 {
        new_ss = new_servers(servers, i)
	all_err = false
        break
      }
    }

    if all_err {
      if initial {
        return "", fmt.Errorf("%s, %d, %s", *server, resp.StatusCode, http.StatusText(resp.StatusCode))
      } else {
        log.Printf("retry after %d seconds\n", sec)

        time.Sleep(time.Duration(sec) * time.Second)

        if sec >= 64 {
          sec = 64
        } else {
          sec *= 2
        }
        continue
      }
    }
    servers = new_ss
    break
  }

  body, _ := io.ReadAll(resp.Body)
  token := string(body)

  if token == "" {
    return "", fmt.Errorf("jwt-agent, Authentication error (exited)")
  }

  tmpname := "token.tmp"
  filepath := dir + "/" + tmpname
  file, err := os.OpenFile(filepath, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0600)
  if err != nil {
      return "", err
  }
  defer file.Close()
  _, err = file.WriteString(token)
  if err != nil {
      return "", err
  }

  err = os.Rename(dir + "/" + tmpname, dir + "/" + basename)
  if err != nil {
      return "", err
  }

  if initial {
    fmt.Printf("Output JWT to %s\n", dir + "/" + basename)
  }

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
    fmt.Println("Usage: jwt-agent-core -s {URL} -l {USER}")
    return
  }

  var passphrase string
  fmt.Scan(&passphrase)

  initial := true
  servers = strings.Split(*server, " ")

  for {
    token, err := getToken(*userId, passphrase, initial)
    if err != nil {
      fmt.Fprintln(os.Stderr, *userId + ": " + err.Error())
      log.Fatalln(*userId + ": " + err.Error())
      panic(err)
    }

    limit, err := parseToken(token)
    if err != nil {
      fmt.Fprintln(os.Stderr, "Invalid Token")
      log.Fatal(*userId + ": Invalid Token")
      log.Fatalln(err)
      panic(err)
    }

    initial = false
    time.Sleep(time.Duration(limit) * time.Second)
  }
}
