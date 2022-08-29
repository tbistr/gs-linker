package link

import (
	"testing"
)

func TestDbConfig_dsn(t *testing.T) {
	tests := map[string]struct {
		UserName string
		Pass     string
		Addr     string
		DbName   string
		want     string
	}{
		"root": {"root", "passwd", "db", "database", "root:passwd@tcp(db:3306)/database?charset=utf8&parseTime=True&loc=Local"},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			conf := &DbConfig{
				UserName: tt.UserName,
				Pass:     tt.Pass,
				Addr:     tt.Addr,
				DbName:   tt.DbName,
			}
			if got := conf.dsn(); got != tt.want {
				t.Errorf("DbConfig.dsn() = %v, want %v", got, tt.want)
			}
		})
	}
}
