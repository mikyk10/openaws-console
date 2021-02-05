package driver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ChromeDriver interface {
	GetChromeVersion() string
	GetDriverVersion() string
	GetLatestDriverURL() string
}

func GetMajorVersion(ver string) string {
	vers := strings.Split(ver, ".")
	if len(vers) == 0 {
		panic("no version(s) found")
	}
	return vers[0]
}

func GetLatestChromeDriverVersion(majv string) string {
	url := fmt.Sprintf("https://chromedriver.storage.googleapis.com/LATEST_RELEASE_%s", majv)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	buf, _ := ioutil.ReadAll(resp.Body)
	return string(buf)
}
