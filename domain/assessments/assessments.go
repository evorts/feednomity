package assessments

import (
	"context"
	"github.com/evorts/feednomity/pkg/database"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type IManager interface {
	FindTemplateDataByKey(ctx context.Context, key string) (*Template, error)
}

type manager struct {
	dbm database.IManager
	cache map[string]Template
}

const (
	yamlFilename = "./forms/360.yaml" // temporary solution to accelerate delivery
)

func NewAssessmentDomain(dbm database.IManager) IManager {
	return &manager{dbm: dbm}
}

func (m *manager) FindTemplateDataByKey(ctx context.Context, key string) (*Template, error) {
	items := m.getDataFromCache()
	if v, ok := items[key]; ok {
		return &v, nil
	}
	return nil, errors.New("not found")
}

func (m *manager) getDataFromCache() map[string]Template {
	if m.cache != nil {
		return m.cache
	}
	m.cache, _ = m.readYaml()
	return m.cache
}

func (m *manager) readYaml() (rs map[string]Template, err error) {
	info, err2 := os.Stat(yamlFilename)
	if os.IsNotExist(err2) {
		return nil, err2
	}
	if info.IsDir() {
		return nil, errors.New("no such file")
	}
	data, err3 := ioutil.ReadFile(yamlFilename)
	if err3 != nil {
		return nil, err3
	}
	err = yaml.Unmarshal(data, &rs)
	return
}
