package internal

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	pb "userservice/server/api"
	"userservice/server/pkg/ent"
	"userservice/server/pkg/ent/user"
	"userservice/server/util"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	database *ent.Client
)

type server struct {
	pb.UnimplementedUserServiceServer
	db     *ent.Client
	logger *zap.Logger
}

func NewLogger(lv zapcore.Level, pretty bool) (*zap.Logger, error) {
	c := zap.NewDevelopmentConfig()
	var opts []zap.Option
	if pretty {
		opts = append(opts, zap.AddStacktrace(zap.ErrorLevel))
	}
	level := zap.NewAtomicLevel()

	if err := level.UnmarshalText([]byte(lv.String())); err != nil {
		return nil, fmt.Errorf("could not parse log level %s", lv.String())
	}
	c.Level = level
	return c.Build(opts...)
}

func (s *server) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterReply, error) {
	already, err := s.db.User.Query().Where(user.Username(in.Username)).Exist(ctx)
	if err != nil {
		s.logger.Sugar().Debug("Fail to get username:", err)
	}
	if already {
		return nil, fmt.Errorf("username is unavailable")
	} else {
		pwd, err := util.HashPwd(in.Password)
		if err != nil {
			s.logger.Sugar().Debug("Fail to hash password:", err)
		}
		res, err := s.db.User.Create().SetUsername(in.Username).SetPassword(pwd).SetName(in.Name).Save(ctx)
		if err != nil {
			s.logger.Sugar().Debug("Fail to register:", err)
			return nil, err
		}
		token, err := util.CreateToken(res.Username)
		if err != nil {
			s.logger.Sugar().Debug("Fail to create token:", err)
			return nil, err
		}
		return &pb.RegisterReply{
			UserId: strconv.Itoa(res.ID),
			Token:  token,
		}, nil
	}
}
func (s *server) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginReply, error) {
	user, err := s.db.User.Query().Where(user.Username(in.Username)).Only(ctx)
	if err != nil {
		return nil, err
	}
	err = util.CheckPwdHash(user.Password, in.Password)
	if err != nil {
		return nil, err
	} else {
		token, err := util.CreateToken(user.Username)
		if err != nil {
			return nil, err
		}
		return &pb.LoginReply{
			UserId: strconv.Itoa(user.ID),
			Token:  token,
		}, nil
	}
}
func (s *server) UpdateProfile(ctx context.Context, in *pb.UpdateProfileRequest) (*pb.UpdateProfileReply, error) {
	return nil, nil
}

func ConnectDB() {
	godotenv.Load()
	DB_HOST := os.Getenv("DB_HOST")
	// DB_OUT_PORT := os.Getenv("DB_OUT_PORT")
	DB_IN_PORT := os.Getenv("DB_IN_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_NAME := os.Getenv("DB_NAME")

	db, err := ent.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=True", DB_USER, DB_PASS, DB_HOST, DB_IN_PORT, DB_NAME), ent.Debug())
	if err != nil {
		log.Fatal("Fail to connect database: ", err)
	}
	if err = db.Schema.Create(context.Background()); err != nil {
		log.Fatal("Fail to create schema resources: ", err)
	}
	database = db
}
func RunGRPC() {
	logger, err := NewLogger(zap.DebugLevel, true)
	if err != nil {
		logger.Sugar().Debug("Fail to initialize log:", err)
	}
	defer logger.Core().Sync()

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		log.Fatal("Fail to listen:", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{
		db:     database,
		logger: logger,
	})
	log.Println("grpc server listening at ", lis.Addr())
	go func() {
		log.Fatal("Fail to serve:", s.Serve(lis))
	}()
	// defer database.Close(): sql database is closed
}
