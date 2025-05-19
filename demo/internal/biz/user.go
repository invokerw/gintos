package biz

import (
	"context"
	"github.com/gin-gonic/gin"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/data/ent"
	"github/invokerw/gintos/demo/internal/data/ent/user"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/demo/internal/pkg/trans"
	"github/invokerw/gintos/demo/internal/pkg/upload"
	"github/invokerw/gintos/log"
	"io"
	"strings"
)

type UserUsecase struct {
	repo UserRepo
	oss  upload.OSS
	log  *log.Helper
}

func NewUserUsecase(repo UserRepo, oss upload.OSS, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, oss: oss, log: log.NewHelper(log.With(logger, "usecase", "user"))}
}

func (uc *UserUsecase) CreateUser(ctx *gin.Context, user *common.User, ignorePassword bool) (*common.User, error) {
	u, err := uc.repo.CreateUsers(ctx, []*common.User{user})
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(u, ignorePassword), nil
}

func (uc *UserUsecase) CreateUsers(ctx *gin.Context, users []*common.User, ignorePassword bool) (*common.User, error) {
	u, err := uc.repo.CreateUsers(ctx, users)
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(u, ignorePassword), nil
}

func (uc *UserUsecase) GetUser(ctx *gin.Context, username string, ignorePassword bool) (*common.User, error) {
	u, err := uc.repo.GetUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(u, ignorePassword), nil
}

func (uc *UserUsecase) GetUserByID(ctx *gin.Context, id uint64, ignorePassword bool) (*common.User, error) {
	u, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(u, ignorePassword), nil
}

func (uc *UserUsecase) DeleteUsers(ctx *gin.Context, names []string) error {
	err := uc.repo.DeleteUsers(ctx, names)
	if err != nil {
		return err
	}
	return nil
}

func (uc *UserUsecase) UpdateUsers(ctx context.Context, users []*common.User, ignorePassword bool) ([]*common.User, error) {
	us, err := uc.repo.UpdateUsers(ctx, users)
	if err != nil {
		return nil, err
	}
	return uc.convertToUsers(us, ignorePassword), nil
}

func (uc *UserUsecase) GetUserList(ctx *gin.Context, req *admin.GetUserListRequest, ignorePassword bool) ([]*common.User, error) {
	users, err := uc.repo.GetUserList(ctx, req)
	if err != nil {
		return nil, err
	}
	var userList []*common.User
	for _, u := range users {
		userList = append(userList, uc.convertToUser(u, ignorePassword))
	}
	return userList, nil
}

func (uc *UserUsecase) GetUserCount(ctx context.Context) (int, error) {
	count, err := uc.repo.Count(ctx)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (uc *UserUsecase) UpdateUserAvatar(ctx *gin.Context, id uint64, avatarData io.Reader, ext string) (*common.User, error) {
	if ext != "png" && ext != "jpg" && ext != "jpeg" {
		return nil, errs.ErrAvatarExtWrong
	}
	u, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	url, _, err := uc.oss.UploadFile(ctx, FILE_CATEGORY_USER_AVATAR, avatarData, *u.Username+"_avatar."+ext)
	if err != nil {
		return nil, err
	}
	uc.log.Infof("upload avatar %s", url)
	if u.Avatar != nil && *u.Avatar != "" {
		old := strings.Split(*u.Avatar, "/")
		if len(old) > 0 {
			key := old[len(old)-1]
			uc.log.Infof("delete old avatar %s", key)
			err = uc.oss.DeleteFile(ctx, FILE_CATEGORY_USER_AVATAR, key)
			if err != nil {
				uc.log.Errorf("delete old avatar %s failed: %v", key, err)
			}
		}
	}
	us, err := uc.repo.UpdateUsers(ctx, []*common.User{
		{
			Id:     trans.Ptr(u.ID),
			Avatar: trans.Ptr(url),
		},
	})
	if err != nil {
		return nil, err
	}
	return uc.convertToUser(us[0], true), nil
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

func (uc *UserUsecase) convertToUsers(us []*ent.User, ignorePassword bool) []*common.User {
	ret := make([]*common.User, 0, len(us))
	for _, u := range us {
		ret = append(ret, uc.convertToUser(u, ignorePassword))
	}
	return ret
}

func (uc *UserUsecase) convertToUser(u *ent.User, ignorePassword bool) *common.User {
	roleName := ""
	if u.RoleName != nil {
		roleName = *u.RoleName
	}
	pass := u.Password
	if ignorePassword {
		pass = trans.Ptr("********")
	}
	return &common.User{
		Id:            trans.Uint64(u.ID),
		RoleName:      u.RoleName,
		CreateBy:      u.CreateBy,
		UpdateBy:      u.UpdateBy,
		Username:      u.Username,
		Nickname:      u.Nickname,
		Password:      pass,
		Avatar:        u.Avatar,
		Email:         u.Email,
		Phone:         u.Phone,
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
