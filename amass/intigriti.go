// Copyright FUCK OFF. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

package amass

import (
	"fmt"

	"github.com/Ergotu/Amass/amass/utils"
)

// Intigriti is the AmassService that handles access to the Intigriti data source.
type Intigriti struct {
	BaseAmassService

	SourceType string
}

// NewIntigriti returns he object initialized, but not yet started.
func NewIntigriti(e *Enumeration) *Intigriti {
	h := &Intigriti{SourceType: API}

	h.BaseAmassService = *NewBaseAmassService(e, "Intigriti", h)
	return h
}

// OnStart implements the AmassService interface
func (h *Intigriti) OnStart() error {
	h.BaseAmassService.OnStart()

	go h.startRootDomains()
	go h.processRequests()
	return nil
}

// OnStop implements the AmassService interface
func (h *Intigriti) OnStop() error {
	h.BaseAmassService.OnStop()
	return nil
}

func (h *Intigriti) processRequests() {
	for {
		select {
		case <-h.PauseChan():
			<-h.ResumeChan()
		case <-h.Quit():
			return
		case <-h.RequestChan():
			// This data source just throws away the checked DNS names
			h.SetActive()
		}
	}
}

func (h *Intigriti) startRootDomains() {
	// Look at each domain provided by the config
	for _, domain := range h.Enum().Config.Domains() {
		h.executeQuery(domain)
	}
}

func (h *Intigriti) executeQuery(domain string) {
	url := h.getURL(domain)
	page, err := utils.RequestWebPage(url, nil, nil, "", "")
	if err != nil {
		h.Enum().Log.Printf("%s: %s: %v", h.String(), url, err)
		return
	}

	h.SetActive()
	re := h.Enum().Config.DomainRegex(domain)
	for _, sd := range re.FindAllString(page, -1) {
		h.Enum().NewNameEvent(&AmassRequest{
			Name:   cleanName(sd),
			Domain: domain,
			Tag:    h.SourceType,
			Source: h.String(),
		})
	}
}

func (h *Intigriti) getURL(domain string) string {
	format := "http://localhost:8080/subdomains/%s"

	return fmt.Sprintf(format, domain)
}
