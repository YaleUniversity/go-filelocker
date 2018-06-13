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
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

var messageSubject, messageBody, expireIn string
var recipientList []string

// sendCmd represents the command to send a message
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a secure message",
	RunE: func(cmd *cobra.Command, args []string) error {
		e, err := time.ParseDuration(expireIn)
		if err != nil {
			return errors.Wrap(err, "unable to parse expiration")
		}
		resp, err := filelockerClient.NewSecureMessage(messageSubject, messageBody, recipientList, time.Now().Add(e))
		if err != nil {
			return errors.Wrap(err, "unable to send secure message")
		}

		if len(resp.InfoMessages) > 0 {
			for _, m := range resp.InfoMessages {
				fmt.Println(m)
			}
		}

		if len(resp.ErrorMessages) > 0 {
			for _, m := range resp.ErrorMessages {
				fmt.Println(m)
			}
			return errors.Wrap(err, "error sending secure message")
		}

		return nil
	},
	TraverseChildren: true,
}

func init() {
	sendCmd.PersistentFlags().StringVarP(&messageSubject, "subject", "s", "Secure Mesaage", "The message subject")
	sendCmd.PersistentFlags().StringVarP(&messageBody, "body", "b", "", "The message body")
	sendCmd.PersistentFlags().StringVarP(&expireIn, "expireIn", "e", "720h", "The message expiration time from now (https://golang.org/pkg/time/#ParseDuration)")
	// sendCmd.PersistentFlags().StringVarP(&expireOn, "expireOn", "o", "", "The message expiration date from ")
	sendCmd.Flags().StringArrayVarP(&recipientList, "recipient", "r", []string{}, "Message recipient(s)")
	RootCmd.AddCommand(sendCmd)
}
