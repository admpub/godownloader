package httpclient
import (
	"net/http"
	"log"
	"errors"
	"strconv"
	"os"
)
const FlushDiskSize = 1024*1024
func CheckMultipart(urls string) (bool, error) {
	r, err := http.NewRequest("GET", urls, nil)
	if err!=nil {
		return false, err
	}
	r.Header.Add("Range", "bytes=0-0")
	cl := http.Client{}
	resp, err := cl.Do(r)
	if err!=nil {
		log.Printf("error: can't check multipart support assume no %v \n", err)
		return false, err
	}
	if resp.StatusCode!=206 {
		return false, errors.New("error: file not found or moved status: "+ resp.Status)
	}
	if (resp.ContentLength==1) {
		log.Printf("info: file size is %d bytes \n", resp.ContentLength)
		return true, nil
	}
	return false, nil
}

func GetSize(urls string) (int64, error) {
	cl := http.Client{}
	resp, err := cl.Head(urls)
	if err!=nil {
		log.Printf("error: when try get file size %v \n", err)
		return 0, err
	}
	if resp.StatusCode!=200 {
		log.Printf("error: file not found or moved status:", resp.StatusCode)
		return 0, errors.New("error: file not found or moved")
	}
	log.Printf("info: file size is %d bytes \n", resp.ContentLength)
	return resp.ContentLength, nil
}

type DownloadProgress struct {
	from int64
	to   int64
	pos  int64
}
type PartialDownloader struct {
	dp     DownloadProgress
	client http.Client
	req    *http.Response
	url    string
	file   *os.File
}
func (pd *PartialDownloader) Init(url string, file  *os.File, from int64, pos int64, to int64) {
	pd.file=file
	pd.url=url
	pd.dp.from=from
	pd.dp.to=to
	pd.dp.pos=pos
}

func (pd PartialDownloader) GetProgress() DownloadProgress {
	return pd.dp
}

func (pd *PartialDownloader)BeforeRun() error {
	nos, err := CheckMultipart(pd.url)
	if !nos {
		return errors.New("error: server unsupport part support")
	}
	if err!=nil {
		return err
	}
	//create new req
	r, err := http.NewRequest("GET", pd.url, nil)
	if err!=nil {
		return err
	}
	r.Header.Add("Range", "bytes="+strconv.FormatInt(pd.dp.from, 10)+"-"+strconv.FormatInt(pd.dp.to, 10))
	//ok we construct query
	//try send request
	resp, err := pd.client.Do(r)
	if err!=nil {
		log.Printf("error: error download part file%v \n", err)
		return err
	}
	//check response
	if resp.StatusCode!=206 {
		log.Printf("error: file not found or moved status:", resp.StatusCode)
		return errors.New("error: file not found or moved")
	}
	pd.req=resp
	return nil
}
func (pd *PartialDownloader)AfterStop() error {
	return nil
}
func (pd *PartialDownloader)DownloadSergment() (bool, error) {
	//write flush data to disk
	 buffer:=make([]byte, FlushDiskSize, FlushDiskSize)
	count,err:= pd.req.Body.Read(buffer)
	if (err!=nil) {
		pd.file.Sync()
		return false, err
	}
	realc, err := pd.file.WriteAt(buffer[:count], pd.dp.pos)
	if err!=nil {
		return false, err
	}

	pd.dp.pos=pd.dp.pos+int64(realc)
	if (pd.dp.pos==pd.dp.to) {
		//ok download part complete normal
		pd.file.Sync()
		return false, nil
	}
	//not full download next segment
	return true, nil
}
func (pd *PartialDownloader) DoWork() (bool, error) {
	return pd.DownloadSergment()
}
