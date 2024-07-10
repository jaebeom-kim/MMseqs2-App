package main

import (
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"
	"sort"
	"strings"
)

type MetabuliClassifyJob struct {
	Size     int      `json:"size" validate:"required"`
	Database []string `json:"database" validate:"required"`
	Mode     string   `json:"mode"`
	Query    string   `json:"q1"`
	Query2   string   `json:"q2"`
	Outdir   string   `json:"outdir"`
	Jobid    string   `json:"jobid"`
}

func (r MetabuliClassifyJob) Hash() Id {
	h := sha256.New224()
	h.Write(([]byte)(JobMetabuliClassify))
	h.Write([]byte(r.Query))
	h.Write([]byte(r.Mode))

	sort.Strings(r.Database)

	for _, value := range r.Database {
		h.Write([]byte(value))
	}

	bs := h.Sum(nil)
	return Id(base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bs))
}

func (r MetabuliClassifyJob) Rank() float64 {
	return float64(r.Size * max(len(r.Database), 1))
}

func (r MetabuliClassifyJob) WritePDB(path string) error {
	err := os.WriteFile(path, []byte(r.Query), 0644)
	if err != nil {
		return err
	}
	return nil
}

func NewMetabuliClassifyJobRequest(
	query string,
	query2 string,
	dbs []string,
	outdir string,
	jobid string,
	validDbs []Params,
	mode string,
	resultPath string,
	email string) (JobRequest, error) {

	job := MetabuliClassifyJob{
		max(strings.Count(query, ">"), 1),
		dbs,
		mode,
		query,
		query2,
		outdir,
		jobid,
	}

	request := JobRequest{
		job.Hash(),
		StatusPending,
		JobMetabuliClassify,
		job,
		email,
	}

	ids := make([]string, len(validDbs))
	for i, item := range validDbs {
		log.Println(item.Path)
		ids[i] = item.Path
	}

	// for _, item := range job.Database {
	// 	idx := isIn(item, ids)
	// 	if idx == -1 {
	// 		return request, errors.New("selected databases are not valid")
	// 	}
	// }
	return request, nil
}
