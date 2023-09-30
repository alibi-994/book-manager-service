package db

import (
	"log"
	"time"
)

type Book struct {
	ID              int       `gorm:"primaryKey"`
	Name            string    `gorm:"varchar(50)"`
	Author          Author    `gorm:"embedded"`
	Category        string    `gorm:"varchar(50)"`
	Volume          int       `gorm:"tinyint"`
	PublishedAt     time.Time `gorm:"dattimeoffset(7)"`
	Summary         string    `gorm:"varchar(100)"`
	TableOfContents []string  `gorm:"-"`
	Publisher       string    `gorm:"varchar(50)"`
}

type Author struct {
	FirstName   string    `gorm:"varchar(50)"`
	LastName    string    `gorm:"varchar(50)"`
	Birthday    time.Time `gorm:"datetimeoffset(7)"`
	Nationality string    `gorm:"varchar(50)"`
}

func (gdb *GormDB) CreteNewBook(b *Book) error {
	if err := gdb.db.Create(b).Error; err != nil {
		return err
	}
	return nil
}

func (gdb GormDB) GetAllBooks() ([]Book, error) {
	var books []Book
	err := gdb.db.Model(&Book{}).Find(&books).Error
	if err != nil {
		log.Print("Could not find books")
		return nil, err
	}
	return books, nil
}

func (gdb *GormDB) GetBookByID(id int) (*Book, error) {
	var book Book
	err := gdb.db.Model(&Book{}).Where(&Book{ID: id}).First(&book).Error
	if err != nil {
		log.Print("Could not find the book")
		return &Book{}, err
	}
	return &book, nil
}

func (gdb *GormDB) UpdateBook(id int, updateQuery map[string]interface{}) error {
	keys := make([]string, 0, len(updateQuery))
	for key := range updateQuery {
		keys = append(keys, key)
	}

	for _, key := range keys {
		gdb.db.Model(&Book{}).Where(&Book{ID: id}).Update(key, updateQuery[key])
	}
	return nil
}

func (gdb *GormDB) DeleteBook(id int) error {
	err := gdb.db.Delete(&Book{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
