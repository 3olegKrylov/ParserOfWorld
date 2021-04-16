package model

import (
	"time"
)

type UserData struct {
	Id                   int32
	LinkAccount          string    //ссылка на аккаунт
	Title                string    //ник аккаунта
	SubTitle             string    //имя аккаунта
	Comment              string    //комментарий / описание аккаунта
	Mail                 string    //почта в комменте
	Telegram             string    //телеграм в комменте
	Instagram            string    //инстаграм в комменте
	Linkes               string    //ссылки в комменте
	Following            int32     //количество подписок
	Followers            int32     //количество подписчиков
	Likes                int32     //~ количество лайков
	ShowTotal            int32     //~ количество показов всего
	AverageNumberOfShows int32     //~ среднее кол-во просмотров
	LastActionTime       time.Time    //время последнего поста
	LanguageAccount      string    //имя аккаунта
	ParsingTime          time.Time //время парсинга аккаунта
}
