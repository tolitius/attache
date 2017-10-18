package attache

import (
	"reflect"
	"testing"
)

func TestOutFail(t *testing.T) {

	spec := ConsulSpec{Address: "answer-is-not-always:42"}

	_, err := ConsulToMap(spec, "/hubble")

	if err == nil {
		t.Errorf("expected an \"no such host\" error, but did not see one")
	}

}

func TestInFail(t *testing.T) {

	spec := ConsulSpec{Address: "answer-is-not-always:42"}

	toConsul := make(map[string]string)
	toConsul["foo"] = "bar"

	_, err := MapToConsul(spec, toConsul)

	if err == nil {
		t.Errorf("expected an \"no such host\" error, but did not see one")
	}

}

func TestInAndOut(t *testing.T) {

	spec := ConsulSpec{Address: "localhost:8500"}

	toConsul := make(map[string]string)

	toConsul["hubble/store"] = "spacecraft://tape"
	toConsul["hubble/camera/mode"] = "color"
	toConsul["hubble/mission/target"] = "Horsehead Nebula"

	duration, _ := MapToConsul(spec, toConsul)
	fromConsul, _ := ConsulToMap(spec, "/hubble")

	t.Logf("writing to Consul took: %v", duration)

	if !reflect.DeepEqual(toConsul, fromConsul) {
		t.Errorf("expected %+v, but read only %+v form consul", toConsul, fromConsul)
	}
}
