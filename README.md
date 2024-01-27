#### To run: 

Include a ServiceAccountKey.json in the root directly with access to a firestore DB (or add your own DB logic)
Add WarcraftLogs credentials for ClientID and ClientSecret to internal/logs/warcraftlogs/client.go

#### Endpoints:

### `/generate/encounter` 
Calls WarcraftLogs for an encounter included in the request_body `encounter_id` which then runs analysis on the response which is then written to firestore. A field `persist` can also be provided which if `true` will continue to run analysis even when the API reaches an encounter it has already analysed before. This is useful if the encounter_id is being analysed for the first couple of times.

### `/encounter/{encounterName}/healer-frequency` 
Calls the firestore table created in the previous request and returns analysis about the healer-frequency. 

#### Notes:

The concept of the API is to run the `generate` endpoints once a day as a cron job to update the analysis. And to call the `encounter` endpoints whenever a request is made to view the analysis, for example, from a web-app. 

The reason for this is to keep the requests to WarcraftLogs to a minimum as it's expensive.
