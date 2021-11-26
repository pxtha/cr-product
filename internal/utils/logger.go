package utils

import (
	"cr-product/conf"
	"fmt"
	"os"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	//log levels
	FATAL_LOG = "FATAL"
	ERROR_LOG = "ERROR"
	WARN_LOG  = "WARN"
	INFO_LOG  = "INFO"
	DEBUG_LOG = "DEBUG"
)

func InitLogger() {
	conf.SetEnv()
	if conf.LoadEnv().Environment == conf.LoadEnv().EnvDev {
		log.SetFormatter(&log.TextFormatter{
			TimestampFormat: time.RFC3339,
			DisableColors:   false,
			FieldMap: log.FieldMap{
				log.FieldKeyTime: "logrus_time",
				log.FieldKeyMsg:  "message",
			},
		})
		log.SetLevel(log.TraceLevel) //or above
	} else {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: log.FieldMap{
				log.FieldKeyTime: "logrus_time",
				log.FieldKeyMsg:  "message",
			},
		})
		log.SetOutput(os.Stdout)
		log.SetLevel(log.TraceLevel) //or above
	}
}

func Log(level string, message interface{}, err error, messageId string) {
	infoLogType := ""
	switch message.(type) {
	case string:
		//separate message and error by spaces
		if message != "" && err != nil {
			message = message.(string) + " "
		}
		infoLogType = "job_info"
	}

	now := time.Now()
	unixNow := now.Unix()
	loc := time.FixedZone("UTC+7", 1*13*16)
	//Timezone, err := time.LoadLocation(loc)
	if err != nil {
		fmt.Println(err)
	}
	readableTime := now.In(loc).Format(time.RFC3339)

	switch level {
	case INFO_LOG:
		switch infoLogType {
		//case "download_info":
		//	downloadInfo := message.(old_data.DownloadInfo)
		//	downloadInfo.DownloadSize = int64(len(downloadInfo.DownloadData))
		//	downloadInfo.DownloadData = ""
		//	//downloadInfoStr, _ := json.Marshal(downloadInfo)  //+ " - DownloadInfo: " + string(downloadInfoStr)
		//	log.WithFields(
		//		log.Fields{
		//			"id":            messageId,
		//			"unix_time":     unixNow,
		//			"readable_time": readableTime,
		//			"download_info": downloadInfo,
		//		}).Info("") //.Info("Worker - Download time - " + downloadInfo.TargetAddress + " - took " + DoubleToStr(downloadInfo.TotalDownloadTime) + " seconds")
		//case "job_info":
		//	jobInfo := message.(old_data.JobLogInfo)
		//	jobInfoInfoStr, _ := json.Marshal(jobInfo)
		//	log.WithFields(
		//		log.Fields{
		//			"id":            messageId,
		//			"unix_time":     unixNow,
		//			"readable_time": readableTime,
		//			"job_info":      jobInfo,
		//		}).Info("Worker - Job lifetime " + jobInfo.Status + " - took " + DoubleToStr(jobInfo.ProcessTime) +
		//		" seconds - JobLogInfo: " + string(jobInfoInfoStr))
		default:
			log.WithFields(
				log.Fields{
					"id":            messageId,
					"unix_time":     unixNow,
					"readable_time": readableTime,
				}).Info(message)
		}
	case DEBUG_LOG:
		log.WithFields(
			log.Fields{
				"id":            messageId,
				"unix_time":     unixNow,
				"readable_time": readableTime,
			}).Debug(message)
	case WARN_LOG:
		log.WithFields(
			log.Fields{
				"id":            messageId,
				"unix_time":     unixNow,
				"readable_time": readableTime,
			}).Warn(message)
	case ERROR_LOG:
		//file, line, funcName := trace()
		log.WithFields(
			log.Fields{
				"id":            messageId,
				"unix_time":     unixNow,
				"readable_time": readableTime,
			}).Error(message, err)
	case FATAL_LOG:
		log.WithFields(
			log.Fields{
				"id":            messageId,
				"unix_time":     unixNow,
				"readable_time": readableTime,
			}).Log(log.FatalLevel, message, err) //if using .Fatal(message, err) directly -> logger auto quit application

	}
}

func trace() (string, int, string) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.File, frame.Line, frame.Function
}
func FailOnError(err error, message string, messageId string) {
	if err != nil {
		Log(FATAL_LOG, message, err, messageId)
		//force to exit
		os.Exit(10)
	}
}

