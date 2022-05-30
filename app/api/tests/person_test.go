package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ribgsilva/person-api/app/api/handlers"
	"github.com/ribgsilva/person-api/business/v1/person"
	"github.com/ribgsilva/person-api/platform/env"
	"github.com/ribgsilva/person-api/platform/logger"
	"github.com/ribgsilva/person-api/sys"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type PersonTests struct {
	app http.Handler
}

func TestPerson(t *testing.T) {
	log, err := logger.New("Person-API-Tests")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// =======================================================================================================
	// Setup configs
	sys.Configs.Mongo.Database = env.OrDefault(log, "MONGO_DATABASE", "tests-Person")
	sys.Configs.Mongo.ConnectionURL = env.OrDefault(log, "MONGO_CONNECTION_URL", "mongodb://person_app:person_app_pass@localhost:27017")
	sys.Configs.Mongo.ConnectionTimeout = env.DurationDefault(log, "MONGO_CONNECTION_TIMEOUT", "10s")
	sys.Configs.Mongo.DisconnectTimeout = env.DurationDefault(log, "MONGO_DISCONNECT_TIMEOUT", "5s")
	sys.Configs.Mongo.OperationTimeout = env.DurationDefault(log, "MONGO_OPERATION_TIMEOUT", "10s")
	sys.Configs.Mongo.PingTimeout = env.DurationDefault(log, "MONGO_PING_TIMEOUT", "2s")

	// =======================================================================================================
	// Setup resources
	// logger
	sys.S.Log = log

	// mongo

	// doing in a func, so I can use defer to cancel the contexts
	if err := func() error {
		mongoCtx, mongoCancel := context.WithTimeout(context.Background(), sys.Configs.Mongo.ConnectionTimeout)
		defer mongoCancel()
		if client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(sys.Configs.Mongo.ConnectionURL)); err != nil {
			return fmt.Errorf("could not connect to mongo: %w", err)
		} else {
			pingCtx, pingCancel := context.WithTimeout(context.Background(), sys.Configs.Mongo.PingTimeout)
			defer pingCancel()
			if err = client.Ping(pingCtx, readpref.Primary()); err != nil {
				return fmt.Errorf("could not connect to mongo: %w", err)
			}
			sys.S.Mongo = client
		}
		return nil
	}(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		sdCtx, sdCancel := context.WithTimeout(context.Background(), sys.Configs.Mongo.DisconnectTimeout)
		defer sdCancel()
		if err := sys.S.Mongo.Database(sys.Configs.Mongo.Database).Drop(sdCtx); err != nil {
			log.Error(err)
		}
		if err := sys.S.Mongo.Disconnect(sdCtx); err != nil {
			log.Error(err)
		}
	}()

	// =======================================================================================================
	// Setup router
	engine := gin.Default()

	handlers.MapHandlers(engine)

	tests := PersonTests{
		engine,
	}

	// =======================================================================================================
	// Tun tests

	t.Run("crudUsers", tests.crudPerson)
}

func (pt *PersonTests) crudPerson(t *testing.T) {
	id := pt.postPerson201(t)

	pt.getPerson200(t, id)
	pt.searchPerson200(t, id)

	pt.putPerson200(t, id)
	pt.deleteUser200(t, id)
	pt.getPerson404(t, id)
}

