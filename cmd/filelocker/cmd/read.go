// Copyright © 2018 Yale University
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/YaleUniversity/go-filelocker/pkg/filelocker"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

var allMessages, markRead bool

// readCmd represents the command to read messages
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Reads secure messages from filelocker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !allMessages {
			resp, err := filelockerClient.SecureMessagesCount()
			if err != nil {
				return errors.Wrap(err, "unable to get message count")
			}

			fmt.Println("New Messages:", resp.Count)
			return nil
		}

		resp, err := filelockerClient.SecureMessages()
		if err != nil {
			return errors.Wrap(err, "unable to read secure message")
		}

		if asJSON {
			out, jsonErr := messagesToJSON(resp)
			if jsonErr != nil {
				return errors.Wrap(jsonErr, "unable to read secure message")
			}
			fmt.Println(string(out))
			return nil
		}

		for _, m := range resp.Messages[0] {
			fmt.Printf("ID: %d | Expiration: %s | Subject: %s | Body: %s\n", m.ID, m.Expiration, m.Subject, m.Body)
		}
		return nil
	},
	TraverseChildren: true,
}

func init() {
	readCmd.PersistentFlags().BoolVarP(&allMessages, "all", "a", false, "Get all messages instead of listing a count of new messages")
	readCmd.PersistentFlags().BoolVarP(&markRead, "mark", "m", false, "Mark secure messages as read")
	RootCmd.AddCommand(readCmd)
}

func messagesToJSON(resp *filelocker.SecureMessagesResponse) ([]byte, error) {
	list := struct {
		Messages []filelocker.SecureMessage
		Info     []string
		Error    []string
	}{
		Info:  resp.InfoMessages,
		Error: resp.InfoMessages,
	}

	// clean up message response from filelocker since its ugly:
	// {
	//     "data": [
	//         [
	//             {
	//                 "body": "Likes cats",
	//                 "creationDatetime": "06/11/2018",
	//                 "expirationDatetime": "07/11/2018",
	//                 "id": 777,
	//                 "ownerId": "Penny",
	//                 "messageRecipients": [
	//                     "Brain"
	//                 ],
	//                 "subject": "Doctor Claw",
	//                 "viewedDatetime": "06/11/2018"
	//             }
	//         ],
	//         [
	//             {
	//                 "body": "sekret message",
	//                 "creationDatetime": "06/11/2018",
	//                 "expirationDatetime": "07/11/2018",
	//                 "id": 778,
	//                 "ownerId": "Quimby",
	//                 "messageRecipients": [
	//                     "Gadget"
	//                 ],
	//                 "subject": "this message will self destruct",
	//                 "viewedDatetime": ""
	//             }
	//         ]
	//     ],
	//     "fMessages": [],
	//	   "sMessages": []
	// }
	for _, m := range resp.Messages[0] {
		list.Messages = append(list.Messages, m)
	}

	out, jsonErr := json.MarshalIndent(list, "", "    ")
	if jsonErr != nil {
		return []byte{}, errors.Wrap(jsonErr, "unable to marshal respons into JSON")
	}

	return out, nil
}
