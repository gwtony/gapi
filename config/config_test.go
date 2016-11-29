package config

import (
	"os"
	"io/ioutil"
	"path/filepath"
	"testing"
	"github.com/gwtony/gapi/test"
)


func testReadConf(t *testing.T, data string) *Config {
	conf := &Config{}
	tempDir, err := ioutil.TempDir("", "test_log")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	path := filepath.Join(tempDir, "test_conf")
	err = ioutil.WriteFile(path, []byte(test.TestNapiConf), 0644)
	if err != nil {
		t.Fatalf("writeFile: %v", err)
	}

	err = conf.ReadConf(path)
	if err != nil {
		t.Fatal("Test read conf failed")
	}
	t.Log("Test read conf ok")
	return conf
}

func TestReadConfOk(t *testing.T) {
	testReadConf(t, test.TestNapiConf)
}

func TestParseConfOk(t *testing.T) {
	c := testReadConf(t, test.TestNapiConf)
	c.ParseConf()
}
