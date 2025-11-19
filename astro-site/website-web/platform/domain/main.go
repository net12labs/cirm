package domain

import (
	dom "github.com/net12labs/cirm/ops/domain"
)

func Domain() *Dom {
	return dom.Platform_WebSite_Web
}

type Dom = dom.WebSiteWebDomain
