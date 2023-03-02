package user

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func TestNewStore(t *testing.T) {
	db, Close := createMemDB(t)
	defer Close()
	NewStore(db)
}

func Test_store_Create(t *testing.T) {
	db, c := createMemDB(t)
	defer c()

	t.Run("add user@example.com", func(t *testing.T) {
		store := NewStore(db)
		err := store.Create(&User{Email: "user@example.com"})
		if err != nil {
			t.Fatal(err)
		}

		createdUser := &User{}
		store.DB().Find(createdUser)

		if createdUser.UUID == "" {
			t.Fatal("the created user uuid should not be empty")
		}
	})
}

func Test_store_User(t *testing.T) {
	db, c := createMemDB(t)
	defer c()

	userExample := User{
		Email: "user@example.com",
	}

	user2Example := User{
		Email: "user2@example.com",
	}

	s := NewStore(db)
	err := s.Create(&user2Example)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		u User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *User
		wantErr error
	}{
		{
			name:    "get user@example.com",
			fields:  fields{db: db},
			args:    args{u: userExample},
			want:    &userExample,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "get user2@example.com",
			fields:  fields{db: db},
			args:    args{u: user2Example},
			want:    &user2Example,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store{
				db: tt.fields.db,
			}

			got, err := s.User(&tt.args.u)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("User() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if err == nil && got.Email != tt.want.Email {
				t.Fatalf("User() got Email = %v, want Email %v", got.Email, tt.want.Email)
			}
		})
	}
}

func Test_store_UserByEmail(t *testing.T) {
	db, c := createMemDB(t)
	defer c()

	userExample := User{
		Email: "user@example.com",
	}

	user2Example := User{
		Email: "user2@example.com",
	}

	s := NewStore(db)
	err := s.Create(&user2Example)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    User
		wantErr error
	}{
		{
			name:    "get user@example.com",
			fields:  fields{db: db},
			args:    args{email: "user@example.com"},
			want:    userExample,
			wantErr: gorm.ErrRecordNotFound,
		},
		{
			name:    "get admin@example.com",
			fields:  fields{db: db},
			args:    args{email: "user@example.com"},
			want:    userExample,
			wantErr: ErrUserNotFound,
		},
		{
			name:    "get user2@example.com",
			fields:  fields{db: db},
			args:    args{email: "user2@example.com"},
			want:    user2Example,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store{
				db: tt.fields.db,
			}

			got, err := s.UserByEmail(tt.args.email)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UserByEmail() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

			}

			if err == nil && got.Email != tt.want.Email {
				t.Fatalf("UserByEmail() got Email = %v, want Email %v", got.Email, tt.want.Email)
			}
		})
	}
}

func Test_store_UserByID(t *testing.T) {
	db, c := createMemDB(t)
	defer c()

	userExample := User{
		Email: "user@example.com",
	}

	user2Example := User{
		Email: "user2@example.com",
	}

	s := NewStore(db)
	err := s.Create(&userExample)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    User
		wantErr error
	}{
		{
			name:    "get id 1",
			fields:  fields{db: db},
			args:    args{1},
			want:    userExample,
			wantErr: nil,
		},
		{
			name:    "get id 2",
			fields:  fields{db: db},
			args:    args{2},
			want:    user2Example,
			wantErr: ErrUserNotFound,
		},
		{
			name:    "get id 2",
			fields:  fields{db: db},
			args:    args{2},
			want:    user2Example,
			wantErr: gorm.ErrRecordNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store{
				db: tt.fields.db,
			}

			got, err := s.UserByID(tt.args.id)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UserByEmail() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

			}

			if err == nil && got.Email != tt.want.Email {
				t.Fatalf("UserByEmail() got Email = %v, want Email %v", got.Email, tt.want.Email)
			}
		})
	}
}

func Test_store_UserByName(t *testing.T) {
	db, c := createMemDB(t)
	defer c()

	userExample := User{
		FirstName: "Foo",
		LastName:  "Bar",
		Email:     "user@example.com",
	}

	user2Example := User{
		FirstName: "Foo",
		LastName:  "Baz",
		Email:     "user2@example.com",
	}

	s := NewStore(db)
	err := s.Create(&userExample)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Create(&user2Example)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCount int
		wantErr   error
	}{
		{
			name:      "get name foo",
			fields:    fields{db: db},
			args:      args{"foo"},
			wantCount: 2,
			wantErr:   nil,
		},
		{
			name:      "get name baz",
			fields:    fields{db: db},
			args:      args{"baz"},
			wantCount: 1,
			wantErr:   nil,
		},
		{
			name:      "get name bar bar",
			fields:    fields{db: db},
			args:      args{"bar bar"},
			wantCount: 0,
			wantErr:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store{
				db: tt.fields.db,
			}

			got, err := s.UserByName(tt.args.name)
			if err != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("UserByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

			}

			if err == nil && len(got) != tt.wantCount {
				t.Fatalf("UserByName() got users = %v, want users %v", len(got), tt.wantCount)
			}
		})
	}
}

func Test_store_SetDB(t *testing.T) {
	db1, c := createMemDB(t)
	defer c()

	db2, c := createMemDB(t)
	defer c()

	db3, c := createMemDB(t)
	defer c()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *gorm.DB
	}{
		{
			name:   "SetDB db1 to db2",
			fields: fields{db: db1},
			args:   args{db: db2},
			want:   db2,
		},
		{
			name:   "SetDB db2 to db3",
			fields: fields{db: db2},
			args:   args{db: db3},
			want:   db3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &store{
				db: tt.fields.db,
			}

			if s.db != tt.fields.db {
				t.Fatalf("contaminated. s.db should be:%p, got: %p", tt.fields.db, s.db)
			}

			s.SetDB(tt.args.db)

			if s.db != tt.want {
				t.Fatalf("SetDB() got db = %p, want db %p", s.db, tt.want)
			}
		})
	}
}

func Test_store_Migrate(t *testing.T) {
	db, c := createMemDB(t)
	defer c()

	s := NewStore(db)

	err := s.Migrate()
	if err != nil {
		t.Errorf("Migrate returned an error: %v", err)
	}
}

func createMemDB(t testing.TB) (*gorm.DB, func()) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Discard,
	})

	if err != nil {
		t.Fatal(err)
	}

	err = db.AutoMigrate(User{})
	if err != nil {
		t.Fatal(err)
	}

	Close := func() {
		d, err := db.DB()
		if err != nil {
			t.Fatal(err)
		}

		err = d.Close()
		if err != nil {
			t.Fatal(err)
		}
	}

	return db, Close
}
