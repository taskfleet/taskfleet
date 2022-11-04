package db

import (
	"github.com/jmoiron/sqlx"
	"go.taskfleet.io/packages/postgres"
)

// InstanceIterator is an iterator that yields compute instances from a listing query.
type InstanceIterator struct {
	db   *postgres.Connection
	rows *sqlx.Rows
	err  error
}

// InstanceResult wraps a result from an iterator over compute instances.
type InstanceResult struct {
	Instance *Instance
	Err      error
}

// ForEach iterates over all items of this iterator and executes the callback for each item. If
// retrieving an item fails or the callback function returns an error, the error is returned
// immediately. If no error is returned, the number of items that was iterated over is also
// returned.
func (i *InstanceIterator) ForEach(callback func(*Instance) error) (int, error) {
	count := 0
	for item := range i.Iter() {
		if item.Err != nil {
			return 0, item.Err
		}
		if err := callback(item.Instance); err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

// Iter returns a channel delivering all of the iterator's items. If an error occurs, it is
// delivered over the channel and the channel is subsequently closed. This method must NOT be
// called concurrently and must only be called once. Either of these operations causes undefined
// behavior.
func (i *InstanceIterator) Iter() <-chan InstanceResult {
	ch := make(chan InstanceResult)
	go func() {
		if i.err != nil {
			ch <- InstanceResult{Err: i.err}
			close(ch)
			return
		}

		for i.rows.Next() {
			var value Instance
			if err := i.rows.StructScan(&value); err != nil {
				ch <- InstanceResult{Err: err}
			}
			value.db = i.db
			ch <- InstanceResult{Instance: &value}
		}
		if err := i.rows.Close(); err != nil {
			ch <- InstanceResult{Err: err}
		}
		close(ch)
	}()
	return ch
}

// Collect collects all queried instances and returns an error if anything fails. This method must
// not be called more than once and it must not be called if `Iter` has already been called and
// vice-versa.
func (i *InstanceIterator) Collect() ([]*Instance, error) {
	if i.err != nil {
		return nil, i.err
	}

	result := []*Instance{}
	for value := range i.Iter() {
		if value.Err != nil {
			return nil, value.Err
		}
		result = append(result, value.Instance)
	}
	return result, nil
}
