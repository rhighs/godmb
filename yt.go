package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

const (
	REGEXP_PREFIX string = `<meta name="title" content=,omitempty"`
	LOG_YTCMD     bool   = true
)

type YTDLPOut struct {
	ID      string `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	Formats []struct {
		FormatID   string  `json:"format_id,omitempty"`
		FormatNote string  `json:"format_note,omitempty"`
		Ext        string  `json:"ext,omitempty"`
		Protocol   string  `json:"protocol,omitempty"`
		Acodec     string  `json:"acodec,omitempty"`
		Vcodec     string  `json:"vcodec,omitempty"`
		URL        string  `json:"url,omitempty"`
		Width      int     `json:"width,omitempty"`
		Height     int     `json:"height,omitempty"`
		Fps        float64 `json:"fps,omitempty"`
		Rows       int     `json:"rows,omitempty"`
		Columns    int     `json:"columns,omitempty"`
		Fragments  []struct {
			URL      string  `json:"url,omitempty"`
			Duration float64 `json:"duration,omitempty"`
		} `json:"fragments,omitempty"`
		Resolution  string  `json:"resolution,omitempty"`
		AspectRatio float64 `json:"aspect_ratio,omitempty"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent,omitempty"`
			Accept         string `json:"Accept,omitempty"`
			AcceptLanguage string `json:"Accept-Language,omitempty"`
			SecFetchMode   string `json:"Sec-Fetch-Mode,omitempty"`
		} `json:"http_headers,omitempty"`
		AudioExt           string      `json:"audio_ext,omitempty"`
		VideoExt           string      `json:"video_ext,omitempty"`
		Format             string      `json:"format,omitempty"`
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
			HTTPChunkSize int `json:"http_chunk_size,omitempty"`
		} `json:"downloader_options,omitempty"`
		Container      string  `json:"container,omitempty"`
		Vbr            float64 `json:"vbr,omitempty"`
		FilesizeApprox int     `json:"filesize_approx,omitempty"`
	} `json:"formats,omitempty"`
	Thumbnails []struct {
		URL        string `json:"url,omitempty"`
		Preference int    `json:"preference,omitempty"`
		ID         string `json:"id,omitempty"`
		Height     int    `json:"height,omitempty"`
		Width      int    `json:"width,omitempty"`
		Resolution string `json:"resolution,omitempty"`
	} `json:"thumbnails,omitempty"`
	Thumbnail         string      `json:"thumbnail,omitempty"`
	Description       string      `json:"description,omitempty"`
	Uploader          string      `json:"uploader,omitempty"`
	UploaderID        string      `json:"uploader_id,omitempty"`
	UploaderURL       string      `json:"uploader_url,omitempty"`
	ChannelID         string      `json:"channel_id,omitempty"`
	ChannelURL        string      `json:"channel_url,omitempty"`
	Duration          int         `json:"duration,omitempty"`
	ViewCount         int         `json:"view_count,omitempty"`
	AverageRating     interface{} `json:"average_rating,omitempty"`
	AgeLimit          int         `json:"age_limit,omitempty"`
	WebpageURL        string      `json:"webpage_url,omitempty"`
	Categories        []string    `json:"categories,omitempty"`
	Tags              []string    `json:"tags,omitempty"`
	PlayableInEmbed   bool        `json:"playable_in_embed,omitempty"`
	LiveStatus        string      `json:"live_status,omitempty"`
	ReleaseTimestamp  interface{} `json:"release_timestamp,omitempty"`
	FormatSortFields  []string    `json:"_format_sort_fields,omitempty"`
	AutomaticCaptions struct {
		Af []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"af,omitempty"`
		Ak []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ak,omitempty"`
		Sq []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sq,omitempty"`
		Am []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"am,omitempty"`
		Ar []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ar,omitempty"`
		Hy []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"hy,omitempty"`
		As []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"as,omitempty"`
		Ay []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ay,omitempty"`
		Az []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"az,omitempty"`
		Bn []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"bn,omitempty"`
		Eu []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"eu,omitempty"`
		Be []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"be,omitempty"`
		Bho []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"bho,omitempty"`
		Bs []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"bs,omitempty"`
		Bg []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"bg,omitempty"`
		My []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"my,omitempty"`
		Ca []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ca,omitempty"`
		Ceb []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ceb,omitempty"`
		ZhHans []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"zh-Hans,omitempty"`
		ZhHant []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"zh-Hant,omitempty"`
		Co []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"co,omitempty"`
		Hr []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"hr,omitempty"`
		Cs []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"cs,omitempty"`
		Da []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"da,omitempty"`
		Dv []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"dv,omitempty"`
		Nl []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"nl,omitempty"`
		EnOrig []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"en-orig,omitempty"`
		En []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"en,omitempty"`
		Eo []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"eo,omitempty"`
		Et []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"et,omitempty"`
		Ee []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ee,omitempty"`
		Fil []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"fil,omitempty"`
		Fi []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"fi,omitempty"`
		Fr []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"fr,omitempty"`
		Gl []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"gl,omitempty"`
		Lg []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"lg,omitempty"`
		Ka []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ka,omitempty"`
		De []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"de,omitempty"`
		El []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"el,omitempty"`
		Gn []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"gn,omitempty"`
		Gu []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"gu,omitempty"`
		Ht []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ht,omitempty"`
		Ha []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ha,omitempty"`
		Haw []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"haw,omitempty"`
		Iw []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"iw,omitempty"`
		Hi []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"hi,omitempty"`
		Hmn []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"hmn,omitempty"`
		Hu []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"hu,omitempty"`
		Is []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"is,omitempty"`
		Ig []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ig,omitempty"`
		ID []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"id,omitempty"`
		Ga []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ga,omitempty"`
		It []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"it,omitempty"`
		Ja []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ja,omitempty"`
		Jv []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"jv,omitempty"`
		Kn []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"kn,omitempty"`
		Kk []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"kk,omitempty"`
		Km []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"km,omitempty"`
		Rw []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"rw,omitempty"`
		Ko []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ko,omitempty"`
		Kri []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"kri,omitempty"`
		Ku []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ku,omitempty"`
		Ky []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ky,omitempty"`
		Lo []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"lo,omitempty"`
		La []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"la,omitempty"`
		Lv []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"lv,omitempty"`
		Ln []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ln,omitempty"`
		Lt []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"lt,omitempty"`
		Lb []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"lb,omitempty"`
		Mk []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"mk,omitempty"`
		Mg []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"mg,omitempty"`
		Ms []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ms,omitempty"`
		Ml []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ml,omitempty"`
		Mt []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"mt,omitempty"`
		Mi []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"mi,omitempty"`
		Mr []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"mr,omitempty"`
		Mn []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"mn,omitempty"`
		Ne []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ne,omitempty"`
		Nso []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"nso,omitempty"`
		No []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"no,omitempty"`
		Ny []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ny,omitempty"`
		Or []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"or,omitempty"`
		Om []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"om,omitempty"`
		Ps []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ps,omitempty"`
		Fa []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"fa,omitempty"`
		Pl []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"pl,omitempty"`
		Pt []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"pt,omitempty"`
		Pa []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"pa,omitempty"`
		Qu []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"qu,omitempty"`
		Ro []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ro,omitempty"`
		Ru []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ru,omitempty"`
		Sm []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sm,omitempty"`
		Sa []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sa,omitempty"`
		Gd []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"gd,omitempty"`
		Sr []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sr,omitempty"`
		Sn []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sn,omitempty"`
		Sd []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sd,omitempty"`
		Si []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"si,omitempty"`
		Sk []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sk,omitempty"`
		Sl []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sl,omitempty"`
		So []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"so,omitempty"`
		St []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"st,omitempty"`
		Es []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"es,omitempty"`
		Su []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"su,omitempty"`
		Sw []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sw,omitempty"`
		Sv []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"sv,omitempty"`
		Tg []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"tg,omitempty"`
		Ta []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ta,omitempty"`
		Tt []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"tt,omitempty"`
		Te []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"te,omitempty"`
		Th []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"th,omitempty"`
		Ti []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ti,omitempty"`
		Ts []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ts,omitempty"`
		Tr []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"tr,omitempty"`
		Tk []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"tk,omitempty"`
		Uk []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"uk,omitempty"`
		Und []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"und,omitempty"`
		Ur []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ur,omitempty"`
		Ug []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"ug,omitempty"`
		Uz []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"uz,omitempty"`
		Vi []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"vi,omitempty"`
		Cy []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"cy,omitempty"`
		Fy []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"fy,omitempty"`
		Xh []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"xh,omitempty"`
		Yi []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"yi,omitempty"`
		Yo []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"yo,omitempty"`
		Zu []struct {
			Ext  string `json:"ext,omitempty"`
			URL  string `json:"url,omitempty"`
			Name string `json:"name,omitempty"`
		} `json:"zu,omitempty"`
	} `json:"automatic_captions,omitempty"`
	Subtitles            struct{}    `json:"subtitles,omitempty"`
	CommentCount         int         `json:"comment_count,omitempty"`
	Chapters             interface{} `json:"chapters,omitempty"`
	LikeCount            int         `json:"like_count,omitempty"`
	Channel              string      `json:"channel,omitempty"`
	ChannelFollowerCount int         `json:"channel_follower_count,omitempty"`
	UploadDate           string      `json:"upload_date,omitempty"`
	Availability         string      `json:"availability,omitempty"`
	OriginalURL          string      `json:"original_url,omitempty"`
	WebpageURLBasename   string      `json:"webpage_url_basename,omitempty"`
	WebpageURLDomain     string      `json:"webpage_url_domain,omitempty"`
	Extractor            string      `json:"extractor,omitempty"`
	ExtractorKey         string      `json:"extractor_key,omitempty"`
	Playlist             interface{} `json:"playlist,omitempty"`
	PlaylistIndex        interface{} `json:"playlist_index,omitempty"`
	DisplayID            string      `json:"display_id,omitempty"`
	Fulltitle            string      `json:"fulltitle,omitempty"`
	DurationString       string      `json:"duration_string,omitempty"`
	IsLive               bool        `json:"is_live,omitempty"`
	WasLive              bool        `json:"was_live,omitempty"`
	RequestedSubtitles   interface{} `json:"requested_subtitles,omitempty"`
	HasDrm               interface{} `json:"_has_drm,omitempty"`
	RequestedDownloads   []struct {
		RequestedFormats []struct {
			Asr                interface{} `json:"asr,omitempty"`
			Filesize           int         `json:"filesize,omitempty"`
			FormatID           string      `json:"format_id,omitempty"`
			FormatNote         string      `json:"format_note,omitempty"`
			SourcePreference   int         `json:"source_preference,omitempty"`
			Fps                float64     `json:"fps,omitempty"`
			AudioChannels      interface{} `json:"audio_channels,omitempty"`
			Height             int         `json:"height,omitempty"`
			Quality            float64     `json:"quality,omitempty"`
			HasDrm             bool        `json:"has_drm,omitempty"`
			Tbr                float64     `json:"tbr,omitempty"`
			URL                string      `json:"url,omitempty"`
			Width              int         `json:"width,omitempty"`
			Language           interface{} `json:"language,omitempty"`
			LanguagePreference int         `json:"language_preference,omitempty"`
			Preference         interface{} `json:"preference,omitempty"`
			Ext                string      `json:"ext,omitempty"`
			Vcodec             string      `json:"vcodec,omitempty"`
			Acodec             string      `json:"acodec,omitempty"`
			DynamicRange       string      `json:"dynamic_range,omitempty"`
			Vbr                float64     `json:"vbr,omitempty"`
			DownloaderOptions  struct {
				HTTPChunkSize int `json:"http_chunk_size,omitempty"`
			} `json:"downloader_options,omitempty"`
			Container   string  `json:"container,omitempty"`
			Protocol    string  `json:"protocol,omitempty"`
			Resolution  string  `json:"resolution,omitempty"`
			AspectRatio float64 `json:"aspect_ratio,omitempty"`
			HTTPHeaders struct {
				UserAgent      string `json:"User-Agent,omitempty"`
				Accept         string `json:"Accept,omitempty"`
				AcceptLanguage string `json:"Accept-Language,omitempty"`
				SecFetchMode   string `json:"Sec-Fetch-Mode,omitempty"`
			} `json:"http_headers,omitempty"`
			VideoExt string  `json:"video_ext,omitempty"`
			AudioExt string  `json:"audio_ext,omitempty"`
			Format   string  `json:"format,omitempty"`
			Abr      float64 `json:"abr,omitempty"`
		} `json:"requested_formats,omitempty"`
		Format               string  `json:"format,omitempty"`
		FormatID             string  `json:"format_id,omitempty"`
		Ext                  string  `json:"ext,omitempty"`
		Protocol             string  `json:"protocol,omitempty"`
		FormatNote           string  `json:"format_note,omitempty"`
		FilesizeApprox       int     `json:"filesize_approx,omitempty"`
		Tbr                  float64 `json:"tbr,omitempty"`
		Width                int     `json:"width,omitempty"`
		Height               int     `json:"height,omitempty"`
		Resolution           string  `json:"resolution,omitempty"`
		Fps                  float64 `json:"fps,omitempty"`
		DynamicRange         string  `json:"dynamic_range,omitempty"`
		Vcodec               string  `json:"vcodec,omitempty"`
		Vbr                  float64 `json:"vbr,omitempty"`
		AspectRatio          float64 `json:"aspect_ratio,omitempty"`
		Acodec               string  `json:"acodec,omitempty"`
		Abr                  float64 `json:"abr,omitempty"`
		Asr                  int     `json:"asr,omitempty"`
		AudioChannels        int     `json:"audio_channels,omitempty"`
		Epoch                int     `json:"epoch,omitempty"`
		Filename             string  `json:"_filename,omitempty"`
		WriteDownloadArchive bool    `json:"__write_download_archive,omitempty"`
	} `json:"requested_downloads,omitempty"`
	RequestedFormats []struct {
		Asr                interface{} `json:"asr,omitempty"`
		Filesize           int         `json:"filesize,omitempty"`
		FormatID           string      `json:"format_id,omitempty"`
		FormatNote         string      `json:"format_note,omitempty"`
		SourcePreference   int         `json:"source_preference,omitempty"`
		Fps                float64     `json:"fps,omitempty"`
		AudioChannels      interface{} `json:"audio_channels,omitempty"`
		Height             int         `json:"height,omitempty"`
		Quality            float64     `json:"quality,omitempty"`
		HasDrm             bool        `json:"has_drm,omitempty"`
		Tbr                float64     `json:"tbr,omitempty"`
		URL                string      `json:"url,omitempty"`
		Width              int         `json:"width,omitempty"`
		Language           interface{} `json:"language,omitempty"`
		LanguagePreference int         `json:"language_preference,omitempty"`
		Preference         interface{} `json:"preference,omitempty"`
		Ext                string      `json:"ext,omitempty"`
		Vcodec             string      `json:"vcodec,omitempty"`
		Acodec             string      `json:"acodec,omitempty"`
		DynamicRange       string      `json:"dynamic_range,omitempty"`
		Vbr                float64     `json:"vbr,omitempty"`
		DownloaderOptions  struct {
			HTTPChunkSize int `json:"http_chunk_size,omitempty"`
		} `json:"downloader_options,omitempty"`
		Container   string  `json:"container,omitempty"`
		Protocol    string  `json:"protocol,omitempty"`
		Resolution  string  `json:"resolution,omitempty"`
		AspectRatio float64 `json:"aspect_ratio,omitempty"`
		HTTPHeaders struct {
			UserAgent      string `json:"User-Agent,omitempty"`
			Accept         string `json:"Accept,omitempty"`
			AcceptLanguage string `json:"Accept-Language,omitempty"`
			SecFetchMode   string `json:"Sec-Fetch-Mode,omitempty"`
		} `json:"http_headers,omitempty"`
		VideoExt string  `json:"video_ext,omitempty"`
		AudioExt string  `json:"audio_ext,omitempty"`
		Format   string  `json:"format,omitempty"`
		Abr      float64 `json:"abr,omitempty"`
	} `json:"requested_formats,omitempty"`
	Format         string      `json:"format,omitempty"`
	FormatID       string      `json:"format_id,omitempty"`
	Ext            string      `json:"ext,omitempty"`
	Protocol       string      `json:"protocol,omitempty"`
	Language       interface{} `json:"language,omitempty"`
	FormatNote     string      `json:"format_note,omitempty"`
	FilesizeApprox int         `json:"filesize_approx,omitempty"`
	Tbr            float64     `json:"tbr,omitempty"`
	Width          int         `json:"width,omitempty"`
	Height         int         `json:"height,omitempty"`
	Resolution     string      `json:"resolution,omitempty"`
	Fps            float64     `json:"fps,omitempty"`
	DynamicRange   string      `json:"dynamic_range,omitempty"`
	Vcodec         string      `json:"vcodec,omitempty"`
	Vbr            float64     `json:"vbr,omitempty"`
	StretchedRatio interface{} `json:"stretched_ratio,omitempty"`
	AspectRatio    float64     `json:"aspect_ratio,omitempty"`
	Acodec         string      `json:"acodec,omitempty"`
	Abr            float64     `json:"abr,omitempty"`
	Asr            int         `json:"asr,omitempty"`
	AudioChannels  int         `json:"audio_channels,omitempty"`
	Epoch          int         `json:"epoch,omitempty"`
	Type           string      `json:"_type,omitempty"`
	Version        struct {
		Version        string      `json:"version,omitempty"`
		CurrentGitHead interface{} `json:"current_git_head,omitempty"`
		ReleaseGitHead string      `json:"release_git_head,omitempty"`
		Repository     string      `json:"repository,omitempty"`
	} `json:"_version,omitempty"`
}

