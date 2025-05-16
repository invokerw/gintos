package data

import (
	"context"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/data/ent"
	"github/invokerw/gintos/demo/internal/data/ent/role"
	"github/invokerw/gintos/demo/internal/data/ent/user"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/log"
	"time"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "user_repo")),
	}
}

func (r *userRepo) convertToStatus(status *common.UserStatus) *user.Status {
	if status == nil {
		return nil
	}
	ret := user.Status(status.String())
	if err := user.StatusValidator(ret); err != nil {
		return nil
	}
	return &ret
}

func (r *userRepo) convertToGender(g *common.UserGender) *user.Gender {
	if g == nil {
		return nil
	}
	ret := user.Gender(g.String())
	if err := user.GenderValidator(ret); err != nil {
		return nil
	}
	return &ret
}

func (r *userRepo) convertToAuthority(a *common.UserAuthority) *user.Authority {
	if a == nil {
		return nil
	}
	ret := user.Authority(a.String())
	if err := user.AuthorityValidator(ret); err != nil {
		return nil
	}
	return &ret
}

func (r *userRepo) CreateUser(ctx context.Context, in *common.User) (*ent.User, error) {
	if in == nil || in.UserName == nil || in.Password == nil {
		return nil, errs.DBErrInvalidParam
	}
	roleName := in.GetRoleName()
	if roleName == "" {
		return nil, errs.DBErrInvalidParam
	}
	now := time.Now()
	var u *ent.User
	var err error
	err = WithTx(ctx, r.data.db, func(tx *ent.Tx) error {
		u, err = tx.User.Query().Where(user.Username(in.GetUserName())).Only(ctx)
		if u != nil {
			return errs.DBErrUserExist
		}
		var roleData *ent.Role
		roleData, err = tx.Role.Query().Where(role.Name(roleName)).Only(ctx)
		if err != nil {
			return err
		}
		uc := tx.User.Create().
			SetUsername(*in.UserName).
			SetPassword(*in.Password).
			SetNillableNickName(in.NickName).
			SetNillableCreateBy(in.CreateBy).
			SetCreateTime(now).
			SetUpdateTime(now).
			SetNillableRemark(in.Remark).
			SetNillableStatus(r.convertToStatus(in.Status)).
			SetNillableEmail(in.Email).
			SetNillablePhone(in.Phone).
			SetNillableAvatar(in.Avatar).
			SetNillableGender(r.convertToGender(in.Gender)).
			SetNillableAuthority(r.convertToAuthority(in.Authority)).
			SetRole(roleData).
			SetNillableLastLoginTime(in.LastLoginTime)

		u, err = uc.Save(ctx)
		if err != nil {
			return errs.DBErrEntError.Wrap(err)
		}
		return nil
	})

	return u, nil
}

func (r *userRepo) UpdateUsers(ctx context.Context, ins []*common.User) ([]*ent.User, error) {
	for _, in := range ins {
		if in == nil || in.UserName == nil {
			return nil, errs.DBErrInvalidParam
		}
	}

	ret := make([]*ent.User, 0, len(ins))
	err := WithTx(ctx, r.data.db, func(tx *ent.Tx) error {
		for _, in := range ins {
			var roleId *uint64
			roleName := in.GetRoleName()
			if roleName != "" {
				roleData, err := tx.Role.Query().Where(role.Name(roleName)).Only(ctx)
				if err == nil {
					roleId = &roleData.ID
				}
			}
			u, err := tx.User.Query().Where(user.Username(in.GetUserName())).Only(ctx)
			if err != nil {
				return errs.DBErrEntError.Wrap(err)
			}
			uc := u.Update().
				SetNillableNickName(in.NickName).
				SetNillableCreateBy(in.CreateBy).
				SetNillableRemark(in.Remark).
				SetNillableStatus(r.convertToStatus(in.Status)).
				SetNillableEmail(in.Email).
				SetNillablePhone(in.Phone).
				SetNillableAvatar(in.Avatar).
				SetNillableGender(r.convertToGender(in.Gender)).
				SetNillableAuthority(r.convertToAuthority(in.Authority)).
				SetNillableRoleID(roleId).
				SetNillableLastLoginTime(in.LastLoginTime).
				SetUpdateTime(time.Now())
			u, err = uc.Save(ctx)
			if err != nil {
				return errs.DBErrEntError.Wrap(err)
			}
			ret = append(ret, u)
		}
		return nil
	})

	if err != nil {
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return ret, nil
}

func (r *userRepo) DeleteUsers(ctx context.Context, names []string) error {
	var err error
	err = WithTx(ctx, r.data.db, func(tx *ent.Tx) error {
		var u *ent.User
		for _, name := range names {
			u, err = tx.User.Query().Where(user.Username(name)).Only(ctx)
			if err != nil {
				return errs.DBErrEntError.Wrap(err)
			}
			err = tx.User.DeleteOne(u).Exec(ctx)
			if err != nil {
				return errs.DBErrEntError.Wrap(err)
			}
		}
		return nil
	})
	if err != nil {
		return errs.DBErrEntError.Wrap(err)
	}
	return err
}

func (r *userRepo) GetUser(ctx context.Context, username string) (*ent.User, error) {
	u, err := r.data.db.User.Query().Where(user.Username(username)).WithRole().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.DBErrUserNotFound
		}
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return u, nil
}

func (r *userRepo) GetUserList(ctx context.Context, req *admin.GetUserListRequest) ([]*ent.User, error) {
	if req == nil {
		return nil, errs.DBErrInvalidParam
	}
	q := r.data.db.User.Query().Offset(int(req.Page.Offset)).Limit(int(req.Page.PageSize)).WithRole()
	if req.Username != nil && *req.Username != "" {
		q.Where(user.UsernameContains(*req.Username))
	}
	if req.Phone != nil && *req.Phone != "" {
		q.Where(user.PhoneContains(*req.Phone))
	}
	if req.Email != nil && *req.Email != "" {
		q.Where(user.EmailContains(*req.Email))
	}
	if req.Status != nil {
		s := user.Status(*req.Status)
		if err := user.StatusValidator(s); err == nil {
			q.Where(user.StatusEQ(s))
		}
	}

	users, err := q.All(ctx)
	if err != nil {
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return users, nil
}

func (r *userRepo) GetUserByID(ctx context.Context, id uint64) (*ent.User, error) {
	u, err := r.data.db.User.Query().Where(user.ID(id)).WithRole().Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.DBErrUserNotFound
		}
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return u, nil
}

func (r *userRepo) Count(ctx context.Context) (int, error) {
	count, err := r.data.db.User.Query().Count(ctx)
	if err != nil {
		return 0, errs.DBErrEntError.Wrap(err)
	}
	return count, nil
}
