package tg

import (
	"log-proj/pkg/clients/tg"
)

const (
	atisuCall = "atisu"
)

var (
	keybord = [][]tg.InlineKeyboardButton{
		{
			{
				Text:         "Подобрать машину на Ati.SU",
				CallbackData: atisuCall,
			},
		},
	}
)
