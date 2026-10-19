package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	pluginruntime "github.com/QuantumNous/new-api/pkg/jsplugin"
)

// seedanceTaskStatusRoute is the Doubao plugin's declared native status
// route behind the vendor-compatible GET /api/v3/contents/generations/tasks/:task_id.
const seedanceTaskStatusRoute = "/doubao/api/v3/contents/generations/tasks/:task_id"

// SeedanceTaskStatusPath is the public status path of a Seedance task.
func SeedanceTaskStatusPath(taskID string) string {
	return "/api/v3/contents/generations/tasks/" + taskID
}

// RenderTaskPluginBody runs a native plugin renderer over task views and
// returns the exact response bytes. GET status routes write these bytes
// as-is, and task callbacks snapshot the same bytes, so a callback body is
// byte-identical to what GET returned at that moment.
func RenderTaskPluginBody(ctx context.Context, plugin *pluginruntime.LoadedPlugin, renderer string, requestContext map[string]any, input any) ([]byte, error) {
	if plugin == nil || plugin.Engine == nil {
		return nil, errors.New("task plugin is unavailable")
	}
	encoded, err := common.Marshal(input)
	if err != nil {
		return nil, err
	}
	result, err := plugin.Engine.CallPath(ctx, "native", []string{renderer}, requestContext, pluginruntime.RawJSON(encoded))
	if err != nil {
		return nil, err
	}
	return common.Marshal(result)
}

// RenderSeedanceTaskStatus renders a Seedance video task exactly as the
// vendor-compatible GET status route does for its owner.
func RenderSeedanceTaskStatus(ctx context.Context, task *model.Task) ([]byte, error) {
	if task == nil {
		return nil, errors.New("task is nil")
	}
	binding, ok := pluginruntime.DefaultRegistry.Generation().LookupDeclaredRoute(http.MethodGet, seedanceTaskStatusRoute)
	if !ok || binding.Plugin == nil || binding.Route.Render == "" {
		return nil, errors.New("seedance task status route is unavailable")
	}
	view, err := BuildTaskPluginView(task)
	if err != nil {
		return nil, err
	}
	requestContext := pluginruntime.RouteRequestContext{
		Path:   SeedanceTaskStatusPath(task.TaskID),
		Method: http.MethodGet,
		Params: map[string]string{"task_id": task.TaskID},
	}
	return RenderTaskPluginBody(ctx, binding.Plugin, binding.Route.Render, requestContext.JSValue(), view)
}
