Usage
===
Emails
---
Create the group in Google Admin with the email for the client to send tickets to.  
Add contact@lamproslabs.com to the group.  
Create a filter in the contact@lamproslabs.com account to forward emails to the support project in Asana when the to: matches the software support email address.  
The support email must be setup as a user in Asana and *TURN OFF EMAIL NOTIFICATIONS* or you will start an infinite email loop!  
Client email addresses should be added to the Asana project.  

If are in the gmail support group, you should not immediately respond to the initial email.  Wait until you get the software support notification and take the link to the Asana task and leave a comment there to respond.  Now everyone will be on the email list.  
contact@lamproslabs.com and the support email addresses should not get Asana notifications.  
Not all the support agents need to be added to the gmail support group.  Troy gets emails from contact@lamproslabs.com anyway.  Agents will get emails and texts if you add them to the projects.json file.  

Add Support Agents
---
Simply, modify the projects.json file and upload to the digissance.com server under `/home/michael/go/src`  
The email should be the support group email and the id is the Asana project id, found in the URL of any project.  I think it's clear what to do with the agents portion.  

Notes
---
At this time it does not check for if a COMMENT was added by a customer that could indicate the task is urgent.  
The urgent tag will only be set automatically if someone emails to the support email with a subject or body containing any form of the word 'urgent'  
However, the urgent response texts/emails will start to get sent if the urgent tag is set manually.  

Setup
===
install golang and cd into the folder.  
running `go build` should require you to setup Google OAuth for the gmail library. (Follow the instructions from the command line.)  
[this may help](https://developers.google.com/gmail/api/auth/web-server)  
the credentials.json and token.json will need to be uploaded with the binary.  
These are referenced by their full path on the server.  

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
When `Environment` is set to `prod` the app will "release the hounds securely" (starting up with SSL).  
This requires you to have the fullchain.pem and privkey.pem in the same folder, or update the code in httpController to have the full path to those files.  
The AsanaSupportId will have to be changed if the support project is changed.  
Also contains the twilio number from the twilio account used to send texts.

Deployment
---
It's currently running on a Linux AMI machine as a service.  
This setup will be different based on what kind of machine you're running.  
If you can get sudo access to the server (currently, Troy's digissance server) where it's hosted you can `service lampros-support stop` or `service lampros-support start` or `service lampros-support restart`  
The binary file is uploaded to `/home/michael/go/src/` you must stop the service to upload and then start again.  
Log can be found: `/var/log/lampros-support.log`  
