package main

import (
	"context"
	"log"
	"net"

	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
	pb "github.com/howcoWeb/grpServiceUsingGo/protos/inventoryRequest"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedInventoryServiceServer
}

func (s *server) GetInventory(ctx context.Context, req *pb.InventoryRequest) (*pb.InventoryResult, error) {
	log.Printf("Received InventoryRequest: %+v", req)

	connString := "Server=s-us-dw02-dev\\qa;Database=ODSNG;integrated security=SSPI;TrustServerCertificate=True;"
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
