package greyhound

import "log"
import "io/ioutil"
import "regexp"

func NewGreyhoundSearch() *GreyhoundSearch {
	return &GreyhoundSearch{make(map[string]*SearchIndex)}
}

type GreyhoundSearch struct {
	Projects map[string]*SearchIndex
}

func (gs *GreyhoundSearch) AddProject(name, path string, exclusions []string) {
	log.Printf("Searching in path %s...", path)
	regexExclusions := make([]*regexp.Regexp, len(exclusions), len(exclusions))
	for p, v := range exclusions {
		var err error
		regexExclusions[p], err = regexp.Compile(v)
		if err != nil {
			log.Print(err)
		}
	}
	gs.Projects[name] = NewSearchIndex(name, path, regexExclusions)
	log.Print(gs.ListProjects())
}

// lists projects
func (gs *GreyhoundSearch) ListProjects() []string {
	project_names := make([]string, len(gs.Projects), len(gs.Projects))
	i := 0
	for k, _ := range gs.Projects {
		project_names[i] = k
		i++
	}
	return project_names
}

// get file contents
func (gs *GreyhoundSearch) ViewFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Unable to read file! ", err)
	}
	return string(content)
}

// return a search result for a projectName query
func (gs *GreyhoundSearch) SearchFile(projectName, query string) []string {
	_, hasKey := gs.Projects[projectName]
	if hasKey {
		return gs.Projects[projectName].FileResults(query)
	}
	return []string{"no results found"}
}

// return a search result for a projectName query
func (gs *GreyhoundSearch) SearchCode(projectName, query string) []string {
	_, hasKey := gs.Projects[projectName]
	if hasKey {
		return gs.Projects[projectName].CodeResults(query)
	}
	return []string{"no results found"}
}
