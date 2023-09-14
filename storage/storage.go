package storage

import (
	"database/sql"
	"fmt"
	"log"
	"tgbotik/errs"

	_ "github.com/lib/pq"
)

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

// Return the *PostgreDB with params
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

// Connect to the database with params
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

// Check exist in database. If exist, returns true
func (s *PostgreDB) Exists(table string, colArgs []string, values []interface{}) (bool, error) {
	req := fmt.Sprintf("select * from %s where ", table)
	if len(colArgs) == 1 {
		req += fmt.Sprintf("%s = $1", colArgs[0])
	} else {
		req += colArgs[0] + " = $1"
		for idx, col := range colArgs[1:] {
			req += fmt.Sprintf(" and %s = $%d", col, idx+2)
		}
	}
	res, err := s.DB.Query(req, values...)
	errs.LogError(err)

	var i int
	for res.Next() {
		i += 1
	}
	if i == 0 {
		return false, nil
	}
	return true, nil
}

// Save the user into a Users table. If user is already exist, activate him in DB
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
		log.Printf("Saving user %d to database", chatID)
		return errs.CheckError(err)
	} else if !stat {
		return s.UpdateUserStatus(chatID, true)
	}

	return err
}

// Save anime into a Animes table with status false
func (s *PostgreDB) SaveAnime(animeName, lastSeriaText, lastSeriaHref string, countSeries int) error {
	res, err := s.DB.Query("select * from Animes where AnimeName = $1", animeName)
	defer res.Close()
	if err != nil {
		return err
	}
	rowCount := countRows(res)
	if rowCount == 0 {
		_, err := s.DB.Exec("insert into Animes (AnimeName, CountSeries, Status, LastSeriaText, LastSeriaHref) values ($1, $2, $3, $4, $5)",
			animeName, countSeries, false, lastSeriaText, lastSeriaHref)
		log.Print("Saving anime to database")
		return errs.CheckError(err)
	}
	return fmt.Errorf("This Anime is already exists: %s", animeName)
}

// Create a new row in Subscribers table
func (s *PostgreDB) Subscribe(chatID, animeID int) error {
	var args []interface{}
	args = append(args, chatID)
	args = append(args, animeID)
	exist, err := s.Exists("Subscribers", []string{"ChatID", "AnimeID"}, args)
	errs.LogError(err)
	if !exist {
		_, err = s.DB.Exec("insert into Subscribers (ChatID, AnimeID) values ($1, $2)", chatID, animeID)
	}

	return errs.CheckError(err)
}

// Removes the row in Subscribers table
func (s *PostgreDB) Unsubscribe(chatID, animeID int) error {
	_, err := s.DB.Exec("delete from Subscribers where ChatID = $1 and AnimeID = $2", chatID, animeID)
	if err == sql.ErrNoRows {
		return nil
	}
	return errs.CheckError(err)
}

// Returns a list of new series from database
func (s *PostgreDB) GetNewSeries() ([]NewSeria, error) {
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

// Update columns LastSeriaText and LastSeriaHref for animeID in Animes table
func (s *PostgreDB) SetNewSeries(animeID int, seriaText, seriaHref string) error {
	req := "update Animes set LastSeriaText = $1 and LastSeriaHref = $2 and Status = $3 where ID = $4"
	_, err := s.DB.Exec(req, seriaText, seriaHref, true, animeID)
	return errs.CheckError(err)
}

// Update column Status
func (s *PostgreDB) SetStatus(table string, id int, status bool) error {
	var args []interface{}
	args = append(args, id)
	exists, err := s.Exists(table, []string{"ID"}, args)
	log.Print(errs.CheckError(err))
	if exists {
		req := fmt.Sprintf("update %s set Status = $1 where ID = $2", table)
		_, err := s.DB.Exec(req, status, id)
		return errs.CheckError(err)
	}
	return fmt.Errorf("Wrong table or id, row isn't exists")
}

// Returns list of users who subscribed on anime with animeID
func (s *PostgreDB) GetSubscribers(animeID int) ([]int, error) {
	var args []interface{}
	args = append(args, animeID)
	exist, err := s.Exists("Subscribers", []string{"AnimeID"}, args)
	errs.LogError(err)
	if !exist {
		return nil, fmt.Errorf("Subscribers on this anime aren't exist")
	}
	req := fmt.Sprintf("select ChatID from Subscribers where AnimeID = $1")
	rows, err := s.DB.Query(req, animeID)
	errs.LogError(err)
	var res []int
	var val int
	for rows.Next() {
		rows.Scan(&val)
		user, err := s.DB.Query("select ID from Users where Status = true and ID = $1", val)
		errs.LogError(err)
		for user.Next() {
			res = append(res, val)
		}
	}
	return res, nil
}

// Return ID and count of series of anime
func (s *PostgreDB) CountSeries(animeName string) (int, int, error) {
	var args []interface{}
	args = append(args, animeName)
	exist, err := s.Exists("Animes", []string{"AnimeName"}, args)
	errs.LogError(err)

	if !exist {
		return 0, 0, fmt.Errorf("Can't fetch CountSeries - %s doesn't exist", animeName)
	}
	rows, err := s.DB.Query("select ID, CountSeries from Animes where AnimeName = $1", animeName)
	errs.LogError(err)
	var id, val int
	for rows.Next() {
		rows.Scan(&id, &val)
	}
	return id, val, nil
}

// Returns list of all saved anime such as (list_ID, list_animeName, error)
func (s *PostgreDB) Animes() ([]int, []string, error) {
	rows, err := s.DB.Query("select ID, AnimeName from Animes")
	errs.LogError(err)
	var resID []int
	var resName []string
	var name string
	var id int
	for rows.Next() {
		rows.Scan(&id, &name)
		resID = append(resID, id)
		resName = append(resName, name)
	}
	return resID, resName, nil
}

// Update users status
func (s *PostgreDB) UpdateUserStatus(chatID int, status bool) error {
	err := s.SetStatus("Users", chatID, status)
	log.Printf("Updating status of user %d in database", chatID)
	return errs.CheckError(err)
}

// Count rows in sql.Rows
func countRows(res *sql.Rows) int {
	var i int
	for res.Next() {
		i += 1
	}
	return i
}
