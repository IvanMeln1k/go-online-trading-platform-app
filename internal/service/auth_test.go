package service

/*

func TestAuth_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	usersRepo := mock_repository.NewMockUsers(ctrl)
	sessionsRepo := mock_repository.NewMockSessions(ctrl)
	authService := NewAuthService(usersRepo, sessionsRepo)

	type args struct {
		user domain.User
	}

	type mockBehavior func(args args, id int)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         int
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				user: domain.User{
					Username: "username",
					Name:     "name",
					Email:    "email",
					Password: "password",
				},
			},
			mockBehavior: func(args args, id int) {
				args.user.Password = authService.hashPassword(args.user.Password)
				usersRepo.EXPECT().GetByEmail(context.Background(), args.user.Email).
					Return(domain.User{}, repository.ErrUserNotFound)
				usersRepo.EXPECT().GetByUserName(context.Background(), args.user.Username).
					Return(domain.User{}, repository.ErrUserNotFound)
				usersRepo.EXPECT().Create(context.Background(), args.user).
					Return(id, nil)
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "email already in user",
			args: args{
				user: domain.User{
					Username: "username",
					Name:     "name",
					Email:    "email",
					Password: "password",
				},
			},
			mockBehavior: func(args args, id int) {
				usersRepo.EXPECT().GetByEmail(context.Background(), args.user.Email).
					Return(domain.User{}, nil)
			},
			want:    0,
			wantErr: ErrEmailAlreadyInUse,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.args, test.want)

			id, err := authService.SignUp(context.Background(), test.args.user)

			assert.ErrorIs(t, test.wantErr, err)
			assert.Equal(t, test.want, id)
		})
	}
}

func TestAuth_SignIn(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	usersRepo := mock_repository.NewMockUsers(ctrl)
	sessionsRepo := mock_repository.NewMockSessions(ctrl)
	authService := NewAuthService(usersRepo, sessionsRepo)

	type args struct {
		email    string
		password string
	}

	type mockBehavior func(args args, tokens domain.Tokens)

	tests := []struct {
		name         string
		args         args
		mockBehavior mockBehavior
		want         domain.Tokens
		wantErr      error
	}{
		{
			name: "ok",
			args: args{
				email:    "email",
				password: "password",
			},
			mockBehavior: func(args args, tokens domain.Tokens) {
				usersRepo.EXPECT().GetByEmail(context.Background(), args.email).
					Return(domain.User{
						Id:       1,
						Username: "username",
						Name:     "name",
						Email:    "email",
						Password: (&AuthService{}).hashPassword("password"),
					})
				refreshToken, _ := (&AuthService{}).createRefreshToken()
				sessionsRepo.EXPECT().Create(context.Background(),
					(&AuthService{}).createSession(1, refreshToken))
			},
		},
	}
}

*/
