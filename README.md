Setup
===
install golang and cd into the folder. </br>
running `go build` should require you to setup Google OAuth for the gmail library. (Follow the instructions from the command line.)
[this may help](https://developers.google.com/gmail/api/auth/web-server)</br>
the credentials.json and token.json will need to be uploaded with the binary. </br>
These are referenced by their full path on the server.</br>

Emails
---
Create the group in Google Admin with the email for the client to send tickets to.</br>
Add the team members and contact@lamproslabs.com to the group.</br>
Create a filter in the contact@lamproslabs.com account to forward emails to the support project in Asana.</br>
The support email must be setup as a user in Asana and *TURN OFF EMAIL NOTIFICATIONS* or you will start an infinite email loop!</br>
Client email addresses should be added to the Asana project.</br>

Credentials
---
You will need a credentials.go file in the controllers folder which looks something like this:</br>

`package controllers`

`const AsanaAccessToken = ""`

`const EmailAddress = ""`

`const Password = ""`

`const TwilioSecret = ""`

`const TwilioVoiceSID = ""`

`const TwilioSID = ""`

`const TwilioAUTH = ""`

where `EmailAddress` and `Password` are for an email account that will send mail notifications.</br>

Constants
---
Contains the recips, which is everyone who should recieve the email notifications.</br>
When `Environment` is set to `prod` the app will "release the hounds securely" (starting up with SSL).</br>
This requires you to have the fullchain.pem and privkey.pem in the same folder, or update the code in httpController to have the full path to those files.</br>
Constants also contains the recipients email list (same as the software support group) and the phone numbers we'll use to send twilio notifications.</br>

Deployment
---
It's currently running on a Linux AMI machine as a service. </br>
This setup will be different based on what kind of machine you're running. </br>
If you can get sudo access to the server where it's hosted you can `service lampros-support stop` or `service lampros-support start` or `service lampros-support restart` </br>
The binary file is uploaded to `/home/michael/go/src/` you must stop the service to upload and then start again. </br>
Log can be found: `/var/log/lampros-support.log`</br>

Usage
---
You should not immediately respond to the initial email.  Wait until you get the software support notification and take the link to the Asana task and leave a comment there to respond.  Now everyone will be on the email list.  contact@lamproslabs.com and the support email address should not get Asana notifications.

