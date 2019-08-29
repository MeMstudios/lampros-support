Description
===
This is a software support middle-ware system that integrates with Gmail, Asana, and Twilio.  
The challenge was to provide our clients with software support without paying an arm and a leg for text notifications with Zendesk.  An added bonus: also keeping our tasks in one project management system.  
The client sends emails to a single email address which forwards a task to Asana.  Once the task is added, Asana sends an event to our server.  The software reads the message and title, applying an urgent tag if it contains: urgent, asap, or important.  After the urgent tag is added (could also be added manually) it will send a notification email, followed by text messages on a schedule until you leave a comment in Asana.  

There used to be more Gmail functionality to get the sender of an email and add them to the task in Asana.  However, the client complained about getting too many Asana emails.  Now, the only thing it does with Gmail is check for email messages matching the software support address in the inbox and remove them from *unread*.  I deleted the logic in the main function.  However, there are still functions in the controllers to add followers to a project, and get senders from a Gmail message.  
Most of the Gmail functionality is external setup, described below.  

Usage
===
Emails
---
Create the group in Google Admin with the email for the client to send tickets to.  
Add contact@lamproslabs.com to the group.  
Create a filter in the contact@lamproslabs.com account to forward emails to the support project in Asana when the to: matches the software support email address and from: matches the client's email domain.  
The support email must be setup as a user in Asana and *TURN OFF EMAIL NOTIFICATIONS* or you will start an infinite email loop!  

All Agents should be added to the Gmail support group as well.  
You should respond to the incoming email by removing the software support address and cc'ing all important parties.  If you don't remove the support address when responding to a ticket, the client will reply, sending another ticket to Asana.  

Add Support Agents
---
Simply, modify the projects.json file and upload to the digissance.com server under `/home/michael/go/src`  
The email should be the support group email and the id is the Asana project id, found in the URL of any project.  I think it's clear what to do with the agents portion.  

The Webhook
---
If you are adding a new support project, you will need to activate a webhook with asana.  You will just need the Asana API access token from our account.  I used postman to POST to https://app.asana.com/api/1.0/webhooks with form-data in the body containing: resource: new_project_id, target: https://our/webhook/endpoint which should be obscured and kept secret.

Notes
---
The urgent tag will only be set automatically if someone emails to the support email with a subject or body containing any form of the word 'urgent', 'asap', or 'important'  
However, the urgent response texts/emails will start to get sent if the urgent tag is set manually.  
Any user with @lamproslabs.com email address should be able to leave a comment to stop the urgent texts.  

Development Setup
===
Install [Golang](https://golang.org/doc/install) and cd into the folder.  

Google Oauth
---
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
When `Environment` is set to `prod` the app will "release the hounds securely" (starting up with SSL). Set it to anything else for local development.  
Prod mode requires you to have the fullchain.pem and privkey.pem in the same folder, or update the code in httpController to have the full path to those files.  
Also contains the twilio number from the twilio account used to send texts.

Golang Notes
---
Please note: all functions that start with lowercase are private to their package.  
You should be able to extend this or fork it interact with most of Asana's API.  
Gmail API is a different pattern because we are using the Golang SDK for it.  
You can do anything with the Gmail API as well, just check out the [documentation](https://godoc.org/google.golang.org/api/gmail/v1)

Asana Notes
---
Recently, Asana changed their webhooks format without warning.  The documentation is still out of date.  So, I have the code I used to debug their request commented out in httpController.  The original format contained IDs as ints and GIDs as strings.  From what I can tell they are moving away from sending ints in JSON and going to all strings, which makes sense, becuase we had to convert ints to strings to call their APIs with the IDs we got.  But now this code has references to Ids and Gids which was my way of distinguishing between strings and ints but now they are mixed.

Deployment
---
It's currently running on a Linux AMI as a service.  
This setup will be different based on what kind of machine you're running.  
If you can get sudo access to the server (currently, Troy's digissance server) where it's hosted you can `service lampros-support stop` or `service lampros-support start` or `service lampros-support restart`  
The binary file is uploaded to `/home/michael/go/src/` you must stop the service to upload and then start again.  
Log can be found: `/var/log/lampros-support.log`  
