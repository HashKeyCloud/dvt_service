package email

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

func (e *Email) SendRegisterValidatorErrorMail(pubkey, OperatorsIds string, errinfo error) error {
	content := fmt.Sprintf(RegisterValidatorErrorTemplate, pubkey, OperatorsIds, errinfo.Error())
	log.Info().Str("Type", "RegisterValidatorError").
		Str("pubkey", pubkey).
		Str("operatorsIds", OperatorsIds).
		Msg("Send Mail")
	return e.sendToMail(RegisterValidatorErrorSubject, content)
}

const (
	RegisterValidatorErrorSubject  = "Notice: DVT RegisterValidator Error"
	RegisterValidatorErrorTemplate = `<div style="font-size: 30px;font-weight: bolder">Notice: Voluntary Exit Fail</div>
    <div style="margin-top: 30px">
            Pubkey: %s
		<br/>
    		OperatorsIds: %s
		<br/>
			Error: %s
    </div>`
)
