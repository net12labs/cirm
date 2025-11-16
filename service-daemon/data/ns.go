package data

import (
	"github.com/net12labs/cirm/service-daemon/data/admin"
	"github.com/net12labs/cirm/service-daemon/data/service"
)

var Service = service.NewService()
var Admin = admin.NewAdmin()

// so for the basic user logic:
/*
context entity - such as organization, platform, user
user (class)
user has profiles (type)
profiles are of specific TYPE - admin, provider, user, platform, root etc.
and the there are profile kinds - and kinds have permissions, capabilities etc.
kinds are used inside a specific organization or context/entity

a profile can only be used in a specific context - such as provider/admn/root etc.

so from  the session cookie - what we know the most is the actual user - but all the other details need to be passed with a request

so -- as just a user - I can do certain things (like manage my personal top level account)
as a user profile - I can do more - for example shop and buy things, join entities
as a kind - I gain attributes and capabilities


*/
