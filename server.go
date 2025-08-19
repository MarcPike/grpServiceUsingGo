package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/howcoWeb/grpServiceUsingGo/protos/inventory"
	"google.golang.org/grpc"
	"database/sqlgithub.com/denisenkom/go-mssqldb"
)	

type server struct {
	pb.UnimplementedInventoryServiceServer
}

func (s *server) GetInventory(ctx context.Context, req *pb.InventoryRequest) (*pb.InventoryResult, error) {
	log.Printf("Received InventoryRequest: %+v", req)		

	connString := "Server=s-us-dw02-dev\qa;Database=ODSNG;integrated security=SSPI;TrustServerCertificate=True;"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		// handle error
	}
	defer db.Close()

	coids := req.GetCoids()
	warehouses := req.GetWarehouses()
	receivedDateMin := req.GetReceivedDateMin()
	receivedDateMax := req.GetReceivedDateMax()
	displayCurrency := req.GetDisplayCurrency()
	displayWeightUOM := req.GetDisplayWeightUOM()

	rows, err := db.Query("EXEC QNG.uspGetInventoryAnalysisNew @COIDs, @DisplayCurrency, @DisplayWeightUOM, @ReceivedDateMin, @ReceivedDateMax, @Warehouses", coids, displayCurrency, displayWeightUOM, receivedDateMin, receivedDateMax, warehouses)
	if err != nil {
		// handle error
	}
	defer rows.Close()


	var items []*pb.InventoryItem
	for rows.Next() {
		var item pb.InventoryItem	

		if err := rows.Scan(&item.Coid, &item.ExRateSummary, &item.Branch, &item.StkItemId, &item.ProductCode, &item.ProductCategory, &item.ProductSize, &item.ProductCondition, &item.Quantity.Units, &item.Quantity.Nanos, &item.UnitPrice.Units, &item.UnitPrice.Nanos, &item.TotalValue.Units, &item.TotalValue.Nanos, &item.InvoiceDate.Year, &item.InvoiceDate.Month, &item.InvoiceDate.Day, &item.FiscalPeriod.Year, &item.FiscalPeriod.Month, &item.FiscalPeriod.Day, &item.StockStatus, &item.StockStatusDate.Year, &item.StockStatusDate.Month, &item.StockStatusDate.Day, &item.StockStatusComment); err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}	

		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("Error reading rows: %v", err)
	}	
	return &pb.InventoryResult{Items: items}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")	
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}	
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, &server{})	
	log.Println("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}	
}
	// // Here you would typically fetch the inventory data based on the request
	// items := []*pb.InventoryItem{
	// 	{		
	// 		Coid:            "coid1",
	// 		ExRateSummary:   "Exchange Rate Summary 1",
	// 		Branch:          "Branch 1",	
	// 		StkItemId:       101,
	// 		ProductCode:     "P001",
	// 		ProductCategory: "Category 1",
	// 		ProductSize:     "Large",	
	// 		ProductCondition: "New",
	// 		Quantity: &pb.DecimalValue{
	// 			Units: 10,
	// 			Nanos: 0,
	// 		},
	// 		UnitPrice: &pb.DecimalValue{
	// 			Units: 100,
	// 			Nanos: 0,
	// 		},
	// 		TotalValue: &pb.DecimalValue{
	// 			Units: 1000,
	// 			Nanos: 0,
	// 		},
	// 		InvoiceDate: &pb.Date{
	// 			Year:  2023,
	// 			Month: 10,
	// 			Day:   1,
	// 		},
	// 		FiscalPeriod: &pb.Date{
	// 			Year:  2023,
	// 			Month: 1,
	// 			Day:   1,
	// 		},
	// 		StockStatus: "In Stock",
	// 		StockStatusDate: &pb.Date{
	// 			Year:  2023,
	// 			Month: 10,
	// 			Day:   1,
	// 		},
	// 		StockStatusComment: "Available for sale",
	// 	}
	// }

	// return &pb.InventoryResult{Items: items}, nil
	
