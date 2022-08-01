package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	server     = flag.String("server", "can3.canaly.co.jp", "JWT Server")
	lock       = flag.String("lock", "/tmp/jwt-agent.pid", "Process ID")
	userId     string
	passphrase string
)

func init() {
	var oldPid string

	fp, err := os.Open(*lock)
	if err == nil {
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			oldPid = scanner.Text()
		}

		fmt.Printf("old pid:%s\n", oldPid)
		targetPid, err := strconv.Atoi(oldPid)
		process, err := os.FindProcess(targetPid)
		if err == nil {
			err = process.Kill()
			if err != nil {
				panic(err)
			}
		}
	}

	pid := os.Getpid()

	file, err := os.Create(*lock)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))
	if err != nil {
		panic(err)
	}
}

func getToken(userId string, passphrase string) (string, error) {
	endpoint := fmt.Sprintf("https://%s/jwt-server/jwt", *server)
	fmt.Println(endpoint)

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
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	token := string(body)

	return token, nil
}

func updateToken(userId string, passphrase string) (string, error) {
	fmt.Println("updateToken")

	return "", nil
}

func main() {
	flag.Parse()

	fmt.Print("userId:")
	fmt.Scan(&userId)

	fmt.Print("passphrase:")
	fmt.Scan(&passphrase)

	fmt.Printf("%s:%s\n", userId, passphrase)

	token, err := getToken(userId, passphrase)

	if err != nil {
		panic(err)
	}

	fmt.Println(token)

}
