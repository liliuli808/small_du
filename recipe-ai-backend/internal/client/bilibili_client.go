package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"recipe-ai-backend/internal/model"
	"recipe-ai-backend/internal/pkg/logger"
	"recipe-ai-backend/internal/pkg/parser"
	"sort"
	"strings"
)

// BilibiliClient B站API客户端
type BilibiliClient struct {
	httpClient *HTTPClient
	baseURL    string
}

// NewBilibiliClient 创建B站客户端
func NewBilibiliClient(httpClient *HTTPClient) *BilibiliClient {
	return &BilibiliClient{
		httpClient: httpClient,
		baseURL:    "https://api.bilibili.com",
	}
}

// ExtractBVID 从URL提取BV号
func (c *BilibiliClient) ExtractBVID(rawURL string) (string, error) {
	return parser.ParseBVID(rawURL)
}

// GetVideoInfo 获取视频信息
func (c *BilibiliClient) GetVideoInfo(ctx context.Context, bvid string) (*model.VideoInfo, error) {
	apiURL := fmt.Sprintf("%s/x/web-interface/view?bvid=%s", c.baseURL, bvid)

	resp, err := c.httpClient.Get(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求视频信息失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			AID      int64  `json:"aid"`
			BVID     string `json:"bvid"`
			Title    string `json:"title"`
			Desc     string `json:"desc"`
			Duration int    `json:"duration"`
			Owner    struct {
				Name string `json:"name"`
			} `json:"owner"`
			Pages []struct {
				CID      int64  `json:"cid"`
				Page     int    `json:"page"`
				Part     string `json:"part"`
				Duration int    `json:"duration"`
			} `json:"pages"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析视频信息失败: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("B站API错误: code=%d, msg=%s", result.Code, result.Message)
	}

	info := &model.VideoInfo{
		AID:         result.Data.AID,
		BVID:        result.Data.BVID,
		Title:       result.Data.Title,
		Description: result.Data.Desc,
		OwnerName:   result.Data.Owner.Name,
		Duration:    result.Data.Duration,
		Pages:       make([]model.VideoPage, 0, len(result.Data.Pages)),
	}

	for _, p := range result.Data.Pages {
		info.Pages = append(info.Pages, model.VideoPage{
			CID:      p.CID,
			Page:     p.Page,
			Part:     p.Part,
			Duration: p.Duration,
		})
	}

	// 默认使用第一P的CID
	if len(info.Pages) > 0 {
		info.CID = info.Pages[0].CID
	}

	logger.Info("获取视频信息成功", logger.String("bvid", bvid), logger.String("title", info.Title))
	return info, nil
}

// TryGetSubtitle 尝试获取字幕
func (c *BilibiliClient) TryGetSubtitle(ctx context.Context, videoInfo *model.VideoInfo) (string, bool) {
	if videoInfo.CID == 0 {
		return "", false
	}

	apiURL := fmt.Sprintf("%s/x/player/wbi/v2?cid=%d&bvid=%s", c.baseURL, videoInfo.CID, videoInfo.BVID)

	resp, err := c.httpClient.Get(ctx, apiURL)
	if err != nil {
		logger.ErrorLog("获取字幕列表失败", logger.Error(err))
		return "", false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLog("读取字幕列表失败", logger.Error(err))
		return "", false
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Subtitle struct {
				Subtitles []struct {
					ID          int64  `json:"id"`
					Lan         string `json:"lan"`
					LanDoc      string `json:"lan_doc"`
					SubtitleURL string `json:"subtitle_url"`
					AIStatus    int    `json:"ai_status"`
				} `json:"subtitles"`
			} `json:"subtitle"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		logger.ErrorLog("解析字幕列表失败", logger.Error(err))
		return "", false
	}

	if result.Code != 0 || len(result.Data.Subtitle.Subtitles) == 0 {
		logger.Info("视频无字幕", logger.String("bvid", videoInfo.BVID))
		return "", false
	}

	// 选择最佳字幕
	subtitle := c.selectBestSubtitle(result.Data.Subtitle.Subtitles)
	if subtitle == nil {
		return "", false
	}

	// 获取字幕内容
	text, err := c.fetchSubtitleContent(ctx, subtitle.SubtitleURL)
	if err != nil {
		logger.ErrorLog("获取字幕内容失败", logger.Error(err))
		return "", false
	}

	logger.Info("获取字幕成功",
		logger.String("bvid", videoInfo.BVID),
		logger.String("lan", subtitle.Lan),
		logger.Int("segments", len(strings.Split(text, "\n"))))

	return text, true
}

