package flow

import (
	"example/bot/telegram-ai-bot/model"
)

func Init() model.Flow {
	return model.Flow{
		"start": model.Usecase{
			"createProfile": model.Chain{
				0: model.Action{
					Handler: model.PromptProfileNameHandler,
					Message: model.Message{
						Text:    "Enter the name of the bot",
						Buttons: []model.Button{},
					},
				},
				1: model.Action{
					Handler: model.PromptInstructionHandler,
					Message: model.Message{
						Text:    "Enter the instruction for the bot",
						Buttons: []model.Button{},
					},
				},
				2: model.Action{
					Handler: model.StoreAIModelHandler,
					Message: model.Message{
						Text: "Select the AI model",
						Buttons: []model.Button{{
							Name: "Normal",
							CallbackData: model.CallbackData{
								CommandKey: "start",
								Case:       "createProfile",
								Step:       2,
								Payload:    "mixtral-8x7b-instruct",
							},
						},
							{
								Name: "Creative",
								CallbackData: model.CallbackData{
									CommandKey: "start",
									Case:       "createProfile",
									Step:       2,
									Payload:    "llama-2-70b-chat",
								},
							}},
					},
				},
			},
		},
		"profile": model.Usecase{
			"options": model.Chain{
				0: model.Action{
					Handler: model.ProfileOptionsHandler,
					Message: model.Message{
						Text: "What do you want to do?",
						Buttons: []model.Button{
							{
								Name: "Create new profile",
								CallbackData: model.CallbackData{
									CommandKey: "start",
									Case:       "createProfile",
									Step:       0,
									Payload:    "",
								},
							},
							{
								Name: "View existing profiles",
								CallbackData: model.CallbackData{
									CommandKey: "profile",
									Case:       "options",
									Step:       2,
									Payload:    "",
								},
							},
						},
					},
				},
				2: model.Action{
					Handler: model.ViewProfilesHandler,
					Message: model.Message{
						Text:    "View",
						Buttons: []model.Button{},
					},
				},
			},
		},
	}
}

// Assuming model has this method
