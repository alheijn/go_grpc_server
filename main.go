package main

import (
	"context"
	_ "fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	pb "go_grpc_server/ecommerce/ordermanagement" // Adjust the import path
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement the OrderManagementServer interface.
type server struct {
	// This is a requirement for forward compatibility.
	pb.UnimplementedOrderManagementServer
}

// CreateOrder implements the CreateOrder RPC method.
func (s *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("Received order creation request for customer: %s", req.CustomerId)

	// --- Simple Data Manipulation ---
	// 1. Generate a unique order ID.
	orderID := uuid.New().String()

	// 2. Calculate the total price.
	var totalPrice float64
	for _, item := range req.Items {
		totalPrice += item.PricePerUnit * float64(item.Quantity)
	}

	// 3. Create the Order object to be returned.
	newOrder := &pb.Order{
		OrderId:            orderID,
		CustomerId:         req.CustomerId,
		Items:              req.Items,
		ShippingAddress:    req.GetShippingAddress(),
		TotalPrice:         totalPrice,
		Status:             pb.Status_PENDING,
		CreatedAtTimestamp: time.Now().Unix(),
	}

	log.Printf("Order created successfully with ID: %s and Total Price: %.2f", newOrder.OrderId, newOrder.TotalPrice)
	return newOrder, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderManagementServer(s, &server{})

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
