package uptimerobot

import "context"

type UptimeRobotClient interface {
	MakeApiRequest(ctx context.Context, methodName string, params map[string]string) ([]byte, error)
}
