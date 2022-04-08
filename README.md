# openaws-console

Automation tool for AWS console login.

## Overview

Assuming you are working on multiple AWS management consoles simultaneously across different account IDs, you might open many browser windows for each account(s) and type lots of credentials or switch IAM role(s) back and forth. Either way, it's repetitive and cumbersome task. This tool can help you and ease your fatigue to open up AWS management console across multiple accounts without having you to enter ID/PW.

# Prerequisites

- macOS
- Go 1.17 or above
- Google Chrome
- Selenium webdriver for Chrome (the tool will automatically install for you.)
- awscli

## Install

Just runing the following command will get you the compiled binary under `~/go/bin/`.
You probably need to configure PATH to execute beforehand.

```
go install github.com/mikyk10/openaws-console@master
```

## Configuring

This tool accepts `console_account` `console_username` `console_password` in your `~/.aws/config`. Enter appropriate credential and save. The original awscli tool will ignore these configurations and nothing will happens to awscli tool.

```
[profile your-aws-profile]
region = ap-northeast-1
aws_access_key_id = AKIA****************
aws_secret_access_key = ********************************
console_account = 999999999999
console_username = foobar
console_password = "your_password"
```

## Usage

```
openaws-console [aws profile name]

# For alfred user
openaws-console --alfred
```

## Disclaimer

The software is expressly provided “AS IS.” THE AUTHOR MAKES NO WARRANTY OF ANY KIND, EXPRESS, IMPLIED, IN FACT OR ARISING BY OPERATION OF LAW, INCLUDING, WITHOUT LIMITATION, THE IMPLIED WARRANTY OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT AND DATA ACCURACY.

## Contributing

Contributions welcome! Please read the [contributing guidelines](CONTRIBUTING.md) first.


## License

[GPLv3](LICENSE)
