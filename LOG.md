### Dec 31, 2025

#### DONE ####
I've written the base code to have a messenger chat bot.
1. Start with creating an app in developers.facebook.com. The current app is called `'golang chat'`.
2. Create a golang app to provide the webhook.
   1. Meta will call something like: `GET /webhook?hub.mode=subscribe&hub.verify_token=...&hub.challenge=...`
   2. If hub.verify_token matches my verify token, respond 200 with the raw hub.challenge string. 
   3. If not, respond 403.
3. In the FB App settings, under Messenger/Messenger API Settings, connect the page to the app and indicate the 
4. webhook subscriptions. Generate an access token from the page to be able to send messages to the page.

#### TODO: ####
Model the dialog. These are AI chat suggestions:[Gemini](https://gemini.google.com/share/09e86e93f129) and
[ChatGPT](https://chatgpt.com/share/6954e9e6-7404-800c-8ef7-b5be54888e58[ChatGPT).
