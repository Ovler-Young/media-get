package bilibili

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/foamzou/audio-get/args"
	"github.com/foamzou/audio-get/consts"
	"github.com/foamzou/audio-get/logger"
	"github.com/foamzou/audio-get/meta"
	"github.com/foamzou/audio-get/utils"
)

type Core struct {
	Opts *args.Options
}

const Album = "Bilibili"

func (c *Core) IsMusicPlatform() bool {
	return false
}

func (c *Core) Domains() []string {
	return []string{"bilibili.com", "b23.tv"}
}

func (c *Core) GetSourceName() string {
	return consts.SourceNameBilibili
}

func (c *Core) FetchMetaAndResourceInfo() (mediaMeta *meta.MediaMeta, err error) {
	if strings.Contains(c.Opts.Url, "b23.tv") {
		if redirectUrl, err := utils.GetLocation(consts.SourceNameBilibili, c.Opts.Url, map[string]string{
			"user-agent": consts.UAAndroid,
		}); err == nil {
			c.Opts.Url = redirectUrl
		}
	}
	html, err := fetchHtml(c.Opts.Url)
	if err != nil {
		return
	}

	// audio resource
	// audio resource
	matchStr, err := utils.RegexSingleMatch(html, `window.__playinfo__=(.+?)<\/script`)
	if err != nil {
		logger.Warn("fetch playinfo failed, try to fetch from api", err)
		return c.fetchBytesFromApi()
	}
	resource := &AudioResource{}
	err = json.Unmarshal([]byte(matchStr), resource)
	if err != nil {
		return
	}

	// audio meta
	mediaMeta = &meta.MediaMeta{
		Duration:     resource.Data.Dash.Duration,
		ResourceType: consts.ResourceTypeVideo,
	}
	metaJson := utils.RegexSingleMatchIgnoreError(html, `__INITIAL_STATE__=(.+?);\(function`, "{}")
	audioMeta := &AudioMeta{}
	err = json.Unmarshal([]byte(metaJson), audioMeta)
	if err != nil {
		logger.Warn("fetch meta json failed", err)
		mediaMeta.Title = utils.RegexSingleMatchIgnoreError(html, `<h1 title="(.+?)"`, utils.Md5(c.Opts.Url))
	} else {
		mediaMeta.Title = audioMeta.VideoData.Title
		mediaMeta.Description = audioMeta.VideoData.Description
	}
	if len(resource.Data.Dash.Audios) == 0 {
		return nil, errors.New("no audio data")
	}
	audio := resource.Data.Dash.Audios[0]
	mediaMeta.Artist = getSinger(audioMeta)
	mediaMeta.Album = Album
	mediaMeta.CoverUrl = audioMeta.VideoData.Pic
	mediaMeta.Headers = map[string]string{
		"user-agent": consts.UAMac,
		"referer":    c.Opts.Url,
	}
	mediaMeta.Audios = append(mediaMeta.Audios, meta.Audio{
		Url:     audio.BaseUrl,
		BitRate: consts.BitRate128,
	})

	for _, bilibiliVideo := range resource.Data.Dash.Videos {
		mediaMeta.Videos = append(mediaMeta.Videos, meta.Video{
			Url:            bilibiliVideo.BaseUrl,
			Width:          bilibiliVideo.Width,
			Height:         bilibiliVideo.Height,
			Ratio:          getRatioById(bilibiliVideo.Id),
			NeedExtraAudio: true,
		})
	}

	return mediaMeta, nil
}

