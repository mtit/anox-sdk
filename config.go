package sdk

import "strconv"

type ConfigAccessor struct {
	client *Client
}

func (c *ConfigAccessor) GetString(key string) string        { return c.client.GetConfig(key) }
func (c *ConfigAccessor) GetServiceString(key string) string { return c.client.GetServiceConfig(key) }
func (c *ConfigAccessor) GetGlobalString(key string) string  { return c.client.GetGlobalConfig(key) }

func (c *ConfigAccessor) GetInt(key string, defaultVal int) int {
	val := c.client.GetConfig(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func (c *ConfigAccessor) GetInt64(key string, defaultVal int64) int64 {
	val := c.client.GetConfig(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return i
}

func (c *ConfigAccessor) GetFloat64(key string, defaultVal float64) float64 {
	val := c.client.GetConfig(key)
	if val == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

func (c *ConfigAccessor) GetBool(key string, defaultVal bool) bool {
	val := c.client.GetConfig(key)
	if val == "" {
		return defaultVal
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultVal
	}
	return b
}

func (c *Client) Config() *ConfigAccessor {
	return &ConfigAccessor{client: c}
}
