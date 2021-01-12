package preferences

import (
	"fmt"
	"github.com/iac-io/myiac/internal/util"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var prefsFilename = "/tmp/.myiac/testPrefs"

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {

}

func teardown() {
	log.Printf("Cleaning up prefs file")
	_ = os.Remove(prefsFilename)
}

func TestCreatesPrefsFile(t *testing.T) {
	prefs := NewConfig(prefsFilename)
	assert.NotEmpty(t, prefs)
	assert.FileExists(t, prefsFilename)
}

func TestSetsProperty(t *testing.T) {
	prefs := NewConfig(prefsFilename)
	prefs.Set("testProperty", "testValue")
	prefsFileContent, _ := util.ReadFileToString(prefsFilename)
	fmt.Printf("Contents %v", prefsFileContent)
	assert.Contains(t, prefsFileContent, "testProperty")
	assert.Contains(t, prefsFileContent, "testValue")
}

func TestGetsProperty(t *testing.T) {
	prefs := NewConfig(prefsFilename)
	prefs.Set("testProperty", "testValue")
	assert.Equal(t, "testValue", prefs.Get("testProperty"))
}

func TestChangesExistingProperty(t *testing.T) {
	prefs := NewConfig(prefsFilename)
	prefs.Set("testProperty", "testValue")
	prefs.Set("testProperty", "newValue")
	val := prefs.Get("testProperty")
	assert.Equal(t, "newValue", val)
}

func TestDeletesProperty(t *testing.T) {
	prefs := NewConfig(prefsFilename)
	prefs.Set("testProperty", "testValue")
	prefs.Del("testProperty")
	val := prefs.Get("testProperty")
	assert.Equal(t, "", val)
}

func TestExistsWithInexistentProperty(t *testing.T) {
	prefs := NewConfig(prefsFilename)
	val := prefs.Get("noProperty")
	assert.Equal(t, "", val)
}

func TestDoesNotErrorOnInexistentProperty(t *testing.T) {
	prefs := NewConfig(prefsFilename)
	assertDoesNotPanic(t, func() {
		prefs.Get("")
	})
}

func TestCreatesInexistentFile(t *testing.T) {
	prefs := NewConfig("~/.inexistent/notExists")
	prefs.Set("testProperty", "testValue")
	val := prefs.Get("testProperty")
	assert.Equal(t, "testValue", val)
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func assertDoesNotPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code panicked")
		}
	}()
	f()
}
