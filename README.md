# autoscaler-service
This is a basic autoscaler service written in Golang.

## How to Build It
#### To test, run one of the following:
```
go test .
go test . -v
```

#### To run quickly for dev purposeswith the following command:
```
go run .
```

#### To build and install (to the Go install path):
```
go build
go install
```

## Coding Decisions
I opted to keep the overall structure relatively simple and contained from the start. However, I also knew I wanted to separate the code in a way that would clarify the boundaries without adding unnecessary complexity. This would also serve to contain the work required if the api or the underlying scaling logic were to ever change.

As far as the actual operation of the autoscaler, I thought that using multiple goroutines made the most sense here but, having never used Go before, it took some reading as well as some time spent experimenting with them. I started out with a more succinct approach, nesting goroutines. It meant fewer lines of code, but the readability didn't feel as clean and using multiple goroutines with channels would be a safer, clearer route.

The first goroutine in ```autoscale()``` runs the timer responsible for the 60 second interval of the DiscoverDeployments API call. It also sends the discovered deployments along to the activeDeployments channel where the second goroutine can pick it up. It then waits for the specified 30 second interval before sending the same deployments again to the activeDeployments channel, triggering the second goroutine again. The second goroutine in ```autoscale()```, upon receiving from the activeDeployments channel,  then immediately iterates through the deployments, spawning a goroutine for each that will check its status and, if necessary, scale it.

Another aspect I spent time considering my options for was the timing of the goroutines. The implementation I went with is rather "synchronized." The discover goroutine calls out to the api every minute and then handles sending to the ```activeDeployments``` channel once, waits 30 seconds, and then seconds it again. All of the timing is handled in the first goroutine and the second goroutine doesn't do anything unless something has been sent to the ```activeDeployments``` channel by the first.

Alternative implementations could pull that portion out into its own goroutine or you could implement it such that the discovery goroutine runs on its own 60 second interval while the scale goroutine runs independently on its own 30 second interval. Ultimately, I opted for the implementation I went with as I didn't see any special benefit to the alternatives and refactoring it to another approach wouldn't be too difficult if deemed necessary.

Currently, the autoscaler is set up to log errors in ```autoscaler.log``` with a basic setup of the built-in log package. In addition to the unit tests in the repo, I also created a mock server with Postman to test against. Some of the things I would likely look at next for this would be end to end testing, improving logging with log levels, and implementing some timeouts on api calls and logging for that.

In all, I spent somewhere around 4.5 hours in total on this project. Around 1.5-2 hours was spent on the goroutines, timing, etc in ```autoscaler.go```. Around 1.5 hours spent on writing tests, designing, implementing, and debugging the delta calculation. The remainder of the time spent on the project was split on the various parts, refactoring, adding doc comments, etc.