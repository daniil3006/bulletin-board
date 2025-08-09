package main

import (
	"bulletin-board/domain"
	"bulletin-board/storage"
	"fmt"
)

func main() {
	ad := domain.Ad{
		Title:       "Телефон",
		Description: "Iphone 16",
		Price:       50000,
		Contact:     "88234553535",
	}
	fs := storage.FileStore{}
	fs.NewBasePath("data.json")
	ad, err := fs.Create(ad)
	if err != nil {
		fmt.Println("ошибка", err)
	} else {
		fmt.Println(ad)
	}
	//
	//err := fs.Update(ad)
	//if err != nil {
	//	fmt.Println("ошибка", err)
	//}

	//items, err := fs.List()
	//if err != nil {
	//	fmt.Println("ошибка", err)
	//}
	//fmt.Println(items)
}
