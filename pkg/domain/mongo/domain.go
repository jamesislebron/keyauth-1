package mongo

import (
	"context"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/types/ftime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/keyauth/common/types"
	"github.com/infraboard/keyauth/pkg/domain"
)

func (s *service) CreateDomain(ownerID string, req *domain.CreateDomainRequst) (*domain.Domain, error) {
	d, err := domain.New(ownerID, req)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}
	if _, err := s.col.InsertOne(context.TODO(), d); err != nil {
		return nil, exception.NewInternalServerError("inserted a domain document error, %s", err)
	}

	return d, nil
}

func (s *service) DescriptionDomain(req *domain.DescriptDomainRequest) (*domain.Domain, error) {
	r, err := newDescDomainRequest(req)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	d := domain.NewDefault()
	if err := s.col.FindOne(context.TODO(), r.FindFilter()).Decode(d); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("domain %s not found", req)
		}

		return nil, exception.NewInternalServerError("find domain %s error, %s", req.Name, err)
	}

	return d, nil
}

func (s *service) QueryDomain(req *domain.QueryDomainRequest) (*domain.Set, error) {
	r := newQueryDomainRequest(req)
	resp, err := s.col.Find(context.TODO(), r.FindFilter(), r.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find domain error, error is %s", err)
	}

	domainSet := domain.NewDomainSet(req.PageRequest)
	// 循环
	for resp.Next(context.TODO()) {
		d := new(domain.Domain)
		if err := resp.Decode(d); err != nil {
			return nil, exception.NewInternalServerError("decode domain error, error is %s", err)
		}

		domainSet.Add(d)
	}

	// count
	count, err := s.col.CountDocuments(context.TODO(), r.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get device count error, error is %s", err)
	}
	domainSet.Total = count

	return domainSet, nil
}

func (s *service) UpdateDomain(req *domain.UpdateDomainRequest) (*domain.Domain, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	d, err := s.DescriptionDomain(domain.NewDescriptDomainRequestWithName(req.Name))
	if err != nil {
		return nil, err
	}
	switch req.UpdateMode {
	case types.PutUpdateMode:
		*d.CreateDomainRequst = *req.CreateDomainRequst
	case types.PatchUpdateMode:
		d.CreateDomainRequst.Patch(req.CreateDomainRequst)
	default:
		return nil, exception.NewBadRequest("unknown update mode: %s", req.UpdateMode)
	}

	d.UpdateAt = ftime.Now()
	_, err = s.col.UpdateOne(context.TODO(), bson.M{"_id": d.Name}, bson.M{"$set": d})
	if err != nil {
		return nil, exception.NewInternalServerError("update domain(%s) error, %s", d.Name, err)
	}

	return d, nil
}

func (s *service) DeleteDomain(id string) error {
	result, err := s.col.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return exception.NewInternalServerError("delete domain(%s) error, %s", id, err)
	}

	if result.DeletedCount == 0 {
		return exception.NewNotFound("domain %s not found", id)
	}

	return nil
}
