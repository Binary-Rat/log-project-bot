package tg

import (
	"log-proj/pkg/clients/tg"
)

const (
	atisuCall = "atisu"
)

var (
	atiSearchKeyboard = [][]tg.KeyboardButton{
		{
			{
				Text:         "Подобрать машину на Ati.SU",
				CallbackData: atisuCall,
			},
		},
	}
	startKeyboard = [][]tg.KeyboardButton{
		{
			{
				Text: "/start",
			},
			{
				Text: "/help",
			},
			{
				Text: "/calc",
			},
		},
	}
)
