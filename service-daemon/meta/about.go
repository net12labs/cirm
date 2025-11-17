package meta

/*

users come here for a reason:
-- improve their network performance (in some mission/project/job they do that we do not know about)
-- save money etc.
-- provide convenience


it can be a eason or a whole reason group

so as a provider - we only post a specific path which we resolve

THERE IS A REASON WHY THEY PICK ONE LINK OVER ANOTHER
and then the providers also have a reason why they promote one link over another (link=method, combination)

--- baobab always gives resons for people to shop
-- we can let them set up pointers and tags etc.
*/

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

/*

So the ACL ing logic will be very important
we can have users only join a specific project, or even task
we can use include and exclude rules

Also - different types of shells for different tasks

A lot of tasks could be virtual tasks -- where you need to edit an image - but the agent puts it in a filder and opens gimp
or redirects you to a temporary station where you have the heavy tools installed

*/
