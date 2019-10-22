package api

import (
	"errors"
	"fmt"
	"time"
)

func (api *API) waitUntilPluginUninstalled(instance_id int, name string) (map[string]interface{}, error) {
	time.Sleep(10 * time.Second)
	for {
		response, err := api.ReadPlugin(instance_id, name)

		if err != nil {
			return nil, err
		}
		if len(response) == 0 {
			return response, nil
		}

		time.Sleep(10 * time.Second)
	}
}

func (api *API) EnablePluginCommunity(instance_id int, name string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	params := &PluginParams{Name: name}
	path := fmt.Sprintf("/api/instances/%d/plugins/community", instance_id)
	response, err := api.sling.Post(path).BodyForm(params).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, errors.New(fmt.Sprintf("EnablePluginCommunity failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return api.waitUntilPluginChanged(instance_id, name, true)
}

func (api *API) ReadPluginCommunity(instance_id int, plugin_name string) (map[string]interface{}, error) {
	var data []map[string]interface{}
	data, err := api.ReadPluginsCommunity(instance_id)
	if err != nil {
		return nil, err
	}

	for _, plugin := range data {
		if plugin["name"] == plugin_name {
			return plugin, nil
		}
	}

	return nil, nil
}

func (api *API) ReadPluginsCommunity(instance_id int) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/plugins/community", instance_id)
	response, err := api.sling.Get(path).Receive(&data, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("ReadPluginsCommunity failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return data, nil
}

func (api *API) UpdatePluginCommunity(instance_id int, params map[string]interface{}) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	pluginParams := &PluginParams{Name: params["name"].(string), Enabled: params["enabled"].(bool)}
	path := fmt.Sprintf("/api/instances/%d/plugins/community", instance_id)
	response, err := api.sling.Put(path).BodyForm(pluginParams).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, errors.New(fmt.Sprintf("UpdatePluginCommunity failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return api.waitUntilPluginChanged(instance_id, params["name"].(string), params["enabled"].(bool))
}

func (api *API) DisablePluginCommunity(instance_id int, name string) (map[string]interface{}, error) {
	failed := make(map[string]interface{})
	path := fmt.Sprintf("/api/instances/%d/plugins/community/%s", instance_id, name)
	response, err := api.sling.Delete(path).Receive(nil, &failed)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 204 {
		return nil, errors.New(fmt.Sprintf("DisablePluginCommunity failed, status: %v, message: %s", response.StatusCode, failed))
	}

	return api.waitUntilPluginUninstalled(instance_id, name)
}
