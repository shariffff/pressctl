package ansible

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pressctl/cli/pkg/models"
)

func TestInventoryAndExtraVarsArePrivate(t *testing.T) {
	ig := &InventoryGenerator{outputDir: t.TempDir()}
	inv, err := ig.Generate(models.Server{Name: "s1", IP: "192.0.2.10"}, "press test", map[string]interface{}{"mysql_pressctlbot_password": "dummy-secret"})
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Dir(inv)
	for path, want := range map[string]os.FileMode{dir: 0o700, inv: 0o600} {
		st, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if st.Mode().Perm() != want {
			t.Errorf("%s mode = %v, want %v", path, st.Mode().Perm(), want)
		}
	}

	vars, err := writeExtraVars(inv, map[string]interface{}{"wordsail_connection_key": "dummy-key"})
	if err != nil {
		t.Fatal(err)
	}
	if st, _ := os.Stat(vars); st.Mode().Perm() != 0o600 {
		t.Errorf("extra vars mode = %v, want 0600", st.Mode().Perm())
	}
	if b, _ := os.ReadFile(vars); !strings.Contains(string(b), "dummy-key") {
		t.Errorf("extra vars file does not hold the vars: %s", b)
	}

	if err := ig.Cleanup(inv); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("run directory %s still exists after Cleanup", dir)
	}
}
