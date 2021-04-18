package model

import (
	"time"
)

type UserData struct {
	Id                int32
	LinkAccount       string    //ссылка на аккаунт
	Title             string    //ник аккаунта
	SubTitle          string    //имя аккаунта
	Comment           string    //комментарий / описание аккаунта
	Mail              string    //почта в комменте
	Telegram          string    //телеграм в комменте
	Instagram         string    //инстаграм в комменте
	Linkes            string    //ссылки в комменте
	LanguageAccount   string    //имя аккаунта
	Following         int32     //~ количество подписок
	Followers         int32     //~ количество подписчиков
	Likes             int32     //~ количество лайков
	LastPostShowTotal int32     //количество показов всего
	AverageShows      int32     //среднее кол-во просмотров
	MedianShows       int32     //медиана просмотрова
	TotalPosts        int32     //количество постов аккаунта
	LastActionTime    time.Time //время последнего поста
	ParsingTime       time.Time //время парсинга аккаунта
}
