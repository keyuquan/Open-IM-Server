package open_im_sdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v6"
	"github.com/tencentyun/cos-go-sdk-v5"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"time"
)

func (u *UserRelated) tencentOssCredentials() (*paramsTencentOssCredentialResp, error) {
	resp, err := post2Api(tencentCloudStorageCredentialRouter, paramsTencentOssCredentialReq{OperationID: operationIDGenerator()}, u.token)
	if err != nil {
		return nil, err
	}

	var ossResp paramsTencentOssCredentialResp
	_ = json.Unmarshal(resp, &ossResp)

	if ossResp.ErrCode != 0 {
		return nil, errors.New(ossResp.ErrMsg)
	}

	return &ossResp, nil
}

func getMinClient() (*minio.Client, error) {
	endpoint := "1.14.194.38:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false
	bucketName := "OpenIM"
	location := "chengdu"

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			sdkLog("We already own", bucketName)
		} else {
			sdkLog("MakeBucket failed, ", err.Error())
			return nil, err
		}
	}
	sdkLog("created ok, ", bucketName)

	return minioClient, nil
}

func getImgContentTypeSuffix(filePath string) (string, string, error) {
	suffix := path.Ext(filePath)
	if len(suffix) == 0 {
		sdkLog("file name failed, ", filePath)
		return "", "", errors.New("file name failed")
	}
	contentType := "image/" + suffix[1:]
	return contentType, suffix, nil
}

func uploadImageMin(filePath string, callback SendMsgCallBack) (string, string, error) {
	return uploadObjectMin(filePath, "img", callback)
}

func uploadSoundMin(filePath string, back SendMsgCallBack) (string, string, error) {
	return uploadObjectMin(filePath, "", back)
}

func uploadFileMin(filePath string, back SendMsgCallBack) (string, string, error) {
	return uploadObjectMin(filePath, "", back)
}

func uploadVideoMin(videoPath, snapshotPath string, back SendMsgCallBack) (string, string, string, string, error) {
	snapshotURL, snapshotUUID, err := uploadObjectMin(snapshotPath, "img", nil)
	if err != nil {
		back.OnError(ErrCodeConversation, err.Error())
		return "", "", "", "", err
	}
	videoURL, videoUUID, err := uploadObjectMin(videoPath, "", back)
	return snapshotURL, snapshotUUID, videoURL, videoUUID, err
}

func uploadObjectMin(filePath string, objectType string, callback SendMsgCallBack) (string, string, error) {
	minioClient, err := getMinClient()
	if err != nil {
		sdkLog("getMinClient failed, ", err.Error())
		if callback != nil {
			callback.OnError(ErrCodeConversation, err.Error())
		}
		return "", "", err
	}

	contentType, suffix, err := getImgContentTypeSuffix(filePath)
	if err != nil {
		sdkLog("getImgContentTypeSuffix failed, ", err.Error())
		if callback != nil {
			callback.OnError(ErrCodeConversation, err.Error())
		}
		return "", "", err
	}

	if objectType != "img" {
		contentType = ""
	}
	bucketName := "OpenIM"
	newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Uint64(), suffix)
	objectName := newName

	if callback != nil {
		go func() {
			n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				sdkLog("FPutObject failed ", bucketName, objectName, filePath, err.Error())
				callback.OnError(ErrCodeConversation, err.Error())
				return
			}
			callback.OnProgress(100)
			callback.OnSuccess("ok")
			sdkLog("upload file: ", filePath, " size: ", n)
		}()
	} else {
		n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			sdkLog("FPutObject failed ", bucketName, objectName, filePath, err.Error())
			return "", "", err
		}
		sdkLog("upload file: ", filePath, " size: ", n)
	}

	reqParams := make(url.Values)
	reqParams.Set("response-content-type", contentType)

	presignedURL, err := minioClient.PresignedGetObject(bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		sdkLog("PresignedGetObject failed, ", err.Error())
		if callback != nil {
			callback.OnError(ErrCodeConversation, err.Error())
		}
		return "", "", err
	}
	return presignedURL.String(), newName, nil
}

func (ur *UserRelated) uploadImage(filePath string, back SendMsgCallBack) (string, string, error) {
	ossResp, err := ur.tencentOssCredentials()
	if err != nil {
		sdkLog("tencentOssCredentials", err.Error())
		return "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, err := url.Parse(dir)
	if err != nil {
		sdkLog("Parse", err.Error())
		return "", "", err
	}
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {
		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(filePath)
		if len(suffix) == 0 {
			return "", "", errors.New("file fail")
		}
		newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		contentType := "image/" + suffix[1:]

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				ContentType: contentType,
				Listener:    lis,
			},
		}
		_, err := client.Object.PutFromFile(context.Background(), newName, filePath, opt)
		if err != nil {
			sdkLog("file:", filePath, err.Error())
			return "", "", err
		}

		targetFileUrl := dir + "/" + newName
		return targetFileUrl, newName, nil
	}

	return "", "", errors.New("client == nil")
}

