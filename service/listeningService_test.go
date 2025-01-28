package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Real-Dev-Squad/discord-service/config"
	"github.com/Real-Dev-Squad/discord-service/dtos"
	"github.com/Real-Dev-Squad/discord-service/queue"
	"github.com/Real-Dev-Squad/discord-service/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestListeningService(t *testing.T) {
	originalSendMessage := queue.SendMessage
	defer func() {
		queue.SendMessage = originalSendMessage
	}()
	config.AppConfig.MAX_RETRIES = 1
	options := &discordgo.ApplicationCommandInteractionDataOption{
		Value: true,
	}

	mockData := &dtos.Data{
		GuildId: "876543210987654321",
		ApplicationCommandInteractionData: discordgo.ApplicationCommandInteractionData{
			Name: "listening",
			Options: []*discordgo.ApplicationCommandInteractionDataOption{
				options,
			},
		},
	}
	t.Run("should return 'You are already set to listen.' if nickname contains suffix and value is true", func(t *testing.T) {
		data := dtos.DataPacket{
			UserID: "userID",
			MetaData: map[string]string{
				"nickname": "testNick" + utils.NICKNAME_SUFFIX,
				"value":    "true",
			},
		}
		body, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", "/listening", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		discordMessage := &dtos.DiscordMessage{
			Data: mockData,
			Member: &discordgo.Member{
				Nick: fmt.Sprintf("joy-gupta-1%s", utils.NICKNAME_SUFFIX),
				User: &discordgo.User{
					ID: "1",
				},
			},
		}

		commandService := &CommandService{discordMessage: discordMessage}
		commandService.ListeningService(rr, req)

		assert.Contains(t, rr.Body.String(), "You are already set to listen.")
	})

	t.Run("should return 'Your nickname remains unchanged.' if nickname contains suffix and value is true", func(t *testing.T) {
		data := dtos.DataPacket{
			UserID: "userID",
			MetaData: map[string]string{
				"nickname": "testNick",
				"value":    "false",
			},
		}
		body, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", "/listening", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		options.Value = false
		discordMessage := &dtos.DiscordMessage{
			Data: mockData,
			Member: &discordgo.Member{
				Nick: fmt.Sprintf("joy-gupta-1"),
				User: &discordgo.User{
					ID: "1",
				},
			},
		}

		commandService := &CommandService{discordMessage: discordMessage}
		commandService.ListeningService(rr, req)

		fmt.Println("sssssssss", rr.Body.String())
		assert.Contains(t, rr.Body.String(), "Your nickname remains unchanged.")
	})

	t.Run("should pass if nickname does not contain suffix and value is true", func(t *testing.T) {
		originalFunc := queue.SendMessage
		defer func() { queue.SendMessage = originalFunc }()
		queue.SendMessage = func(message []byte) error {
			return nil
		}
		data := dtos.DataPacket{
			UserID: "userID",
			MetaData: map[string]string{
				"nickname": "testNick",
				"value":    "true",
			},
		}
		body, _ := json.Marshal(data)
		req, _ := http.NewRequest("POST", "/listening", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		options.Value = true
		discordMessage := &dtos.DiscordMessage{
			Data: mockData,
			Member: &discordgo.Member{
				Nick: fmt.Sprintf("joy-gupta-1"),
				User: &discordgo.User{
					ID: "1",
				},
			},
		}

		commandService := &CommandService{discordMessage: discordMessage}
		commandService.ListeningService(rr, req)

		fmt.Println("sssssssss", rr.Body.String())
		assert.Contains(t, rr.Body.String(), "Your nickname will be updated shortly.")
	})

	// t.Run("should return 'Your nickname remains unchanged.' if nickname does not contain suffix and value is false", func(t *testing.T) {
	// 	data := dtos.DataPacket{
	// 		UserID: "userID",
	// 		MetaData: map[string]string{
	// 			"nickname": "testNick",
	// 			"value":    "false",
	// 		},
	// 	}
	// 	body, _ := json.Marshal(data)
	// 	req, _ := http.NewRequest("POST", "/listening", bytes.NewBuffer(body))
	// 	rr := httptest.NewRecorder()

	// 	options := discordgo.ApplicationCommandInteractionDataOption{
	// 		Value: false,
	// 	}
	// 	discordMessage := &discordgo.InteractionCreate{
	// 		Interaction: &discordgo.Interaction{
	// 			Data: discordgo.ApplicationCommandInteractionData{
	// 				Options: []*discordgo.ApplicationCommandInteractionDataOption{&options},
	// 			},
	// 			Member: &discordgo.Member{
	// 				Nick: "testNick",
	// 				User: &discordgo.User{ID: "userID"},
	// 			},
	// 		},
	// 	}

	// 	commandService := &CommandService{discordMessage: discordMessage}
	// 	commandService.ListeningService(rr, req)

	// 	assert.Equal(t, http.StatusOK, rr.Code)
	// 	assert.Contains(t, rr.Body.String(), "Your nickname remains unchanged.")
	// })

	// t.Run("should return 'Your nickname will be updated shortly.' and send message to queue if requires update", func(t *testing.T) {
	// 	data := dtos.DataPacket{
	// 		UserID: "userID",
	// 		MetaData: map[string]string{
	// 			"nickname": "testNick",
	// 			"value":    "true",
	// 		},
	// 	}
	// 	body, _ := json.Marshal(data)
	// 	req, _ := http.NewRequest("POST", "/listening", bytes.NewBuffer(body))
	// 	rr := httptest.NewRecorder()

	// 	options := discordgo.ApplicationCommandInteractionDataOption{
	// 		Value: true,
	// 	}
	// 	discordMessage := &discordgo.InteractionCreate{
	// 		Interaction: &discordgo.Interaction{
	// 			Data: discordgo.ApplicationCommandInteractionData{
	// 				Options: []*discordgo.ApplicationCommandInteractionDataOption{&options},
	// 			},
	// 			Member: &discordgo.Member{
	// 				Nick: "testNick",
	// 				User: &discordgo.User{ID: "userID"},
	// 			},
	// 		},
	// 	}

	// 	commandService := &CommandService{discordMessage: discordMessage}
	// 	mockQueue.On("SendMessage", mock.Anything).Return(nil).Once()
	// 	mockSuccess.On("NewDiscordResponse", mock.Anything, "Success", mock.Anything).Once()

	// 	commandService.ListeningService(rr, req)

	// 	assert.Equal(t, http.StatusOK, rr.Code)
	// 	assert.Contains(t, rr.Body.String(), "Your nickname will be updated shortly.")
	// 	mockQueue.AssertExpectations(t)
	// 	mockSuccess.AssertExpectations(t)
	// })

	// t.Run("should return 'Failed to update your nickname.' if sending message to queue fails", func(t *testing.T) {
	// 	data := dtos.DataPacket{
	// 		UserID: "userID",
	// 		MetaData: map[string]string{
	// 			"nickname": "testNick",
	// 			"value":    "true",
	// 		},
	// 	}
	// 	body, _ := json.Marshal(data)
	// 	req, _ := http.NewRequest("POST", "/listening", bytes.NewBuffer(body))
	// 	rr := httptest.NewRecorder()

	// 	options := discordgo.ApplicationCommandInteractionDataOption{
	// 		Value: true,
	// 	}
	// 	discordMessage := &discordgo.InteractionCreate{
	// 		Interaction: &discordgo.Interaction{
	// 			Data: discordgo.ApplicationCommandInteractionData{
	// 				Options: []*discordgo.ApplicationCommandInteractionDataOption{&options},
	// 			},
	// 			Member: &discordgo.Member{
	// 				Nick: "testNick",
	// 				User: &discordgo.User{ID: "userID"},
	// 			},
	// 		},
	// 	}

	// 	commandService := &CommandService{discordMessage: discordMessage}
	// 	mockQueue.On("SendMessage", mock.Anything).Return(errors.New("failed to send message")).Once()

	// 	commandService.ListeningService(rr, req)

	// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	// 	assert.Contains(t, rr.Body.String(), "Failed to update your nickname.")
	// 	mockQueue.AssertExpectations(t)
	// })
}
