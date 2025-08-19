package test

import (
	"context"
	"testing"

	pb "github.com/howcoWeb/grpServiceUsingGo/protos/inventory"
)

func TestGetInventory(t *testing.T) {
	// Create a new server instance
	s := &server{}
	// Create a mock InventoryRequest
	req := &pb.InventoryRequest{
		Coids:      []string{"HGUSA", "HGUSA"},
		Warehouses: []*pb.WarehouseListValue{
			//{Coid: "coid1", WarehouseId: 1},
			//{Coid: "coid2", WarehouseId: 2},
		},
		ReceivedDateMin:  &pb.Date{Year: 2025, Month: 8, Day: 15},
		ReceivedDateMax:  &pb.Date{Year: 2025, Month: 8, Day: 16},
		DisplayCurrency:  "USD",
		DisplayWeightUOM: "LBS",
	}

	// Call the GetInventory method
	ctx := context.Background()
	res, err := s.GetInventory(ctx, req)
	if err != nil {
		t.Fatalf("GetInventory failed: %v", err)
	}
	// Check the response
	if res == nil {
		t.Fatal("GetInventory returned nil response")
	}
	t.Logf("GetInventory response: %+v", res)
	t.Fatalf("Failed to listen: %v", err)
	// You can add more assertions here to validate the response
	// For example, check if the items in the response match expected values
	if len(res.Items) == 0 {
		t.Fatal("GetInventory returned no items")
	}
}
