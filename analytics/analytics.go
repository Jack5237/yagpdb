package analytics

import (
	"sync"

	"github.com/botlabs-gg/yagpdb/v2/common"
	"github.com/botlabs-gg/yagpdb/v2/common/config"
	"github.com/mediocregopher/radix/v3"
)

// Logger for the analytics plugin
var logger = common.GetPluginLogger(&Plugin{})

// Plugin represents the analytics plugin structure
// stopWorkers is used to signal worker shutdown
 type Plugin struct {
	stopWorkers chan *sync.WaitGroup
}

// PluginInfo returns metadata about the plugin
func (p *Plugin) PluginInfo() *common.PluginInfo {
	return &common.PluginInfo{
		Name:     "Analytics",
		SysName:  "analytics",
		Category: common.PluginCategoryCore,
	}
}

// RegisterPlugin initializes and registers the analytics plugin
func RegisterPlugin() {
	common.RegisterPlugin(&Plugin{
		stopWorkers: make(chan *sync.WaitGroup),
	})
	common.InitSchemas("analytics", dbSchemas...)
}

// RecordActiveUnit logs activity for a specific plugin feature in a guild
func RecordActiveUnit(guildID int64, plugin common.Plugin, analyticName string) {
	if err := recordActiveUnit(guildID, plugin, analyticName); err != nil {
		logger.WithError(err).
			WithField("guild", guildID).
			WithField("plugin", plugin.PluginInfo().SysName).
			WithField("analytic", analyticName).
			Error("Failed updating analytic in Redis")
	}
}

// Config option to enable or disable analytics tracking
var confEnableAnalytics = config.RegisterOption("yagpdb.enable_analytics", "Enable usage analytics tracking", false)

// recordActiveUnit increments analytics tracking data in Redis for a specific guild and plugin
func recordActiveUnit(guildID int64, plugin common.Plugin, analyticName string) error {
	if !confEnableAnalytics.GetBool() {
		return nil // Analytics tracking is disabled
	}

	redisKey := "analytics_active_units." + plugin.PluginInfo().SysName + "." + analyticName
	if err := common.RedisPool.Do(radix.FlatCmd(nil, "HINCRBY", redisKey, guildID, 1)); err != nil {
		return err // Return error if Redis command fails
	}

	return nil // Success
}
