package grpc

import (
	pb "github.com/fishus/go-advanced-gophkeeper/internal/adapter/handler/proto"
	"github.com/fishus/go-advanced-gophkeeper/internal/core/port"
)

type server struct {
	pb.UnimplementedVaultServer
	userService port.UserService
	authService port.AuthService
}

func NewServer(userService port.UserService, authService port.AuthService) *server {
	server := server{
		userService: userService,
		authService: authService,
	}

	return &server
}

//func (s *VaultServer) AddVault(ctx context.Context, in *pb.AddVaultRequest) (*pb.AddVaultResponse, error) {
//	var response pb.AddVaultResponse
//
//	fmt.Printf("User: %#v\n", in.User.Email)
//
//	for _, data := range in.Storage.Data {
//		fmt.Printf("ID: %#v\n", data.Id)
//		fmt.Printf("Data: %#v\n", data.Data)
//		if data.GetCreds() != nil {
//			fmt.Println("Creds:")
//			fmt.Printf("Login: %#v\n", data.GetCreds().Login)
//			fmt.Printf("Password: %#v\n", data.GetCreds().Password)
//		}
//		if data.GetCard() != nil {
//			fmt.Println("Card:")
//			fmt.Printf("Name: %#v\n", data.GetCard().Name)
//			fmt.Printf("Number: %#v\n", data.GetCard().Number)
//			fmt.Printf("Code: %#v\n", data.GetCard().Code)
//		}
//		fmt.Println("------------")
//	}
//
//fmt.Printf("CreatedAt: %#v\n", in.User.CreatedAt.AsTime().Format(time.RFC3339))

//	return &response, nil
//}
