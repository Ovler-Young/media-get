package bilibili

type AudioResource struct {
	Data struct {
		Dash struct {
			Duration int     `json:"duration"`
			Videos   []Video `json:"video"`
			Audios   []struct {
				BaseUrl string `json:"baseUrl"`
			} `json:"audio"`
		} `json:"dash"`
	} `json:"data"`
}

type Video struct {
	Id      int    `json:"id"`
	BaseUrl string `json:"baseUrl"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
}

type AudioMeta struct {
	VideoData struct {
		Title       string `json:"title"`
		Duration    int    `json:"duration"`
		Description string `json:"desc"`
		Pic         string `json:"pic"`
		Owner       struct {
			Mid  int    `json:"mid"`
			Name string `json:"name"`
			Face string `json:"face"`
		} `json:"owner"`
		Staff []struct {
			Mid   int    `json:"mid"`
			Title string `json:"title"`
			Name  string `json:"name"`
		} `json:"staff"`
	} `json:"videoData"`
}

type BilibiliWebInterfaceView struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Bvid  string `json:"bvid"`
		Aid   int    `json:"aid"`
		Cid   int    `json:"cid"`
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Pic   string `json:"pic"`
		Owner struct {
			Mid  int    `json:"mid"`
			Name string `json:"name"`
			Face string `json:"face"`
		} `json:"owner"`
		Staff []struct {
			Mid   int    `json:"mid"`
			Title string `json:"title"`
			Name  string `json:"name"`
		} `json:"staff"`
		Pages []struct {
			Cid      int    `json:"cid"`
			Page     int    `json:"page"`
			From     string `json:"from"`
			Part     string `json:"part"`
			Duration int    `json:"duration"`
			Vid      string `json:"vid"`
			Weblink  string `json:"weblink"`
			Dimension struct {
				Width  int `json:"width"`
				Height int `json:"height"`
				Rotate int `json:"rotate"`
			} `json:"dimension"`
		} `json:"pages"`
	} `json:"data"`
}

type BilibiliPlayUrlResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		From              string   `json:"from"`
		Result            string   `json:"result"`
		Message           string   `json:"message"`
		Quality           int      `json:"quality"`
		Format            string   `json:"format"`
		Timelength        int      `json:"timelength"`
		AcceptFormat      string   `json:"accept_format"`
		AcceptDescription []string `json:"accept_description"`
		AcceptQuality     []int    `json:"accept_quality"`
		VideoCodecid      int      `json:"video_codecid"`
		SeekParam         string   `json:"seek_param"`
		SeekType          string   `json:"seek_type"`
		Durl              []struct {
			Order     int    `json:"order"`
			Length    int    `json:"length"`
			Size      int    `json:"size"`
			Ahead     string `json:"ahead"`
			Vhead     string `json:"vhead"`
			Url       string `json:"url"`
			BackupUrl []string `json:"backup_url"`
		} `json:"durl"`
		Dash struct {
			Duration      int     `json:"duration"`
			MinBufferTime float64 `json:"minBufferTime"`
			MinBuffer     float64 `json:"min_buffer"`
			Video         []struct {
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
			} `json:"video"`
			Audio []struct {
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
			} `json:"audio"`
			Dolby struct {
				Type  int `json:"type"`
				Audio []struct {
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
				} `json:"audio"`
			} `json:"dolby"`
			Flac interface{} `json:"flac"`
		} `json:"dash"`
		SupportFormats []struct {
			Quality        int    `json:"quality"`
			Format         string `json:"format"`
			NewDescription string `json:"new_description"`
			DisplayDesc    string `json:"display_desc"`
			Superscript    string `json:"superscript"`
			Codecs         []interface{}  `json:"codecs"`
		} `json:"support_formats"`
		HighFormat interface{} `json:"high_format"`
		LastPlayTime int `json:"last_play_time"`
		LastPlayCid int `json:"last_play_cid"`
	} `json:"data"`
}
