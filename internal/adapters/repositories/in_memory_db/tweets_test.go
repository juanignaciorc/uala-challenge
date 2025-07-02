package in_memory_db

import (
	"context"
	"github.com/google/uuid"
	"github.com/juanignaciorc/microbloggin-pltf/internal/domain"
	mock_ports "github.com/juanignaciorc/microbloggin-pltf/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

const uuidMock = "77dae0ef-658c-44c6-803f-f849854a7033"

func TestInMemoryDB_CreateTweet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mock_ports.NewMockInMemoryDBTweetsInterface(ctrl)

	tweet := domain.Tweet{
		ID:      uuid.MustParse(uuidMock),
		UserID:  uuid.MustParse(uuidMock),
		Message: "test tweet",
	}

	mockDB.EXPECT().
		CreateTweet(gomock.Any(), gomock.Eq(tweet)).
		Return(tweet, nil).
		Times(1)

	createdTweet, err := mockDB.CreateTweet(context.Background(), tweet)

	assert.NoError(t, err)
	assert.Equal(t, tweet.ID, createdTweet.ID)
	assert.Equal(t, tweet.Message, createdTweet.Message)

}
