package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func GetReservation(c *gin.Context) {

	reserveId := c.Param("reserveId")

	db, err := sql.Open("mysql", "root:BNSoft2021@@tcp(125.7.152.117:30336)/bnspace")
	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Query("SELECT * FROM TB_MEETING_RESERVATION WHERE reserve_id = ?", reserveId)
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count) // interface는 java interface와 같다
	valuePtrs := make([]interface{}, count)

	final_result := map[int]map[string]string{}
	result_id := 0

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		db.Close()
		recover()
	}()

	for rows.Next() {

		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...) // scan에 포인터를 넘겨서 포인터에 값을 직접 넣는다.

		tmp_struct := map[string]string{}

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte) // []byte 타입으로 interface를 가져옴 ( 인터페이스는 모든 자료형이 들어감 ) ( 값, 타입맞는지 Bool )
			if ok {
				v = string(b)
				tmp_struct[col] = fmt.Sprintf("%s", v)
			} else {
				v = val
				tmp_struct[col] = fmt.Sprintf("%d", v)
			}

		}

		final_result[result_id] = tmp_struct
		result_id++

	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Success",
		"result":    final_result,
		"reserveId": reserveId,
	})

}

func GetReservationList(c *gin.Context) {

	db, err := sql.Open("mysql", "root:BNSoft2021@@tcp(125.7.152.117:30336)/bnspace")
	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Query("SELECT * FROM TB_MEETING_RESERVATION")

	log.Print(err)

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	final_result := []map[string]string{}

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		db.Close()
		recover()
	}()

	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]

		}
		rows.Scan(valuePtrs...)

		tmp_struct := map[string]string{}

		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
				tmp_struct[col] = fmt.Sprintf("%s", v)
			} else {
				v = val
				tmp_struct[col] = fmt.Sprintf("%d", v)
			}

		}

		final_result = append(final_result, tmp_struct)

	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Success",
		"connect success": db,
		"reserves":        final_result,
	})

}

func AddReservation(c *gin.Context) {
	meetingStartTime := c.DefaultPostForm("meetingStartTime", "current_timestamp()") // 회의 시작시간
	meetingEndTime := c.DefaultPostForm("meetingEndTime", "current_timestamp()")     // 회의 종료시간
	meetingPlaceCode := c.DefaultPostForm("meetingPlaceCode", "1")                   // 회의 장소 코드
	meetingDepartmentCode := c.DefaultPostForm("meetingDepartmentCode", "SSD")       // 예약 부서 코드
	reserveUserName := c.DefaultPostForm("reserveUserName", "")                      // 예약자명
	meetingContent := c.DefaultPostForm("meetingContent", "")                        // 회의 내용 ( 메모 )
	gubun := c.DefaultPostForm("gubun", "0")                                         // 회의 구분
	regUser := reserveUserName
	modUser := reserveUserName
	var reserveId int64 = 0
	db, err := sql.Open("mysql", "root:BNSoft2021@@tcp(125.7.152.117:30336)/bnspace")
	if err != nil {
		log.Panic(err)
	}
	var cnt int
	_ = db.QueryRow("SELECT COUNT(*) FROM TB_MEETING_RESERVATION WHERE meeting_end_time > ? AND meeting_start_time  < ? AND meeting_place_code = ?", meetingStartTime, meetingEndTime, meetingPlaceCode).Scan(&cnt)
	if cnt > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Fail",
			"err":     err,
		})
	} else {
		result, err2 := db.Exec("INSERT INTO TB_MEETING_RESERVATION (meeting_start_time, meeting_end_time, meeting_place_code, meeting_department_code, reserve_user_name, meeting_content, reg_user, mod_user, gubun) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)", meetingStartTime, meetingEndTime, meetingPlaceCode, meetingDepartmentCode, reserveUserName, meetingContent, regUser, modUser, gubun)
		if err2 != nil {
			log.Panic(err2)
		} else {
			id, err := result.LastInsertId()
			if err != nil {
				log.Panic(err)
			} else {
				reserveId = id
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"message":               "Success",
			"meetingStartTime":      meetingStartTime,
			"meetingEndTime":        meetingEndTime,
			"meetingPlaceCode":      meetingPlaceCode,
			"meetingDepartmentCode": meetingDepartmentCode,
			"reserveUserName":       reserveUserName,
			"meetingContent":        meetingContent,
			"gubun":                 gubun,
			"regUser":               regUser,
			"modUser":               modUser,
			"result":                result,
			"reserveId":             reserveId,
			"err":                   err,
		})
	}

	defer func() {
		db.Close()
		recover()
	}()

}

