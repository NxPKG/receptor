package controlsvc

import (
	"context"
	"fmt"

	"github.com/distronode/receptor/internal/version"
	"github.com/distronode/receptor/pkg/utils"
)

type (
	StatusCommandType struct{}
	StatusCommand     struct {
		requestedFields []string
	}
)

func (t *StatusCommandType) InitFromString(params string) (ControlCommand, error) {
	if params != "" {
		return nil, fmt.Errorf("status command does not take parameters")
	}
	c := &StatusCommand{}

	return c, nil
}

func (t *StatusCommandType) InitFromJSON(config map[string]interface{}) (ControlCommand, error) {
	requestedFields, ok := config["requested_fields"]
	var requestedFieldsStr []string
	if ok {
		requestedFieldsStr = make([]string, 0)
		for _, v := range requestedFields.([]interface{}) {
			vStr, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("each element of requested_fields must be a string")
			}
			requestedFieldsStr = append(requestedFieldsStr, vStr)
		}
	} else {
		requestedFieldsStr = nil
	}
	c := &StatusCommand{
		requestedFields: requestedFieldsStr,
	}

	return c, nil
}

func (c *StatusCommand) ControlFunc(_ context.Context, nc NetceptorForControlCommand, _ ControlFuncOperations) (map[string]interface{}, error) {
	status := nc.Status()
	statusGetters := make(map[string]func() interface{})
	statusGetters["Version"] = func() interface{} { return version.Version }
	statusGetters["SystemCPUCount"] = func() interface{} { return utils.GetSysCPUCount() }
	statusGetters["SystemMemoryMiB"] = func() interface{} { return utils.GetSysMemoryMiB() }
	statusGetters["NodeID"] = func() interface{} { return status.NodeID }
	statusGetters["Connections"] = func() interface{} { return status.Connections }
	statusGetters["RoutingTable"] = func() interface{} { return status.RoutingTable }
	statusGetters["Advertisements"] = func() interface{} { return status.Advertisements }
	statusGetters["KnownConnectionCosts"] = func() interface{} { return status.KnownConnectionCosts }
	cfr := make(map[string]interface{})
	if c.requestedFields == nil { // if nil, fill it with the keys in statusGetters
		for field := range statusGetters {
			c.requestedFields = append(c.requestedFields, field)
		}
	}
	for _, field := range c.requestedFields {
		getter, ok := statusGetters[field]
		if ok {
			cfr[field] = getter()
		}
	}

	return cfr, nil
}
