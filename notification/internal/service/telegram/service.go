package telegram

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"go.uber.org/zap"

	"github.com/nimbodex/microservices-factory/notification/internal/client/http/telegram"
	"github.com/nimbodex/microservices-factory/notification/internal/model"
)

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type service struct {
	client    *telegram.Client
	logger    Logger
	chatID    string
	templates map[string]*template.Template
}

func NewService(client *telegram.Client, logger Logger, chatID string) *service {
	templates := make(map[string]*template.Template)

	// Шаблон для уведомления об оплате
	paidTemplate := template.Must(template.New("paid").Parse(`
🚀 <b>Заказ оплачен!</b>

📋 Номер заказа: <code>{{.OrderUUID}}</code>
👤 Пользователь: <code>{{.UserUUID}}</code>
💳 Способ оплаты: {{.PaymentMethod}}
🔗 Транзакция: <code>{{.TransactionUUID}}</code>

Ожидайте уведомления о готовности заказа.
	`))

	// Шаблон для уведомления о сборке
	assembledTemplate := template.Must(template.New("assembled").Parse(`
✅ <b>Заказ собран!</b>

📋 Номер заказа: <code>{{.OrderUUID}}</code>
👤 Пользователь: <code>{{.UserUUID}}</code>
⏱️ Время сборки: {{.BuildTimeSec}} сек.

Ваш заказ готов к отправке!
	`))

	templates["paid"] = paidTemplate
	templates["assembled"] = assembledTemplate

	return &service{
		client:    client,
		logger:    logger,
		chatID:    chatID,
		templates: templates,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, event *model.OrderPaidEvent) error {
	tmpl, exists := s.templates["paid"]
	if !exists {
		return fmt.Errorf("template 'paid' not found")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, event); err != nil {
		s.logger.Error(ctx, "Failed to execute paid template", zap.Error(err))
		return err
	}

	if err := s.client.SendMessage(ctx, s.chatID, buf.String()); err != nil {
		s.logger.Error(ctx, "Failed to send paid notification", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Order paid notification sent",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
	)

	return nil
}

func (s *service) SendShipAssembledNotification(ctx context.Context, event *model.ShipAssembledEvent) error {
	tmpl, exists := s.templates["assembled"]
	if !exists {
		return fmt.Errorf("template 'assembled' not found")
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, event); err != nil {
		s.logger.Error(ctx, "Failed to execute assembled template", zap.Error(err))
		return err
	}

	if err := s.client.SendMessage(ctx, s.chatID, buf.String()); err != nil {
		s.logger.Error(ctx, "Failed to send assembled notification", zap.Error(err))
		return err
	}

	s.logger.Info(ctx, "Ship assembled notification sent",
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	return nil
}
