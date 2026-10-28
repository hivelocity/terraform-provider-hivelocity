package hivelocity

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	swagger "github.com/hivelocity/terraform-provider-hivelocity/hivelocity-client-go"
)

// apiErrorDiag wraps an error from the generated client (or anything else) into
// diagnostics, including the API's response body when available.
func apiErrorDiag(err error, format string, a ...interface{}) diag.Diagnostics {
	return diag.FromErr(formatSwaggerError(err, format, a...))
}

// isNotFound reports whether the HTTP response is a 404.
func isNotFound(resp *http.Response) bool {
	return resp != nil && resp.StatusCode == http.StatusNotFound
}

// networkTaskFailedError is returned when an asynchronous network task ends in
// the Failed state.
type networkTaskFailedError struct {
	task *swagger.NetworkTaskDump
}

func (e networkTaskFailedError) Error() string {
	return fmt.Sprintf("network task %s failed (device %d): %s", e.task.TaskId, e.task.DeviceId, e.task.MetaData)
}

// waitForNetworkTask polls /network/status/{taskId} until the task succeeds,
// fails, or the timeout elapses. Many networking endpoints (VLAN, ports, IP
// routing, VPS operations) return a task instead of a final object.
func (hv *Client) waitForNetworkTask(ctx context.Context, taskId string, timeout time.Duration) (*swagger.NetworkTaskDump, error) {
	if taskId == "" {
		return nil, fmt.Errorf("API returned an empty network task ID")
	}

	stateChangeConf := &retry.StateChangeConf{
		Pending: []string{"Pending", ""},
		Target:  []string{"Success"},
		Refresh: func() (interface{}, string, error) {
			task, _, err := hv.client.NetworkApi.GetNetworkTaskIdResource(hv.auth, taskId, nil)
			if err != nil {
				return nil, "", formatSwaggerError(err, "GET /network/status/%s", taskId)
			}
			if task.Result == "Failed" {
				return &task, task.Result, networkTaskFailedError{task: &task}
			}
			return &task, task.Result, nil
		},
		Timeout:                   timeout,
		Delay:                     5 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 1,
	}

	r, err := stateChangeConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}
	return r.(*swagger.NetworkTaskDump), nil
}

// Int32ListFromSet converts a schema.Set of ints into []int32.
func Int32ListFromSet(s *schema.Set) []int32 {
	int32List := make([]int32, s.Len())
	for i, n := range s.List() {
		int32List[i] = int32(n.(int))
	}
	return int32List
}

// SetFromInt32List converts []int32 into a schema.Set of ints.
func SetFromInt32List(items []int32) *schema.Set {
	intList := make([]interface{}, len(items))
	for i, n := range items {
		intList[i] = int(n)
	}
	return schema.NewSet(schema.HashInt, intList)
}

// expandInt32List converts a schema list or set attribute value into []int32.
// The result is never nil so callers can safely send it as a JSON array.
func expandInt32List(v interface{}) []int32 {
	out := make([]int32, 0)
	switch items := v.(type) {
	case *schema.Set:
		for _, n := range items.List() {
			out = append(out, int32(n.(int)))
		}
	case []interface{}:
		for _, n := range items {
			out = append(out, int32(n.(int)))
		}
	}
	return out
}

// expandStringList converts a schema list or set attribute value into []string.
// The result is never nil so callers can safely send it as a JSON array.
func expandStringList(v interface{}) []string {
	out := make([]string, 0)
	switch items := v.(type) {
	case *schema.Set:
		for _, s := range items.List() {
			out = append(out, s.(string))
		}
	case []interface{}:
		for _, s := range items {
			out = append(out, s.(string))
		}
	}
	return out
}

// flattenInt32List converts []int32 into []int for d.Set.
func flattenInt32List(items []int32) []int {
	out := make([]int, len(items))
	for i, n := range items {
		out[i] = int(n)
	}
	return out
}

// flattenStringList makes sure a nil slice becomes an empty one for d.Set.
func flattenStringList(items []string) []string {
	if items == nil {
		return []string{}
	}
	return items
}

// setResourceData sets every key in values, returning the first error.
func setResourceData(d *schema.ResourceData, values map[string]interface{}) error {
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := d.Set(k, values[k]); err != nil {
			return fmt.Errorf("setting %q: %w", k, err)
		}
	}
	return nil
}

// expandStringMap converts a schema map attribute into map[string]string.
func expandStringMap(v interface{}) map[string]string {
	out := map[string]string{}
	if m, ok := v.(map[string]interface{}); ok {
		for k, val := range m {
			out[k] = fmt.Sprintf("%v", val)
		}
	}
	return out
}
