package service

import (
	"os"
	"testing"
	"time"
)

func TestNewOrderController(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	if oc.nextOrderID != 1 {
		t.Errorf("Expected nextOrderID to be 1, got %d", oc.nextOrderID)
	}
	if oc.nextBotID != 1 {
		t.Errorf("Expected nextBotID to be 1, got %d", oc.nextBotID)
	}
}

func TestAddNormalOrder(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddNormalOrder()

	if len(oc.orders) != 1 {
		t.Errorf("Expected 1 order, got %d", len(oc.orders))
	}
	if oc.orders[0].Type != NormalOrder {
		t.Errorf("Expected NormalOrder type, got %v", oc.orders[0].Type)
	}
	if oc.orders[0].ID != 1 {
		t.Errorf("Expected order ID 1, got %d", oc.orders[0].ID)
	}
}

func TestAddVIPOrder(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddNormalOrder()
	oc.AddVIPOrder()

	if len(oc.orders) != 2 {
		t.Errorf("Expected 2 orders, got %d", len(oc.orders))
	}
	if oc.orders[0].Type != VIPOrder {
		t.Errorf("Expected first order to be VIP, got %v", oc.orders[0].Type)
	}
	if oc.orders[1].Type != NormalOrder {
		t.Errorf("Expected second order to be Normal, got %v", oc.orders[1].Type)
	}
}

func TestVIPOrderPriority(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddNormalOrder()
	oc.AddNormalOrder()
	oc.AddVIPOrder()
	oc.AddVIPOrder()

	expectedOrder := []OrderType{VIPOrder, VIPOrder, NormalOrder, NormalOrder}
	for i, expected := range expectedOrder {
		if oc.orders[i].Type != expected {
			t.Errorf("Order at position %d: expected %v, got %v", i, expected, oc.orders[i].Type)
		}
	}
}

func TestAddBot(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddBot()

	if len(oc.bots) != 1 {
		t.Errorf("Expected 1 bot, got %d", len(oc.bots))
	}
	if oc.bots[0].ID != 1 {
		t.Errorf("Expected bot ID 1, got %d", oc.bots[0].ID)
	}
}

func TestRemoveBot(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddBot()
	oc.AddBot()
	oc.RemoveBot()

	if len(oc.bots) != 1 {
		t.Errorf("Expected 1 bot after removal, got %d", len(oc.bots))
	}
}

func TestRemoveBotWithProcessingOrder(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddNormalOrder()
	oc.AddBot()

	time.Sleep(500 * time.Millisecond)

	oc.RemoveBot()

	oc.mu.RLock()
	ordersInQueue := len(oc.orders)
	oc.mu.RUnlock()

	if ordersInQueue != 1 {
		t.Errorf("Expected order to be returned to queue, got %d orders", ordersInQueue)
	}

	time.Sleep(100 * time.Millisecond)
}

func TestOrderProcessing(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddNormalOrder()
	oc.AddBot()

	time.Sleep(11 * time.Second)

	oc.mu.RLock()
	ordersRemaining := len(oc.orders)
	oc.mu.RUnlock()

	if ordersRemaining != 0 {
		t.Errorf("Expected order to be processed, but %d orders remain", ordersRemaining)
	}
}

func TestMultipleBotsProcessing(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddNormalOrder()
	oc.AddNormalOrder()
	oc.AddBot()
	oc.AddBot()

	time.Sleep(11 * time.Second)

	if len(oc.orders) != 0 {
		t.Errorf("Expected all orders to be processed, but %d orders remain", len(oc.orders))
	}
}

func TestUniqueOrderNumbers(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	oc.AddNormalOrder()
	oc.AddVIPOrder()
	oc.AddNormalOrder()

	ids := make(map[int]bool)
	for _, order := range oc.orders {
		if ids[order.ID] {
			t.Errorf("Duplicate order ID found: %d", order.ID)
		}
		ids[order.ID] = true
	}
}

func TestGetStatus(t *testing.T) {
	oc, err := NewOrderController("/tmp/test_result.txt")
	if err != nil {
		t.Fatalf("Failed to create OrderController: %v", err)
	}
	defer oc.Close()
	defer os.Remove("/tmp/test_result.txt")

	status := oc.GetStatus()
	expected := "status: bot: [0/0], order: []"
	if status != expected {
		t.Errorf("Expected status '%s', got '%s'", expected, status)
	}
}
