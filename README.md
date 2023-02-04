# autoscaler-service
This is a basic autoscaler service written in Golang.

Note: This is my first time using Go.

### How to Build It
```go run .```

### Coding Decisions
I opted to keep the overall structure relatively simple and contained from the start, but I also knew I wanted to separate out the code contained in ```api.go``` and ```deploymentservice.go``` at least into separate files. This both clarifies the boundaries in the code as well as contains the work required if the api the autoscaler hits ever changes.

Using multiple goroutines made the most sense here but, having never used Go before, it took some reading as well as some trial and error. I started out with a more succinct approach, nesting goroutines. It meant fewer lines of code, but the readability didn't feel as clean and using multiple goroutines with channels would be a safer, clearer route.

The first goroutine in ```autoscale()``` runs the timer responsible for the 60 second interval of the DiscoverDeployments API call. It also sends the discovered deployments along to the activeDeployments channel where the second goroutine can pick it up. It then waits for the specified 30 second interval before sending the same deployments again to the activeDeployments channel, triggering the second goroutine again.

The second goroutine in ```autoscale()``` will receive from the activeDeployments channel and immediately iterate through the deployments, spawning a goroutine for each that will check its status and, if necessary, scale it.

In ```autoscaler.go```, I opted to pull the autoscaler logic itself out of main and into its own function for clarity when reading the code.

The ```deploymentservice.go``` file contains code handling the calls to ```api.go``` as well as the logic behind the scaling delta calculation. This code does assume that we want to check every case where the targetFreePct doesn't match the deployment's current free percentage of servers, even if the free percentage is only an extremely tiny percentage above the target and unlikely to benefit from scaling. A way to configure a limit would be easy to implement, however.

```api.go``` contains all of the calls out to the deployment api and handles reading, deserialization/serialization, and errors from the api.

```helpers.go``` contains some helpful wrapper functions that were the result of refactoring some repetitive error handling in ```api.go```.

Currently, the autoscaler is set up to log errors in ```autoscaler.log``` with a basic setup of the built-in log package and continue on.

In addition to tests in the repo, I also created a mock server with Postman to test against.

### TODO:
- Write Tests
- Ensure proper error handling for api calls
  - optional error field