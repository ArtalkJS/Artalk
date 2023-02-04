package email

import (
	"encoding/json"
	"sync"

	"github.com/ArtalkJS/Artalk/internal/config"
	"github.com/ArtalkJS/Artalk/internal/query"
	"github.com/sirupsen/logrus"
)

// Email Queue
var emailCh chan Email
var emailMutex sync.Mutex

// Initialize Email Queue
func InitQueue(emailConf *config.EmailConf) {
	emailMutex.Lock()
	defer emailMutex.Unlock()

	// email queue only need init once
	if emailCh != nil {
		return
	}

	// init email queue
	emailCh = make(chan Email, emailConf.Queue.BufferSize)

	if config.Instance.Debug {
		logrus.Debug("[Email] Email Queue initialize complete")
	}

	go func() {
		for {
			email := <-emailCh
			sender := NewSender(emailConf.SendType)

			if config.Instance.Debug {
				emailJson, _ := json.Marshal(email)
				logrus.Debug("[Email] Sending an email: ", string(emailJson))
			}

			// send email
			isOK := sender.Send(email)

			if !isOK {
				logrus.Errorf("[Email] Failed send email to addr: %s", email.ToAddr)
				continue
			}

			// send success
			if email.LinkedNotify != nil {
				if err := query.NotifySetEmailed(email.LinkedNotify); err != nil { // flag associated comment as emailed
					logrus.Errorf("[Email] Flag email delivery status for associated comment failed: %s", err)
				}
			}
		}
	}()
}

// Add an email to the sending queue
func AddToQueue(email Email) {
	emailCh <- email
}
