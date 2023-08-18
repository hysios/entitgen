// Code generated by entitgen. DO NOT EDIT.
package out

import (
	"time"

	pb "github.com/hysios/entitgen/example/gen/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Member struct {
	ID           uint
	UserID       uint
	EnterpriseID uint
	User         *User
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ToProto converts the model to protobuf type.
func (m *Member) ToProto() *pb.Member {
	return &pb.Member{
		Id:           uint32(m.ID),
		UserId:       uint32(m.UserID),
		EnterpriseId: uint32(m.EnterpriseID),
		User:         m.User.ToProto(),
		CreatedAt:    timestamppb.New(m.CreatedAt),
		UpdatedAt:    timestamppb.New(m.UpdatedAt),
	}
}

// FromProto converts the protobuf type to model.
func (m *Member) FromProto(pMember *pb.Member) *Member {
	return &Member{
		ID:           uint(pMember.Id),
		UserID:       uint(pMember.UserId),
		EnterpriseID: uint(pMember.EnterpriseId),
		User:         (*User)(nil).FromProto(pMember.User),
		CreatedAt:    pMember.CreatedAt.AsTime(),
		UpdatedAt:    pMember.UpdatedAt.AsTime(),
	}
}

func MemberFromProto(pMember *pb.Member) *Member {
	return (*Member)(nil).FromProto(pMember)
}