func (ur *UserRelated) uploadSound(filePath string, back SendMsgCallBack) (string, string, error) {
	ossResp, err := ur.tencentOssCredentials()
	if err != nil {
		sdkLog(err.Error())
		return "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {

		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(filePath)
		if len(suffix) == 0 {
			return "", "", errors.New("file fail")
		}
		newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		//contentType := "image/" + suffix[1:]

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				//ContentType: contentType,
				Listener: lis,
			},
		}

		_, err := client.Object.PutFromFile(context.Background(), newName, filePath, opt)
		if err != nil {
			sdkLog("PutFromFile", err.Error())
			return "", "", err
		}

		targetFile := dir + "/" + newName
		return targetFile, newName, nil
	}
	sdkLog("client == nil")
	return "", "", errors.New("client == nil")
}

func (ur *UserRelated) uploadFile(filePath string, back SendMsgCallBack) (string, string, error) {
	ossResp, err := ur.tencentOssCredentials()
	if err != nil {
		sdkLog(err.Error())
		return "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {

		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(filePath)
		if len(suffix) == 0 {
			return "", "", errors.New("file fail")
		}
		newName := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		//contentType := "image/" + suffix[1:]

		opt := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				//ContentType: contentType,
				Listener: lis,
			},
		}

		_, err := client.Object.PutFromFile(context.Background(), newName, filePath, opt)
		if err != nil {
			sdkLog(err.Error())
			return "", "", err
		}

		targetFile := dir + "/" + newName
		return targetFile, newName, nil
	}

	return "", "", errors.New("client == nil")
}

func (ur *UserRelated) uploadVideo(videoPath, snapshotPath string, back SendMsgCallBack) (string, string, string, string, error) {
	sdkLog("input args:", videoPath, snapshotPath)
	ossResp, err := ur.tencentOssCredentials()
	if err != nil {
		sdkLog("tencentOssCredentials err:", err.Error())
		return "", "", "", "", err
	}

	dir := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", ossResp.Bucket, ossResp.Region)
	u, _ := url.Parse(dir)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     ossResp.Data.Credentials.TmpSecretId,
			SecretKey:    ossResp.Data.Credentials.TmpSecretKey,
			SessionToken: ossResp.Data.Credentials.Token,
		},
	})
	if client != nil {
		var newNameSnapshot, targetSnapshot string
		if len(snapshotPath) > 0 {
			//-----first------
			suffix := path.Ext(snapshotPath)
			if len(suffix) == 0 {
				sdkLog("suffix =0 Snapshot err:")
				return "", "", "", "", errors.New("file fail")
			}
			newNameSnapshot := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
			contentTypeSnapshot := "image/" + suffix[1:]

			opt1 := &cos.ObjectPutOptions{
				ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
					ContentType: contentTypeSnapshot,
				},
			}

			_, err := client.Object.PutFromFile(context.Background(), newNameSnapshot, snapshotPath, opt1)
			if err != nil {
				sdkLog("PutFromFile Snapshot err:", err.Error())
				return "", "", "", "", err
			}
			targetSnapshot = dir + "/" + newNameSnapshot
		}

		//-----second------
		var lis = &selfListener{}
		lis.SendMsgCallBack = back

		suffix := path.Ext(videoPath)
		if len(suffix) == 0 {
			sdkLog("suffix =0  Video err:")
			return "", "", "", "", errors.New("file fail")
		}
		newNameVideo := fmt.Sprintf("%d-%d%s", time.Now().UnixNano(), rand.Int(), suffix)
		//contentType := "image/" + suffix[1:]

		opt2 := &cos.ObjectPutOptions{
			ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
				//ContentType: contentType,
				Listener: lis,
			},
		}

		_, err = client.Object.PutFromFile(context.Background(), newNameVideo, videoPath, opt2)
		if err != nil {
			sdkLog("PutFromFile Video err:", err.Error())
			return "", "", "", "", err
		}

		targetVideo := dir + "/" + newNameVideo

		sdkLog("ok", videoPath, snapshotPath, targetSnapshot, targetVideo)

		return targetSnapshot, newNameSnapshot, targetVideo, newNameVideo, nil
	}
	sdkLog("client == nil")
	return "", "", "", "", errors.New("client == nil")
}

type selfListener struct {
	SendMsgCallBack
}

func (l *selfListener) ProgressChangedCallback(event *cos.ProgressEvent) {
	switch event.EventType {
	case cos.ProgressDataEvent:
		if event.ConsumedBytes == event.TotalBytes {
			l.SendMsgCallBack.OnProgress(int((event.ConsumedBytes - 1) * 100 / event.TotalBytes))
		} else {
			l.SendMsgCallBack.OnProgress(int(event.ConsumedBytes * 100 / event.TotalBytes))
		}
		log(fmt.Sprintf("\r[ConsumedBytes/TotalBytes: %d/%d, %d%%]", event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes))

	case cos.ProgressFailedEvent:
		sdkLog(fmt.Sprintf("\nTransfer Failed: %v", event.Err))
	}
}
