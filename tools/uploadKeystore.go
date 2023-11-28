package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/pterm/pterm"
)

const uploadKeystore = "uploadKeystore"

func uploadCommandAction() {
	api, _ := url.ParseRequestURI(apiUrl)
	apiRouter := api.JoinPath("ssv").JoinPath("upload").String()

	pterm.Info.Println("DVT-Service api rawUrl is set to ", apiRouter)

	keystoreFolder, _ := pterm.DefaultInteractiveTextInput.
		WithDefaultText("Input keystore folder path").
		WithMultiLine(false).
		Show()
	if len(keystoreFolder) == 0 {
		pterm.Error.Println("path is empty")
		return
	}

	dirEntry, err := os.ReadDir(keystoreFolder)
	if err != nil {
		pterm.Error.Printf("read keystore folder error %s\n", err)
		return
	}
	pterm.Success.Println("Check keystore folder success")

	var keys []*KeystoreV4
	packageProgressbar, _ := pterm.DefaultProgressbar.WithTotal(len(dirEntry)).WithTitle("Package keystore info...").Start()
	for _, entry := range dirEntry {
		if strings.Contains(entry.Name(), "keystore") {
			filePath := fmt.Sprintf("%s/%s", keystoreFolder, entry.Name())
			fileBytes, err := os.ReadFile(filePath)
			if err != nil {
				pterm.Error.Println(fmt.Sprintf("Read keystore file fail: %s", err))
				packageProgressbar.Stop()
				return
			}

			var file KeystoreV4
			if err := sonic.Unmarshal(fileBytes, &file); err != nil {
				pterm.Error.Println(fmt.Sprintf("Unmarshal keystore file fail: %s", err))
				packageProgressbar.Stop()
				return
			}
			pterm.Success.Println(fmt.Sprintf("Keystore %s packaged", file.Pubkey))
			keys = append(keys, &file)
		} else {
			pterm.Warning.Println(fmt.Sprintf("Skip by not keystore file: %s", entry.Name()))
		}
		packageProgressbar.Increment()
	}
	pterm.Success.Println("Package keystore info success!")

	password, _ := pterm.DefaultInteractiveTextInput.
		WithDefaultText("Enter keystore password").
		WithMask("*").
		WithMultiLine(false).
		Show()
	if len(password) < 8 {
		pterm.Error.Println("password too short")
		return
	}

	postSpinnerInfo, _ := pterm.DefaultSpinner.Start("Post request...")
	request := &UploadKeystore{
		Keys:     keys,
		Password: password,
	}

	marshal, _ := sonic.Marshal(request)

	if response, err := httpPost(apiRouter, marshal); err != nil {
		postSpinnerInfo.Fail("Post upload error: ", err)
		return
	} else {
		var res ResultResponse[string]
		sonic.Unmarshal(response, &res)
		if res.Code == 200 {
			postSpinnerInfo.Success("Keystore info upload success!")
		} else {
			postSpinnerInfo.Fail("Post upload error: ", res.Data)
		}
	}
}

func httpPost(url string, body []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
