Setup
===
install golang and cd into the folder.  
running `go build` should require you to setup Google OAuth for the gmail library. (Follow the instructions from the command line.)  
[this may help](https://developers.google.com/gmail/api/auth/web-server)  
the credentials.json and token.json will need to be uploaded with the binary.  
These are referenced by their full path on the server.  

Emails
---
Create the group in Google Admin with the email for the client to send tickets to.  
Add contact@lamproslabs.com to the group.  
Create a filter in the contact@lamproslabs.com account to forward emails to the support project in Asana when the to: matches the software support email address.  
The support email must be setup as a user in Asana and *TURN OFF EMAIL NOTIFICATIONS* or you will start an infinite email loop!  
Client email addresses should be added to the Asana project.  

Credentials
---
You will need a credentials.go file in the controllers folder which looks something like this:  

`package controllers`

`const AsanaAccessToken = ""`

`const EmailAddress = ""`

`const Password = ""`

`const TwilioSecret = ""`

`const TwilioVoiceSID = ""`

`const TwilioSID = ""`

`const TwilioAUTH = ""`

where `EmailAddress` and `Password` are for an email account that will send mail notifications.  

Constants
---
Contains the recips, which is everyone who should recieve the email notifications.  
When `Environment` is set to `prod` the app will "release the hounds securely" (starting up with SSL).  
This requires you to have the fullchain.pem and privkey.pem in the same folder, or update the code in httpController to have the full path to those files.  
Constants also contains the recipients email list (same as the software support group) and the phone numbers we'll use to send twilio notifications.  

Deployment
---
It's currently running on a Linux AMI machine as a service.  
This setup will be different based on what kind of machine you're running.  
If you can get sudo access to the server where it's hosted you can `service lampros-support stop` or `service lampros-support start` or `service lampros-support restart`  
The binary file is uploaded to `/home/michael/go/src/` you must stop the service to upload and then start again.  
Log can be found: `/var/log/lampros-support.log`  

Usage
---
You should not immediately respond to the initial email.  Wait until you get the software support notification and take the link to the Asana task and leave a comment there to respond.  Now everyone will be on the email list.  contact@lamproslabs.com and the support email address should not get Asana notifications.  
Not all the support agents need to be added to the support group.  Troy gets emails from contact@lamproslabs.com anyway.  
POST to the /add-agent endpoint to add an email and phone number send json body like:  
`{"email": "email@address.com", "phone": "+1234567890"}`
and set a header called: `Api-Key` using the OurApiKey in constants.go  
This is just in memory, so if you start and stop the service you'll have to add agents again, so just save a request in postman.

