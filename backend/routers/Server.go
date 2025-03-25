package routers

import (
	"awesomeProject/controllers"
	"net/http"

	api "awesomeProject/openAPIGenerated"
	"github.com/gin-gonic/gin"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Server implements the ServerInterface
type Server struct {
	userController *controllers.UserController
}

// NewServer creates a new server instance with the user controller
func NewServer(
	userController *controllers.UserController,
) *Server {
	return &Server{
		userController: userController,
	}
}

// PostApiAgendaInvites creates a new agenda invite
func (s *Server) PostApiAgendaInvites(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// DeleteApiAgendaInvitesId deletes an agenda invite by ID
func (s *Server) DeleteApiAgendaInvitesId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetApiAgendaInvitesId gets an agenda invite by ID
func (s *Server) GetApiAgendaInvitesId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// PutApiAgendaInvitesId updates an agenda invite by ID
func (s *Server) PutApiAgendaInvitesId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetApiAgendaItems queries agenda items
func (s *Server) GetApiAgendaItems(c *gin.Context, params api.GetApiAgendaItemsParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// PostApiAgendaItems creates or updates multiple agenda items
func (s *Server) PostApiAgendaItems(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// DeleteApiAgendaItemsId deletes one or multiple agenda items
func (s *Server) DeleteApiAgendaItemsId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetApiAgendaItemsId gets an agenda item by ID
func (s *Server) GetApiAgendaItemsId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetApiAgendaSources gets a list of agenda sources
func (s *Server) GetApiAgendaSources(c *gin.Context, params api.GetApiAgendaSourcesParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// PostApiAgendaSources creates a new agenda source
func (s *Server) PostApiAgendaSources(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// DeleteApiAgendaSourcesId deletes an agenda source by ID
func (s *Server) DeleteApiAgendaSourcesId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// GetApiAgendaSourcesId gets an agenda source by ID
func (s *Server) GetApiAgendaSourcesId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// PutApiAgendaSourcesId updates an agenda source by ID
func (s *Server) PutApiAgendaSourcesId(c *gin.Context, id openapi_types.UUID) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}

// PostApiRegister registers a new user
func (s *Server) PostApiRegister(c *gin.Context) {
	s.userController.CreateUser(c)
}

// PutApiUsersId updates a user's account details
func (s *Server) PutApiUsersId(c *gin.Context, id openapi_types.UUID) {
	s.userController.UpdateUser(c, id)
}

// GetApiViewAgendaInviteId provides a publicly available view of a user agenda
func (s *Server) GetApiViewAgendaInviteId(c *gin.Context, id openapi_types.UUID, params api.GetApiViewAgendaInviteIdParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented"})
}
