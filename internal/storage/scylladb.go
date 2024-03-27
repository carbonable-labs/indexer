// Package storage provides a ScyllaDB storage implementation for managing key-value pairs.
package storage

import (
	"github.com/gocql/gocql"
)

// ScyllaStorage represents a storage handler for ScyllaDB.
type ScyllaStorage struct {
	handle *gocql.Session
}

// Get retrieves the value associated with the given key from ScyllaDB.
func (s *ScyllaStorage) Get(id gocql.UUID) ([]byte, error) {
	var value []byte
	if err := s.handle.Query("SELECT value FROM Blocks WHERE key = ?", id).Scan(&value); err != nil {
		return nil, err
	}
	return value, nil
}

// Has checks if the given key exists in ScyllaDB.
func (s *ScyllaStorage) Has(id gocql.UUID) (bool, error) {
	var value []byte
	if err := s.handle.Query("SELECT value FROM Blocks WHERE key = ?", id).Scan(&value); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Set sets the value for the given key in ScyllaDB.
func (s *ScyllaStorage) Set(id gocql.UUID, value []byte) error {
	if err := s.handle.Query("INSERT INTO Blocks (key, value) VALUES (?, ?)", id, value).Exec(); err != nil {
		return err
	}
	return nil
}

// Scan retrieves all key-value pairs from ScyllaDB.
func (s *ScyllaStorage) Scan() ([]gocql.UUID, [][]byte, error) {
	iter := s.handle.Query("SELECT key, value FROM Blocks").Iter()
	defer iter.Close()

	var ids []gocql.UUID
	var values [][]byte

	var id gocql.UUID
	var value []byte

	for iter.Scan(&id, &value) {
		ids = append(ids, id)
		values = append(values, value)
	}

	if err := iter.Close(); err != nil {
		return nil, nil, err
	}

	return ids, values, nil
}

// NewScyllaStorage creates a new instance of ScyllaStorage.
func NewScyllaStorage() (*ScyllaStorage, error) {
	cluster = gocql.NewCluster("127.0.0.1")
	//cluster.Hosts = append(cluster.Hosts, "127.0.0.2", "127.0.0.3") // Add more addresses to the cluster.
	cluster.Keyspace = "index_keyspace"
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "admin",
		Password: "admin",
	}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &ScyllaStorage{
		handle: session,
	}, nil
}
