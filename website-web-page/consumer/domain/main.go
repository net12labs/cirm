package domain

import (
	dom "github.com/net12labs/cirm/ops/domain"
)

func Domain() *Dom {
	return dom.Consumer_WebSite_Page
}

type Dom = dom.WebSitePageDomain
