package main

import (
	"fmt"
	"log"

	"github.com/elastic/go-libaudit"
	"github.com/elastic/go-libaudit/auparse"
)

func receive(r *libaudit.AuditClient) error {
	for {
		rawEvent, err := r.Receive(false)
		if err != nil {
			return fmt.Errorf("receive failed: %w", err)
		}

		// Messages from 1300-2999 are valid audit messages.
		if rawEvent.Type < auparse.AUDIT_USER_AUTH ||
			rawEvent.Type > auparse.AUDIT_LAST_USER_MSG2 {
			continue
		}

		log.Printf("type=%v msg=%s\n", rawEvent.Type, rawEvent.Data)
	}
}

func main() {
	var err error
	var client *libaudit.AuditClient
	client, err = libaudit.NewMulticastAuditClient(nil)
	if err != nil {
		log.Fatal("failed to create receive-only audit client: %w", err)
	}
	defer client.Close()
	log.Println("Listening to audit events on Multicast socket ...")
	status, _ := client.GetStatus()
	log.Println(status.Enabled)
	receive(client)
}
