package service

import (
	"HighArch/entity"
	"HighArch/storage"
)

type FriendLinksService struct {
	friendLinksStore storage.FriendLinksStore
}

func NewFriendLinksService(friendLinksStore storage.FriendLinksStore) *FriendLinksService {
	return &FriendLinksService{
		friendLinksStore: friendLinksStore,
	}
}

func (s *FriendLinksService) SetFriendsLink(friendOneUserId string, friendTwoUserId string) error {
	if len(friendOneUserId) <= 0 || len(friendTwoUserId) <= 0 {
		return ErrorValidation
	}
	err := s.friendLinksStore.SetFriends(entity.FriendsLink{
		Friend1UserId: friendOneUserId,
		Friend2UserId: friendTwoUserId,
	})
	if err != nil {
		return ErrorStoreError
	}
	return nil
}

func (s *FriendLinksService) DeleteFriendsLink(friendOneUserId string, friendTwoUserId string) error {
	if len(friendOneUserId) <= 0 || len(friendTwoUserId) <= 0 {
		return ErrorValidation
	}
	err := s.friendLinksStore.DeleteFriends(entity.FriendsLink{
		Friend1UserId: friendOneUserId,
		Friend2UserId: friendTwoUserId,
	})
	if err != nil {
		return ErrorStoreError
	}
	return nil
}
