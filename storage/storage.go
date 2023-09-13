package storage

import (
	"database/sql"
	"fmt"
	"log"
	"tgbotik/errs"

	_ "github.com/lib/pq"
)

type Storage interface {
	Connect() error
	SaveUser(chatID int, userName string, status bool) error
	SaveAnime(animeName, lastSeriaText, lastSeriaHref string, countSeries int) (int64, error)
	Subscribe(chatID, animeID int) error
	Unsubscribe(chatID, animeID int) error
	SetStatus(table string, id int, status bool) error
	NewSeries() ([]NewSeria, error)
	Exists(ID int) bool
}

type PostgreDB struct {
	host     string
	user     string
	password string
	dbname   string
	port     int
	DB       *sql.DB
}

type NewSeria struct {
	ID        int
	AnimeName string
	Text      string
	Href      string
}

func New(host, user, password, dbname string, port int) (*PostgreDB, error) {
	p := PostgreDB{
		host:     host,
		user:     user,
		password: password,
		dbname:   dbname,
		port:     port,
	}
	err := p.Connect()
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *PostgreDB) Connect() error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.host, s.port, s.user, s.password, s.dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return err
	}
	s.DB = db
	log.Printf("Connected to database %s", s.dbname)
	return nil
}

func (s *PostgreDB) Exists(ID int, table string) (bool, error) {
	req := fmt.Sprintf("select * from %s where ID = $1", table)
	res, err := s.DB.Query(req, ID)
	if err != nil {
		return false, err
	}
	var i int
	for res.Next() {
		i += 1
	}
	if i == 0 {
		return false, nil
	}
	return true, nil
}

func (s *PostgreDB) SaveUser(chatID int, userName string, status bool) error {
	res, err := s.DB.Query("select Status from Users where ID = $1", chatID)
	defer res.Close()

	var stat bool
	var rowCount int
	for res.Next() {
		rowCount += 1
		res.Scan(&stat)
	}
	if rowCount == 0 {
		_, err := s.DB.Exec("insert into Users (ID, Username, Status) values ($1, $2, $3)", chatID, userName, status)
		return errs.CheckError(err)
	} else if !stat {
		err = s.SetStatus("Users", chatID, true)
		log.Printf("Saving user %d to database", chatID)
		return errs.CheckError(err)
	}
	return err
}

func (s *PostgreDB) SaveAnime(animeName, lastSeriaText, lastSeriaHref string, countSeries int) error {
	res, err := s.DB.Query("select * from Animes where AnimeName = $1", animeName)
	defer res.Close()
	if err != nil {
		return err
	}
	rowCount := CountRows(res)
	if rowCount == 0 {
		_, err := s.DB.Exec("insert into Animes (AnimeName, CountSeries, Status, LastSeriaText, LastSeriaHref) values ($1, $2, $3, $4, $5)",
			animeName, countSeries, false, lastSeriaText, lastSeriaHref)
		log.Print("Saving anime to database")
		return errs.CheckError(err)
	}
	return fmt.Errorf("This Anime is already exists: %s", animeName)

}

func (s *PostgreDB) Subscribe(chatID, animeID int) error {
	_, err := s.DB.Exec("select * from Subscribers where ChatID = $1 and AnimeID = $2", chatID, animeID)
	if err == sql.ErrNoRows {
		_, err := s.DB.Exec("insert into Subscribers (ChatID, AnimeID) values ($1, $2)", chatID, animeID)
		if err != nil {
			return err
		}
		return nil
	}
	return errs.CheckError(err)
}

func (s *PostgreDB) Unsubscribe(chatID, animeID int) error {
	_, err := s.DB.Exec("delete from Subscribers where ChatID = $1 and AnimeID = $2", chatID, animeID)
	if err == sql.ErrNoRows {
		return nil
	}
	return errs.CheckError(err)
}

func (s *PostgreDB) NewSeries() ([]NewSeria, error) {
	rows, err := s.DB.Query("select ID, AnimeName, LastSeriaText, LastSeriaHref from Animes where Status = true")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	var res []NewSeria
	for rows.Next() {
		var row NewSeria
		err := rows.Scan(&row.ID, &row.AnimeName, &row.Text, &row.Href)
		if err != nil {
			log.Print("Error in NewSeries function: ", err)
			continue
		}
		res = append(res, row)
	}
	return res, nil
}

func (s *PostgreDB) SetStatus(table string, id int, status bool) error {
	exists, err := s.Exists(id, table)
	log.Print(errs.CheckError(err))
	if exists {
		req := fmt.Sprintf("update %s set Status = $1 where ID = $2", table)
		_, err := s.DB.Exec(req, status, id)
		return errs.CheckError(err)
	}
	return fmt.Errorf("Wrong table or id, row isn't exists")
}

func CountRows(res *sql.Rows) int {
	var i int
	for res.Next() {
		i += 1
	}
	return i
}
