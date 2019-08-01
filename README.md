# go-filelocker

Golang library for [Filelocker 2](http://filelocker2.sourceforge.net/)

## Library

```golang
package main

import (
  "fmt"
  "net/http"
  "time"

  "github.com/YaleUniversity/go-filelocker/pkg/filelocker"
)

var userID = "user123"
var apiKey = "aaaabbbbcccccddddeeeeeffff"
var filelockerURL = "https://files.example.edu"

var expireIn = "5d"
var secretSubject = "shhhhh"
var secretMessage = "pssssst! i have a secret!"
var recipients = []string{"user123"}

func main() {
  expire, _ := time.ParseDuration(expireIn)
  filelockerClient, _ := filelocker.NewClient(userID, apiKey, filelockerURL, &http.Client{Timeout: 30 * time.Second})
  resp, _ := filelockerClient.NewSecureMessage(secretSubject, secretMessage, recipients, time.Now().Add(expire))

  if len(resp.InfoMessages) > 0 {
    for _, m := range resp.InfoMessages {
      fmt.Println(m)
    }
  }

  if len(resp.ErrorMessages) > 0 {
    for _, m := range resp.ErrorMessages {
      fmt.Println(m)
    }
    panic("errrrrrrrrr!")
  }

  time.Sleep(5 * time.Second)

  msgs, _ := filelockerClient.SecureMessages()
  for _, m := range msgs.Messages[0] {
    fmt.Printf("ID: %d | Expiration: %s | Subject: %s | Body: %s\n", m.ID, m.Expiration, m.Subject, m.Body)
  }
}


```

## Command Line Interface

```bash
A go cli for interacting with filelocker 2.

Usage:
  filelocker [command]

Available Commands:
  help        Help about any command
  read        Reads secure messages from filelocker
  send        Send a secure message
  version     Displays version information

Flags:
      --config string    filelocker config file -- _not_ the control file (default is $HOME/.filelocker.yaml)
  -h, --help             help for filelocker
  -j, --json             Format the response as JSON where applicable
  -k, --key string       The api key to use for connections to filelocker
  -l, --login string     The userid to use for connections to filelocker
  -t, --timeout string   The filelocker http client timeout (seconds) (default "30s")
  -u, --url string       The base URL to use for connections to filelocker (ie. https://files.example.edu

Use "filelocker [command] --help" for more information about a command.
```

### Examples

**Send a secure message**

```bash
filelocker send -u 'https://files.example.edu' -l mynetid -k xxxxxyyyyyybbbbbbbzzzzzz -s 'test test' -r netid123 -b 'test123 go have fun'
```

**Read all messages**

```bash
filelocker read -u 'https://files.example.edu' -l mynetid -k xxxxxyyyyyybbbbbbbzzzzzz -a
```

**Read all messages as JSON***

```bash
filelocker read -u 'https://files.example.edu' -l mynetid -k xxxxxyyyyyybbbbbbbzzzzzz -a -j
```

```json
{
    "Messages": [
        {
            "body": "pssssst! i have a secret!",
            "creationDatetime": "08/01/2019",
            "expirationDatetime": "08/31/2019",
            "id": 12345,
            "ownerId": "user123",
            "messageRecipients": [
                "mynetid"
            ],
            "subject": "shhhhh",
            "viewedDatetime": "08/01/2019"
        }
    ],
    "Info": [],
    "Error": []
}
```

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

The MIT License (MIT)

Copyright (c) 2019 Yale University

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
