package getserverrooms

import (
	"context"
	"maps"
	"slices"

	"github.com/Arheon/markus-backend/internal/server/domain/repository"
)

type Query struct {
	roomRepo repository.RoomRepository
}

func NewQuery(repo repository.RoomRepository) *Query {
	return &Query{
		roomRepo: repo,
	}
}

func (q *Query) Handle(ctx context.Context, serverID string) (*Result, error) {
	rooms, err := q.roomRepo.GetServerRooms(ctx, serverID)
	if err != nil {
		return nil, err
	}

	var roomsWithoutCategoryResults []ResultRoom
	var roomCategories map[string]ResultRoomCategory = make(map[string]ResultRoomCategory)

	for _, room := range rooms {
		if room.RoomCategoryID == nil {
			roomsWithoutCategoryResults = append(roomsWithoutCategoryResults, ResultRoom{
				RoomID:   room.ID,
				RoomName: room.RoomName,
			})

			continue
		}

		category, exists := roomCategories[*room.RoomCategoryID]
		if !exists {
			category = ResultRoomCategory{
				CategoryID:   *room.RoomCategoryID,
				CategoryName: room.RoomCategory.Name,
				Rooms:        []ResultRoom{},
			}
		}

		category.Rooms = append(category.Rooms, ResultRoom{
			RoomID:   room.ID,
			RoomName: room.RoomName,
		})

		roomCategories[category.CategoryID] = category
	}

	return &Result{
		ServerID:             serverID,
		RoomCategories:       slices.Collect(maps.Values(roomCategories)),
		RoomsWithoutCategory: roomsWithoutCategoryResults,
	}, nil
}
