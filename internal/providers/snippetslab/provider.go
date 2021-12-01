package snippetslab

import (
	"fmt"
	"strings"

	"github.com/lemoony/snippet-kit/internal/model"
	"github.com/lemoony/snippet-kit/internal/utils"
)

type Provider struct {
	system      *utils.System
	libraryPath snippetsLabLibrary
	tagsFilter  []string
}

// Option configures a Provider.
type Option interface {
	apply(p *Provider)
}

// optionFunc wraps a func so that it satisfies the Option interface.
type optionFunc func(provider *Provider)

func (f optionFunc) apply(provider *Provider) {
	f(provider)
}

// WithSystem sets the utils.System instance to be used by Provider.
func WithSystem(system *utils.System) Option {
	return optionFunc(func(p *Provider) {
		p.system = system
	})
}

func WithTagsFilter(tagFilter []string) Option {
	return optionFunc(func(p *Provider) {
		p.tagsFilter = tagFilter
	})
}

func NewProvider(options ...Option) (*Provider, error) {
	provider := &Provider{
		tagsFilter: []string{},
	}

	for _, o := range options {
		o.apply(provider)
	}

	if err := provider.init(); err != nil {
		return nil, err
	}

	return provider, nil
}

func (p Provider) Info() model.ProviderInfo {
	var lines []model.ProviderLine

	if preferencesURL, err := getPreferencesURL(p.system); err != nil {
		lines = append(lines, model.ProviderLine{IsError: true, Key: "SnippetsLab library Path", Value: err.Error()})
	} else {
		lines = append(lines, model.ProviderLine{
			IsError: true, Key: "SnippetsLab preferences path", Value: preferencesURL.Path,
		})
	}

	if libraryURL, err := getLibraryURL(p.system); err != nil {
		lines = append(lines, model.ProviderLine{IsError: true, Key: "SnippetsLab library path", Value: err.Error()})
	} else {
		lines = append(lines, model.ProviderLine{IsError: true, Key: "SnippetsLab library path", Value: string(libraryURL)})
	}

	if tags, err := p.getValidTagUUIDs(); err != nil {
		lines = append(lines, model.ProviderLine{IsError: true, Key: "SnippetsLab tags", Value: err.Error()})
	} else {
		lines = append(lines, model.ProviderLine{
			IsError: true, Key: "SnippetsLab tags", Value: utils.StringOrDefault(strings.Join(tags.Keys(), ","), "None"),
		})
	}

	if snippets, err := p.GetSnippets(); err != nil {
		lines = append(lines, model.ProviderLine{
			IsError: true, Key: "SnippetsLab total number of snippets", Value: err.Error(),
		})
	} else {
		lines = append(lines, model.ProviderLine{
			IsError: true, Key: "SnippetsLab total number of snippets", Value: fmt.Sprintf("%d", len(snippets)),
		})
	}

	return model.ProviderInfo{
		Lines: lines,
	}
}

func (p *Provider) GetSnippets() ([]model.Snippet, error) {
	validTagUUIDs, err := p.getValidTagUUIDs()
	if err != nil {
		return nil, err
	}

	snippets, err := parseSnippets(p.libraryPath)
	if err != nil {
		return nil, err
	}

	if len(validTagUUIDs) == 0 {
		return snippets, nil
	} else {
		var result []model.Snippet
		for _, snippet := range snippets {
			if hasValidTag(snippet.TagUUIDs, validTagUUIDs) {
				result = append(result, snippet)
			}
		}
		return result, nil
	}
}

func (p *Provider) init() error {
	if libPath, err := getLibraryURL(p.system); err != nil {
		return err
	} else {
		p.libraryPath = libPath
	}

	return nil
}

func hasValidTag(snippetTagUUIDS []string, validTagUUIDs utils.StringSet) bool {
	for _, tagUUID := range snippetTagUUIDS {
		if _, ok := validTagUUIDs[tagUUID]; ok {
			return true
		}
	}
	return false
}

func (p *Provider) getValidTagUUIDs() (utils.StringSet, error) {
	tags, err := parseTags(p.libraryPath)
	if err != nil {
		return nil, err
	}

	result := utils.StringSet{}
	for _, validTag := range p.tagsFilter {
		for tagKey, tagValue := range tags {
			if strings.Compare(tagValue, validTag) == 0 {
				result[tagKey] = struct{}{}
			}
		}
	}

	return result, nil
}
