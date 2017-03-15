package controllers

import (
	"encoding/json"
	"errors"
	"github.com/revel/revel"
	"gitlab.bt.bpc.in/DevOps/release/app/models"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type ReleaseController struct {
	*revel.Controller
}

func (c ReleaseController) Index() revel.Result {
	var (
		releases []models.Release
		err      error
	)
	m := make(map[string]string)

	for k := range c.Params.Query {
		m[k] = c.Params.Query.Get(k)
	}

	releases, err = models.GetReleases(m)
	if err != nil {
		errResp := buildErrResponse(err, "500")
		c.Response.Status = 500
		return c.RenderJson(errResp)
	}
	c.Response.Status = 200
	return c.RenderJson(releases)
}

func (c ReleaseController) GetCSV(project, product string) revel.Result {
	var (
		releases []models.Release
		err      error
		s        []string
	)
	m := make(map[string]string)
	for k := range c.Params.Query {
		m[k] = c.Params.Query.Get(k)
	}
	releases, err = models.GetReleases(m)
	if err != nil {
		errResp := buildErrResponse(err, "500")
		c.Response.Status = 500
		return c.RenderJson(errResp)
	}

	for _, v := range releases {
		s = append(s, v.Name)
	}

	c.Response.Status = 200
	return c.RenderText(strings.Join(s, ","))
}

func (c ReleaseController) Show(id string) revel.Result {
	var (
		release   models.Release
		err       error
		releaseID bson.ObjectId
	)

	if id == "" {
		errResp := buildErrResponse(errors.New("Invalid release id format"), "400")
		c.Response.Status = 400
		return c.RenderJson(errResp)
	}

	releaseID, err = convertToObjectIdHex(id)
	if err != nil {
		errResp := buildErrResponse(errors.New("Invalid release id format"), "400")
		c.Response.Status = 400
		return c.RenderJson(errResp)
	}

	release, err = models.GetRelease(releaseID)
	if err != nil {
		errResp := buildErrResponse(err, "500")
		c.Response.Status = 500
		return c.RenderJson(errResp)
	}

	c.Response.Status = 200
	return c.RenderJson(release)
}

func (c ReleaseController) Create() revel.Result {
	var (
		release models.Release
		err     error
	)

	err = json.NewDecoder(c.Request.Body).Decode(&release)
	if err != nil {
		errResp := buildErrResponse(err, "403")
		c.Response.Status = 403
		return c.RenderJson(errResp)
	}

	release, err = models.AddRelease(release)
	if err != nil {
		errResp := buildErrResponse(err, "500")
		c.Response.Status = 500
		return c.RenderJson(errResp)
	}
	c.Response.Status = 201
	return c.RenderJson(release)
}

func (c ReleaseController) Update() revel.Result {
	var (
		release models.Release
		err     error
	)
	err = json.NewDecoder(c.Request.Body).Decode(&release)
	if err != nil {
		errResp := buildErrResponse(err, "400")
		c.Response.Status = 400
		return c.RenderJson(errResp)
	}

	err = release.UpdateRelease()
	if err != nil {
		errResp := buildErrResponse(err, "500")
		c.Response.Status = 500
		return c.RenderJson(errResp)
	}
	return c.RenderJson(release)
}

func (c ReleaseController) Delete(id string) revel.Result {
	var (
		err       error
		release   models.Release
		releaseID bson.ObjectId
	)
	if id == "" {
		errResp := buildErrResponse(errors.New("Invalid release id format"), "400")
		c.Response.Status = 400
		return c.RenderJson(errResp)
	}

	releaseID, err = convertToObjectIdHex(id)
	if err != nil {
		errResp := buildErrResponse(errors.New("Invalid release id format"), "400")
		c.Response.Status = 400
		return c.RenderJson(errResp)
	}

	release, err = models.GetRelease(releaseID)
	if err != nil {
		errResp := buildErrResponse(err, "500")
		c.Response.Status = 500
		return c.RenderJson(errResp)
	}
	err = release.DeleteRelease()
	if err != nil {
		errResp := buildErrResponse(err, "500")
		c.Response.Status = 500
		return c.RenderJson(errResp)
	}
	c.Response.Status = 204
	return c.RenderJson(nil)
}
