package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"workWithDB/DataBase"
	"workWithDB/models"
)

func init(){
	DataBase.TEST = true
}

func TestCreateUser(t *testing.T) {
	//type context struct {
	//	c *gin.Context
	//}
	//var wants []context
	//wants = append(wants)
	tests := []struct {
		name       string
		user       *models.User
		wantStatus int
		//error errorBody
	}{
		{
			name: "error on filling email",
			user: &models.User{
				Name:     "Name",
				Surname:  "Surname",
				Email:    "nam.ru",
				Login:    "Login14",
				Password: "User'sP@ss123",
			},
			wantStatus: 400,
			//error: "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},

		{
			name: "ok",
			user: &models.User{
				Name:     "Name",
				Surname:  "Surname",
				Email:    "name@mail.ru",
				Login:    "Login14",
				Password: "User'sP@ss123",
			},
			wantStatus: 200,
			//error:"",
		},
		{
			name: "this login is already exist",
			user: &models.User{
				Name:     "Name",
				Surname:  "Surname",
				Email:    "name@mail.ru",
				Login:    "Login14",
				Password: "User'sP@ss123",
			},
			wantStatus: 400,
			//error:"this login is already exist",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(t.Name())
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			var req http.Request
			var body []byte
			if tt.user != nil {
				body, _ = json.Marshal(tt.user)
				req.Body = ioutil.NopCloser(strings.NewReader(string(body)))
			}
			c.Request = &req
			CreateUser(c)

			assert.Equal(t, tt.wantStatus, w.Code)
			//assert.Equal(t, tt.error, w.Body)

		})
	}
}
