package notificationsapi

import (
	"log/slog"

	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"google.golang.org/protobuf/types/known/structpb"
)

type service struct {
	notificationsapiv1.UnimplementedNotificationsServer
	notifiee notifications.Notifiee
	logger   *slog.Logger
}

func NewService(notifiee notifications.Notifiee, logger *slog.Logger) notificationsapiv1.NotificationsServer {
	return &service{
		notifiee: notifiee,
		logger:   logger,
	}
}

func (s *service) Subscribe(req *notificationsapiv1.SubscribeRequest, stream notificationsapiv1.Notifications_SubscribeServer) error {
	notificationChan := s.notifiee.Subscribe(req.Topics...)
	defer func() {
		<-s.notifiee.Unsubscribe(notificationChan)
	}()
	s.logger.Info("subscribed to notifications", "topics", req.Topics)
	for {
		select {
		case notification := <-notificationChan:
			s.logger.Info("received notification", "notification", notification)
			val, err := structpb.NewStruct(notification.Value)
			if err != nil {
				s.logger.Error("failed to convert notification value to structpb.Value", "error", err, "notification", notification)
				continue
			}
			err = stream.Send(&notificationsapiv1.Notification{
				Topic: notification.Topic,
				Value: val,
			})
			if err != nil {
				s.logger.Error("failed to send notification", "error", err, "notification", notification)
				return err
			}
			s.logger.Info("sent notification", "notification", notification)
		case <-stream.Context().Done():
			return nil
		}
	}
}
