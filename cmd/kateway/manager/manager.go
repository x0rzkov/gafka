package manager

import (
	"net/http"
)

// Manager is the interface that integrates with pubsub manager UI.
type Manager interface {
	// Name of the manager implementation.
	Name() string

	Start() error
	Stop()

	// AuthAdmin check if an app with the key has admin rights.
	AuthAdmin(appid, pubkey string) (ok bool)

	// OwnTopic checks if an appid owns a topic.
	OwnTopic(appid, pubkey, topic string) error

	AllowSubWithUnregisteredGroup(bool)

	// KafkaTopic returns raw kafka topic name.
	KafkaTopic(appid string, topic string, ver string) string

	// TopicSchema returns the avro schema definition json string.
	TopicSchema(appid, topic, ver string) (string, error)

	// ShadowTopic returns raw kafka topic name of a shadowed topic.
	ShadowTopic(shadow, myAppid, hisAppid, topic, ver, group string) string

	// AuthSub checks if an appid is able to consume message from hisAppid.hisTopic.
	AuthSub(appid, subkey, hisAppid, hisTopic, group string) error

	// LookupCluster locate the cluster name of an appid.
	LookupCluster(appid string) (cluster string, found bool)

	// WebHooks returns all registered webhooks object.
	WebHooks() ([]WebHook, error)

	// ForceRefresh will force manager to refresh the management data at once.
	ForceRefresh()

	// IsShadowedTopic checks if a topic has retry/dead sub/shadow topics.
	IsShadowedTopic(hisAppid, topic, ver, myAppid, group string) bool

	ValidateTopicName(topic string) bool
	ValidateGroupName(header http.Header, group string) bool

	DeadPartitions() map[string]map[int32]struct{}

	Dump() map[string]interface{}

	IsDryrunTopic(appid, topic, ver string) bool
	MarkTopicDryrun(appid, topic, ver string)
	ClearDryrunTopics()
}

var Default Manager
