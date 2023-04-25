package main

// Generated via https://transform.tools/json-to-go
type YTDLPOut struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Formats []struct {
		FormatID   string  `json:"format_id"`
		FormatNote string  `json:"format_note"`
		Ext        string  `json:"ext"`
		Protocol   string  `json:"protocol"`
		Acodec     string  `json:"acodec"`
		Vcodec     string  `json:"vcodec"`
		URL        string  `json:"url"`
		Width      int     `json:"width"`
		Height     int     `json:"height"`
		Fps        float64 `json:"fps"`
		Rows       int     `json:"rows,omitempty"`
		Columns    int     `json:"columns,omitempty"`
		Fragments  []struct {
			URL      string  `json:"url"`
			Duration float64 `json:"duration"`
		} `json:"fragments,omitempty"`
		Resolution  string  `json:"resolution"`
		AspectRatio float64 `json:"aspect_ratio"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent"`
			Accept         string `json:"Accept"`
			AcceptLanguage string `json:"Accept-Language"`
			SecFetchMode   string `json:"Sec-Fetch-Mode"`
		} `json:"http_headers"`
		AudioExt           string      `json:"audio_ext"`
		VideoExt           string      `json:"video_ext"`
		Format             string      `json:"format"`
		Asr                int         `json:"asr,omitempty"`
		Filesize           int         `json:"filesize,omitempty"`
		SourcePreference   int         `json:"source_preference,omitempty"`
		AudioChannels      int         `json:"audio_channels,omitempty"`
		Quality            float64     `json:"quality,omitempty"`
		HasDrm             bool        `json:"has_drm,omitempty"`
		Tbr                float64     `json:"tbr,omitempty"`
		Language           interface{} `json:"language,omitempty"`
		LanguagePreference int         `json:"language_preference,omitempty"`
		Preference         interface{} `json:"preference,omitempty"`
		DynamicRange       interface{} `json:"dynamic_range,omitempty"`
		Abr                float64     `json:"abr,omitempty"`
		DownloaderOptions  struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options,omitempty"`
		Container      string  `json:"container,omitempty"`
		Vbr            float64 `json:"vbr,omitempty"`
		FilesizeApprox int     `json:"filesize_approx,omitempty"`
	} `json:"formats"`
	Thumbnails []struct {
		URL        string `json:"url"`
		Preference int    `json:"preference"`
		ID         string `json:"id"`
		Height     int    `json:"height,omitempty"`
		Width      int    `json:"width,omitempty"`
		Resolution string `json:"resolution,omitempty"`
	} `json:"thumbnails"`
	Thumbnail         string      `json:"thumbnail"`
	Description       string      `json:"description"`
	Uploader          string      `json:"uploader"`
	UploaderID        string      `json:"uploader_id"`
	UploaderURL       string      `json:"uploader_url"`
	ChannelID         string      `json:"channel_id"`
	ChannelURL        string      `json:"channel_url"`
	Duration          int         `json:"duration"`
	ViewCount         int         `json:"view_count"`
	AverageRating     interface{} `json:"average_rating"`
	AgeLimit          int         `json:"age_limit"`
	WebpageURL        string      `json:"webpage_url"`
	Categories        []string    `json:"categories"`
	Tags              []string    `json:"tags"`
	PlayableInEmbed   bool        `json:"playable_in_embed"`
	LiveStatus        string      `json:"live_status"`
	ReleaseTimestamp  interface{} `json:"release_timestamp"`
	FormatSortFields  []string    `json:"_format_sort_fields"`
	AutomaticCaptions struct {
		Af []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"af"`
		Ak []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ak"`
		Sq []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sq"`
		Am []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"am"`
		Ar []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ar"`
		Hy []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"hy"`
		As []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"as"`
		Ay []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ay"`
		Az []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"az"`
		Bn []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"bn"`
		Eu []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"eu"`
		Be []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"be"`
		Bho []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"bho"`
		Bs []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"bs"`
		Bg []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"bg"`
		My []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"my"`
		Ca []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ca"`
		Ceb []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ceb"`
		ZhHans []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"zh-Hans"`
		ZhHant []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"zh-Hant"`
		Co []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"co"`
		Hr []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"hr"`
		Cs []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"cs"`
		Da []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"da"`
		Dv []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"dv"`
		Nl []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"nl"`
		EnOrig []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"en-orig"`
		En []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"en"`
		Eo []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"eo"`
		Et []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"et"`
		Ee []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ee"`
		Fil []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"fil"`
		Fi []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"fi"`
		Fr []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"fr"`
		Gl []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"gl"`
		Lg []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"lg"`
		Ka []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ka"`
		De []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"de"`
		El []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"el"`
		Gn []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"gn"`
		Gu []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"gu"`
		Ht []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ht"`
		Ha []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ha"`
		Haw []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"haw"`
		Iw []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"iw"`
		Hi []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"hi"`
		Hmn []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"hmn"`
		Hu []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"hu"`
		Is []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"is"`
		Ig []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ig"`
		ID []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"id"`
		Ga []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ga"`
		It []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"it"`
		Ja []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ja"`
		Jv []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"jv"`
		Kn []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"kn"`
		Kk []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"kk"`
		Km []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"km"`
		Rw []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"rw"`
		Ko []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ko"`
		Kri []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"kri"`
		Ku []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ku"`
		Ky []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ky"`
		Lo []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"lo"`
		La []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"la"`
		Lv []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"lv"`
		Ln []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ln"`
		Lt []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"lt"`
		Lb []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"lb"`
		Mk []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"mk"`
		Mg []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"mg"`
		Ms []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ms"`
		Ml []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ml"`
		Mt []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"mt"`
		Mi []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"mi"`
		Mr []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"mr"`
		Mn []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"mn"`
		Ne []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ne"`
		Nso []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"nso"`
		No []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"no"`
		Ny []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ny"`
		Or []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"or"`
		Om []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"om"`
		Ps []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ps"`
		Fa []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"fa"`
		Pl []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"pl"`
		Pt []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"pt"`
		Pa []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"pa"`
		Qu []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"qu"`
		Ro []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ro"`
		Ru []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ru"`
		Sm []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sm"`
		Sa []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sa"`
		Gd []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"gd"`
		Sr []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sr"`
		Sn []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sn"`
		Sd []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sd"`
		Si []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"si"`
		Sk []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sk"`
		Sl []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sl"`
		So []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"so"`
		St []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"st"`
		Es []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"es"`
		Su []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"su"`
		Sw []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sw"`
		Sv []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"sv"`
		Tg []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"tg"`
		Ta []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ta"`
		Tt []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"tt"`
		Te []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"te"`
		Th []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"th"`
		Ti []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ti"`
		Ts []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ts"`
		Tr []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"tr"`
		Tk []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"tk"`
		Uk []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"uk"`
		Und []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"und"`
		Ur []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ur"`
		Ug []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"ug"`
		Uz []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"uz"`
		Vi []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"vi"`
		Cy []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"cy"`
		Fy []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"fy"`
		Xh []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"xh"`
		Yi []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"yi"`
		Yo []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"yo"`
		Zu []struct {
			Ext  string `json:"ext"`
			URL  string `json:"url"`
			Name string `json:"name"`
		} `json:"zu"`
	} `json:"automatic_captions"`
	Subtitles            struct{}    `json:"subtitles"`
	CommentCount         int         `json:"comment_count"`
	Chapters             interface{} `json:"chapters"`
	LikeCount            int         `json:"like_count"`
	Channel              string      `json:"channel"`
	ChannelFollowerCount int         `json:"channel_follower_count"`
	UploadDate           string      `json:"upload_date"`
	Availability         string      `json:"availability"`
	OriginalURL          string      `json:"original_url"`
	WebpageURLBasename   string      `json:"webpage_url_basename"`
	WebpageURLDomain     string      `json:"webpage_url_domain"`
	Extractor            string      `json:"extractor"`
	ExtractorKey         string      `json:"extractor_key"`
	Playlist             interface{} `json:"playlist"`
	PlaylistIndex        interface{} `json:"playlist_index"`
	DisplayID            string      `json:"display_id"`
	Fulltitle            string      `json:"fulltitle"`
	DurationString       string      `json:"duration_string"`
	IsLive               bool        `json:"is_live"`
	WasLive              bool        `json:"was_live"`
	RequestedSubtitles   interface{} `json:"requested_subtitles"`
	HasDrm               interface{} `json:"_has_drm"`
	RequestedDownloads   []struct {
		RequestedFormats []struct {
			Asr                interface{} `json:"asr"`
			Filesize           int         `json:"filesize"`
			FormatID           string      `json:"format_id"`
			FormatNote         string      `json:"format_note"`
			SourcePreference   int         `json:"source_preference"`
			Fps                float64     `json:"fps"`
			AudioChannels      interface{} `json:"audio_channels"`
			Height             int         `json:"height"`
			Quality            float64     `json:"quality"`
			HasDrm             bool        `json:"has_drm"`
			Tbr                float64     `json:"tbr"`
			URL                string      `json:"url"`
			Width              int         `json:"width"`
			Language           interface{} `json:"language"`
			LanguagePreference int         `json:"language_preference"`
			Preference         interface{} `json:"preference"`
			Ext                string      `json:"ext"`
			Vcodec             string      `json:"vcodec"`
			Acodec             string      `json:"acodec"`
			DynamicRange       string      `json:"dynamic_range"`
			Vbr                float64     `json:"vbr,omitempty"`
			DownloaderOptions  struct {
				HTTPChunkSize int `json:"http_chunk_size"`
			} `json:"downloader_options"`
			Container   string  `json:"container"`
			Protocol    string  `json:"protocol"`
			Resolution  string  `json:"resolution"`
			AspectRatio float64 `json:"aspect_ratio"`
			HTTPHeaders struct {
				UserAgent      string `json:"User-Agent"`
				Accept         string `json:"Accept"`
				AcceptLanguage string `json:"Accept-Language"`
				SecFetchMode   string `json:"Sec-Fetch-Mode"`
			} `json:"http_headers"`
			VideoExt string  `json:"video_ext"`
			AudioExt string  `json:"audio_ext"`
			Format   string  `json:"format"`
			Abr      float64 `json:"abr,omitempty"`
		} `json:"requested_formats"`
		Format               string  `json:"format"`
		FormatID             string  `json:"format_id"`
		Ext                  string  `json:"ext"`
		Protocol             string  `json:"protocol"`
		FormatNote           string  `json:"format_note"`
		FilesizeApprox       int     `json:"filesize_approx"`
		Tbr                  float64 `json:"tbr"`
		Width                int     `json:"width"`
		Height               int     `json:"height"`
		Resolution           string  `json:"resolution"`
		Fps                  float64 `json:"fps"`
		DynamicRange         string  `json:"dynamic_range"`
		Vcodec               string  `json:"vcodec"`
		Vbr                  float64 `json:"vbr"`
		AspectRatio          float64 `json:"aspect_ratio"`
		Acodec               string  `json:"acodec"`
		Abr                  float64 `json:"abr"`
		Asr                  int     `json:"asr"`
		AudioChannels        int     `json:"audio_channels"`
		Epoch                int     `json:"epoch"`
		Filename             string  `json:"_filename"`
		WriteDownloadArchive bool    `json:"__write_download_archive"`
	} `json:"requested_downloads"`
	RequestedFormats []struct {
		Asr                interface{} `json:"asr"`
		Filesize           int         `json:"filesize"`
		FormatID           string      `json:"format_id"`
		FormatNote         string      `json:"format_note"`
		SourcePreference   int         `json:"source_preference"`
		Fps                float64     `json:"fps"`
		AudioChannels      interface{} `json:"audio_channels"`
		Height             int         `json:"height"`
		Quality            float64     `json:"quality"`
		HasDrm             bool        `json:"has_drm"`
		Tbr                float64     `json:"tbr"`
		URL                string      `json:"url"`
		Width              int         `json:"width"`
		Language           interface{} `json:"language"`
		LanguagePreference int         `json:"language_preference"`
		Preference         interface{} `json:"preference"`
		Ext                string      `json:"ext"`
		Vcodec             string      `json:"vcodec"`
		Acodec             string      `json:"acodec"`
		DynamicRange       string      `json:"dynamic_range"`
		Vbr                float64     `json:"vbr,omitempty"`
		DownloaderOptions  struct {
			HTTPChunkSize int `json:"http_chunk_size"`
		} `json:"downloader_options"`
		Container   string  `json:"container"`
		Protocol    string  `json:"protocol"`
		Resolution  string  `json:"resolution"`
		AspectRatio float64 `json:"aspect_ratio"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent"`
			Accept         string `json:"Accept"`
			AcceptLanguage string `json:"Accept-Language"`
			SecFetchMode   string `json:"Sec-Fetch-Mode"`
		} `json:"http_headers"`
		VideoExt string  `json:"video_ext"`
		AudioExt string  `json:"audio_ext"`
		Format   string  `json:"format"`
		Abr      float64 `json:"abr,omitempty"`
	} `json:"requested_formats"`
	Format         string      `json:"format"`
	FormatID       string      `json:"format_id"`
	Ext            string      `json:"ext"`
	Protocol       string      `json:"protocol"`
	Language       interface{} `json:"language"`
	FormatNote     string      `json:"format_note"`
	FilesizeApprox int         `json:"filesize_approx"`
	Tbr            float64     `json:"tbr"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Resolution     string      `json:"resolution"`
	Fps            float64     `json:"fps"`
	DynamicRange   string      `json:"dynamic_range"`
	Vcodec         string      `json:"vcodec"`
	Vbr            float64     `json:"vbr"`
	StretchedRatio interface{} `json:"stretched_ratio"`
	AspectRatio    float64     `json:"aspect_ratio"`
	Acodec         string      `json:"acodec"`
	Abr            float64     `json:"abr"`
	Asr            int         `json:"asr"`
	AudioChannels  int         `json:"audio_channels"`
	Epoch          int         `json:"epoch"`
	Type           string      `json:"_type"`
	Version        struct {
		Version        string      `json:"version"`
		CurrentGitHead interface{} `json:"current_git_head"`
		ReleaseGitHead string      `json:"release_git_head"`
		Repository     string      `json:"repository"`
	} `json:"_version"`
}
