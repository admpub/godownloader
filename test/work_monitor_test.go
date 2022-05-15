package dtest

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/admpub/godownloader/monitor"
)

type TestWork struct {
	From, sleep, To int
}

func (tw TestWork) GetProgress() interface{} {
	return tw.From
}
func (tw *TestWork) DoWork(context.Context) (bool, error) {
	time.Sleep(time.Millisecond * 300)
	tw.From += 1
	log.Println("info: exec DoWork", tw.From)
	if tw.From == tw.To {
		log.Println("done")
		return true, nil
	}
	if tw.From > tw.To {
		return false, errors.New("failed")
	}
	return false, nil
}

func (tw *TestWork) BeforeRun() error {
	log.Println("info: exec before run")
	return nil
}
func (tw *TestWork) AfterStop() error {
	log.Println("info: exec after stop")
	return nil
}

func (tw *TestWork) IsPartialDownload() bool {
	return true
}

func (tw *TestWork) ResetProgress() {
	tw.From = 0
}

func TestWorker(t *testing.T) {
	tes := new(monitor.MonitoredWorker)
	itw := &TestWork{From: 0, To: 8, sleep: 300}
	tes.Itw = itw
	ctx := context.Background()
	tes.Start(ctx)
	log.Println(tes.Start(ctx))
	time.Sleep(time.Second * 1)
	if tes.GetState() != monitor.Running {
		t.Error("Expected Running(1)")
		return
	}
	tes.Stop(ctx)
	if tes.GetState() != monitor.Stopped {
		t.Errorf("Expected Stoped(0): %v", tes.GetState().String())
		return
	}
	tes.Start(ctx)
	time.Sleep(time.Second * 9)
	if tes.GetState() != monitor.Completed {
		t.Errorf("Expected Comlete(3): %v", tes.GetState().String())
		return
	}

	tes.Start(ctx)
	time.Sleep(time.Second * 1)
	if tes.GetState() != monitor.Completed {
		t.Errorf("Expected Failed(3): %v", tes.GetState().String())
		return
	}
}