// selectBestSubtitle 选择最佳字幕
func (c *BilibiliClient) selectBestSubtitle(subtitles []struct {
	ID          int64  `json:"id"`
	Lan         string `json:"lan"`
	LanDoc      string `json:"lan_doc"`
	SubtitleURL string `json:"subtitle_url"`
	AIStatus    int    `json:"ai_status"`
}) *model.SubtitleMeta {
	if len(subtitles) == 0 {
		return nil
	}

	// 优先级映射
	priority := map[string]int{
		"zh-CN":   100, // 简体中文人工字幕
		"zh-Hans": 95,  // 简体中文
		"zh":      90,  // 中文
		"ai-zh":   80,  // AI中文
		"zh-TW":   70,  // 繁体中文
		"zh-Hant": 65,  // 繁体中文
	}

	type scoredSubtitle struct {
		meta  model.SubtitleMeta
		score int
	}

	var scored []scoredSubtitle
	for _, s := range subtitles {
		meta := model.SubtitleMeta{
			ID:          s.ID,
			Lan:         s.Lan,
			LanDoc:      s.LanDoc,
			SubtitleURL: s.SubtitleURL,
			AIStatus:    s.AIStatus,
		}

		score := priority[s.Lan]
		if score == 0 {
			score = 10 // 默认低分
		}

		// AI字幕降低优先级
		if s.AIStatus == 1 {
			score -= 20
		}

		scored = append(scored, scoredSubtitle{meta: meta, score: score})
	}

	// 按分数排序
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	best := scored[0].meta
	return &best
}

// fetchSubtitleContent 获取字幕内容
func (c *BilibiliClient) fetchSubtitleContent(ctx context.Context, subtitleURL string) (string, error) {
	resp, err := c.httpClient.Get(ctx, subtitleURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var subtitleFile model.SubtitleFile
	if err := json.NewDecoder(resp.Body).Decode(&subtitleFile); err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, seg := range subtitleFile.Body {
		text := strings.TrimSpace(seg.Content)
		if text == "" {
			continue
		}
		sb.WriteString(fmt.Sprintf("[%.2f-%.2f] %s\n", seg.From, seg.To, text))
	}

	return sb.String(), nil
}

// GetRecipeLikeComments 获取可能包含菜谱的评论
func (c *BilibiliClient) GetRecipeLikeComments(ctx context.Context, aid int64, limit int) ([]model.CommentItem, error) {
	if aid == 0 {
		return nil, fmt.Errorf("aid为空")
	}

	apiURL := fmt.Sprintf("%s/x/v2/reply/main?type=1&oid=%d&ps=20", c.baseURL, aid)

	resp, err := c.httpClient.Get(ctx, apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求评论失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取评论失败: %w", err)
	}

	var result struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Top struct {
				Upper struct {
					Content struct {
						Message string `json:"message"`
					} `json:"content"`
					Member struct {
						Uname string `json:"uname"`
					} `json:"member"`
					RPID int64 `json:"rpid"`
					Like int   `json:"like"`
				} `json:"upper"`
			} `json:"top"`
			Replies []struct {
				RPID    int64 `json:"rpid"`
				Like    int   `json:"like"`
				Content struct {
					Message string `json:"message"`
				} `json:"content"`
				Member struct {
					Uname string `json:"uname"`
				} `json:"member"`
			} `json:"replies"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析评论失败: %w", err)
	}

	if result.Code != 0 {
		return nil, fmt.Errorf("B站评论API错误: code=%d, msg=%s", result.Code, result.Message)
	}

	var comments []model.CommentItem

	// 添加置顶评论
	if result.Data.Top.Upper.Content.Message != "" {
		comments = append(comments, model.CommentItem{
			RPID:    result.Data.Top.Upper.RPID,
			Message: result.Data.Top.Upper.Content.Message,
			Like:    result.Data.Top.Upper.Like,
			IsTop:   true,
			Member:  result.Data.Top.Upper.Member.Uname,
		})
	}

	// 筛选高赞评论
	for _, reply := range result.Data.Replies {
		msg := reply.Content.Message
		if LooksLikeRecipeComment(msg) {
			comments = append(comments, model.CommentItem{
				RPID:    reply.RPID,
				Message: msg,
				Like:    reply.Like,
				IsTop:   false,
				Member:  reply.Member.Uname,
			})
		}

		if len(comments) >= limit+1 { // +1 因为可能包含置顶评论
			break
		}
	}

	logger.Info("获取评论成功",
		logger.Int64("aid", aid),
		logger.Int("count", len(comments)))

	return comments, nil
}

// LooksLikeRecipeComment 判断评论是否像菜谱
func LooksLikeRecipeComment(msg string) bool {
	keywords := []string{
		"食材", "材料", "配方", "步骤", "做法",
		"克", "g", "勺", "分钟", "腌", "焯水",
		"爆香", "小火", "大火", "生抽", "老抽",
		"蚝油", "盐", "糖", "油",
	}

	count := 0
	for _, kw := range keywords {
		if strings.Contains(msg, kw) {
			count++
		}
	}

	return count >= 2
}

func (c *BilibiliClient) IsValidURL(input string) bool {
	return parser.IsValidBilibiliURL(input)
}
