# autoscaler-service



### Requirements
- every 1 minute, get /discover to get json object containing deployments[]
- every 30 seconds, get /status/{id} for each deployment.
  -  for every deployment that needs action taken, post /scale with delta and deploymentid
- exit cleanly on a SIGINT
- long-running daemon