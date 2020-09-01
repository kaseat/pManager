package sberbank

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/kaseat/pManager/gmail"
	"github.com/kaseat/pManager/models"
	"github.com/kaseat/pManager/models/provider"
	"github.com/kaseat/pManager/storage"
)

var isSync int32

// SyncGmail init sberbank report sync
func SyncGmail(login, pid, from, to string) {
	defer atomic.StoreInt32(&isSync, 0)
	if atomic.LoadInt32(&isSync) == 1 {
		err := errors.New("Sync already in process")
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
		return
	}
	atomic.StoreInt32(&isSync, 1)
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Begin sync sberbank operations via Gmail")
	cl := gmail.GetClient()
	srv, err := cl.GetServiceForUser(login)
	if err != nil {
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
		return
	}

	s := storage.GetStorage()
	t, err := s.GetUserLastUpdateTime(login, provider.Sber)
	if err != nil {
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
		return
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
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
		return
	}

	parsedDates := make(map[string]bool)
	operations := make([]models.Operation, 0)

	for _, m := range r.Messages {
		msg, err := srv.Users.Messages.Get("me", m.Id).Do()

		if err != nil {
			fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
			return
		}
		attachmentID := ""
		for _, p := range msg.Payload.Parts {
			if strings.Contains(p.Filename, ".html") {
				attachmentID = p.Body.AttachmentId
			}
		}

		att, err := srv.Users.Messages.Attachments.Get("me", m.Id, attachmentID).Do()
		if err != nil {
			fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
			return
		}

		b, err := base64.URLEncoding.DecodeString(att.Data)
		if err != nil {
			fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
			return
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
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Parse message on", time.Unix(msg.InternalDate/1000, 0), "ok!")
	}

	if len(operations) != 0 {
		sort.Sort(models.OperationSorter(operations))
		_, err = s.AddOperations(pid, operations)
		if err != nil {
			fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
			return
		}
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "save opertions for", login, "to storge OK")
	}

	err = s.AddUserLastUpdateTime(login, provider.Sber, time.Now())
	if err != nil {
		fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Error sync instruments:", err)
		return
	}
	fmt.Println(time.Now().Format("2006-02-01 15:04:05"), "Success sync sberbank operations via Gmail")
}
