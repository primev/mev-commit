package notificationsapi

import (
	"log/slog"

	notificationsapiv1 "github.com/primev/mev-commit/p2p/gen/go/notificationsapi/v1"
	"github.com/primev/mev-commit/p2p/pkg/notifications"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

type Service struct {
	notificationsapiv1.UnimplementedNotificationsServer
	notifiee notifications.Notifiee
	logger   *slog.Logger
}

func NewService(notifiee notifications.Notifiee, logger *slog.Logger) *Service {
	return &Service{
		notifiee: notifiee,
		logger:   logger,
	}
}

func (s *Service) Subscribe(
	req *notificationsapiv1.SubscribeRequest,
	stream notificationsapiv1.Notifications_SubscribeServer,
) error {
	var topics []notifications.Topic
	for _, topic := range req.Topics {
		if notifications.IsTopicValid(notifications.Topic(topic)) {
			topics = append(topics, notifications.Topic(topic))
		} else {
			return status.Errorf(codes.InvalidArgument, "invalid topic: %s", topic)
		}
	}
	notificationChan := s.notifiee.Subscribe(topics...)
	defer func() {
		<-s.notifiee.Unsubscribe(notificationChan)
	}()
	for {
		select {
		case notification := <-notificationChan:
			val, err := structpb.NewStruct(notification.Value())
			if err != nil {
				s.logger.Error("failed to convert notification value to structpb.Value", "error", err, "notification", notification)
				continue
			}
			err = stream.Send(&notificationsapiv1.Notification{
				Topic: string(notification.Topic()),
				Value: val,
			})
			if err != nil {
				s.logger.Error("failed to send notification", "error", err, "notification", notification)
				return err
			}
			s.logger.Debug("sent notification", "notification", notification)
		case <-stream.Context().Done():
			return nil
		}
	}
}
