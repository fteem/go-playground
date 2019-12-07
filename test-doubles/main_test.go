package main

import "testing"

type DummySearcher struct{}

func (ds DummySearcher) Search(people []*Person, firstName, lastName string) *Person {
	return &Person{}
}

func TestFindReturnsError(t *testing.T) {
	phonebook := &Phonebook{}

	want := ErrMissingArgs
	_, got := phonebook.Find(DummySearcher{}, "", "")

	if got != want {
		t.Errorf("Want '%s', got '%s'", want, got)
	}
}

type StubSearcher struct {
	phone string
}

func (ss StubSearcher) Search(people []*Person, firstName, lastName string) *Person {
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     ss.phone,
	}
}

func TestFindReturnsPerson(t *testing.T) {
	fakePhone := "+31 65 222 333"
	phonebook := &Phonebook{}

	phone, _ := phonebook.Find(StubSearcher{fakePhone}, "Jane", "Doe")

	if phone != fakePhone {
		t.Errorf("Want '%s', got '%s'", fakePhone, phone)
	}
}

type SpySearcher struct {
	phone           string
	searchWasCalled bool
}

func (ss *SpySearcher) Search(people []*Person, firstName, lastName string) *Person {
	ss.searchWasCalled = true
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     ss.phone,
	}
}

func TestFindCallsSearchAndReturnsPerson(t *testing.T) {
	fakePhone := "+31 65 222 333"
	phonebook := &Phonebook{}
	spy := &SpySearcher{phone: fakePhone}

	phone, _ := phonebook.Find(spy, "Jane", "Doe")

	if !spy.searchWasCalled {
		t.Errorf("Expected to call 'Search' in 'Find', but it wasn't.")
	}

	if phone != fakePhone {
		t.Errorf("Want '%s', got '%s'", fakePhone, phone)
	}
}

type MockSearcher struct {
	phone         string
	methodsToCall map[string]bool
}

func (ms *MockSearcher) Search(people []*Person, firstName, lastName string) *Person {
	ms.methodsToCall["Search"] = true
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Phone:     ms.phone,
	}
}

func (ms *MockSearcher) ExpectToCall(methodName string) {
	if ms.methodsToCall == nil {
		ms.methodsToCall = make(map[string]bool)
	}
	ms.methodsToCall[methodName] = false
}

func (ms *MockSearcher) Verify(t *testing.T) {
	for methodName, called := range ms.methodsToCall {
		if !called {
			t.Errorf("Expected to call '%s', but it wasn't.", methodName)
		}
	}
}

func TestFindCallsSearchAndReturnsPersonUsingMock(t *testing.T) {
	fakePhone := "+31 65 222 333"
	phonebook := &Phonebook{}
	mock := &MockSearcher{phone: fakePhone}
	mock.ExpectToCall("Search")

	phone, _ := phonebook.Find(mock, "Jane", "Doe")

	if phone != fakePhone {
		t.Errorf("Want '%s', got '%s'", fakePhone, phone)
	}

	mock.Verify(t)
}

type FakeSearcher struct{}

func (fs FakeSearcher) Search(people []*Person, firstName string, lastName string) *Person {
	if len(people) == 0 {
		return nil
	}

	return people[0]
}

func TestFindCallsSearchAndReturnsEmptyStringForNoPerson(t *testing.T) {
	phonebook := &Phonebook{}
	fake := &FakeSearcher{}

	phone, _ := phonebook.Find(fake, "Jane", "Doe")

	if phone != "" {
		t.Errorf("Wanted '', got '%s'", phone)
	}
}
