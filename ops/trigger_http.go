package ops

import (
	"bufio"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/http2"

	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Triggers a backup via HTTP",
	Long:  `Triggers a backup via HTTP`,
	Run:   AttachHandler(remoteBackupTrigger),
}

var customHeaders []string // given in key=value format via command line
var method string
var isHTTP2 bool
var payloadFile string
var skipSSLVerification bool

func remoteBackupTrigger(args []string) error {
	if len(args) != 1 {
		return errors.New("No url was passed")
	}

	url := args[0]
	var inputFile *os.File
	if payloadFile == "stdin" {
		inputFile = os.Stdin
	} else if payloadFile != "" {
		f, err := os.Open(payloadFile)
		if err != nil {
			return err
		}
		inputFile = f
	}
	inputPayload := bufio.NewReader(inputFile)

	// Builds the HTTP Client
	client, err := buildClient()
	if err != nil {
		return err
	}

	// Builds the HTTP Request
	req, err := http.NewRequest(method, url, inputPayload)
	if err != nil {
		return err
	}
	headers := parseHeaders(customHeaders)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Make the HTTP Request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))
	return nil
}

func parseHeaders(input []string) map[string]string {
	output := make(map[string]string)
	for _, header := range input {
		splits := strings.Split(header, "=")
		if len(splits) == 2 {
			output[splits[0]] = splits[1]
		} // TODO - Handle when we want to set a header name along without a value
	}

	return output
}

func buildClient() (*http.Client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipSSLVerification},
	}
	if isHTTP2 {
		err := http2.ConfigureTransport(tr)
		if err != nil {
			return nil, err
		}
	}
	client := http.Client{
		Transport: tr,
	}
	return &client, nil
}

func init() {
	trigger.AddCommand(httpCmd)
	httpCmd.PersistentFlags().StringVarP(&method, "method", "X", "GET", "HTTP Method to invoke on the url")
	httpCmd.PersistentFlags().StringSliceVarP(&customHeaders, "header", "H", []string{}, "HTTP Headers that needs to be passed with the request in key=value format")
	httpCmd.PersistentFlags().BoolVarP(&isHTTP2, "http2", "", false, "Make a HTTP2 request instead of the default HTTP1.1")
	httpCmd.PersistentFlags().StringVarP(&payloadFile, "payload", "D", "", "File contents to be passed as payload in the request. If you're pipeing use \"stdin\".")
	httpCmd.PersistentFlags().BoolVarP(&skipSSLVerification, "skip-ssl-verificatioin", "k", false, "Ignore SSL errors - while connecting to self-signed HTTPS servers")
}
