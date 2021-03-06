package micro

import (
	"errors"

	"github.com/infraboard/mcube/http/request"
)

// Service token管理服务
type Service interface {
	CreateService(req *CreateMicroRequest) (*Micro, error)
	QueryService(req *QueryMicroRequest) (*Set, error)
	DescribeService(req *DescribeMicroRequest) (*Micro, error)
	DeleteService(id string) error
}

// NewQueryMicroRequest 列表查询请求
func NewQueryMicroRequest(pageReq *request.PageRequest) *QueryMicroRequest {
	return &QueryMicroRequest{
		PageRequest: pageReq,
	}
}

// QueryMicroRequest 查询应用列表
type QueryMicroRequest struct {
	*request.PageRequest
}

// NewDescriptServiceRequest new实例
func NewDescriptServiceRequest() *DescribeMicroRequest {
	return &DescribeMicroRequest{}
}

// DescribeMicroRequest 查询应用详情
type DescribeMicroRequest struct {
	Name string
	ID   string
}

// Validate 校验详情查询请求
func (req *DescribeMicroRequest) Validate() error {
	if req.ID == "" && req.Name == "" {
		return errors.New("id, name or service_id is required")
	}

	return nil
}
