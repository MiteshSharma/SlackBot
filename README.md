# Go SlackBot

We are running an slack bot that listens to bot messages in a specified channel to reply to user. This connection is created using RTM stream connection with slack.

Slack Channel --> Message to SlackBot --> Parse Message --> Reply --> Go to slack channel

Slack parameters needed:

1. token: Bot User OAuth Token in OAuth & Permissions
2. channelName: Default channel to be used
3. signingSecret: Available in App Credentials
4. clientId: Available in App Credentials
5. clientSecret: Available in App Credentials

Add slack button: https://api.slack.com/docs/slack-button 

Enable Home tab to show home screen:

<img width="672" alt="Screenshot 2022-09-25 at 4 50 27 PM" src="https://user-images.githubusercontent.com/5562910/192140887-8c61cb0b-8562-4da6-bdf1-d817456d81c0.png">

Add interactivity URL for callbacks in section: Interactivity & Shortcuts

URL from current project: /api/v1/slack/interaction 

Add slash comment in section: Slash Commands

Slash comment API: /api/v1/slack/command (All commands handled in same API)

<img width="661" alt="Screenshot 2022-09-25 at 4 52 01 PM" src="https://user-images.githubusercontent.com/5562910/192140931-d6d92774-8523-4713-92da-5f626be4bd82.png">

Enable events API callback on receiving mentions, home page view etc to take action in Event Subscriptions section:

Event API: /api/v1/slack 

<img width="666" alt="Screenshot 2022-09-25 at 4 52 48 PM" src="https://user-images.githubusercontent.com/5562910/192140957-43b69984-8d80-405a-a871-8cecbf4d666f.png">

Subscribe bot for different events that we want to be called in event API webhook specified:

<img width="661" alt="Screenshot 2022-09-25 at 4 53 02 PM" src="https://user-images.githubusercontent.com/5562910/192140961-cba00e69-bbdb-453b-acce-fd23a2bf95a2.png">
