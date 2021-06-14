// Code generated by goa v3.4.2, DO NOT EDIT.
//
// sunflower gRPC client CLI support package
//
// Command:
// $ goa gen sunflower/pkg/api/design -o pkg/api/

package cli

import (
	"flag"
	"fmt"
	"os"
	scorec "sunflower/pkg/api/gen/grpc/score/client"

	goa "goa.design/goa/v3/pkg"
	grpc "google.golang.org/grpc"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `score (score-list|score-detail)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` score score-list --message '{
      "class": "Nisi aliquam est consequatur quod ea vitae.",
      "cursor": 0,
      "limit": 20,
      "name": "Est vel.",
      "scores": 8275923798665478266,
      "sortField": "score",
      "sortOrder": "asc",
      "subject": "Sed aut labore eum placeat quo corporis."
   }'` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(cc *grpc.ClientConn, opts ...grpc.CallOption) (goa.Endpoint, interface{}, error) {
	var (
		scoreFlags = flag.NewFlagSet("score", flag.ContinueOnError)

		scoreScoreListFlags       = flag.NewFlagSet("score-list", flag.ExitOnError)
		scoreScoreListMessageFlag = scoreScoreListFlags.String("message", "", "")

		scoreScoreDetailFlags       = flag.NewFlagSet("score-detail", flag.ExitOnError)
		scoreScoreDetailMessageFlag = scoreScoreDetailFlags.String("message", "", "")
	)
	scoreFlags.Usage = scoreUsage
	scoreScoreListFlags.Usage = scoreScoreListUsage
	scoreScoreDetailFlags.Usage = scoreScoreDetailUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "score":
			svcf = scoreFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "score":
			switch epn {
			case "score-list":
				epf = scoreScoreListFlags

			case "score-detail":
				epf = scoreScoreDetailFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "score":
			c := scorec.NewClient(cc, opts...)
			switch epn {
			case "score-list":
				endpoint = c.ScoreList()
				data, err = scorec.BuildScoreListPayload(*scoreScoreListMessageFlag)
			case "score-detail":
				endpoint = c.ScoreDetail()
				data, err = scorec.BuildScoreDetailPayload(*scoreScoreDetailMessageFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// scoreUsage displays the usage of the score command and its subcommands.
func scoreUsage() {
	fmt.Fprintf(os.Stderr, `成绩系统
Usage:
    %s [globalflags] score COMMAND [flags]

COMMAND:
    score-list: 成绩列表
    score-detail: 成绩详情

Additional help:
    %s score COMMAND --help
`, os.Args[0], os.Args[0])
}
func scoreScoreListUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] score score-list -message JSON

成绩列表
    -message JSON: 

Example:
    `+os.Args[0]+` score score-list --message '{
      "class": "Nisi aliquam est consequatur quod ea vitae.",
      "cursor": 0,
      "limit": 20,
      "name": "Est vel.",
      "scores": 8275923798665478266,
      "sortField": "score",
      "sortOrder": "asc",
      "subject": "Sed aut labore eum placeat quo corporis."
   }'
`, os.Args[0])
}

func scoreScoreDetailUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] score score-detail -message JSON

成绩详情
    -message JSON: 

Example:
    `+os.Args[0]+` score score-detail --message '{
      "id": 1
   }'
`, os.Args[0])
}