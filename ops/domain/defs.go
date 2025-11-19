package domain

import "strconv"

type Domain struct {
	BasePath string
}

func (d *Domain) Path() string {
	return d.BasePath
}
func (d *Domain) MakePath(subpath string) string {
	return d.BasePath + "/" + subpath
}

type AiAgentDomain struct {
	*Domain
}

func NewAiAgentDomain(dom *Domain, path string) *AiAgentDomain {
	return &AiAgentDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

type AgentDomain struct {
	*Domain
}

func NewAgentDomain(dom *Domain, path string) *AgentDomain {
	return &AgentDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

type WebSiteDomain struct {
	*Domain
}

type PageDomain struct {
	*Domain
}

func NewWebSiteDomain(dom *Domain, path string) *WebSiteDomain {
	return &WebSiteDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

type AiAgentPageDomain struct {
	*Domain
	*PageDomain
}

func NewAiAgentPageDomain(dom *AiAgentDomain, path string) *AiAgentPageDomain {
	p := &AiAgentPageDomain{
		Domain:     &Domain{BasePath: dom.Path() + "/" + path},
		PageDomain: &PageDomain{},
	}
	p.PageDomain.Domain = p.Domain
	return p
}

type WebSitePageDomain struct {
	*Domain
	*PageDomain
}

func NewWebSitePageDomain(dom *WebSiteDomain, path string) *WebSitePageDomain {
	p := &WebSitePageDomain{
		Domain:     &Domain{BasePath: dom.Path() + "/" + path},
		PageDomain: &PageDomain{},
	}
	p.PageDomain.Domain = p.Domain
	return p
}

type AgentPageDomain struct {
	*Domain
	*PageDomain
}

func NewAgentPageDomain(dom *AgentDomain, path string) *AgentPageDomain {
	p := &AgentPageDomain{
		Domain:     &Domain{BasePath: dom.Path() + "/" + path},
		PageDomain: &PageDomain{},
	}
	p.PageDomain.Domain = p.Domain
	return p
}

func (d *PageDomain) WrapHTML(data []byte, path string, id int64) []byte {
	head := []byte("<x-domain path=\"" + d.Path() + "/" + path + "/" + strconv.FormatInt(id, 10) + "\">")
	tail := []byte("</x-domain>")
	result := make([]byte, 0, len(head)+len(data)+len(tail))
	result = append(result, head...)
	result = append(result, data...)
	result = append(result, tail...)
	return result
}

type AiAgentWebDomain struct {
	*Domain
}

func NewAiAgentWebDomain(dom *AiAgentDomain, path string) *AiAgentWebDomain {
	return &AiAgentWebDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

type AgentWebDomain struct {
	*Domain
}

func NewAgentWebDomain(dom *AgentDomain, path string) *AgentWebDomain {
	return &AgentWebDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

type WebSiteWebDomain struct {
	*Domain
}

func NewWebSiteWebDomain(dom *WebSiteDomain, path string) *WebSiteWebDomain {
	return &WebSiteWebDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

// api domains

type AiAgentApiDomain struct {
	*Domain
}

func NewAiAgentApiDomain(dom *AiAgentDomain, path string) *AiAgentApiDomain {
	return &AiAgentApiDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

type AgentApiDomain struct {
	*Domain
}

func NewAgentApiDomain(dom *AgentDomain, path string) *AgentApiDomain {
	return &AgentApiDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}

type WebSiteApiDomain struct {
	*Domain
}

func NewWebSiteApiDomain(dom *WebSiteDomain, path string) *WebSiteApiDomain {
	return &WebSiteApiDomain{
		Domain: &Domain{BasePath: dom.Path() + "/" + path},
	}
}
