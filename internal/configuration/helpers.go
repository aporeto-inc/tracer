package configuration

import (
	"fmt"
	"net/url"
	"os"

	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/pflag"
)

// showHelp show a full help
func showHelp() {

	fmt.Println("Usage:")
	pflag.PrintDefaults()
	fmt.Printf(`

> Display all queries with traces from the last 1h

  ./tracer --since 1h

> Display all queries for a service from the last 1h

  ./tracer --since 1h --service squall

> Display all queries for a service in a given namespace from the last 1h

  ./tracer --since 1h --service squall --namespace /foo/bar

> Display all queries for a service in a given namespace that took more than 2s from the last 1h

  ./tracer --since 1h --service squall --namespace /foo/bar --slower-than 2s

> Display all requests that returns with an error for the past hour

  ./tracer --since 1h --errors-only

> Display all requests that return with a code 200 or 400-422 in the past hour

  ./tracer --since 1h --code 200,400-422

> Display all requests made to /flowreports

  ./tracer --since 1h --url /flowreports

> Display all 400-403 requests on service squall, cid and /issue between two dates

  ./tracer --code 400-403 --service squal --service cid --url /issue --from 2020-10-21T17:56:17Z --to 2020-10-22T17:56:17Z

> Display logs for 2 services between two dates

  ./tracer --log --service squal --service cid --from 2020-10-21T17:56:17Z --to 2020-10-22T17:56:17Z

Some queries are not providing traces (like reports because this is too much for jaeger to handle).
In general errors are logged in the service in debug mode. Use the switch-debug <service name>  command to enable it.
And look at the logs either through Grafana->Explore->Loki or with the k get log <pod_name> command.
`)
	os.Exit(0)
}

// openTrace will open a trace in the browser
func openTrace(u, trace string) {

	toOpen, err := url.Parse(u)
	if err != nil {
		fmt.Println("Invalid url: ", err)
		return
	}

	toOpen.Path = "explore"
	q, _ := url.ParseQuery(toOpen.RawQuery)
	q.Add("orgId", "1")
	q.Add("left", fmt.Sprintf(`["now-1h","now","jaeger-aporeto",{"query":"%s"},{"ui":[true,true,true,"none"]}]`, trace))

	toOpen.RawQuery = q.Encode()

	err = open.Run(toOpen.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Open this URL in your browser:", toOpen)
	}
	os.Exit(0)
}
