package driver

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type chromeDriver struct{}

func (c *chromeDriver) GetChromeVersion() string {
	var chromeVersion string
	{
		cmd := exec.Command("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "--version")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		rawVersion, _ := bufio.NewReader(stdout).ReadString('\n')
		vs := strings.Split(strings.NewReplacer(" ", "", "Google", "", "Chrome", "", "\n", "").Replace(rawVersion), ".")
		chromeVersion = strings.Join(vs[0:3], ".")
	}

	return chromeVersion
}

func (c *chromeDriver) GetDriverVersion() string {
	var driverVersion string
	{
		cmd := exec.Command("chromedriver", "--version")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, "`chromedriver` not found.")
			return ""
		}
		rawVersion, _ := bufio.NewReader(stdout).ReadString('\n')
		vs := strings.Split(strings.NewReplacer(" ", "", "ChromeDriver", "", "\n", "").Replace(rawVersion), ".")
		driverVersion = strings.Join(vs[0:3], ".")
	}
	return driverVersion
}

func (c *chromeDriver) GetLatestDriverURL() string {
	latestVersion := GetLatestChromeDriverVersion(GetMajorVersion(c.GetChromeVersion()))
	url := fmt.Sprintf("https://chromedriver.storage.googleapis.com/%s/chromedriver_mac64.zip", latestVersion)
	return url
}

func NewChromeDriver() ChromeDriver {
	return &chromeDriver{}
}
