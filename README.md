#Setup
install golang and cd into the folder
running `go build` should require you to setup Google OAuth for the gmail library. (Follow the instructions from the command line.)
[this may help](https://developers.google.com/gmail/api/auth/web-server)
the credentials.json and token.json will need to be uploaded with the binary

###Emails
Create the group in Google Admin with the email for the client to send tickets to.
Add the team members and contact@lamproslabs.com to the group.
Create a filter in the contact@lamproslabs.com account to forward emails to the support project in Asana.
The support email must be setup as a user in Asana and ####TURN OFF EMAIL NOTIFICATIONS or you will start an infinite email loop!
Client email addresses should be added to the Asana project.

###Credentials
You will need a credentials.go file in the controllers folder which looks something like this:
`
package controllers

const AsanaAccessToken = ""

const EmailAddress = ""
const Password = ""

const TwilioSecret = ""
const TwilioVoiceSID = ""

const TwilioSID = ""
const TwilioAUTH = ""
`
where `EmailAddress` and `Password` are for an email account that will send mail notifications.

###Constants
Contains the recips, which is everyone who should recieve the email notifications.
When `Environment` is set to `prod` the app will "release the hounds securely" (starting up with SSL).
This requires you to have the fullchain.pem and privkey.pem in the same folder, or update the code in httpController to have the full path to those files.



