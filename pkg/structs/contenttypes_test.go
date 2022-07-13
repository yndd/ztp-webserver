package structs

import "testing"

func TestString2ContentTypes(t *testing.T) {

	for _, s := range []string{"software", "config", "md5hashfile", "script"} {
		_, err := String2ContentTypes(s)
		if err != nil {
			t.Errorf("Error converting string '%s' to ContentType enum", s)
		}
	}

	badContentTypeString := "NonExistingContentType"
	_, err := String2ContentTypes(badContentTypeString)
	if err == nil {
		t.Errorf("Trying to parse the bad ContentType from string '%s' did not yield an error, as expected.", badContentTypeString)
	}

}
