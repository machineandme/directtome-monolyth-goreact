package storage

import (
	"net/url"
	"time"
)

// SystemData is model that store system fields
type SystemData struct {
	IsActive  bool
	IsDeleted bool
	Updated   time.Time
	Created   time.Time
}

// User data model
type User struct {
	SystemData // system fields

	BucketKey []byte // Name of personal bucket

	Name    string
	Email   string
	Phone   string
	Company string

	BillingPlan uint8
	Password    string
	Tokens      map[string]string

	Country   string
	City      string
	LastUA    string
	Visits    uint64
	LastVisit time.Time
	TimeZone  time.Location
}

// Domain data model
type Domain struct {
	SystemData // system fields

	BucketKey []byte // Name of domain bucket

	FeaturesKey []byte // Relation to Features

	Hostname    string
	NotFoundURL url.URL
}

// Redirect data model
type Redirect struct {
	SystemData // system fields

	StatisticKey []byte // Relation to Statistic

	FromURL url.URL
	ToURL   url.URL

	StatusCode   uint8
	CacheControl string
	CanonicalURL url.URL

	ExpireClicks uint64
	ExpireTime   time.Time
}

// Features data model that has flags to enable or disable features
type Features struct {
	FeatureAutoHTTPS bool
	FeatureGeoIP     bool

	FeatureCaptureUA       bool
	FeatureCaptureReferrer bool

	Additional map[string]bool
}

// Statistic data model that have fields filled with statistics
type Statistic struct {
	Clicks         uint64
	ExpiredClicks  uint64
	NotFoundClicks uint64

	SocialShareClicks uint64

	ExtendedStats map[string]map[string]uint64
}
