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
// route: Ark's own GET /api/v3/contents/generations/tasks/:task_id (the
// gateway also serves it under /doubao).
const seedanceTaskStatusRoute = "/api/v3/contents/generations/tasks/:task_id"

// SeedanceTaskStatusPath is the public status path of a Seedance task.
func SeedanceTaskStatusPath(taskID string) string {
	return "/api/v3/contents/generations/tasks/" + taskID
}

// TaskRendererArtifacts projects a finished task's artifacts into the
// gateway addresses a native presenter gets. It is installed by the
// middleware package (it lives in the task plugin adaptor, which imports
// this package).
var TaskRendererArtifacts func(plugin *pluginruntime.LoadedPlugin, task *model.Task) (map[string]any, error)

// TaskPluginPresenterContext is the context a native query presenter runs
// with: the route request context plus the artifacts of one task (empty
// for multi-task queries). The GET query route and task callbacks both use
// it, so their bodies stay identical. A projection error leaves the
// artifacts empty and is returned for logging.
func TaskPluginPresenterContext(plugin *pluginruntime.LoadedPlugin, requestContext map[string]any, task *model.Task) (map[string]any, error) {
	if requestContext == nil {
		requestContext = map[string]any{}
	}
	artifacts := map[string]any{}
	var projectErr error
	if task != nil && TaskRendererArtifacts != nil {
		projected, err := TaskRendererArtifacts(plugin, task)
		if err != nil {
			projectErr = err
		} else if projected != nil {
			artifacts = projected
		}
	}
	requestContext["artifacts"] = artifacts
	return requestContext, projectErr
}

// RenderTaskPluginResult runs a native presenter and returns its value.
func RenderTaskPluginResult(ctx context.Context, plugin *pluginruntime.LoadedPlugin, renderer string, presenterContext map[string]any, input any) (any, error) {
	if plugin == nil || plugin.Engine == nil {
		return nil, errors.New("task plugin is unavailable")
	}
	encoded, err := common.Marshal(input)
	if err != nil {
		return nil, err
	}
	return plugin.Engine.CallPath(ctx, "native", []string{renderer}, presenterContext, pluginruntime.RawJSON(encoded))
}

// RenderTaskPluginBody runs a native presenter and returns the exact body
// bytes the GET query route writes (a reply wrapper contributes its
// __body), so a callback body is byte-identical to what GET returned.
func RenderTaskPluginBody(ctx context.Context, plugin *pluginruntime.LoadedPlugin, renderer string, presenterContext map[string]any, input any) ([]byte, error) {
	result, err := RenderTaskPluginResult(ctx, plugin, renderer, presenterContext, input)
	if err != nil {
		return nil, err
	}
	if wrapper, ok := result.(map[string]any); ok {
		if _, isWrapper := wrapper["__status"]; isWrapper {
			body, hasBody := wrapper["__body"]
			if !hasBody {
				return nil, errors.New("task status presenter returned no body")
			}
			result = body
		}
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
	presenterContext, _ := TaskPluginPresenterContext(binding.Plugin, requestContext.JSValue(), task)
	return RenderTaskPluginBody(ctx, binding.Plugin, binding.Route.Render, presenterContext, view)
}