//func CreateJobLogInfo(
//	job old_data.Job,
//	isSuccess bool,
//	startAt time.Time,
//	crawlerName string,
//	downloadIp string,
//	jobResult old_data.JobResult,
//) old_data.JobLogInfo {
//
//	processTime := time.Since(startAt).Seconds()
//	jobStatusStr := "ok"
//	if isSuccess == false {
//		jobStatusStr = "fail"
//	}
//
//	errorCount := 0
//	warnCount := 0
//	fatalCount := 0
//	for _, error := range jobResult.Data.Errors {
//		switch error.ErrorLevel {
//		case ERROR_LOG:
//			errorCount++
//		case WARN_LOG:
//			warnCount++
//		case FATAL_LOG:
//			fatalCount++
//		}
//	}
//
//	dataCountMap := make(map[string]int)
//
//	//check has raw data
//	rawLen := len(jobResult.Data.RawDataArray)
//
//	//check has public data
//	publicProductLen := len(jobResult.Data.PublicData.PublicProductData)
//	publicShopLen := len(jobResult.Data.PublicData.PublicShopData)
//	publicLen := publicProductLen + publicShopLen
//
//	if rawLen > 0 {
//		switch job.LinkType {
//		case "get_listings":
//			dataCountMap["RawListings"] = rawLen
//		case "get_orders":
//			dataCountMap["RawOrders"] = rawLen
//		case "get_transaction_details":
//			dataCountMap["RawTransactions"] = rawLen
//		case "get_seller_performance":
//			dataCountMap["RawShopPerformances"] = rawLen
//		}
//	} else if publicLen > 0 {
//		//switch job.LinkType {
//		//case LINK_TYPE_PUBLIC_PRODUCT:
//		//	dataCountMap["PublicProductData"] = publicProductLen
//		//case LINK_TYPE_PUBLIC_SHOP:
//		//	dataCountMap["PublicShopData"] = publicShopLen
//		//}
//	} else {
//		switch job.LinkType {
//		case "create_access_token", "refresh_access_token":
//			dataCountMap["Accesstokens"] = len(jobResult.Data.Data.Accesstokens)
//		case "get_categories":
//			dataCountMap["Categories"] = len(jobResult.Data.Data.Categories)
//		case "get_listings":
//			dataCountMap["Listings"] = len(jobResult.Data.Data.Listings)
//			dataCountMap["ListingVariants"] = len(jobResult.Data.Data.ListingVariants)
//		case "get_orders":
//			dataCountMap["Listings"] = len(jobResult.Data.Data.Listings)
//			dataCountMap["ListingVariants"] = len(jobResult.Data.Data.ListingVariants)
//			dataCountMap["Orders"] = len(jobResult.Data.Data.Orders)
//			dataCountMap["OrderItems"] = len(jobResult.Data.Data.OrderItems)
//		case "process_orders_cancel", "process_orders_rts", "process_orders_packed", "process_orders_rts_sof", "process_orders_delivered_sof":
//			dataCountMap["Listings"] = len(jobResult.Data.Data.Listings)
//			dataCountMap["ListingVariants"] = len(jobResult.Data.Data.ListingVariants)
//			dataCountMap["Orders"] = len(jobResult.Data.Data.Orders)
//			dataCountMap["OrderItems"] = len(jobResult.Data.Data.OrderItems)
//		case "update_listings_price", "update_listings_stock", "update_listings_active", "update_listings_images", "update_listings", "create_listing":
//			dataCountMap["Listings"] = len(jobResult.Data.Data.Listings)
//			dataCountMap["ListingVariants"] = len(jobResult.Data.Data.ListingVariants)
//			dataCountMap["DeletedListings"] = len(jobResult.Data.Data.DeletedListings)
//		}
//	}
//
//	jobLogInfo := old_data.JobLogInfo{
//		Status:       jobStatusStr,
//		JobDomain:    job.Domain,
//		Type:         job.LinkType,
//		StartAt:      startAt.Unix(),
//		ProcessTime:  processTime,
//		Retry:        cast.ToInt(job.AdditionalInfo.Flow.Meta.Retry),
//		CrawlerName:  crawlerName,
//		Ip:           downloadIp,
//		ErrorCount:   errorCount,
//		WarnCount:    warnCount,
//		FatalCount:   fatalCount,
//		DataCountMap: dataCountMap,
//	}
//	return jobLogInfo
//}
