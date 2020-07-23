package sberbank

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/kaseat/pManager/gmail"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/storage"
)

// SyncGmail init sberbank report sync
func SyncGmail(login, pid, from, to string) error {
	cl := gmail.GetClient()
	srv, err := cl.GetServiceForUser(login)
	if err != nil {
		return err
	}

	s := storage.GetStorage()
	t, err := s.GetUserLastUpdateTime(login, "sberbank")
	if err != nil {
		return err
	}

	query := "from:broker_rep@sberbank.ru subject:report filename:html"
	if from == "" && to == "" && !t.IsZero() {
		query = fmt.Sprintf("%s after:%s", query, t.Format("2006/01/02"))
	}
	if from != "" {
		query = fmt.Sprintf("%s after:%s", query, from)
	}
	if to != "" {
		query = fmt.Sprintf("%s before:%s", query, to)
	}

	r, err := srv.Users.Messages.List("me").Q(query).Do()
	if err != nil {
		return err
	}

	parsedDates := make(map[string]bool)
	operations := make([]models.Operation, 0)

	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Do()

		if err != nil {
			return err
		}
		attachmentID := ""
		for _, p := range msg.Payload.Parts {
			if strings.Contains(p.Filename, ".html") {
				attachmentID = p.Body.AttachmentId
			}
		}

		att, err := srv.Users.Messages.Attachments.Get("me", m.Id, attachmentID).Do()
		if err != nil {
			return err
		}

		b, err := base64.URLEncoding.DecodeString(att.Data)
		if err != nil {
			return err
		}

		reader := bytes.NewReader(b)
		report := parseReport(reader)
		if !report.IsEmpty {
			if !parsedDates[report.Date] {
				parsedDates[report.Date] = true
				if len(report.Operations) > 0 {
					operations = append(operations, report.Operations...)
				}
				if len(report.Buybacks) > 0 {
					operations = append(operations, report.Buybacks...)
				}
				if len(report.CashFlow) > 0 {
					operations = append(operations, report.CashFlow...)
				}
			}
		}
		fmt.Println("Parse message on", time.Unix(msg.InternalDate/1000, 0), "ok!")
	}

	if len(operations) != 0 {
		sort.Sort(models.OperationSorter(operations))
		_, err = s.AddOperations(pid, operations)
		if err != nil {
			return err
		}
		fmt.Println("save opertions for", login, "to storge OK")
	}

	err = s.AddUserLastUpdateTime(login, "sberbank", time.Now())
	if err != nil {
		return err
	}
	return nil
}
