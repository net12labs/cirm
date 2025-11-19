package domain

// Admin
var Admin = &Domain{BasePath: "/admin"}

var Admin_Agent = NewAgentDomain(Admin, "agent")
var Admin_AiAgent = NewAiAgentDomain(Admin, "ai-agent")
var Admin_WebSite = NewWebSiteDomain(Admin, "site")

var Admin_AiAgent_Page = NewAiAgentPageDomain(Admin_AiAgent, "at")
var Admin_AiAgent_Web = NewAiAgentWebDomain(Admin_AiAgent, "web")
var Admin_AiAgent_Api = NewAiAgentApiDomain(Admin_AiAgent, "api")

var Admin_Agent_Page = NewAgentPageDomain(Admin_Agent, "at")
var Admin_Agent_Web = NewAgentWebDomain(Admin_Agent, "web")
var Admin_Agent_Api = NewAgentApiDomain(Admin_Agent, "api")

var Admin_WebSite_Page = NewWebSitePageDomain(Admin_WebSite, "home")
var Admin_WebSite_Web = NewWebSiteWebDomain(Admin_WebSite, "web")
var Admin_WebSite_Api = NewWebSiteApiDomain(Admin_WebSite, "api")

// Platform
var Platform = &Domain{BasePath: "/platform"}

var Platform_Agent = NewAgentDomain(Platform, "agent")
var Platform_AiAgent = NewAiAgentDomain(Platform, "ai-agent")
var Platform_WebSite = NewWebSiteDomain(Platform, "site")

var Platform_AiAgent_Page = NewAiAgentPageDomain(Platform_AiAgent, "at")
var Platform_AiAgent_Web = NewAiAgentWebDomain(Platform_AiAgent, "web")
var Platform_AiAgent_Api = NewAiAgentApiDomain(Platform_AiAgent, "api")

var Platform_Agent_Page = NewAgentPageDomain(Platform_Agent, "at")
var Platform_Agent_Web = NewAgentWebDomain(Platform_Agent, "web")
var Platform_Agent_Api = NewAgentApiDomain(Platform_Agent, "api")

var Platform_WebSite_Page = NewWebSitePageDomain(Platform_WebSite, "home")
var Platform_WebSite_Web = NewWebSiteWebDomain(Platform_WebSite, "web")
var Platform_WebSite_Api = NewWebSiteApiDomain(Platform_WebSite, "api")

// Consumer
var Consumer = &Domain{BasePath: "/consumer"}

var Consumer_Agent = NewAgentDomain(Consumer, "agent")
var Consumer_AiAgent = NewAiAgentDomain(Consumer, "ai-agent")
var Consumer_WebSite = NewWebSiteDomain(Consumer, "site")

var Consumer_AiAgent_Page = NewAiAgentPageDomain(Consumer_AiAgent, "at")
var Consumer_AiAgent_Web = NewAiAgentWebDomain(Consumer_AiAgent, "web")
var Consumer_AiAgent_Api = NewAiAgentApiDomain(Consumer_AiAgent, "api")

var Consumer_Agent_Page = NewAgentPageDomain(Consumer_Agent, "at")
var Consumer_Agent_Web = NewAgentWebDomain(Consumer_Agent, "web")
var Consumer_Agent_Api = NewAgentApiDomain(Consumer_Agent, "api")

var Consumer_WebSite_Web = NewWebSiteWebDomain(Consumer_WebSite, "web")
var Consumer_WebSite_Api = NewWebSiteApiDomain(Consumer_WebSite, "api")
var Consumer_WebSite_Page = NewWebSitePageDomain(Consumer_WebSite, "home")

// Provider
var Provider = &Domain{BasePath: "/provider"}

var Provider_Agent = NewAgentDomain(Provider, "agent")
var Provider_AiAgent = NewAiAgentDomain(Provider, "ai-agent")
var Provider_WebSite = NewWebSiteDomain(Provider, "site")

var Provider_AiAgent_Page = NewAiAgentPageDomain(Provider_AiAgent, "at")
var Provider_AiAgent_Web = NewAiAgentWebDomain(Provider_AiAgent, "web")
var Provider_AiAgent_Api = NewAiAgentApiDomain(Provider_AiAgent, "api")

var Provider_Agent_Page = NewAgentPageDomain(Provider_Agent, "at")
var Provider_Agent_Web = NewAgentWebDomain(Provider_Agent, "web")
var Provider_Agent_Api = NewAgentApiDomain(Provider_Agent, "api")

var Provider_WebSite_Page = NewWebSitePageDomain(Provider_WebSite, "home")
var Provider_WebSite_Web = NewWebSiteWebDomain(Provider_WebSite, "web")
var Provider_WebSite_Api = NewWebSiteApiDomain(Provider_WebSite, "api")

// Site
var Site = &Domain{BasePath: "/site"}

var Site_Agent = NewAgentDomain(Site, "agent")
var Site_AiAgent = NewAiAgentDomain(Site, "ai-agent")
var Site_WebSite = NewWebSiteDomain(Site, "site")

var Site_AiAgent_Page = NewAiAgentPageDomain(Site_AiAgent, "at")
var Site_AiAgent_Web = NewAiAgentWebDomain(Site_AiAgent, "web")
var Site_AiAgent_Api = NewAiAgentApiDomain(Site_AiAgent, "api")

var Site_Agent_Page = NewAgentPageDomain(Site_Agent, "at")
var Site_Agent_Web = NewAgentWebDomain(Site_Agent, "web")
var Site_Agent_Api = NewAgentApiDomain(Site_Agent, "api")

var Site_WebSite_Page = NewWebSitePageDomain(Site_WebSite, "home")
var Site_WebSite_Web = NewWebSiteWebDomain(Site_WebSite, "web")
var Site_WebSite_Api = NewWebSiteApiDomain(Site_WebSite, "api")
