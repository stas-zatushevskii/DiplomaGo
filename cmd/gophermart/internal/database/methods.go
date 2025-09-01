package database

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/constants"
	customErrors "github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/errors"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/models"
	"github.com/stas-zatushevskii/DiplomaGo/cmd/gophermart/internal/utils"
	"gorm.io/gorm"
	"time"
)

func (d *Database) CreateNewUser(user *models.User) error {
	res := d.GormDB.Create(user)
	if res.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(res.Error, &pgErr) && pgErr.Code == "23505" {
			return customErrors.ErrUserAlreadyExists
		}
	}
	return res.Error
}

func (d *Database) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	res := d.GormDB.Where("id = ?", id).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.ErrUserNotFound
		}
	}
	return &user, res.Error
}

func (d *Database) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	res := d.GormDB.Where("username = ?", username).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.ErrUserNotFound
		}
	}
	return &user, res.Error
}

func (d *Database) IncreaseOrderAccrual(orderNumber string, accrual utils.Money) error {
	res := d.GormDB.Model(&models.Order{}).
		Where("order_number = ?", orderNumber).
		UpdateColumn("accrual", gorm.Expr("accrual + ?", accrual))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return customErrors.ErrOrdersNotFound
	}
	return nil
}

func (d *Database) DecreaseOrderAccrual(orderNumber string, withdraw utils.Money) error {
	res := d.GormDB.Model(&models.Order{}).
		Where("order_number = ?", orderNumber).
		UpdateColumn("accrual", gorm.Expr("accrual - ?", withdraw)).
		UpdateColumn("withdrawn_accrual", gorm.Expr("withdrawn_accrual + ?", withdraw))

	return res.Error
}

func (d *Database) GetUserBalance(userID uint) (utils.Money, error) {
	var balance utils.Money

	res := d.GormDB.Model(&models.Order{}).
		Select("COALESCE(SUM(accrual), 0)").
		Where("user_id = ? AND status = ?", userID, constants.OrderStatusProcessed).
		Scan(&balance)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
	}

	return balance, res.Error
}

func (d *Database) GetUserWithDrawnBalance(userID uint) (utils.Money, error) {
	var balance utils.Money

	res := d.GormDB.Model(&models.Order{}).
		Select("COALESCE(SUM(withdrawn_accrual), 0)").
		Where("user_id = ? AND status = ? AND withdrawn_accrual > 0", userID, constants.OrderStatusProcessed).
		Scan(&balance)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
	}

	return balance, res.Error
}

func (d *Database) AddToOrderHistory(order *models.Order, sum utils.Money) error {
	h := models.OrderHistory{
		OrderNumber: order.OrderNumber,
		Sum:         sum,
		OrderID:     &order.ID,
		ProcessedAt: time.Now().Format(time.RFC3339),
	}

	if err := d.GormDB.Create(&h).Error; err != nil {
		return err
	}

	return nil
}

func (d *Database) GetWithdrawalsHistory(userID uint) ([]models.OrderHistory, error) {
	var history []models.OrderHistory

	err := d.GormDB.
		Model(&models.OrderHistory{}).
		Joins("JOIN orders ON orders.id = order_histories.order_id").
		Where("orders.user_id = ?", userID).
		Where("orders.status = ?", constants.OrderStatusProcessed).
		Order("order_histories.processed_at DESC").
		Find(&history).Error
	if err != nil {
		return nil, err
	}
	return history, nil
}

func (d *Database) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	res := d.GormDB.Where("user_id = ?", userID).Find(&orders)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.ErrOrdersNotFound
		}
	}
	return orders, res.Error
}

func (d *Database) GetOrderByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	res := d.GormDB.Preload("User").Where("order_number = ?", orderNumber).First(&order)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.ErrOrdersNotFound
		}
	}
	return &order, res.Error
}

func (d *Database) GetNewOrProcessingOrders() ([]models.Order, error) {
	var orders []models.Order
	res := d.GormDB.
		Where("status IN ?", []constants.OrderStatus{
			constants.OrderStatusNew,
			constants.OrderStatusProcessing,
		}).
		Find(&orders)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return nil, customErrors.ErrOrdersNotFound
		}
	}
	return orders, res.Error
}

func (d *Database) ChangeOrderStatus(orderNumber string, newStatus constants.OrderStatus) error {
	res := d.GormDB.
		Model(&models.Order{}).
		Where("order_number = ?", orderNumber).
		UpdateColumn("status", newStatus)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return customErrors.ErrUserNotFound
	}
	return nil
}

func (d *Database) CreateNewOrder(orderNumber string, userID uint) (*models.Order, error) {

	var existing models.Order
	err := d.GormDB.
		Select("id", "user_id").
		Where("order_number = ?", orderNumber).
		Take(&existing).Error

	switch {
	case err == nil:
		if existing.UserID != nil && *existing.UserID == userID { // order found and userID == requested UserID
			return nil, customErrors.ErrOrderAlreadyUsed
		}
		return nil, customErrors.ErrOrderAlreadyExist // order found and userID != requested UserID
	case errors.Is(err, gorm.ErrRecordNotFound):
		newOrder := models.Order{
			OrderNumber: orderNumber,
			UserID:      &userID,
			Status:      constants.OrderStatusNew,
			CreatedAt:   time.Now().Format(time.RFC3339),
		}

		if err := d.GormDB.Create(&newOrder).Error; err != nil {
			return nil, err
		}
		return &newOrder, nil
	default:
		return nil, err
	}
}
