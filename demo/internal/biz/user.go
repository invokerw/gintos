package biz

import (
	"context"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/data/ent"
	"github/invokerw/gintos/demo/internal/data/ent/user"
	"github/invokerw/gintos/demo/internal/pkg/trans"
	"github/invokerw/gintos/log"

	"github.com/gin-gonic/gin"
)

type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(log.With(logger, "usecase", "user"))}
}

func (uc *UserUsecase) CreateUser(ctx *gin.Context, user *common.User) (*common.User, error) {
	u, err := uc.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(u), nil
}

func (uc *UserUsecase) GetUser(ctx *gin.Context, username string) (*common.User, error) {
	u, err := uc.repo.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(u), nil
}

func (uc *UserUsecase) GetUserByID(ctx *gin.Context, id uint64) (*common.User, error) {
	u, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(u), nil
}

func (uc *UserUsecase) DeleteUsers(ctx *gin.Context, names []string) error {
	err := uc.repo.DeleteUsers(ctx, names)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) UpdateUsers(ctx context.Context, users []*common.User) ([]*common.User, error) {
	us, err := uc.repo.UpdateUsers(ctx, users)
	if err != nil {
		return nil, err
	}
	return uc.convertToUsers(us), nil
}

func (uc *UserUsecase) GetUserList(ctx *gin.Context, req *admin.GetUserListRequest) ([]*common.User, error) {
	users, err := uc.repo.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}
	var userList []*common.User
	for _, u := range users {
		userList = append(userList, uc.convertToUser(u))
	}
	return userList, nil
}

func (uc *UserUsecase) convertToGender(g *user.Gender) *common.UserGender {
	if g == nil {
		return nil
	}
	find, ok := common.UserGender_value[g.String()]
	if !ok {
		return nil
	}
	return trans.Ptr(common.UserGender(find))
}

func (uc *UserUsecase) convertToStatus(s *user.Status) *common.UserStatus {
	if s == nil {
		return nil
	}
	find, ok := common.UserStatus_value[s.String()]
	if !ok {
		return nil
	}
	return trans.Ptr(common.UserStatus(find))
}

func (uc *UserUsecase) convertToAuthority(a *user.Authority) *common.UserAuthority {
	if a == nil {
		return nil
	}
	find, ok := common.UserAuthority_value[a.String()]
	if !ok {
		return nil
	}
	return trans.Ptr(common.UserAuthority(find))
}

func (uc *UserUsecase) convertToUsers(us []*ent.User) []*common.User {
	ret := make([]*common.User, 0, len(us))
	for _, u := range us {
		ret = append(ret, uc.convertToUser(u))
	}
	return ret
}

func (uc *UserUsecase) convertToUser(u *ent.User) *common.User {
	roleName := u.Edges.Role.Name
	return &common.User{
		Id:            trans.Uint64(u.ID),
		RoleName:      &roleName,
		CreateBy:      u.CreateBy,
		UpdateBy:      u.UpdateBy,
		UserName:      u.Username,
		NickName:      u.NickName,
		Password:      u.Password,
		Avatar:        u.Avatar,
		Email:         u.Email,
		Mobile:        u.Mobile,
		Gender:        uc.convertToGender(u.Gender),
		Remark:        u.Remark,
		LastLoginTime: u.LastLoginTime,
		Status:        uc.convertToStatus(u.Status),
		Authority:     uc.convertToAuthority(u.Authority),
		Roles:         []string{roleName},
		CreateTime:    trans.Ptr(u.CreateTime.Unix()),
		UpdateTime:    trans.Ptr(u.UpdateTime.Unix()),
	}
}
