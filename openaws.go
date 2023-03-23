package main

/* openaws.go
 *
 * Copyright (C) 2023 Mitsutaka Kato
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or (at
 * your option) any later version.
 *
 * This program is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA 02110-1301, USA.
 */

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"os"

	"github.com/manifoldco/promptui"
	"github.com/mikyk10/openaws-console/driver"
	"github.com/sclevine/agouti"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

type Profile struct {
	SourceProfile     string
	AssumingRoleName  string
	AssumingAccountID string
	ConsoleUserName   string
	ConsolePassword   string
	AccountID         string
}

type AlfredJSONItem struct {
	Title    string `json:"title"`
	Arg      string `json:"arg"`
	Subtitle string `json:"subtitle"`
	//Icon string `json:"icon"`
}

type AlfredJSON struct {
	Items []AlfredJSONItem `json:"items"`
}

const url = "https://console.aws.amazon.com/console/home"

func main() {
	cmdMain.Flags().Bool("alfred", false, "For alfred")

	// コマンドを実行
	_, err := cmdMain.ExecuteC()
	if err != nil {
		os.Exit(1)
		return
	}
}

var cmdMain = &cobra.Command{
	RunE: func(cmd *cobra.Command, args []string) error {

		profs := map[string]Profile{}

		// obtain HOME directory
		homedir, _ := os.UserHomeDir()

		var sourceProfile Profile
		var assumingProfile Profile

		// Obtain AWS Profiles
		opt := ini.LoadOptions{
			UnescapeValueDoubleQuotes: true,
		}

		cfg, _ := ini.LoadSources(opt, filepath.Join(homedir, ".aws", "config"))
		sections := cfg.Sections()
		for i := range sections {
			pname := strings.Replace(sections[i].Name(), "profile ", "", 1)
			roleArn := strings.Split(sections[i].Key("role_arn").Value(), ":")

			p := Profile{
				SourceProfile: sections[i].Key("source_profile").Value(),
			}

			if len(roleArn) == 6 {
				p.AssumingRoleName = strings.Replace(roleArn[5], "role/", "", 1)
				p.AssumingAccountID = roleArn[4]
			}

			// for this time, get console username and password from the .aws/config
			//TODO: organize well
			if sections[i].HasKey("console_account") {
				p.AccountID = sections[i].Key("console_account").Value()
				p.ConsoleUserName = sections[i].Key("console_username").Value()
				p.ConsolePassword = sections[i].Key("console_password").Value()
			}

			profs[pname] = p
		}

		var text string

		isAlfred, _ := cmd.Flags().GetBool("alfred")
		if isAlfred {
			alf := AlfredJSON{}

			for key := range profs {
				if len(args) > 0 {
					if !strings.Contains(key, args[0]) {
						continue
					}
				}

				subt := strings.Builder{}
				if profs[key].AccountID != "" {
					subt.WriteString(fmt.Sprintf("AccountID: %s  ", profs[key].AccountID))
				}
				if profs[key].ConsoleUserName != "" {
					subt.WriteString(fmt.Sprintf("Username: %s  ", profs[key].ConsoleUserName))
				}
				if profs[key].SourceProfile != "" {
					subt.WriteString(fmt.Sprintf("SourceProfile: %s ", profs[key].SourceProfile))

					if profs[key].AssumingAccountID != "" {
						subt.WriteString(fmt.Sprintf("[AccountID: %s  Role: %s]", profs[key].AssumingAccountID, profs[key].AssumingRoleName))
					}
				}

				if subt.Len() == 0 {
					continue
				}

				alf.Items = append(alf.Items, AlfredJSONItem{
					Title:    key,
					Subtitle: subt.String(),
					Arg:      key,
				})
			}

			s, _ := json.Marshal(alf)
			fmt.Println(string(s))

			return nil
		}

		if len(args) == 0 {
			pnames := []string{}
			for key := range profs {
				pnames = append(pnames, key)
			}
			sort.Slice(pnames, func(i, j int) bool {
				return pnames[i] < pnames[j]
			})

			prompt := promptui.Select{
				Label:     "Choose AWS Profile: ",
				HideHelp:  true,
				Items:     pnames,
				Size:      10,
				Templates: &promptui.SelectTemplates{Active: "{{.|red|cyan}}", Inactive: "{{.}}", Selected: "{{.}}"},
			}
			i, _, err := prompt.Run()
			if err != nil {
				os.Exit(1)
			}

			text = pnames[i]
		} else {
			text = args[0]
		}

		// Check Chrome installation
		chromeVer := driver.NewChromeDriver().GetChromeVersion()
		driverVer := driver.NewChromeDriver().GetDriverVersion()
		latestDriverURL := driver.NewChromeDriver().GetLatestDriverURL()

		// Print driver/Chrome versions
		fmt.Fprintf(os.Stderr, "Your Google Chrome version: %s\n", chromeVer)
		fmt.Fprintf(os.Stderr, "Your ChromeDriver version: %s\n", driverVer)
		// update the Chrome if needed
		if strings.Compare(driver.GetMajorVersion(chromeVer), driver.GetMajorVersion(driverVer)) != 0 {
			fmt.Fprintf(os.Stderr, "Downloading the latest driver at: %s\n", latestDriverURL)
			resp, err := http.Get(latestDriverURL)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			f, err := ioutil.TempFile(os.TempDir(), "chromedriver")
			if err != nil {
				panic(err)
			}

			io.Copy(f, resp.Body)

			buf, err := exec.Command("/usr/bin/unzip", "-od", "/usr/local/bin/", f.Name()).CombinedOutput()
			if err != nil {
				panic(err)
			}
			fmt.Fprintf(os.Stderr, string(buf))
		}

		prof := strings.Replace(text, "\n", "", -1)
		sourceProfile = profs[prof]
		assumingProfile = profs[prof]

		if assumingProfile.SourceProfile != "" {
			sourceProfile = profs[assumingProfile.SourceProfile]
		}

		// launch Selenium driver
		options := agouti.ChromeOptions(
			"args", []string{
				"--disable-gpu",
				"--disable-extensions",
			})

		// we don't want to close driver after login.
		driver := agouti.ChromeDriver(options)
		driver.Start()

		page, err := driver.NewPage()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return nil
		}

		fmt.Fprintf(os.Stderr, "opening AWS...\n")
		page.Navigate(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return nil
		}

		// enter login credentials...
		fmt.Fprintf(os.Stderr, "entering account ID\n")
		page.Find("#iam_user_radio_button").Click()
		page.Find("#resolving_input").SendKeys(sourceProfile.AccountID)
		page.Find("#next_button").Click()

		fmt.Fprintf(os.Stderr, "entering credentials\n")
		page.Find("#username").SendKeys(sourceProfile.ConsoleUserName)
		page.Find("#password").SendKeys(sourceProfile.ConsolePassword)
		page.Find("#signin_button").Click()

		fmt.Fprintf(os.Stderr, "waiting for login...\n")

		if assumingProfile.SourceProfile == "" {
			return nil
		}

		// wait for browser title change
		for {
			title, _ := page.Title()
			time.Sleep(100 * time.Millisecond)
			if title == "AWS Management Console" {
				time.Sleep(1 * time.Second)
				break
			}
		}

		// visit Switch Role
		fmt.Fprintf(os.Stderr, "ok. assuming a role\n")
		page.FindByID("nav-usernameMenu").Click()
		//page.FindByLink("Switch Roles").Click()
		page.FindByXPath("//a[@data-testid='awsc-switch-roles']").Click()
		page.FindByID("switchrole_firstrun_button").Click()

		// do the switch role
		page.FindByID("account").Fill(profs[prof].AssumingAccountID)
		page.FindByID("roleName").Fill(profs[prof].AssumingRoleName)
		page.FindByID("input_switchrole_button").Click()

		return nil
	},
}
