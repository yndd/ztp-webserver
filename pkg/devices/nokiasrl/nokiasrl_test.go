package nokiasrl

import (
	_ "embed"
	"fmt"
	"testing"
)

func TestNokiaSRL_handleScript(t *testing.T) {

}

func TestNokiaSRL_handleConfig(t *testing.T) {

}

func TestNokiaSRL_SetWebserverSetupper(t *testing.T) {

}

func Test_getTemplatingFunctions(t *testing.T) {
	tf := getTemplatingFunctions()

	expected_functions_defined := []string{"join", "jsonstringify"}

	for _, x := range expected_functions_defined {
		if funcname, exists := tf[x]; !exists {
			t.Errorf("Expected to have '%s' available as a function, but was not.", funcname)
		}
	}
}

func Test_getTemplatingFunctionsJsonstringifyArray(t *testing.T) {

	input := []string{"a", "b", "c"}
	result := jsonStringifyArray(input)

	for x, val := range result {
		expected := fmt.Sprintf("\"%s\"", input[x])
		if val != expected {
			t.Errorf("Expected result to be '%s' but was '%s'", expected, val)
		}
	}

}

func TestNewNokiaSRL(t *testing.T) {
	nnsrl := GetNokiaSRL()
	// check for singleton
	if nnsrl != GetNokiaSRL() {
		t.Errorf("NokiaSRL seems to not be a singleton.")
	}
}
