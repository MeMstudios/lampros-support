package controllers

const Environment = "prod" //checks for prod, anything else starts a normal server

const AsanaBase = "https://app.asana.com/api/1.0"
const SupportProjectID = "1107602815071855"
const UrgentTagGid = "1107602815071859"
const NewTagGid = "1107602815071860"
const PendingTagGid = "1107602815071861"
const ResolvedTagGid = "1107602815071862"
const SupportEmailAddress = "testsoftwaresupport@lamproslabs.com"

//EMAIL RECIPIENTS
var recips = []string{"michael@lamproslabs.com"}

const TwilioNumber = "+18594465953"
const TwiML = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
	"<Response>\n" +
	"<Say voice=\"man\">You have an urgent support ticket that hasn't been responded to.  Please respond or face the consequences!</Say>\n" +
	"<Record maxLength=\"20\" />\n" +
	"</Response>\n"

const TwilioBase = "https://api.twilio.com/2010-04-01/Accounts/" + TwilioSID
