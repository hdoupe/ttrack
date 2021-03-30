# ttrack

Time tracking CLI tool that integrates with Freshbooks.

I use this tool _every single day_, but I only work on it when I'm sufficiently ticked off about a bug or missing feature. If you like to stay in the terminal to log time entries and don't mind a few rough edges, then this may be the tool for you.

## Install

 TODO

## Getting started

1. Get the client id and secret from [me](mailto:henrymdoupe@gmail.com).

2. Sign into freshbooks with this command: `ttrack connect`.

   ```
   $ ttrack connect
   Using client: default

   Go to link:  https://my.freshbooks.com/service/auth/oauth/authorize/?response_type=code&redirect_uri=https://hankdoupe.com/ttrack.html&client_id=9af0623cc6bb6d3717e1c5e73f2f779992ad74e5187e6a1c95e4a651bb2eef0c
   Enter authorization code: 13cafca6bf813e73c039169fe4fc989456bef0214fafe269ea29a824e386438d
   200 OK
   Writing credentials to: /home/hankdoupe/.ttrack.creds.json
   ```

3. Create a file in your home directory named `.ttrack.yaml` with the client ID and secret values:

   ```yaml
   # .ttrack.yaml
   clientID:
   clientSecret:
   ```

4. Create a client and project on Freshbooks if you haven't already.

   - Go to the client page and use the URL get the client ID: https://my.freshbooks.com/#/client/1234 --> 1234 is the client ID.
   - Go to the project page on Freshbooks and use the URL to get the project ID: https://my.freshbooks.com/#/project/5678 --> 5678 is the project ID.

   Connect ttrack with the project using the ID's you just found:

   ```
   $ ttrack clients add --client-id 158233 --project-id 7129723 --nickname my-project
   Using client: default

   Added new client:
   Nickname: my-project
   Client ID: 158233
   Project ID: 7129723
   ```

   List your clients:

   ```
   $ ttrack clients list
   Using client: default

   Nickname: default
   Client ID: 0
   Project ID: 0

   Nickname: my-project
   Client ID: 158233
   Project ID: 7129723
   ```

   Swap to the new client:

   ```
   $ ttrack clients set-current my-project
   Using client: default

   Current client set to: my-project
   ```

5. Create and finish a time entry using current client / project:

   ```
   $ ttrack start -s 2:00PM 'Doing something awesome'
   Using client: my-project

   Creating new time entry...

   Description: Doing something awesome
   Started At: Tue Mar 30 14:30:00 EDT 2021
   Finished At: In progress
   Duration: 0s
   ID: 0 (External ID: 134861033)
   Client ID: 158233
   ```

   ```
   $ ttrack start -s 2:00PM 'Doing something awesome'
   Using client: my-project

   Creating new time entry...

   Description: Doing something awesome
   Started At: Tue Mar 30 14:00:00 EDT 2021
   Finished At: In progress
   Duration: 0s
   ID: 0 (External ID: 134861033)
   Client ID: 158233
   ```

6. List today's time entries:

   ```
   $ ttrack list --since 8:00AM

   # ignore real time entries ;)...

   Description: Doing something awesome
   Started At: Tue Mar 30 14:30:00 EDT 2021
   Finished At: Tue Mar 30 14:43:00 EDT 2021
   Duration: 13m0s
   ID: 478 (External ID: 134861033)
   Client ID: 158233

   Total hours recorded:  1h57m0s
   ```

   or this month's time entries:

   ```
   $ ttrack log --since 2021-03-21

   # ignore real time entries ;)...

   Description: Doing something awesome
   Started At: Tue Mar 30 14:30:00 EDT 2021
   Finished At: Tue Mar 30 14:43:00 EDT 2021
   Duration: 13m0s
   ID: 478 (External ID: 134861033)
   Client ID: 158233

   Total hours recorded:  76h34m0s
   ```
