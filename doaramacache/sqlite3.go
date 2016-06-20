package doaramacache

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"io"
	"io/ioutil"

	"github.com/cnf/structhash"
	"github.com/twpayne/go-doarama"
	"golang.org/x/net/context"
)

type sqlite struct {
	client     *doarama.Client
	db         *sql.DB
	insertStmt *sql.Stmt
	queryStmt  *sql.Stmt
}

// NewSQLite3 returns a new ActivityCreator that caches activities from client
// in dataSourceName.
func NewSQLite3(dataSourceName string, client *doarama.Client) (ActivityCreator, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec("" +
		"CREATE TABLE IF NOT EXISTS activities (\n" +
		"  activity_id INT UNIQUE,\n" +
		"  gpstrack_sha256 STRING(32) UNIQUE,\n" +
		"  activityinfo_sha256 STRING(32) UNIQUE\n" +
		");"); err != nil {
		return nil, err
	}
	insertStmt, err := db.Prepare("" +
		"INSERT INTO activities(activity_id, gpstrack_sha256, activityinfo_sha256)\n" +
		"VALUES (?, ?, ?);")
	if err != nil {
		return nil, err
	}
	queryStmt, err := db.Prepare("" +
		"SELECT activity_id\n" +
		"FROM activities\n" +
		"WHERE gpstrack_sha256 = ? AND activityinfo_sha256 = ?;")
	if err != nil {
		return nil, err
	}
	return &sqlite{
		client:     client,
		db:         db,
		insertStmt: insertStmt,
		queryStmt:  queryStmt,
	}, nil
}

// Close releases any resources.
func (s *sqlite) Close() error {
	if s != nil {
		if err := s.insertStmt.Close(); err != nil {
			return err
		}
		if err := s.queryStmt.Close(); err != nil {
			return err
		}
	}
	return nil
}

// CreateActivityWithInfo creates an activity, re-using a previous activity if
// available.
func (s *sqlite) CreateActivityWithInfo(ctx context.Context, filename string, gpsTrack io.Reader, activityInfo *doarama.ActivityInfo) (*doarama.Activity, error) {
	content, err := ioutil.ReadAll(gpsTrack)
	if err != nil {
		return nil, err
	}
	gpsTrackSha256 := sha256.Sum256(content)
	activityInfoSha256 := sha256.Sum256(structhash.Dump(activityInfo, 0))
	var activityID int
	switch err := s.queryStmt.QueryRow(gpsTrackSha256[:], activityInfoSha256[:]).Scan(&activityID); err {
	case nil:
		return &doarama.Activity{
			Client: s.client,
			ID:     activityID,
		}, nil
	case sql.ErrNoRows:
		activity, err := s.client.CreateActivityWithInfo(ctx, filename, bytes.NewBuffer(content), activityInfo)
		if err != nil {
			return activity, err
		}
		_, err = s.insertStmt.Exec(activity.ID, gpsTrackSha256[:], activityInfoSha256[:])
		return activity, err
	default:
		return nil, err
	}
}
