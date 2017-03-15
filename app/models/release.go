package models

import (
	"gitlab.bt.bpc.in/DevOps/release/app/models/mongodb"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type Release struct {
	ID          bson.ObjectId `json:"id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Product     string        `json:"product" bson:"product"`
	Project     string        `json:"project" bson:"project"`
	Jira        string        `json:"jira" bson:"jira"`
	IsActive    bool          `json:"is_active" bson:"is_active"`
	Released    bool          `json:"released" bson:"released"`
	DateRelease time.Time     `json:"date_release" bson:"date_release"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

func newReleaseCollection() *mongodb.Collection {
	c := mongodb.NewCollectionSession("release")
	index := mgo.Index{
		Key:        []string{"project", "product", "name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.Session.EnsureIndex(index)

	if err != nil {
		panic(err)
	}
	return c
}

// AddRelease insert a new Release into database and returns
// last inserted release on success.
func AddRelease(m Release) (release Release, err error) {
	c := newReleaseCollection()
	defer c.Close()

	m.ID = bson.NewObjectId()
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return m, c.Session.Insert(m)
}

// UpdateRelease update a Release into database and returns
// last nil on success.
func (m Release) UpdateRelease() error {
	c := newReleaseCollection()
	defer c.Close()

	err := c.Session.Update(bson.M{
		"_id": m.ID,
	}, bson.M{
		"$set": bson.M{
			"name": m.Name, "product": m.Product, "project": m.Project, "jira": m.Jira, "is_active": m.IsActive, "released": m.Released, "date_release": m.DateRelease, "updatedAt": time.Now()},
	})
	return err
}

// DeleteRelease Delete Release from database and returns
// last nil on success.
func (m Release) DeleteRelease() error {
	c := newReleaseCollection()
	defer c.Close()

	err := c.Session.Remove(bson.M{"_id": m.ID})
	return err
}

// GetReleases Get all Release from database and returns
// list of Release on success
func GetReleases(args map[string]string) ([]Release, error) {
	var (
		releases []Release
		err      error
		flag     bool
	)
	q := bson.M{}
	for k, v := range args {
		flag, err = strconv.ParseBool(v)
		if err == nil {
			q[k] = flag
		} else {
			q[k] = v
		}
	}

	c := newReleaseCollection()
	defer c.Close()

	err = c.Session.Find(q).Sort("-createdAt").All(&releases)
	return releases, err
}

// GetRelease Get a Release from database and returns
// a Release on success
func GetRelease(id bson.ObjectId) (Release, error) {
	var (
		release Release
		err     error
	)

	c := newReleaseCollection()
	defer c.Close()

	err = c.Session.Find(bson.M{"_id": id}).One(&release)
	return release, err
}
