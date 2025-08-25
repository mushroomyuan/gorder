package decorator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mushroomyuan/gorder/common/logging"
	"github.com/sirupsen/logrus"
)

type QueryLoggingDecorator[C, R any] struct {
	base QueryHandler[C, R]
}

func (q QueryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	body, _ := json.Marshal(cmd)
	fields := logrus.Fields{
		"query":      generateActionName(cmd),
		"query_body": string(body),
	}

	defer func() {
		if err == nil {
			logging.Infof(ctx, fields, "%s", "Query execute successfully!")
		} else {
			logging.Errorf(ctx, fields, "Fail to execute query: %v", err)
		}
	}()
	result, err = q.base.Handle(ctx, cmd)
	return
}

type CommandLoggingDecorator[C, R any] struct {
	base CommandHandler[C, R]
}

func (q CommandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	body, _ := json.Marshal(cmd)
	fields := logrus.Fields{
		"command":      generateActionName(cmd),
		"command_body": string(body),
	}

	defer func() {
		if err == nil {
			logging.Infof(ctx, fields, "%s", "Command execute successfully!")
		} else {
			logging.Errorf(ctx, fields, "Fail to execute command: %v", err)
		}
	}()
	result, err = q.base.Handle(ctx, cmd)
	return
}

func generateActionName(cmd any) string {
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
