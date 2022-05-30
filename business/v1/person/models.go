package person

import "time"

// CreateRequest model info
// @Description Person create request body
type CreateRequest struct {
	Name            string   `json:"name" minLength:"3" maxLength:"100" example:"William" binding:"required,min=3,max=100"`            // Name of the person
	Type            string   `json:"type" enums:"contractor,employee" example:"employee" binding:"required,oneof=contractor employee"` // Type of the person: contractor - employee
	Role            *string  `json:"role" minLength:"3" maxLength:"50" example:"Software Engineer"`                                    // Role of the employee
	ContactDuration *string  `json:"contractDuration" example:"5 months"`                                                              // ContactDuration time of contract
	Tags            []string `json:"tags" example:"c#,c++"`                                                                            // Tags technologies of the person
}

// CreateResponse model info
// @Description Person create request body
type CreateResponse struct {
	Id string `json:"id" minLength:"36" maxLength:"36" example:"3de26feb-5cd5-4d70-81b9-44bf6f74f453" binding:"required,min=36,max=36"` // Id of the person in the system
}

// UpdateRequest model info
// @Description Person update request body
type UpdateRequest struct {
	Name            string   `json:"name" minLength:"3" maxLength:"100" example:"William" binding:"required,min=3,max=100"`            // Name of the person
	Type            string   `json:"type" enums:"contractor,employee" example:"employee" binding:"required,oneof=contractor employee"` // Type of the person: contractor - employee
	Role            *string  `json:"role" minLength:"3" maxLength:"50" example:"Software Engineer"`                                    // Role of the employee
	ContactDuration *string  `json:"contractDuration" example:"5 months"`                                                              // ContactDuration time of contract
	Tags            []string `json:"tags" example:"c#,c++"`                                                                            // Tags technologies of the person
}

// Person model info
// @Description Person represents data about a person
type Person struct {
	Id              string    `json:"id" minLength:"36" maxLength:"36" example:"3de26feb-5cd5-4d70-81b9-44bf6f74f453" binding:"required,min=36,max=36"` // Id of the person in the system
	Name            string    `json:"name" minLength:"3" maxLength:"100" example:"William" binding:"required,min=3,max=100"`                            // Name of the person
	Type            string    `json:"type" enums:"contractor,employee" example:"employee" binding:"required,oneof=contractor employee"`                 // Type of the person: contractor - employee
	Role            *string   `json:"role" minLength:"3" maxLength:"50" example:"Software Engineer" binding:"min=3,max=50"`                             // Role of the employee
	ContactDuration *string   `json:"contractDuration" example:"5 months" binding:"min=0,max=50"`                                                       // ContactDuration time of contract
	Tags            []string  `json:"tags" example:"c#,c++"`                                                                                            // Tags technologies of the person
	UpdatedAt       time.Time `json:"updatedAt" example:"2006-01-02T15:04:05Z"`                                                                         // Date of the last update
	CreatedAt       time.Time `json:"createdAt" example:"2006-01-02T15:04:05Z"`                                                                         // Date of creation
}

// SearchRequest model info
// @Description Search request
type SearchRequest struct {
	Tags []string `form:"tags" example:"c#,c++" binding:"required"` // Tags technologies of the person
}
