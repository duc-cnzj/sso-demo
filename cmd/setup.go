package cmd

import (
	"sso/app/controllers/api"
	"sso/app/models"
	"sso/config/env"
	"sso/repositories/user_repository"
	"sso/server"

	"github.com/spf13/cobra"
)

var (
	repos     *api.AllRepo
	adminUser = struct {
		name  string
		pwd   string
		email string
		roles []string
	}{
		name:  "duc",
		pwd:   "12345",
		email: "1025434218@qq.com",
		roles: []string{"sso"},
	}

	roles = []*models.Role{{
		Text: "sso",
		Name: "sso",
	}}

	permissions = []struct {
		roleName    string
		permissions []*models.Permission
	}{
		{
			roleName: "sso",
			permissions: []*models.Permission{
				{
					Text:    "sso login",
					Name:    "login",
					Project: "sso",
				},
				{
					Text:    "sso create user",
					Name:    "user_create",
					Project: "sso",
				},
				{
					Text:    "sso create user",
					Name:    "user_edit",
					Project: "sso",
				},
				{
					Text:    "sso create user",
					Name:    "user_delete",
					Project: "sso",
				},
			}},
	}
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "初始化用户和权限",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			s   = &server.Server{}
		)
		s.SetRunningInConsole()
		if err = s.Init(envPath, ""); err != nil {
			return
		}

		repos = api.NewAllRepo(s.Env())
		user := newAdmin(s.Env())
		newRoles(s.Env())
		newPermissions(s.Env())
		for _, role := range adminUser.roles {
			r, _ := repos.RoleRepo.FindByName(role)
			if r != nil {
				repos.UserRepo.SyncRoles(user, []uint{r.ID})
			}
		}
	},
}

func newAdmin(env *env.Env) *models.User {
	repo := user_repository.NewUserRepository(env)
	pwd, _ := repo.GeneratePwd(adminUser.pwd)
	admin := &models.User{
		UserName: adminUser.name,
		Email:    adminUser.email,
		Password: pwd,
	}
	u, _ := repos.UserRepo.FindByEmail(adminUser.email)
	if u != nil {
		u.Password = pwd
		u.UserName = adminUser.name
		env.GetDB().Model(&models.User{}).Update(u)
		return u
	} else {
		env.GetDB().Create(admin)
		return admin
	}
}

func newRoles(env *env.Env) []*models.Role {
	var results []*models.Role
	for _, r := range roles {
		role, _ := repos.RoleRepo.FindByName(r.Name)
		if role != nil {
			role.Name = r.Name
			role.Text = r.Text
			env.GetDB().Model(&models.Role{}).Update(role)
			results = append(results, role)
		} else {
			env.GetDB().Create(r)
			results = append(results, r)
		}
	}
	return results
}

func newPermissions(env *env.Env) []*models.Permission {
	var res []*models.Permission
	for _, items := range permissions {
		var result []*models.Permission
		var pids []uint
		role, _ := repos.RoleRepo.FindByName(items.roleName)
		for _, p := range items.permissions {
			perm, _ := repos.PermRepo.FindByName(p.Name)
			if perm != nil {
				perm.Text = p.Text
				env.GetDB().Model(&models.Permission{}).Update(perm)
				result = append(result, perm)
				pids = append(pids, perm.ID)
			} else {
				env.GetDB().Create(p)
				result = append(result, p)
				pids = append(pids, p.ID)
			}
		}
		if role != nil {
			repos.RoleRepo.SyncPermissions(role, pids, nil)
		}
		res = append(res, result...)
	}

	return res
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
