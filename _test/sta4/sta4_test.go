package sta4_test

import (
	"reflect"
	"testing"

	"github.com/leonard475192/go-stations/model"
)

func TestStation4(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		Target       interface{}
		FieldName    string
		WantKinds    []reflect.Kind
		JSONTagValue string
	}{
		"HealthzResponse has Message field": {
			Target:       model.HealthzResponse{},
			FieldName:    "Message",
			WantKinds:    []reflect.Kind{reflect.String},
			JSONTagValue: "message",
		},
	}

	for name, tc := range testcases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tp := reflect.TypeOf(tc.Target)
			f, ok := tp.FieldByName(tc.FieldName)
			if !ok {
				t.Error(tc.FieldName + " field が見つかりません")
				return
			}

			notFound := true
			for _, k := range tc.WantKinds {
				if f.Type.Kind() == k {
					notFound = false
					break
				}
			}
			if notFound {
				t.Errorf(tc.FieldName+" が期待している kind ではありませ, got = %s, want = %s", f.Type.Kind(), tc.WantKinds)
				return
			}

			v, ok := f.Tag.Lookup("json")
			if !ok {
				t.Error("json tag が見つかりません")
				return
			}

			if v != tc.JSONTagValue {
				t.Errorf("json tag の内容が期待している内容ではありません, got = %s, want = %s", v, tc.JSONTagValue)
			}
		})
	}
}