func ModifyReservation(c *gin.Context) {

	reserveId := c.Param("reserveId")                                                // 예약 일련번호
	meetingStartTime := c.DefaultPostForm("meetingStartTime", "current_timestamp()") // 회의 시작시간
	meetingEndTime := c.DefaultPostForm("meetingEndTime", "current_timestamp()")     // 회의 종료시간
	meetingPlaceCode := c.DefaultPostForm("meetingPlaceCode", "01")                  // 회의 장소 코드
	meetingDepartmentCode := c.DefaultPostForm("meetingDepartmentCode", "01")        // 예약 부서 코드
	reserveUserName := c.DefaultPostForm("reserveUserName", "")                      // 예약자명
	meetingContent := c.DefaultPostForm("meetingContent", "")                        // 회의 내용 ( 메모 )
	gubun := c.DefaultPostForm("gubun", "0")                                         // 회의 구분
	modUser := reserveUserName

	if reserveId == "" {
		log.Panic("no have reserveId")
	} // 예약 일련번호

	if reserveId == "" {
		log.Panic("no have reserveID")
	}

	db, err := sql.Open("mysql", "root:BNSoft2021@@tcp(125.7.152.117:30336)/bnspace")
	if err != nil {
		log.Panic(err)
	}

	var cnt int
	_ = db.QueryRow("SELECT COUNT(*) FROM TB_MEETING_RESERVATION WHERE meeting_end_time > ? AND meeting_start_time  < ? AND meeting_place_code = ? AND reserve_id != ?", meetingStartTime, meetingEndTime, meetingPlaceCode, reserveId).Scan(&cnt)

	if cnt > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Fail",
			"err":     err,
		})
	} else {
		result, err := db.Exec("UPDATE TB_MEETING_RESERVATION SET meeting_start_time = ?, meeting_end_time = ?, meeting_place_code = ?, meeting_department_code = ?, reserve_user_name = ?, meeting_content = ?, mod_user = ?, gubun = ? WHERE reserve_id = ?", meetingStartTime, meetingEndTime, meetingPlaceCode, meetingDepartmentCode, reserveUserName, meetingContent, modUser, gubun, reserveId)
		if err != nil {
			log.Panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message":               "Success",
			"meetingStartTime":      meetingStartTime,
			"meetingEndTime":        meetingEndTime,
			"meetingPlaceCode":      meetingPlaceCode,
			"meetingDepartmentCode": meetingDepartmentCode,
			"reserveUserName":       reserveUserName,
			"meetingContent":        meetingContent,
			"modUser":               modUser,
			"gubun":                 gubun,
			"result":                result,
			"err":                   err,
		})
	}

	defer func() {
		db.Close()
		recover()
	}()
}

func DeleteReservation(c *gin.Context) {

	reserveId := c.Param("reserveId") // 예약 일련번호
	if reserveId == "" {
		log.Panic("no have reserveId")
	}

	db, err := sql.Open("mysql", "root:BNSoft2021@@tcp(125.7.152.117:30336)/bnspace")
	if err != nil {
		log.Panic(err)
	}

	result, err := db.Exec("DELETE FROM TB_MEETING_RESERVATION WHERE reserve_id = ?", reserveId)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		db.Close()
		recover()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "delete Success",
		"result":  result,
		"err":     err,
	})
}
