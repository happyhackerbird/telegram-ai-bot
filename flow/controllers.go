package flow

import (
	"example/bot/telegram-ai-bot/model"
	"example/bot/telegram-ai-bot/services"
)

func Init() model.Flow {
	return model.Flow{
		"start": model.Usecase{
			"createProfile": model.Chain{
				0: model.Action{
					Handler: services.PromptProfileNameHandler,
					Message: model.Message{
						Text:    "Enter the name of the bot",
						Buttons: []model.Button{},
					},
				},
				1: model.Action{
					Handler: services.PromptInstructionHandler,
					Message: model.Message{
						Text:    "Enter the instruction for the bot",
						Buttons: []model.Button{},
					},
				},
				2: model.Action{
					Handler: services.FinalizeProfileHandler,
					Message: model.Message{
						Text: "Select the AI model",
						Buttons: []model.Button{{
							Name: "Mixtral 8x7b",
							CallbackData: model.CallbackData{
								CommandKey: "start",
								Case:       "createProfile",
								Step:       2,
								Payload:    "mixtral-8x7b-instruct",
							},
						},
							{
								Name: "Perplexity 70b",
								CallbackData: model.CallbackData{
									CommandKey: "start",
									Case:       "createProfile",
									Step:       2,
									Payload:    "pplx-70b-chat",
								},
							},
							{
								Name: "GPT-4-Turbo",
								CallbackData: model.CallbackData{
									CommandKey: "start",
									Case:       "createProfile",
									Step:       2,
									Payload:    "gpt-4-turbo-preview",
								},
							}},
					},
				},
			},
		},
		"profile": model.Usecase{
			"options": model.Chain{
				0: model.Action{
					Handler: services.ProfileOptionsHandler,
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
									Step:       1,
									Payload:    "",
								},
							},
						},
					},
				},
				1: model.Action{
					Handler: services.ViewProfilesHandler,
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
