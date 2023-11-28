package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pterm/pterm"
)

const setFeeRecipient = "setFeeRecipient"

type setFeeRecipientAddressBody struct {
	FeeRecipient string `json:"fee_recipient"`
}

type resultResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func setFeeRecipientCommandAction() {
	api, _ := url.ParseRequestURI(apiUrl)
	apiRouter := api.JoinPath("ssv").JoinPath("setFeeRecipientAddress").String()

	pterm.Info.Println("DVT-Service api rawUrl is set to ", apiRouter)

	orcAddress, _ := pterm.DefaultInteractiveTextInput.
		WithDefaultText("Input Fee Recipient Address").
		WithMultiLine(false).
		Show()
	address := strings.TrimSpace(orcAddress)
	if len(address) == 0 {
		pterm.Error.Println("address is empty")
		return
	}

	if !common.IsHexAddress(address) {
		pterm.Error.Println("invalid address")
		return
	}

	spinnerLiveText, _ := pterm.DefaultSpinner.Start("Send a request to the dvt service...")

	postbody := &setFeeRecipientAddressBody{
		FeeRecipient: address,
	}

	marshal, _ := sonic.Marshal(postbody)

	resp, err := http.Post(apiRouter, "application/json; charset=utf-8", bytes.NewReader(marshal))
	if err != nil {
		spinnerLiveText.Fail(fmt.Sprintf("Send a request fail: %s", err))
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		spinnerLiveText.Fail(fmt.Sprintf("Read response body error: %s", err))
		return
	}
	var res resultResponse
	err = sonic.Unmarshal(body, &res)
	if err != nil {
		spinnerLiveText.Fail(fmt.Sprintf("Response unmarshal fail! %s", err))
		return
	}

	if res.Code != 200 {
		spinnerLiveText.Fail(fmt.Sprintf("Response Fail! %s", res.Msg))
	} else {
		spinnerLiveText.Success(fmt.Sprintf("Success! %s", res.Msg))
	}

	return
}
