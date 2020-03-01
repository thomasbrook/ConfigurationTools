package client

import (
	"ConfigurationTools/dataSource"
	"log"

	"github.com/lxn/walk"
)

type Organization struct {
	orgId    string
	orgName  string
	parent   *Organization
	children []*Organization
}

func NewOrganization(id string, orgName string, parent *Organization) *Organization {
	return &Organization{orgId: id, orgName: orgName, parent: parent}
}

var _ walk.TreeItem = new(Organization)

func (org *Organization) Text() string {
	return org.orgName
}

func (org *Organization) Parent() walk.TreeItem {
	if org.parent == nil {
		return nil
	}
	return org.parent
}

func (org *Organization) ChildCount() int {
	if org.children == nil {
		if err := org.ResetChildren(); err != nil {
			log.Print(err)
		}
	}
	return len(org.children)
}

func (org *Organization) ChildAt(index int) walk.TreeItem {
	return org.children[index]
}

func (org *Organization) ResetChildren() error {
	org.children = nil

	orgs, err := dataSource.ListSubOrg(org.orgId)
	if err != nil {
		return err
	}

	for i := 0; i < len(orgs); i++ {
		if orgs[i] == nil {
			continue
		}

		org.children = append(org.children, NewOrganization(orgs[i].Id, orgs[i].OrgName, org))
	}

	return nil
}

type OrganizationTreeModel struct {
	walk.TreeModelBase
	roots []*Organization
}

var _ walk.TreeModel = new(OrganizationTreeModel)

func NewOrganizationTreeModel() (*OrganizationTreeModel, error) {
	model := new(OrganizationTreeModel)

	orgs, err := dataSource.ListSubOrg("")
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orgs); i++ {
		if orgs[i] == nil {
			continue
		}
		model.roots = append(model.roots, NewOrganization(orgs[i].Id, orgs[i].OrgName, nil))
	}

	return model, nil
}

func (*OrganizationTreeModel) LazyPopulation() bool {
	return true
}

func (m *OrganizationTreeModel) RootCount() int {
	return len(m.roots)
}

func (m *OrganizationTreeModel) RootAt(index int) walk.TreeItem {
	return m.roots[index]
}
