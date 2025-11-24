package ioc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

// defaultStore 默认存储实现
type defaultStore struct {
	store []*NamespaceStore
}

func (s *defaultStore) Len() int {
	return len(s.store)
}

func (s *defaultStore) Less(i, j int) bool {
	return s.store[i].Priority > s.store[j].Priority
}

func (s *defaultStore) Swap(i, j int) {
	s.store[i], s.store[j] = s.store[j], s.store[i]
}

// Sort 根据空间优先级进行排序
func (s *defaultStore) Sort() {
	sort.Sort(s)
}

// Namespace 获取一个对象存储空间
func (s *defaultStore) Namespace(namespace string) *NamespaceStore {
	for i := range s.store {
		item := s.store[i]
		if item.Namespace == namespace {
			return item
		}
	}

	ns := newNamespaceStore(namespace)
	s.store = append(s.store, ns)
	return ns
}

// LoadFromEnv 加载环境变量配置
func (s *defaultStore) LoadFromEnv(prefix string) error {
	errs := []string{}
	for i := range s.store {
		item := s.store[i]
		err := item.LoadFromEnv(prefix)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return nil
}

// LoadFromFile 从配置文件加载配置
func (s *defaultStore) LoadFromFile(path string) error {
	// 尝试多个可能的路径
	possiblePaths := []string{
		path,
		"config/" + filepath.Base(path),
		"ragljx_go/config/" + filepath.Base(path),
	}

	actualPath := ""
	for _, p := range possiblePaths {
		if isFileExists(p) {
			actualPath = p
			break
		}
	}

	if actualPath == "" {
		return fmt.Errorf("file %s not exist (tried: %v)", path, possiblePaths)
	}

	path = actualPath

	fileType := filepath.Ext(path)
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	errs := []string{}
	for i := range s.store {
		item := s.store[i]
		err := item.LoadFromFileContent(content, fileType)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return nil
}

// InitIocObject 初始化托管的所有对象
func (s *defaultStore) InitIocObject() error {
	s.Sort()
	for i := range s.store {
		item := s.store[i]
		err := item.Init()
		if err != nil {
			return fmt.Errorf("[%s] %s", item.Namespace, err)
		}
	}
	return nil
}

// Stop 倒序遍历s.store，关闭对象
func (s *defaultStore) Stop(ctx context.Context) {
	for i := len(s.store) - 1; i >= 0; i-- {
		item := s.store[i]
		item.Close(ctx)
	}
}

// NamespaceStore 命名空间存储
type NamespaceStore struct {
	Namespace string
	Priority  int
	Items     []*ObjectWrapper
}

func newNamespaceStore(namespace string) *NamespaceStore {
	return &NamespaceStore{
		Namespace: namespace,
		Items:     []*ObjectWrapper{},
	}
}

func (s *NamespaceStore) SetPriority(v int) *NamespaceStore {
	s.Priority = v
	return s
}

func (s *NamespaceStore) Registry(v Object) {
	obj := NewObjectWrapper(v)
	old, index := s.getWithIndex(obj.Name)
	if old == nil {
		s.Items = append(s.Items, obj)
		return
	}
	// 已存在则替换
	s.setWithIndex(index, obj)
}

func (s *NamespaceStore) Get(name string) Object {
	obj, _ := s.getWithIndex(name)
	if obj == nil {
		return nil
	}
	return obj.Value
}

func (s *NamespaceStore) setWithIndex(index int, obj *ObjectWrapper) {
	s.Items[index] = obj
}

func (s *NamespaceStore) getWithIndex(name string) (*ObjectWrapper, int) {
	for i := range s.Items {
		obj := s.Items[i]
		if obj.Name == name {
			return obj, i
		}
	}
	return nil, -1
}

func (s *NamespaceStore) List() (names []string) {
	for i := range s.Items {
		item := s.Items[i]
		names = append(names, item.Name)
	}
	return
}

func (s *NamespaceStore) Len() int {
	return len(s.Items)
}

func (s *NamespaceStore) Less(i, j int) bool {
	return s.Items[i].Priority > s.Items[j].Priority
}

func (s *NamespaceStore) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

// Sort 根据对象的优先级进行排序
func (s *NamespaceStore) Sort() {
	sort.Sort(s)
}

func (s *NamespaceStore) Init() error {
	s.Sort()
	for i := range s.Items {
		obj := s.Items[i]
		err := obj.Value.Init()
		if err != nil {
			return fmt.Errorf("init object %s error, %s", obj.Name, err)
		}
		log.Printf("init app %s[priority: %d] ok.", obj.Value.Name(), obj.Value.Priority())
	}
	return nil
}

// Close 倒序关闭
func (s *NamespaceStore) Close(ctx context.Context) {
	for i := len(s.Items) - 1; i >= 0; i-- {
		obj := s.Items[i]
		obj.Value.Close(ctx)
	}
}

// LoadFromEnv 从环境变量中加载对象配置
func (s *NamespaceStore) LoadFromEnv(prefix string) error {
	errs := []string{}
	for _, w := range s.Items {
		prefixList := strings.ToUpper(w.Name) + "_"
		if prefix != "" {
			prefixList = fmt.Sprintf("%s_%s", strings.ToUpper(prefix), prefixList)
		}
		err := env.Parse(w.Value, env.Options{
			Prefix: prefixList,
		})
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return nil
}

// LoadFromFileContent 从文件内容加载配置
func (s *NamespaceStore) LoadFromFileContent(fileContent []byte, fileType string) error {
	fileData := make(map[string]interface{})

	switch fileType {
	case ".yml", ".yaml":
		if err := yaml.Unmarshal(fileContent, &fileData); err != nil {
			return fmt.Errorf("yaml decode error: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(fileContent, &fileData); err != nil {
			return fmt.Errorf("json decode error: %w", err)
		}
	default:
		return fmt.Errorf("unsupported format: %s", fileType)
	}

	// 简单的配置加载，将配置数据映射到对象
	for _, w := range s.Items {
		if configData, exists := fileData[w.Value.Name()]; exists {
			// 这里简化处理，实际应该使用mapstructure等库
			data, _ := json.Marshal(configData)
			json.Unmarshal(data, w.Value)
		}
	}
	return nil
}

func isFileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || !os.IsNotExist(err)
}