func (c *Core) fetchBytesFromApi() (mediaMeta *meta.MediaMeta, err error) {
	bvid := utils.RegexSingleMatchIgnoreError(c.Opts.Url, `(BV[a-zA-Z0-9]+)`, "")
	if bvid == "" {
		return nil, errors.New("bvid not found")
	}

	// fetch meta
	metaUrl := "https://api.bilibili.com/x/web-interface/view?bvid=" + bvid
	metaJson, err := utils.HttpGet(consts.SourceNameBilibili, metaUrl, map[string]string{
		"user-agent": consts.UAMac,
	})
	if err != nil {
		return nil, err
	}
	var metaData BilibiliWebInterfaceView
	err = json.Unmarshal([]byte(metaJson), &metaData)
	if err != nil {
		return nil, err
	}
	if metaData.Code != 0 {
		return nil, errors.New(metaData.Message)
	}

	// fetch resource
	playUrl := "https://api.bilibili.com/x/player/playurl?bvid=" + bvid + "&cid=" + strconv.Itoa(metaData.Data.Cid) + "&qn=80&fnval=4048"
	playJson, err := utils.HttpGet(consts.SourceNameBilibili, playUrl, map[string]string{
		"user-agent": consts.UAMac,
		"referer":    c.Opts.Url,
	})
	if err != nil {
		return nil, err
	}
	var playData BilibiliPlayUrlResponse
	err = json.Unmarshal([]byte(playJson), &playData)
	if err != nil {
		return nil, err
	}
	if playData.Code != 0 {
		return nil, errors.New(playData.Message)
	}

	// construct meta
	mediaMeta = &meta.MediaMeta{
		Title:        metaData.Data.Title,
		Description:  metaData.Data.Desc,
		Duration:     playData.Data.Dash.Duration,
		CoverUrl:     metaData.Data.Pic,
		ResourceType: consts.ResourceTypeVideo,
		Album:        Album,
		Headers: map[string]string{
			"user-agent": consts.UAMac,
			"referer":    c.Opts.Url,
		},
	}

	// artist
	if len(metaData.Data.Staff) == 0 {
		mediaMeta.Artist = metaData.Data.Owner.Name
	} else {
		var names []string
		for _, staff := range metaData.Data.Staff {
			names = append(names, staff.Name)
		}
		mediaMeta.Artist = strings.Join(names, ", ")
	}

	// audios
	var bestAudio *struct {
		Id           int      `json:"id"`
		BaseUrl      string   `json:"baseUrl"`
		BackupUrl    []string `json:"backupUrl"`
		Bandwidth    int      `json:"bandwidth"`
		MimeType     string   `json:"mimeType"`
		Codecid      int      `json:"codecid"`
		Codecs       string   `json:"codecs"`
		Width        int      `json:"width"`
		Height       int      `json:"height"`
		FrameRate    string   `json:"frameRate"`
		Sar          string   `json:"sar"`
		StartWithSap int      `json:"startWithSap"`
		SegmentBase  struct {
			Initialization string `json:"Initialization"`
			IndexRange     string `json:"indexRange"`
		} `json:"SegmentBase"`
		Codecid2 int `json:"codecid"`
	}
	maxBandwidth := 0

	for i := range playData.Data.Dash.Audio {
		audio := &playData.Data.Dash.Audio[i]
		if audio.Bandwidth > maxBandwidth {
			maxBandwidth = audio.Bandwidth
			bestAudio = audio
		}
	}

	if bestAudio != nil {
		mediaMeta.Audios = append(mediaMeta.Audios, meta.Audio{
			Url:     bestAudio.BaseUrl,
			BitRate: bestAudio.Bandwidth / 1000,
		})
	}
	
	// videos
	for _, video := range playData.Data.Dash.Video {
		mediaMeta.Videos = append(mediaMeta.Videos, meta.Video{
			Url:            video.BaseUrl,
			Width:          video.Width,
			Height:         video.Height,
			Ratio:          getRatioById(video.Id),
			NeedExtraAudio: true,
		})
	}
	
	return mediaMeta, nil
}

func getRatioById(id int) string {
	switch id {
	case 16:
		return consts.Ratio360
	case 32:
		return consts.Ratio480
	case 64:
		return consts.Ratio720
	case 80:
		return consts.Ratio1080
	case 112:
		return consts.Ratio1080Plus
	default:
		return consts.RatioUnknown
	}
}

func getSinger(audioMeta *AudioMeta) string {
	var name string
	if len(audioMeta.VideoData.Staff) == 0 {
		name = audioMeta.VideoData.Owner.Name
	} else {
		var names []string
		for _, staff := range audioMeta.VideoData.Staff {
			names = append(names, staff.Name)
		}
		name = strings.Join(names, ", ")
	}
	if strings.TrimSpace(name) == "" {
		return "unknown"
	}

	return name
}

func fetchHtml(url string) (string, error) {
	return utils.HttpGet(consts.SourceNameBilibili, url, map[string]string{
		"user-agent": consts.UAMac,
	})
}
