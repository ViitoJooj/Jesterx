package supabase

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ViitoJooj/Jesterx/pkg/dotenv"
)

func TestConn(t *testing.T) {
	originalOpen := openDB
	originalPing := pingDB
	originalURI := dotenv.Supabase_uri
	originalDB := DB

	t.Cleanup(func() {
		openDB = originalOpen
		pingDB = originalPing
		dotenv.Supabase_uri = originalURI
		DB = originalDB
	})

	tests := []struct {
		name      string
		openErr   error
		pingErr   error
		wantPanic bool
		wantDBSet bool
	}{
		{
			name:      "connects successfully",
			wantPanic: false,
			wantDBSet: true,
		},
		{
			name:      "panics when opening database fails",
			openErr:   errors.New("open failed"),
			wantPanic: true,
			wantDBSet: false,
		},
		{
			name:      "panics when ping fails",
			pingErr:   errors.New("ping failed"),
			wantPanic: true,
			wantDBSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			openDB = func(driverName, dataSourceName string) (*sql.DB, error) {
				if tt.openErr != nil {
					return nil, tt.openErr
				}

				return &sql.DB{}, nil
			}

			pingDB = func(db *sql.DB) error {
				return tt.pingErr
			}

			didPanic := false
			func() {
				defer func() {
					if recover() != nil {
						didPanic = true
					}
				}()

				Conn()
			}()

			if didPanic != tt.wantPanic {
				t.Fatalf("panic = %v, want %v", didPanic, tt.wantPanic)
			}

			if (DB != nil) != tt.wantDBSet {
				t.Fatalf("DB set = %v, want %v", DB != nil, tt.wantDBSet)
			}
		})
	}
}
