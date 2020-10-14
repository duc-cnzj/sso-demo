package cmd

import (
	"fmt"
	"sso/app/middlewares/jwt"
	"sso/repositories/user_repository"
	"sso/server"

	"github.com/spf13/cobra"
)

var id uint

var genJwtCmd = &cobra.Command{
	Use:   "genJwt",
	Short: "生成 jwt token.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err error
			s   = &server.Server{}
		)
		s.SetRunningInConsole()
		if err = s.Init(envPath, ""); err != nil {
			return
		}
		env := s.Env()
		userRepo := user_repository.NewUserRepository(env)
		u, _ := userRepo.FindById(id)
		token, _ := jwt.GenerateToken(u, env)
		fmt.Println(token)
		//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7ImlkIjoxLCJ1c2VyX25hbWUiOiJhYmMiLCJlbWFpbCI6IjFAcXEuY29tIiwibGFzdF9sb2dpbl9hdCI6bnVsbCwiY3JlYXRlZF9hdCI6IjAwMDEtMDEtMDFUMDA6MDA6MDBaIiwidXBkYXRlZF9hdCI6IjIwMjAtMDgtMjZUMTA6MzA6NTYrMDg6MDAiLCJkZWxldGVkX2F0IjpudWxsLCJwZXJtaXNzaW9ucyI6bnVsbCwicm9sZXMiOm51bGx9LCJleHAiOjE1MDAwLCJqdGkiOiIxIn0.hQwhEM6Rc7MUfUZzGRH7DijtNCupfzbgw-IIfs3NFTM
	},
}

func init() {
	rootCmd.AddCommand(genJwtCmd)
	genJwtCmd.Flags().UintVar(&id, "id", 0, "--id=1")
}
