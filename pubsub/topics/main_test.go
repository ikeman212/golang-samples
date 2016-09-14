// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/golang-samples/internal/testutil"
)

const topicID = "golang-samples-topic-example"

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestCreate(t *testing.T) {
	c := setup(t)
	if err := create(c, topicID); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}
	ok, err := c.Topic(topicID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want topic = %q", topicID)
	}
}

func TestList(t *testing.T) {
	c := setup(t)

	testutil.Flaky(t, 10, time.Second, func(r *testutil.R) {
		topics, err := list(c)
		if err != nil {
			r.Failf("failed to list topics: %v", err)
		}

		for _, t := range topics {
			if t.ID() == topicID {
				return // PASS
			}
		}

		topicNames := make([]string, len(topics))
		for i, t := range topics {
			topicNames[i] = t.ID()
		}
		r.Failf("got %+v; want a list with topic = %q", topicNames, topicID)
	})
}

func TestPublish(t *testing.T) {
	// Nothing much to do here, unless we are consuming.
	// TODO(jbd): Merge topics and subscriptions programs maybe?
	c := setup(t)
	if err := publish(c, topicID, "hello world"); err != nil {
		t.Errorf("failed to publish message: %v", err)
	}
}

func TestDelete(t *testing.T) {
	c := setup(t)
	if err := delete(c, topicID); err != nil {
		t.Fatalf("failed to delete subscription (%q): %v", topicID, err)
	}
	ok, err := c.Topic(topicID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if sub exists: %v", err)
	}
	if ok {
		t.Fatalf("got sub = %q; want none", topicID)
	}
}
