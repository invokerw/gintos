package data

import (
	"context"
	"github/invokerw/gintos/demo/api/v1/admin"
	"github/invokerw/gintos/demo/api/v1/common"
	"github/invokerw/gintos/demo/internal/biz"
	"github/invokerw/gintos/demo/internal/data/ent"
	"github/invokerw/gintos/demo/internal/data/ent/role"
	"github/invokerw/gintos/demo/internal/errs"
	"github/invokerw/gintos/log"
	"time"
)

type roleRepo struct {
	data *Data
	log  *log.Helper
}

// NewRoleRepo .
func NewRoleRepo(data *Data, logger log.Logger) biz.RoleRepo {
	return &roleRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "data", "role_repo")),
	}
}

func (r *roleRepo) CreateRole(ctx context.Context, in *common.Role) (*ent.Role, error) {
	if in == nil || in.Name == nil {
		return nil, errs.DBErrInvalidParam
	}
	now := time.Now()
	var u *ent.Role
	var err error
	err = WithTx(ctx, r.data.db, func(tx *ent.Tx) error {
		u, err = tx.Role.Query().Where(role.Name(in.GetName())).Only(ctx)
		if u != nil {
			return errs.DBErrRoleExist
		}
		uc := tx.Role.Create().
			SetName(in.GetName()).
			SetNillableDesc(in.Desc).
			SetNillableParentID(in.ParentId).
			SetNillableSortID(in.SortId).
			SetCreateTime(now).
			SetUpdateTime(now)

		u, err = uc.Save(ctx)
		if err != nil {
			return errs.DBErrEntError.Wrap(err)
		}
		return nil
	})
	if err != nil {
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return u, nil
}

func (r *roleRepo) UpdateRoles(ctx context.Context, roles []*common.Role) ([]*ent.Role, error) {
	for _, in := range roles {
		if in == nil || in.Name == nil {
			return nil, errs.DBErrInvalidParam
		}
	}
	ret := make([]*ent.Role, 0, len(roles))
	err := WithTx(ctx, r.data.db, func(tx *ent.Tx) error {
		for _, in := range roles {
			u, err := tx.Role.Query().Where(role.Name(in.GetName())).Only(ctx)
			if err != nil {
				return err
			}
			uc := u.Update().
				SetNillableDesc(in.Desc).
				SetNillableParentID(in.ParentId).
				SetNillableSortID(in.SortId).
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

func (r *roleRepo) DeleteRoles(ctx context.Context, names []string) error {
	var err error
	err = WithTx(ctx, r.data.db, func(tx *ent.Tx) error {
		var u *ent.Role
		for _, name := range names {
			u, err = tx.Role.Query().Where(role.Name(name)).Only(ctx)
			if err != nil {
				return errs.DBErrEntError.Wrap(err)
			}
			err = tx.Role.DeleteOne(u).Exec(ctx)
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

func (r *roleRepo) GetRole(ctx context.Context, name string) (*ent.Role, error) {
	u, err := r.data.db.Role.Query().Where(role.Name(name)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.DBErrRoleNotFound
		}
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return u, nil
}

func (r *roleRepo) GetRoleList(ctx context.Context, req *admin.GetRoleListRequest) ([]*ent.Role, error) {
	if req == nil {
		return nil, errs.DBErrInvalidParam
	}
	q := r.data.db.Role.Query().Offset(int(req.Page.Offset)).Limit(int(req.Page.PageSize))
	name := req.GetName()
	if name != "" {
		q.Where(role.NameContains(name))
	}
	rs, err := q.All(ctx)
	if err != nil {
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return rs, nil
}

func (r *roleRepo) GetRoleByID(ctx context.Context, id uint64) (*ent.Role, error) {
	u, err := r.data.db.Role.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errs.DBErrRoleNotFound
		}
		return nil, errs.DBErrEntError.Wrap(err)
	}
	return u, nil
}

func (r *roleRepo) Count(ctx context.Context) (int, error) {
	count, err := r.data.db.Role.Query().Count(ctx)
	if err != nil {
		return 0, errs.DBErrEntError.Wrap(err)
	}
	return count, nil
}