var YTDLPPath string = ""

func SetYtdlpPath(path string) {
	YTDLPPath = path
	c := exec.Command(YTDLPPath)
	if err := c.Start(); err != nil {
		log.Fatalf("[COMMAND_ERR] yt-dlp start failed with reason: %v", err)
	}

	if err := c.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() != 2 {
			// Exit code 2 means an argument-less command invocation
			log.Fatalf("[COMMAND_ERR] yt-dlp exited with code %d Error: %s", exitErr.ExitCode(), exitErr.Error())
		}
	}
}

func GetYtdlpPath() string {
	return YTDLPPath
}

var ytPrefixes = []string{
	"https://www.youtu.be",
	"https://m.youtu.be",
	"https://youtu.be",
	"http://www.youtu.be",
	"http://m.youtu.be",
	"http://youtu.be",
	"https://www.youtube",
	"https://m.youtube",
	"https://youtube",
	"http://www.youtube",
	"http://m.youtube",
	"http://youtube",
}

var videoTitleRegexp *regexp.Regexp

func init() {
	sr := `<meta name="title" content=".*">`
	r, err := regexp.Compile(sr)
	if err != nil {
		panic("Couldn't compile regexp: " + sr)
	}
	videoTitleRegexp = r
}

func IsYoutubeUrl(url string) bool {
	for _, pre := range ytPrefixes {
		if strings.HasPrefix(url, pre) {
			return true
		}
	}
	return false
}

