package data

import (
	"github.com/zbitech/controller/pkg/model"
)

var (
	Project1 = model.Project{Name: "project1", Owner: Owner1_UserId}
	Project2 = model.Project{Name: "project2", Owner: Owner1_UserId}
	Project3 = model.Project{Name: "project3", Owner: Owner2_UserId}
	Project4 = model.Project{Name: "project4", Owner: Owner2_UserId}

	Projects       = []model.Project{Project1, Project2, Project3, Project4}
	Owner1Projects = []model.Project{Project1, Project2}
	Owner2Projects = []model.Project{Project3, Project4}
)

func AppendProjects(projects []model.Project, _projects ...model.Project) []model.Project {
	return append(projects, _projects...)
}

func CreateProjects(count int, props map[string]interface{}) []model.Project {

	var projects = make([]model.Project, count)
	for index := range projects {
		projects[index] = *CreateProject(props)
	}

	return projects
}

func CreateProject(props map[string]interface{}) *model.Project {
	return &model.Project{
		Name:  getProperty(props, "name", randomString(10)).(string),
		Owner: getProperty(props, "owner", randomString(10)).(string),
		//Network: getProperty(props, "network", randomValue(networkTypes)).(model.NetworkType),
		//TeamId:  getProperty(props, "team", randomString(50)).(string),
	}
}
