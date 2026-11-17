package platform

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestGetLicensesForUser_ReturnsTheUsersLicenses(t *testing.T) {
	client, _ := newStubbedClient(t, `{"data":{"user":{"licenses":[{"id":3,"type":"pixo-builder","name":"Pixo Builder"}]}}}`)

	licenses, err := client.GetLicensesForUser(context.Background(), 12)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []License{{ID: 3, Type: "pixo-builder", Name: "Pixo Builder"}}
	if !reflect.DeepEqual(licenses, expected) {
		t.Errorf("expected %+v, got %+v", expected, licenses)
	}
}

func TestGetLicensesForUser_AsksForTheGivenUsersLicenses(t *testing.T) {
	client, asked := newStubbedClient(t, `{"data":{"user":{"licenses":[]}}}`)

	if _, err := client.GetLicensesForUser(context.Background(), 12); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(asked.Query, "licenses { id type name }") {
		t.Errorf("expected the query to ask for the user's licenses, got %s", asked.Query)
	}
	if asked.Variables["id"] != float64(12) {
		t.Errorf("expected the query to ask for user 12, got %v", asked.Variables["id"])
	}
}

func TestGetLicensesForUser_ReturnsTheErrorWhenThePlatformRefuses(t *testing.T) {
	client, _ := newStubbedClient(t, `{"errors":[{"message":"forbidden"}],"data":null}`)

	if _, err := client.GetLicensesForUser(context.Background(), 12); err == nil {
		t.Fatal("expected an error")
	}
}
