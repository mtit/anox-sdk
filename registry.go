package sdk

type Registry struct {
	client *Client
}

func (r *Registry) Register() error   { return nil }
func (r *Registry) Deregister() error { return nil }
func (r *Registry) GetInstanceID() string {
	if r.client == nil {
		return ""
	}
	return r.client.GetInstanceID()
}