func GetYTVideoTitle(URL string) (string, error) {
	resp, _ := http.Get(URL)
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	matched := videoTitleRegexp.FindAllStringSubmatch(string(bodyBytes), 1)
	if len(matched[0]) == 0 {
		return "", errors.New("Could find any title at: " + URL)
	}
	str := matched[0][0]
	idx := strings.Index(str, ">")
	str = str[len(REGEXP_PREFIX) : idx-1]
	return str, nil
}

func YoutubeMediaUrl(videoUrl string) (string, error) {
	args := []string{
		"--dump-single-json",
		"--no-warnings",
		"--call-home",
		"--youtube-skip-dash-manifest",
		"--rm-cache-dir",
		videoUrl,
	}

	cmd := exec.Command(
		YTDLPPath,
		args...,
	)

	if LOG_YTCMD {
		log.Println("[YTDL_CMD_USED]:", YTDLPPath, strings.Join(args, " "))
	}

	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var ytdlOutput YTDLPOut
	if err := json.Unmarshal(stdout, &ytdlOutput); err != nil {
		log.Println(err)
		return "", err
	}

	for _, format := range ytdlOutput.Formats {
		if format.Vcodec == "none" && format.Acodec == "opus" {
			return format.URL, nil
		}
	}

	err = fmt.Errorf("no media url found")
	log.Println(err)
	return "", err
}

func UpdateYTDLP() {
    args := []string{
        "--update",
    }

    cmd := exec.Command(
        YTDLPPath,
        args...,
    )

	if LOG_YTCMD {
		log.Println("[YTDL_CMD_UPDATE]:", YTDLPPath, strings.Join(args, " "))
	}

    if err := cmd.Start(); err != nil {
		log.Println("[YTDL_CMD_UPDATE_ERR]:", YTDLPPath, strings.Join(args, " "))
    }

    if err := cmd.Wait(); err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            log.Println("[YTDL_CMD_UPDATE_EXITCODE]:",
                YTDLPPath, strings.Join(args, " "), "code: ", exitErr.ExitCode())
        }
    }
}

