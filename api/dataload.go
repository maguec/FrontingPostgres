package api

import (
	"fmt"
	"time"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
  "go.uber.org/zap"
)

func generateProfile(id int) Profile {

	profile := Profile{
		ID:          id,
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		Email:       gofakeit.Email(),
		Phone:       gofakeit.Phone(),
		SSN:         gofakeit.SSN(),
		Title:       gofakeit.JobTitle(),
		Company:     gofakeit.Company(),
		SecondaryId: fmt.Sprintf("user%d", id),
	}
	return profile
}

func genProfiles(count int) []Profile {
	var profiles []Profile
	for i := 1; i <= count; i++ {
		profiles = append(profiles, generateProfile(i))
	}
	return profiles
}

func loadProfiles(count int, db *gorm.DB) error {
	var w []Profile
	var err error
	//var p bytes.Buffer
	//ctx := context.Background()
	profiles := genProfiles(count)
	db.AutoMigrate(&Profile{})
	for i := 0; i < len(profiles); i++ {
		w = append(w, profiles[i])
		if len(w) == 500 {
			//err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&w).Error
			err = db.Clauses(hints.CommentAfter("returning", "route='/load',module='api.Dataload'")).Clauses(clause.OnConflict{UpdateAll: true}).Create(&w).Error
			if err != nil {
				return err
			}
			w = nil
		}
		if len(w) > 0 {
			//	err = db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&w).Error
			err = db.Clauses(hints.CommentAfter("returning", "route='/load',module='api.Dataload'")).Clauses(clause.OnConflict{UpdateAll: true}).Create(&w).Error
		}
	}
	return err
}

func Dataload(c *gin.Context) {
	record_count := c.MustGet("datasize").(int)
  logger := c.MustGet("logger").(*zap.SugaredLogger)
  start := time.Now()
	err := loadProfiles(
		record_count,
		c.MustGet("db").(*gorm.DB),
	)
	if err != nil {
    logger.Errorw("dataload", "error", err, "elapsed", time.Since(start).Milliseconds(), "record_count", record_count)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
  logger.Infow("dataload", "elapsed", time.Since(start).String(), "record_count", record_count)
	c.JSON(http.StatusOK, gin.H{
		"profiles_loaded": record_count,
	})
}
