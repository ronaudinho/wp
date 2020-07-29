// integration test with sqlite db
package sqlite_test

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/ronaudinho/wp/internal/model"
	"github.com/ronaudinho/wp/internal/repository/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

var (
	testdbfile = "./testmessage.db"
	testSQRepo *sqlite.Repository
)

var createTests = []struct {
	name   string
	in     model.Message
	errnil bool
}{
	{
		name:   "ok",
		in:     model.Message{Message: strptr("brainfcuk is fun")},
		errnil: true,
	},
	{
		name:   "ok",
		in:     model.Message{Message: strptr("hi")},
		errnil: true,
	},
}

var getTests = []struct {
	name   string
	in     model.Message
	want   []string
	errnil bool
}{
	{
		name: "ok",
		in:   model.Message{},
		want: []string{
			"brainfcuk is fun",
			"hi",
		},
		errnil: true,
	},
	{
		name: "get after create ok",
		in:   model.Message{Message: strptr("really?")},
		want: []string{
			"brainfcuk is fun",
			"hi",
			"really?",
		},
		errnil: true,
	},
}

func strptr(s string) *string {
	return &s
}

func TestMain(m *testing.M) {
	testdb, err := sql.Open("sqlite3", testdbfile)
	if err != nil {
		log.Fatal(err)
	}
	err = sqlite.InitDB(testdb)
	if err != nil {
		log.Fatal(err)
	}
	testSQRepo = sqlite.New(testdb)
	code := m.Run()
	testdb.Close()
	os.Remove(testdbfile)
	os.Exit(code)
}

func TestRepository_CreateMessage(t *testing.T) {
	for _, tt := range createTests {
		t.Run(tt.name, func(t *testing.T) {
			err := testSQRepo.CreateMessage(tt.in)
			if tt.errnil != (err == nil) {
				t.Errorf("got errnil %v, want errnil %v", err != nil, tt.errnil)
			}
		})
	}
}

func TestRepository_GetMessages(t *testing.T) {
	for _, tt := range getTests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.in, model.Message{}) {
				testSQRepo.CreateMessage(tt.in)
			}
			msg, err := testSQRepo.GetMessages()
			if tt.errnil != (err == nil) {
				t.Errorf("got errnil %v, want errnil %v", err != nil, tt.errnil)
			}
			var s []string
			for _, m := range msg {
				s = append(s, *m.Message)
			}
			if !reflect.DeepEqual(s, tt.want) {
				t.Errorf("got %v, want %v", s, tt.want)
			}
		})
	}
}