func (pt *PersonTests) putPerson200(t *testing.T, id string) {
	role := "Software Engineer Sr"
	nu := person.CreateRequest{
		Name:            "Gabriel More Names",
		Type:            "employee",
		Role:            &role,
		ContactDuration: nil,
		Tags:            []string{"JAVA", "GoLang"},
	}

	body, err := json.Marshal(&nu)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPut, "/v1/persons/"+id, bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	pt.app.ServeHTTP(w, r)

	var resp person.Person
	if w.Code != http.StatusOK {
		t.Fatalf("Test putPerson200: Should receive a status code of 200 for the response : %v", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Test putPerson200: Should be able to unmarshal the response : %v", err)
	}

	if resp.Id != id {
		t.Fatalf("Test putPerson200: Should have received the same id inside the response: %v", resp)
	}

	if resp.Name != "Gabriel More Names" {
		t.Fatalf("Test putPerson200: Should have received the same name inside the response: %v", resp)
	}

	if resp.Role == nil || *resp.Role != "Software Engineer Sr" {
		t.Fatalf("Test putPerson200: Should have received the same role inside the response: %v", resp)
	}
}

func (pt *PersonTests) searchPerson200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/persons?tags=GoLang&tags=JAVA", nil)
	w := httptest.NewRecorder()

	pt.app.ServeHTTP(w, r)

	var resp []person.Person
	if w.Code != http.StatusOK {
		t.Fatalf("Test searchPerson200: Should receive a status code of 200 for the response : %v", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Test searchPerson200: Should be able to unmarshal the response : %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("Test searchPerson200: Should have received one result in the response: %v", resp)
	}
	if resp[0].Id != id {
		t.Fatalf("Test searchPerson200: Should have received the same id in the response: %v", resp)
	}
	if resp[0].Name != "Gabriel" {
		t.Fatalf("Test searchPerson200: Should have received \"Gabriel\" as name in the response: %v", resp)
	}
	if resp[0].Type != "employee" {
		t.Fatalf("Test searchPerson200: Should have received \"employee\" as type in the response: %v", resp)
	}
	if resp[0].Role == nil || *resp[0].Role != "Software Engineer" {
		t.Fatalf("Test searchPerson200: Should have received \"Software Engineer\" as role in the response: %v", resp)
	}
	if resp[0].ContactDuration != nil {
		t.Fatalf("Test searchPerson200: Should have received \"nil\" as contractDuration in the response: %v", resp)
	}
	if len(resp[0].Tags) != 2 {
		t.Fatalf("Test searchPerson200: Should have received \"[\"JAVA\", \"GoLang\"]\" as Tags in the response: %v", resp)
	}
}

func (pt *PersonTests) getPerson200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/persons/"+id, nil)
	w := httptest.NewRecorder()

	pt.app.ServeHTTP(w, r)

	var resp person.Person
	if w.Code != http.StatusOK {
		t.Fatalf("Test getPerson200: Should receive a status code of 200 for the response : %v", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Test getPerson200: Should be able to unmarshal the response : %v", err)
	}

	if resp.Id != id {
		t.Fatalf("Test getPerson200: Should have received the same id in the response: %v", resp)
	}
	if resp.Name != "Gabriel" {
		t.Fatalf("Test getPerson200: Should have received \"Gabriel\" as name in the response: %v", resp)
	}
	if resp.Type != "employee" {
		t.Fatalf("Test getPerson200: Should have received \"employee\" as type in the response: %v", resp)
	}
	if resp.Role == nil || *resp.Role != "Software Engineer" {
		t.Fatalf("Test getPerson200: Should have received \"Software Engineer\" as role in the response: %v", resp)
	}
	if resp.ContactDuration != nil {
		t.Fatalf("Test getPerson200: Should have received \"nil\" as contractDuration in the response: %v", resp)
	}
	if len(resp.Tags) != 2 {
		t.Fatalf("Test getPerson200: Should have received \"[\"JAVA\", \"GoLang\"]\" as Tags in the response: %v", resp)
	}
}

func (pt *PersonTests) getPerson404(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodGet, "/v1/persons/"+id, nil)
	w := httptest.NewRecorder()

	pt.app.ServeHTTP(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Test getPerson404: Should receive a status code of 404 for the response : %v", w.Code)
	}
}

func (pt *PersonTests) postPerson201(t *testing.T) string {
	role := "Software Engineer"
	nu := person.CreateRequest{
		Name:            "Gabriel",
		Type:            "employee",
		Role:            &role,
		ContactDuration: nil,
		Tags:            []string{"JAVA", "GoLang"},
	}

	body, err := json.Marshal(&nu)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/v1/persons", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	pt.app.ServeHTTP(w, r)

	var resp person.CreateResponse
	if w.Code != http.StatusCreated {
		t.Fatalf("Test postPerson201: Should receive a status code of 201 for the response : %v", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Test postPerson201: Should be able to unmarshal the response : %v", err)
	}

	if resp.Id == "" {
		t.Fatalf("Test postPerson201: Should have received an id inside response: %v", resp)
	}

	return resp.Id
}

func (pt *PersonTests) deleteUser200(t *testing.T, id string) {
	r := httptest.NewRequest(http.MethodDelete, "/v1/persons/"+id, nil)
	w := httptest.NewRecorder()

	pt.app.ServeHTTP(w, r)

	var resp person.Person
	if w.Code != http.StatusOK {
		t.Fatalf("Test deleteUser200: Should receive a status code of 200 for the response : %v", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Test deleteUser200: Should be able to unmarshal the response : %v", err)
	}

	if resp.Id != id {
		t.Fatalf("Test deleteUser200: Should have received the same id in the response: %v", resp)
	}
	if resp.Name != "Gabriel More Names" {
		t.Fatalf("Test deleteUser200: Should have received \"Gabriel More Names\" as name in the response: %v", resp)
	}
	if resp.Type != "employee" {
		t.Fatalf("Test deleteUser200: Should have received \"employee\" as type in the response: %v", resp)
	}
	if resp.Role == nil || *resp.Role != "Software Engineer Sr" {
		t.Fatalf("Test deleteUser200: Should have received \"Software Engineer Sr\" as role in the response: %v", resp)
	}
	if resp.ContactDuration != nil {
		t.Fatalf("Test deleteUser200: Should have received \"nil\" as contractDuration in the response: %v", resp)
	}
	if len(resp.Tags) != 2 {
		t.Fatalf("Test deleteUser200: Should have received \"[\"JAVA\", \"GoLang\"]\" as Tags in the response: %v", resp)
	}
}
