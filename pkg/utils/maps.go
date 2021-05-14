package utils

type MapStringInterface map[string]interface{}

func (msi MapStringInterface) ToMapString() map[string]string {
	rs := make(map[string]string, 0)
	for k, v := range msi {
		if vc, ok := v.(string); ok {
			rs[k] = vc
		}
	}
	return rs
}
