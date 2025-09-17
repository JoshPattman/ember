package ember

import (
	"testing"
)

type TestCase interface {
	Name() string
	Test() error
}

func RunTests(t *testing.T, tests []TestCase) {
	for _, testCase := range tests {
		t.Run(testCase.Name(), func(t *testing.T) {
			err := testCase.Test()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
