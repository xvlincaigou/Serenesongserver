package utils

import (
	"Serenesongserver/config"
	"Serenesongserver/models"
	"context"
	"fmt"
	"strings"

	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	RecommendedCi      models.Ci
	RecommendedPicPath string
)

func RandomSelectCi() error {
	pipeline := []bson.D{
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: 1}}}},
	}

	cursor, err := config.MongoClient.Database("serenesong").Collection("Ci").Aggregate(context.Background(), pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	if cursor.Next(context.Background()) {
		return cursor.Decode(&RecommendedCi)
	}

	return errors.New("no document found")
}

type WanxRequest struct {
	Model string `json:"model"`
	Input struct {
		Prompt string `json:"prompt"`
	} `json:"input"`
	Parameters struct {
		Style string `json:"style"`
		Size  string `json:"size"`
		N     int    `json:"n"`
	} `json:"parameters"`
}

type TaskResponse struct {
	Output struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
	} `json:"output"`
	RequestID string `json:"request_id"`
}

type TaskResult struct {
	RequestID string `json:"request_id"`
	Output    struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"`
		Results    []struct {
			URL string `json:"url"`
		} `json:"results"`
		TaskMetrics struct {
			Total     int `json:"TOTAL"`
			Succeeded int `json:"SUCCEEDED"`
			Failed    int `json:"FAILED"`
		} `json:"task_metrics"`
	} `json:"output"`
	Usage struct {
		ImageCount int `json:"image_count"`
	} `json:"usage"`
}

func createImageTask(apiKey, prompt string) (string, error) {
	reqBody := WanxRequest{
		Model: config.Model,
	}
	reqBody.Input.Prompt = prompt
	reqBody.Parameters.Style = "<auto>"
	reqBody.Parameters.Size = "1024*1024"
	reqBody.Parameters.N = 1

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST",
		config.GeneratePicURL,
		bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-DashScope-Async", "enable")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	var taskResp TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&taskResp); err != nil {
		return "", err
	}

	return taskResp.Output.TaskID, nil
}

func queryTaskResult(apiKey, taskID string) (*TaskResult, error) {
	url := fmt.Sprintf("%s%s", config.CheckPicURL, taskID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("查询任务失败，状态码: %d", resp.StatusCode)
	}

	var result TaskResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func downloadImage(url, filePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载图片失败，状态码: %d", resp.StatusCode)
	}

	// 创建目录（如果不存在）
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func GenerateAndDownloadImage() error {
	// 1.获取随机的词
	if err := RandomSelectCi(); err != nil {
		return fmt.Errorf("获取词失败: %w", err)
	}

	// 2. 提取词的文字
	prompt := strings.Join(RecommendedCi.Content, "\n")

	// 3. 创建任务
	taskID, err := createImageTask(config.ApiKey, prompt)
	if err != nil {
		return err
	}

	// 4. 轮询获取结果
	var result *TaskResult
	for i := 0; i < 30; i++ { // 最多等待90秒
		result, err = queryTaskResult(config.ApiKey, taskID)
		if err != nil {
			return fmt.Errorf("查询任务结果失败: %w", err)
		}

		if result.Output.TaskStatus == "SUCCEEDED" {
			break
		}

		if result.Output.TaskStatus == "FAILED" {
			return errors.New("任务执行失败")
		}

		time.Sleep(3 * time.Second) // 每3秒查询一次
	}

	if result == nil || result.Output.TaskStatus != "SUCCEEDED" {
		return errors.New("任务超时")
	}

	// 5. 下载图片
	for i, img := range result.Output.Results {
		fileName := fmt.Sprintf("image_%s_%d.png", taskID, i)
		filePath := filepath.Join(config.PicFolder, fileName)

		if err := downloadImage(img.URL, filePath); err != nil {
			return fmt.Errorf("下载图片失败: %w", err)
		}
	}

	RecommendedPicPath = filepath.Join(config.PicFolder, fmt.Sprintf("image_%s_0.png", taskID))

	return nil
}

func GenerateAndDownloadImageWrapper() {
	if err := GenerateAndDownloadImage(); err != nil {
		log.Println(err)
	}
}
