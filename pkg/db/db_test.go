package db

import (
	"fmt"
	"testing"
)

type databaseOpName int

const (
	set databaseOpName = iota
	get
	deleteByID
	deleteByVal
	begin
	commit
	rollback
)

func (don databaseOpName) String() string {
	switch don {
	case set:
		return "set"
	case get:
		return "get"
	case deleteByID:
		return "deleteByID"
	case deleteByVal:
		return "deleteByVal"
	case begin:
		return "begin"
	case commit:
		return "commit"
	case rollback:
		return "rollback"
	default:
		return "unrecognized"
	}
}

type databaseOp struct {
	name        databaseOpName
	argKey      int
	argVal      string
	expectedVal string
	expectError bool
}

func execTestCase(testCaseOps []databaseOp, t *testing.T) {
	db := NewDatabase()

	for i, op := range testCaseOps {
		switch op.name {
		case set:
			db.Set(op.argKey, op.argVal)
		case get:
			val := db.Get(op.argKey)
			if val != op.expectedVal {
				t.Errorf("op #%d: %q by id %d value %q did not match expected value %q",
					i+1, op.name, op.argKey, val, op.expectedVal)
			}
		case deleteByID:
			db.DeleteByID(op.argKey)
		case deleteByVal:
			db.DeleteByValue(op.argVal)
		case begin:
			err := db.Begin()
			if op.expectError != (err != nil) {
				t.Errorf("op #%d: %q error expected %t, error received %t",
					i+1, op.name, op.expectError, err != nil)
			}
		case commit:
			err := db.Commit()
			if op.expectError != (err != nil) {
				t.Errorf("op #%d: %q error expected %t, error received %t",
					i+1, op.name, op.expectError, err != nil)
			}
		case rollback:
			err := db.Rollback()
			if op.expectError != (err != nil) {
				t.Errorf("op #%d: %q error expected %t, error received %t",
					i+1, op.name, op.expectError, err != nil)
			}
		default:
			t.Errorf("unrecognized test case command #%d: %s", i+1, op.name)
		}
	}
}

func TestDatabase(t *testing.T) {
	testCases := [][]databaseOp{
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: get, argKey: 1, expectedVal: "foo"},
			{name: get, argKey: 0, expectedVal: noVal},
			{name: set, argKey: 1, argVal: "bar"},
			{name: get, argKey: 1, expectedVal: "bar"},
		},
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: set, argKey: 2, argVal: "foo"},
			{name: deleteByID, argKey: 1},
			{name: get, argKey: 1, expectedVal: noVal},
			{name: get, argKey: 2, expectedVal: "foo"},
		},
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: set, argKey: 2, argVal: "foo"},
			{name: deleteByVal, argVal: "foo"},
			{name: get, argKey: 1, expectedVal: noVal},
			{name: get, argKey: 2, expectedVal: noVal},
		},
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: set, argKey: 2, argVal: "foo"},
			{name: begin, expectError: false},
			{name: set, argKey: 2, argVal: "bar"},
			{name: set, argKey: 3, argVal: "buzz"},
			{name: get, argKey: 1, expectedVal: "foo"},
			{name: get, argKey: 2, expectedVal: "bar"},
			{name: get, argKey: 3, expectedVal: "buzz"},
			{name: commit, expectError: false},
			{name: get, argKey: 1, expectedVal: "foo"},
			{name: get, argKey: 2, expectedVal: "bar"},
			{name: get, argKey: 3, expectedVal: "buzz"},
		},
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: set, argKey: 2, argVal: "foo"},
			{name: begin, expectError: false},
			{name: set, argKey: 2, argVal: "bar"},
			{name: set, argKey: 3, argVal: "buzz"},
			{name: get, argKey: 1, expectedVal: "foo"},
			{name: get, argKey: 2, expectedVal: "bar"},
			{name: get, argKey: 3, expectedVal: "buzz"},
			{name: rollback, expectError: false},
			{name: get, argKey: 1, expectedVal: "foo"},
			{name: get, argKey: 2, expectedVal: "foo"},
			{name: get, argKey: 3, expectedVal: noVal},
		},
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: set, argKey: 2, argVal: "foo"},
			{name: begin, expectError: false},
			{name: set, argKey: 3, argVal: "buzz"},
			{name: deleteByID, argKey: 3},
			{name: deleteByID, argKey: 2},
			{name: get, argKey: 2, expectedVal: noVal},
			{name: get, argKey: 3, expectedVal: noVal},
			{name: rollback, expectError: false},
			{name: get, argKey: 1, expectedVal: "foo"},
			{name: get, argKey: 2, expectedVal: "foo"},
			{name: get, argKey: 3, expectedVal: noVal},
		},
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: set, argKey: 2, argVal: "foo"},
			{name: set, argKey: 3, argVal: "buzz"},
			{name: begin, expectError: false},
			{name: deleteByID, argKey: 3},
			{name: deleteByID, argKey: 2},
			{name: get, argKey: 2, expectedVal: noVal},
			{name: get, argKey: 3, expectedVal: noVal},
			{name: rollback, expectError: false},
			{name: get, argKey: 1, expectedVal: "foo"},
			{name: get, argKey: 2, expectedVal: "foo"},
			{name: get, argKey: 3, expectedVal: "buzz"},
		},
		{
			{name: set, argKey: 1, argVal: "foo"},
			{name: set, argKey: 2, argVal: "foo"},
			{name: set, argKey: 3, argVal: "bar"},
			{name: begin, expectError: false},
			{name: set, argKey: 2, argVal: "bar"},
			{name: set, argKey: 3, argVal: "foo"},
			{name: set, argKey: 4, argVal: "foo"},
			{name: deleteByVal, argVal: "foo"},
			{name: get, argKey: 1, expectedVal: noVal},
			{name: get, argKey: 2, expectedVal: "bar"},
			{name: get, argKey: 3, expectedVal: "bar"},
			{name: get, argKey: 4, expectedVal: noVal},
			{name: commit, expectError: false},
			{name: get, argKey: 1, expectedVal: noVal},
			{name: get, argKey: 2, expectedVal: "bar"},
			{name: get, argKey: 3, expectedVal: "bar"},
			{name: get, argKey: 4, expectedVal: noVal},
		},
	}

	for i, tc := range testCases {
		tcName := fmt.Sprintf("TestDB%d", i+1)
		t.Run(tcName, func(t *testing.T) {
			execTestCase(tc, t)
		})
	}
}
