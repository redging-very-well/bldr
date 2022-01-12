package pipeline

import (
	"fmt"

	"github.com/Redgwell/bldr/internal/config"
	"github.com/Redgwell/bldr/internal/models"
	"github.com/Redgwell/bldr/pkg/contexts"
	"github.com/Redgwell/bldr/pkg/plugins"
	"github.com/sirupsen/logrus"
)

type PluginPipeline struct {
	config          *config.Base
	contextProvider *contexts.ContextProvider
	logger          *logrus.Logger
	mode            string
	plugins         []plugins.PluginDefinition
}

func (p *PluginPipeline) addPlugin(plugin plugins.PluginDefinition) {
	p.plugins = append(p.plugins, plugin)
}

func NewPluginPipeline(logger *logrus.Logger, baseConfig *config.Base) *PluginPipeline {
	pipeline := &PluginPipeline{
		config:          baseConfig,
		contextProvider: contexts.NewContextProvider(logger),
		logger:          logger,
		mode:            "build",
	}
	return pipeline
}

func (p *PluginPipeline) InitialisePipeline(registry *plugins.Registry) error {

	pipelineCfg, err := models.LoadPipelineConfig(p.config.Pipeline.Path)
	if err != nil {
		// log/blow up
		return err
	}

	for i, t := range pipelineCfg.Targets {
		p.logger.Infof("Initialising target: %v using %v\n", t.Name, t.Build.Plugin)
		// Load the PluginDefinition using the plugin registry for now
		// Later, we could potentially support go plugins

		pluginInstance, err := registry.CreateInstance(t.Build.Plugin)
		if err != nil {
			return err
		}

		yamlPath := fmt.Sprintf("$.targets[%d].%s", i, p.mode)
		pluginConfig := pipelineCfg.LoadPluginConfig(yamlPath)

		err = pluginInstance.SetConfig(p.logger, t.Name, pluginConfig)
		if err != nil {
			return err
		}
		p.addPlugin(pluginInstance)
	}

	return nil
}

func (p *PluginPipeline) Run() error {

	for _, plugin := range p.plugins {
		err := plugin.Execute(p.contextProvider)

		if err != nil {

			// log/blow up
			p.logger.Errorf("failed running %v: %v", plugin, err)
			return err
		}
	}

	return nil
}
