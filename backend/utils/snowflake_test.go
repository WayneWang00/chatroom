package utils

import "testing"

func TestGetId(t *testing.T) {
	worker, err := NewWorker(1)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Logf("new worker:%+v\n", worker)

	id := worker.GetId()
	t.Log("id:", id)
}
