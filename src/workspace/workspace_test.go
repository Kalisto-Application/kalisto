package workspace

import (
	"fmt"
	"kalisto/src/models"
	"testing"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WorkspaceSuite struct {
	suite.Suite

	w *Workspace
}

func (s *WorkspaceSuite) SetupTest() {
	store := NewMockStore(s.T())
	store.EXPECT().SaveWorkspaces(mock.Anything).Return(nil)
	store.EXPECT().Workspace().Return(nil, nil)

	w, err := New(store)
	s.Require().NoError(err)
	s.w = w
}

func (s *WorkspaceSuite) TestWorkspace() {
	var err error
	// save
	w1 := models.Workspace{
		Name:     "w1",
		BasePath: "path",
	}
	w1, err = s.w.Save(w1)
	s.Require().NoError(err)

	// save
	w2 := models.Workspace{
		Name:     "w2",
		BasePath: "path",
	}
	w2, err = s.w.Save(w2)
	s.Require().NoError(err)

	// save
	w3 := models.Workspace{
		Name:     "w3",
		BasePath: "path",
	}
	w3, err = s.w.Save(w3)
	s.Require().NoError(err)

	// list
	list := s.w.List()
	s.EqualValues(list[0].ID, w3.ID)
	s.EqualValues(list[1].ID, w2.ID)
	s.EqualValues(list[2].ID, w1.ID)

	// find
	w3Found, err := s.w.Find(w3.ID)
	w3.LastUsage = w3Found.LastUsage
	s.Require().NoError(err)
	s.EqualValues(w3Found, w3)

	// rename
	err = s.w.Rename(w1.ID, "w11")
	s.Require().NoError(err)
	w1Found, err := s.w.Find(w1.ID)
	s.Require().NoError(err)
	s.EqualValues(w1Found.Name, "w11")

	// delete
	err = s.w.Delete(w1.ID)
	s.Require().NoError(err)

	// not found
	_, err = s.w.Find(w1.ID)
	s.ErrorIs(err, models.ErrWorkspaceNotFound)

	// list
	// w3, w2
	list = s.w.List()
	s.EqualValues(list[0].ID, w3.ID)
	s.EqualValues(list[1].ID, w2.ID)
	fmt.Println(list[0].LastUsage.String())
	fmt.Println(list[1].LastUsage.String())
}

func TestWorkspace(t *testing.T) {
	suite.Run(t, new(WorkspaceSuite))
}
